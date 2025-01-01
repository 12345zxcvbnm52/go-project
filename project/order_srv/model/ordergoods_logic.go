package model

import (
	gb "order_srv/global"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderGoodsResult struct {
	Data  []*OrderGoods
	Total int64
}

func (u *OrderGoods) FindByOrderId(orderId uint32, Tx ...*gorm.DB) (*OrderGoodsResult, error) {
	LocDB := gb.DB
	res := &OrderGoodsResult{}
	if len(Tx) >= 1 {
		LocDB = Tx[0]
	}
	r := LocDB.Model(&OrderGoods{}).Where("order_id = ?", orderId).Find(&res.Data)
	if r.RowsAffected == 0 {
		zap.S().Errorw("订单内商品详细内容查找失败", "msg", "该订单不存在对应的商品信息")
		return nil, ErrOrderGoodsNotFound
	}
	if r.Error != nil {
		zap.S().Errorw("订单内商品详细内容查找失败", "msg", r.Error.Error())
		return nil, ErrInternalWrong

	}
	res.Total = int64(len(res.Data))
	return res, nil
}

func (u *OrderGoods) Insert(orderGoods []*OrderGoods, tx ...*gorm.DB) error {
	LocDB := gb.DB
	if len(tx) >= 1 {
		LocDB = tx[0]
	}
	if len(orderGoods) == 0 {
		return ErrInvalidArgument
	}
	if res := LocDB.Model(&OrderGoods{}).CreateInBatches(orderGoods, 100); res.Error != nil {
		zap.S().Errorw("批量插入订单内商品失败", "msg", res.Error.Error())
		if res.Error == gorm.ErrDuplicatedKey {
			return ErrOrderGoodsDuplicated
		}
		return ErrInternalWrong
	}
	return nil
}
