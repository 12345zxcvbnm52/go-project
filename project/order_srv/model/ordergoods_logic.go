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

func (u *OrderGoods) FindByOrderId() (*OrderGoodsResult, error) {
	res := &OrderGoodsResult{}
	LocDB := gb.DB.Where("order_id = ?", u.OrderId)
	LocDB.Count(&res.Total)
	r := LocDB.Find(res.Data)
	if r.RowsAffected == 0 {
		return nil, ErrOrderGoodsNotFound
	}
	if r.Error != nil {
		return nil, ErrInternalWrong
	}
	return res, nil
}

// 这个代码完全不安全,应该考虑事务
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
		if len(tx) >= 1 {
			LocDB.Rollback()
		}
		if res.Error == gorm.ErrDuplicatedKey {
			return ErrOrderGoodsDuplicated
		}
		return ErrInternalWrong
	}
	return nil
}
