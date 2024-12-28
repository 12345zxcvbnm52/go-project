package model

import (
	gb "order_srv/global"

	"go.uber.org/zap"
)

type OrderGoodsResult struct {
	Data  []*OrderGoods
	Total int64
}

func (u *OrderGoods) FindByOrderId(orderId uint32) (*OrderGoodsResult, error) {
	res := &OrderGoodsResult{}
	LocDB := gb.DB.Where("order_id = ?", u.OrderId)
	LocDB.Count(&res.Total)
	if res := LocDB.Find(res.Data); res.Error != nil {
		return nil, ErrInternalWrong
	}
	return res, nil
}

// 这个代码完全不安全,应该考虑事务
func (u *OrderGoods) Insert(orderGoods []*OrderGoods) error {
	end := len(orderGoods)
	if end > 100 {
		for i := 0; i < end; i += 100 {
			var j int = 100
			if end-i <= 100 {
				j = end - i
			}
			if res := gb.DB.Model(&OrderGoods{}).CreateInBatches(orderGoods[i:], j); res.Error != nil {
				zap.S().Errorw("批量插入失败", "msg", res.Error.Error())
			}
		}
	}
	return nil
}
