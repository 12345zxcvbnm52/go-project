package inventorydata

import (
	"context"
	"fmt"
	"kenshop/pkg/errors"
	"kenshop/pkg/redlock"
	proto "kenshop/proto/inventory"
	model "kenshop/service/inventory/internal/model"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// InventoryDataService是提供Inventory底层相关数据操作的接口
type InventoryDataService interface {
	//后续可以考虑给每个库存添加id(同一个id则增加库存量,否则根据地址创建),address(地址)
	CreateStockDB(ctx context.Context, in *proto.CreateInventoryReq) (*emptypb.Empty, error)

	SetStockDB(ctx context.Context, in *proto.SetInventoryReq) (*emptypb.Empty, error)

	GetStockInfoDB(ctx context.Context, in *proto.InventoryInfoReq) (*proto.InventoryInfoRes, error)

	DecrStockDB(ctx context.Context, in *proto.UpdateStockReq) (*emptypb.Empty, error)

	IncrStockDB(ctx context.Context, in *proto.UpdateStockReq) (*emptypb.Empty, error)

	RebackStockDB(ctx context.Context, in *proto.RebackStockReq) (*emptypb.Empty, error)
}

// Inventory服务中的Data层,是数据操作的具体逻辑
type GormInventoryData struct {
	DB *gorm.DB
	//OpenTransaction bool
	Redlock *redlock.RedLock
	Logger  *otelzap.Logger
}

func MustNewGormInventoryData(db *gorm.DB, redlock *redlock.RedLock) *GormInventoryData {
	s := &GormInventoryData{DB: db, Redlock: redlock}
	return s
}

var _ InventoryDataService = (*GormInventoryData)(nil)

// func (s *GormInventoryData) WithTransaction(opts ...*sql.TxOptions) *GormInventoryData {
// 	return &GormInventoryData{
// 		OpenTransaction: true,
// 		DB:              s.DB.Begin(opts...),
// 	}
// }

// 后续可以考虑给每个库存添加id(同一个id则增加库存量,否则根据地址创建),address(地址)
func (s *GormInventoryData) CreateStockDB(ctx context.Context, in *proto.CreateInventoryReq) (*emptypb.Empty, error) {
	stock := &model.Inventory{}
	lock := s.Redlock.NewMutex(fmt.Sprintf("goods_id:%d", in.GoodsId))

	if err := lock.Lock(); err != nil {
		return nil, GormInventoryErrHandle(err)
	}

	defer func() {
		if _, err := lock.Unlock(); err != nil {
			s.Logger.Sugar().Errorf("释放锁失败, err= %+v", err)
		}
	}()

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			s.Logger.Sugar().Errorf("gorm遇到异常无法正常执行, err=%+v", r)
		}
	}()

	res := tx.Model(&model.Inventory{}).Where("goods_id = ?", in.GoodsId).First(stock)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			stock.GoodsId = in.GoodsId
			stock.GoodsNum = in.GoodsNum
			res = tx.Create(stock)
			if res.Error != nil {
				tx.Rollback()
				return nil, GormInventoryErrHandle(res.Error)
			}

			tx.Commit()
			return &emptypb.Empty{}, nil
		} else {
			tx.Rollback()
			return nil, GormInventoryErrHandle(res.Error)
		}
	}

	// 更新库存
	stock.GoodsNum += in.GoodsNum
	res = tx.Model(&model.Inventory{}).Where("goods_id = ?", stock.GoodsId).Updates(stock)
	if res.Error != nil {
		tx.Rollback()
		return nil, GormInventoryErrHandle(res.Error)
	}

	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *GormInventoryData) SetStockDB(ctx context.Context, in *proto.SetInventoryReq) (*emptypb.Empty, error) {
	stock := &model.Inventory{}
	res := s.DB.Model(&model.Inventory{}).Where("goods_id = ?", in.GoodsId).First(stock)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			stock.GoodsId = in.GoodsId
			stock.GoodsNum = in.GoodsNum
			res = s.DB.Model(&model.Inventory{}).Create(stock)
			if res.Error != nil {
				return nil, GormInventoryErrHandle(res.Error)
			}
		} else {
			return nil, GormInventoryErrHandle(res.Error)
		}
	}
	stock.GoodsNum = in.GoodsNum
	tx := s.DB.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			return
		}
	}()
	res = tx.Model(&model.Inventory{}).Where("goods_id = ?", stock.GoodsId).Updates(stock)
	if res.Error != nil {
		tx.Rollback()
		return nil, GormInventoryErrHandle(res.Error)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *GormInventoryData) GetStockInfoDB(ctx context.Context, in *proto.InventoryInfoReq) (*proto.InventoryInfoRes, error) {
	stock := &model.Inventory{GoodsId: in.GoodsId}
	res := s.DB.Model(&model.Inventory{}).Where("goods_id = ?", stock.GoodsId).First(stock)
	if res.Error != nil {
		return nil, GormInventoryErrHandle(res.Error)
	}
	return &proto.InventoryInfoRes{GoodsId: stock.GoodsId, GoodsNum: stock.GoodsNum}, nil
}

