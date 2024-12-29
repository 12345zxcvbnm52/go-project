package handler

import (
	"context"
	"encoding/json"
	"fmt"
	gb "inventory_srv/global"
	"inventory_srv/model"
	pb "inventory_srv/proto"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type InventoryServer struct {
	pb.UnimplementedInventoryServer
}

func (is *InventoryServer) SetStock(ctx context.Context, req *pb.WriteInvtReq) (*emptypb.Empty, error) {
	u := &model.Inventory{
		GoodsId:  req.GoodsId,
		GoodsNum: req.GoodsNum,
	}
	if err := u.InsertOne(); err != nil {
		if err == model.ErrDuplicated {
			err = u.UpdateOneByGoodsId()
		}
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// func (is *InventoryServer) FindStock(ctx context.Context, req *pb.StockInfoReq) (*pb.StockInfoRes, error) {
// 	u := &model.Inventory{
// 		GoodsId: req.GoodsId,
// 	}
// 	if err := u.FindOneByGoodsId(); err != nil {
// 		return nil, err
// 	}
// 	return &pb.StockInfoRes{
// 		GoodsId:  u.GoodsId,
// 		GoodsNum: u.GoodsNum,
// 	}, nil
// }

// 这里的两个锁还没有排除连接出错的可能,可以考虑后期再封装一层或循环tryLock
func (is *InventoryServer) DecrStock(ctx context.Context, req *pb.DecrStockReq) (*emptypb.Empty, error) {
	return decrStockByLock(req)
}

func (is *InventoryServer) IncrStock(ctx context.Context, req *pb.IncrStockReq) (*emptypb.Empty, error) {
	tx := gb.DB.Begin()
	for _, v := range req.IncrData {
		inventory := &model.Inventory{}
		if err := inventory.FindOneByGoodsId(); err != nil {
			tx.Rollback()
			return nil, err
		}
		inventory.GoodsNum += v.GoodsNum
		if res := tx.Where(&model.Inventory{GoodsId: inventory.GoodsId}).Update("goods_num", inventory.GoodsNum); res.Error != nil {
			tx.Rollback()
			return nil, model.ErrInternalWrong
		}
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

// 这里的逻辑是残缺的:没有考虑到失败时应当重新减少库存
// func decrStockByLockByOptLock(ctx context.Context, req *pb.DecrStockReq) (*emptypb.Empty, error) {
// 	tx := gb.DB.Begin()
// 	ids := make([]int32, len(req.DecrData))
// 	for i, v := range req.DecrData {
// 		ids[i] = v.GoodsId
// 	}
// 	logic := &model.Inventory{}
// 	res, err := logic.FindByGoodsIds(ids...)
// 	if err != nil || int(res.Total) != len(req.DecrData) {
// 		tx.Rollback()
// 		return nil, err
// 	}
// 	for i, v := range res.Data {
// 		if v.GoodsNum < req.DecrData[i].GoodsNum {
// 			zap.S().Errorw("一次事务执行失败", "msg", "库存不足")
// 			tx.Rollback()
// 			return nil, model.ErrInvalidArgument
// 		}
// 		v.Version++
// 		v.GoodsNum -= req.DecrData[i].GoodsNum
// 		r := tx.Model(&model.Inventory{}).Select("goods_num", "version").
// 			Where("goods_id = ? AND version = ?", v.GoodsId, v.Version).Update("goods_num", v.GoodsNum)
// 		if r.Error != nil {
// 			zap.S().Errorw("一次事务执行失败", "msg", r.Error.Error())
// 			tx.Rollback()
// 			return nil, model.ErrInternalWrong
// 		}
// 	}
// 	tx.Commit()
// 	return &emptypb.Empty{}, nil
// }

func decrStockByLock(req *pb.DecrStockReq) (*emptypb.Empty, error) {
	records := &model.InvtRecord{
		OrderSign: req.OrderSign,
		Status:    model.StatusWaitingReback,
	}
	record := []model.GoodsRecord{}
	tx := gb.DB.Begin()
	for _, v := range req.DecrData {
		record = append(record, model.GoodsRecord{
			GoodsId:  v.GoodsId,
			GoodsNum: v.GoodsNum,
		})
		mtx := gb.RedLock.NewMutex(fmt.Sprintf("goods_id:%d", v.GoodsId))
		if err := mtx.Lock(); err != nil {
			zap.S().Errorw("获取分布式锁出错", "msg", err.Error())
			tx.Rollback()
			return nil, model.ErrInternalWrong
		}

		u := &model.Inventory{GoodsId: v.GoodsId}
		//因为这里使用了分布式锁,直接用DB查也没问题
		if err := u.FindOneByGoodsId(); err != nil {
			tx.Rollback()
			mtx.Unlock()
			return nil, err
		}

		if v.GoodsNum > u.GoodsNum {
			zap.S().Errorw("一次事务执行失败", "msg", "库存不足")
			tx.Rollback()
			mtx.Unlock()
			return nil, model.ErrLackInventory
		}
		u.GoodsNum -= v.GoodsNum
		if r := tx.Model(&model.Inventory{}).Where("goods_id = ? ", u.GoodsId).Update("goods_num", u.GoodsNum); r.Error != nil {
			zap.S().Errorw("一次事务执行失败", "msg", r.Error.Error())
			tx.Rollback()
			mtx.Unlock()
			return nil, model.ErrInternalWrong
		}

		if ok, err := mtx.Unlock(); !ok || err != nil {
			zap.S().Errorw("释放redis分布式锁异常", "msg", err.Error())
			tx.Rollback()
			return nil, model.ErrInternalWrong
		}
	}
	records.GoodsRecord = record
	if res := tx.Model(&model.InvtRecord{}).Create(&records); res.Error != nil {
		var err error
		if res.Error == gorm.ErrDuplicatedKey {
			err = model.ErrDuplicated
		} else {
			err = model.ErrInternalWrong
		}
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

// 这个函数的错误处理还不完整
// 这里是使用的sql数据库记录的Reback记录,如果为了追求速度可以使用redis记录Reback记录,(但是提升也不大就是了)
func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderSignReader struct {
		OrderSign string
	}
	for i := range msgs {
		//用于读出OrderSign
		order := &OrderSignReader{}
		json.Unmarshal(msgs[i].Body, order)
		tx := gb.DB.Begin()
		records := &model.InvtRecord{}
		//把等待Reback的记录查出来并且在最后Reback
		if res := tx.Model(&model.InvtRecord{}).
			Where("order_sign = ? and status = ?", order.OrderSign, model.StatusWaitingReback).First(records); res.Error != nil {
			return consumer.ConsumeSuccess, nil
		}
		for _, v := range records.GoodsRecord {
			res := tx.Model(&model.Inventory{}).Where("order_sign = ?", order.OrderSign).Update("goods_num", gorm.Expr("goods_num + ?", v.GoodsNum))
			if res.Error != nil {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}
		//Reback后把这条记录改为已经Reback
		if res := tx.Model(&model.Inventory{}).Where("order_sign = ?", order.OrderSign).Update("status", model.StatusAlreadyReback); res.Error != nil {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
	}
	return consumer.ConsumeSuccess, nil
}
