package handler

import (
	"context"
	"fmt"
	gb "inventory_srv/global"
	"inventory_srv/model"
	pb "inventory_srv/proto"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (is *InventoryServer) FindStock(ctx context.Context, req *pb.StockInfoReq) (*pb.StockInfoRes, error) {
	u := &model.Inventory{
		GoodsId: req.GoodsId,
	}
	if err := u.FindOneByGoodsId(); err != nil {
		return nil, err
	}
	return &pb.StockInfoRes{
		GoodsId:  u.GoodsId,
		GoodsNum: u.GoodsNum,
	}, nil
}

// 这里的两个锁还没有排除连接出错的可能,可以考虑后期再封装一层或循环tryLock
func (is *InventoryServer) DecrStock(ctx context.Context, req *pb.DecrStockReq) (*emptypb.Empty, error) {
	return decrStockByLock(ctx, req)
}

func (is *InventoryServer) IncrStock(ctx context.Context, req *pb.IncrStockReq) (*emptypb.Empty, error) {
	// tx := gb.DB.Begin()
	// ids := make([]int32, len(req.IncrData))
	// for i, v := range req.IncrData {
	// 	ids[i] = v.GoodsId
	// }
	// logic := &model.Inventory{}
	// if err := gb.RedLock.Lock(); err != nil {
	// 	zap.S().Errorw("获取分布式锁出错", "msg", err.Error())
	// }
	// res, err := logic.FindByGoodsIds(ids...)
	// if err != nil || int(res.Total) != len(req.IncrData) {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	// if false {
	// 	for i, v := range res.Data {
	// 		v.Version++
	// 		v.GoodsNum += req.IncrData[i].GoodsNum
	// 		r := tx.Model(&model.Inventory{}).Select("goods_num", "version").
	// 			Where("goods_id = ? AND version = ?", v.GoodsId, v.Version).Update("goods_num", v.GoodsNum)
	// 		if r.Error != nil {
	// 			zap.S().Errorw("一次事务执行失败", "msg", r.Error.Error())
	// 			tx.Rollback()
	// 			return nil, model.ErrInternalWrong
	// 		}
	// 	}
	// } else {
	// 	for i, v := range res.Data {
	// 		v.Version++
	// 		v.GoodsNum += req.IncrData[i].GoodsNum
	// 		r := tx.Model(&model.Inventory{}).Select("goods_num").Where("goods_id = ?", v.GoodsId).Update("goods_num", v.GoodsNum)
	// 		if r.Error != nil {
	// 			zap.S().Errorw("一次事务执行失败", "msg", r.Error.Error())
	// 			tx.Rollback()
	// 			gb.RedLock.Unlock()
	// 			return nil, model.ErrInternalWrong
	// 		}
	// 	}
	// 	gb.RedLock.Unlock()
	// }
	// tx.Commit()
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

func decrStockByLock(ctx context.Context, req *pb.DecrStockReq) (*emptypb.Empty, error) {
	tx := gb.DB.Begin()
	for _, v := range req.DecrData {
		mtx := gb.RedLock.NewMutex(fmt.Sprintf("goods_id:%d", v.GoodsId))
		if err := mtx.Lock(); err != nil {
			zap.S().Errorw("获取分布式锁出错", "msg", err.Error())
		}
		u := &model.Inventory{GoodsId: v.GoodsId}

		if err := u.FindOneByGoodsId(); err != nil {
			tx.Rollback()
			mtx.Unlock()
			return nil, err
		}

		if v.GoodsNum > u.GoodsNum {
			zap.S().Errorw("一次事务执行失败", "msg", "库存不足")
			tx.Rollback()
			mtx.Unlock()
			return nil, model.ErrInvalidArgument
		}
		u.GoodsNum -= v.GoodsNum
		r := tx.Model(&model.Inventory{}).Select("goods_num").Where("goods_id = ? ", u.GoodsId).Update("goods_num", u.GoodsNum)

		if r.Error != nil {
			zap.S().Errorw("一次事务执行失败", "msg", r.Error.Error())
			tx.Rollback()
			mtx.Unlock()
			return nil, model.ErrInternalWrong
		}
		tx.Commit()
		mtx.Unlock()
	}
	return &emptypb.Empty{}, nil
}