func (s *GormInventoryData) DecrStockDB(ctx context.Context, in *proto.UpdateStockReq) (*emptypb.Empty, error) {
	records := &model.OrderDecrRecord{
		OrderSign: in.OrderSign,
		Status:    model.StatusWaitingReback,
	}
	record := []*model.InventoryDecrRecord{}
	tx := s.DB.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			s.Logger.Sugar().Errorf("gorm遇到异常无法正常执行, err=%+v", err)
			return
		}
	}()

	for _, v := range in.DecrData {
		record = append(record, &model.InventoryDecrRecord{
			GoodsId:  v.GoodsId,
			GoodsNum: v.GoodsNum,
		})

		lock := s.Redlock.NewMutex(fmt.Sprintf("goods_id:%d", v.GoodsId))
		if err := lock.Lock(); err != nil {
			tx.Rollback()
			return nil, GormInventoryErrHandle(err)
		}

		stock := &model.Inventory{GoodsId: v.GoodsId}
		if err := tx.Model(&model.Inventory{}).Where("goods_id = ?", stock.GoodsId).First(stock).Error; err != nil {
			tx.Rollback()
			lock.Unlock()
			return nil, GormInventoryErrHandle(err)
		}

		if v.GoodsNum > stock.GoodsNum {
			tx.Rollback()
			lock.Unlock()
			err := errors.Errorf("库存不足,无法扣减,库存id为:%d,库存现有存量为:%d,欲扣减库存量为:%d", v.GoodsId, stock.GoodsNum, v.GoodsNum)
			return nil, GormInventoryErrHandle(errors.WithCoder(err, errors.CodeLackInventory, ""))
		}

		stock.GoodsNum -= v.GoodsNum
		res := tx.Model(&model.Inventory{}).Where("goods_id = ? ", stock.GoodsId).Update("goods_num", stock.GoodsNum)
		if res.Error != nil {
			tx.Rollback()
			lock.Unlock()
			return nil, GormInventoryErrHandle(res.Error)
		}

		if _, err := lock.Unlock(); err != nil {
			tx.Rollback()
			return nil, GormInventoryErrHandle(err)
		}

	}
	records.GoodsRecord = record
	if res := tx.Model(&model.OrderDecrRecord{}).Create(&records); res.Error != nil {
		tx.Rollback()
		return nil, GormInventoryErrHandle(res.Error)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *GormInventoryData) IncrStockDB(ctx context.Context, in *proto.UpdateStockReq) (*emptypb.Empty, error) {
	tx := s.DB.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			s.Logger.Sugar().Errorf("gorm遇到异常无法正常执行, err=%+v", err)
			return
		}
	}()
	for _, v := range in.DecrData {
		stock := &model.Inventory{}
		res := tx.Model(&model.Inventory{}).Where("goods_id = ?", v.GoodsId).First(stock)
		if res.Error != nil {
			tx.Rollback()
			return nil, GormInventoryErrHandle(res.Error)
		}
		stock.GoodsNum += v.GoodsNum
		res = tx.Model(&model.Inventory{}).Where("goods_id = ?", v.GoodsId).Update("goods_num", stock.GoodsNum)
		if res.Error != nil {
			tx.Rollback()
			return nil, GormInventoryErrHandle(res.Error)
		}
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

// 库存归还
func (s *GormInventoryData) RebackStockDB(ctx context.Context, in *proto.RebackStockReq) (*emptypb.Empty, error) {
	u := &model.OrderDecrRecord{OrderSign: in.OrderSign}
	tx := s.DB.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			s.Logger.Sugar().Errorf("gorm遇到异常无法正常执行, err=%+v", err)
			return
		}
	}()

	//因为只有内部会查订单扣减库存信息表,所以不会有严重的并发问题
	res := tx.Model(&model.OrderDecrRecord{}).Where("order_sign = ?", in.OrderSign).First(u)
	if res.Error != nil {
		tx.Rollback()
		return nil, GormInventoryErrHandle(res.Error)
	}
	if u.Status == model.StatusAlreadyReback {
		tx.Commit()
		return &emptypb.Empty{}, nil
	}
	for _, v := range u.GoodsRecord {
		lock := s.Redlock.NewMutex(fmt.Sprintf("goods_id:%d", v.GoodsId))
		if err := lock.Lock(); err != nil {
			tx.Rollback()
			return nil, GormInventoryErrHandle(err)
		}

		stock := &model.Inventory{GoodsId: v.GoodsId}
		res = tx.Model(&model.Inventory{}).Where("goods_id = ?", v.GoodsId).First(stock)
		if res.Error != nil {
			tx.Rollback()
			//暂时先不考虑锁释放失败和事务rollback完全失败(在分布式中概率微小)
			lock.Unlock()
			return nil, GormInventoryErrHandle(res.Error)
		}
		stock.GoodsNum += v.GoodsNum
		res = tx.Model(&model.Inventory{}).Where("goods_id = ?", v.GoodsId).Updates(stock)
		if res.Error != nil {
			tx.Rollback()
			lock.Unlock()
			return nil, GormInventoryErrHandle(res.Error)
		}
		lock.Unlock()
	}
	u.Status = model.StatusAlreadyReback
	res = tx.Model(&model.OrderDecrRecord{}).Where("order_sign = ?", in.OrderSign).Updates(u)
	if res.Error != nil {
		tx.Rollback()
		return nil, GormInventoryErrHandle(res.Error)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *GormInventoryData) Reback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	orderSign := ""
	for _, msg := range msgs {
		//用于读出OrderSign
		orderSign = string(msg.Body)
		tx := s.DB.Begin()
		defer func() {
			if err := recover(); err != nil {
				tx.Rollback()
				s.Logger.Sugar().Errorf("gorm遇到异常无法正常执行, err=%+v", err)
				return
			}
		}()

		records := &model.OrderDecrRecord{}
		records.OrderSign = orderSign
		//把等待Reback的记录查出来并且在最后Reback
		res := tx.Model(&model.OrderDecrRecord{}).
			Where("order_sign = ? and status = ?", orderSign, model.StatusWaitingReback).First(records)
		if res.Error != nil {
			return consumer.ConsumeSuccess, nil
		}

		for _, v := range records.GoodsRecord {
			res := tx.Model(&model.Inventory{}).Where("goods_id = ?", v.GoodsId).
				Update("goods_num", gorm.Expr("goods_num + ?", v.GoodsNum))
			if res.Error != nil {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}
		//Reback后把这条记录改为已经Reback
		if res := tx.Model(&model.OrderDecrRecord{}).Where("order_sign = ?", orderSign).
			Update("status", model.StatusAlreadyReback); res.Error != nil {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
	}
	return consumer.ConsumeSuccess, nil
}
