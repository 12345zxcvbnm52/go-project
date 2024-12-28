package model

import (
	gb "order_srv/global"
	"time"

	"gorm.io/gorm"
)

type OrderResult struct {
	Data  []*Order
	Total int64
}

type OrderFindOption struct {
	AfterTime  time.Time
	BeforeTime time.Time
	PagesNum   int32
	PageSize   int32
}

func (u *Order) FindByOpt(opt *OrderFindOption) (*OrderResult, error) {
	res := &OrderResult{}
	LocDB := gb.DB.Model(&Order{})
	if !opt.AfterTime.IsZero() {
		LocDB = LocDB.Where("created_at >= ?", opt.AfterTime)
	}
	if !opt.BeforeTime.IsZero() {
		LocDB = LocDB.Where("created_at <= ?", opt.BeforeTime)
	}
	LocDB.Count(&res.Total)
	if opt.PageSize > 0 {
		LocDB = LocDB.Scopes(Paginate(int(opt.PagesNum), int(opt.PageSize)))
	}
	r := LocDB.Find(res.Data)
	if r.Error != nil {
		return nil, ErrInternalWrong
	}
	return res, nil
}

func (u *Order) InsertOne() error {
	if res := gb.DB.Model(&Order{}).Create(u); res.Error != nil {
		if res.Error == gorm.ErrDuplicatedKey {
			return ErrOrderDuplicated
		}
		return ErrInternalWrong
	}
	return nil
}

// 隐式
func (u *Order) FindOne() error {
	if res := gb.DB.Model(&Order{}).Where("user_id = ?", u.UserId).First(u, u.ID); res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return ErrOrderNotFound
		}
		return ErrInternalWrong
	}
	return nil
}
