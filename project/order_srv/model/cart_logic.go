package model

import (
	gb "order_srv/global"

	"gorm.io/gorm"
)

type CartResult struct {
	Data  []*Cart
	Total int64
}

// 可以添加排序选项
type CartFindOption struct {
	UserId   uint32
	Selected bool
	PagesNum int32
	PageSize int32
}

func (u *Cart) FindByOpt(opt *CartFindOption) (*CartResult, error) {
	LocDB := gb.DB.Model(&Cart{}).Where("user_id = ?", opt.UserId)
	if opt.Selected {
		LocDB = LocDB.Where("selected = ?", 1)
	}
	res := &CartResult{}
	LocDB.Count(&res.Total)
	if opt.PageSize > 0 {
		LocDB = LocDB.Scopes(Paginate(int(opt.PagesNum), int(opt.PageSize)))
	}
	r := LocDB.Find(res.Data)
	if r.RowsAffected == 0 {
		return nil, ErrCartNotFound
	}
	if r.Error != nil {
		return nil, ErrInternalWrong
	}
	return res, nil
}

func (u *Cart) FindByUserId(userId uint32) (*CartResult, error) {
	item := []*Cart{}
	res := gb.DB.Where("user_id = ?", userId).Omit("user_id").Find(item)
	if res.RowsAffected == 0 {
		return nil, ErrCartNotFound
	}
	if res.Error != nil {
		return nil, ErrInternalWrong
	}
	return &CartResult{
		Total: res.RowsAffected,
		Data:  item,
	}, nil
}

func (u *Cart) FindOneById() error {
	res := gb.DB.Model(&Cart{}).First(u, u.ID)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return ErrCartNotFound
		}
		return ErrInternalWrong
	}
	return nil
}

// 通过用户和商品id来找到购物车的一个商品
func (u *Cart) FindOneByUserGoodsIds() error {
	res := gb.DB.Model(&Cart{}).Where("goods_id = ? AND user_id = ?", u.GoodsId, u.UserId).Find(u)
	if res.RowsAffected == 0 {
		return ErrCartNotFound
	}
	if res.Error != nil {
		return ErrInternalWrong
	}
	return nil
}

func (u *Cart) UpdateOneById() error {
	res := gb.DB.Model(&Cart{}).Updates(u)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return ErrCartNotFound
		}
		return ErrInternalWrong
	}
	return nil
}

func (u *Cart) InsertOne() error {
	model := &Cart{}
	if err := model.FindOneByUserGoodsIds(); err != nil {
		if err == ErrCartNotFound {
			res := gb.DB.Model(&Cart{}).Create(u)
			if res.Error != nil {
				return ErrInternalWrong
			}
			return nil
		} else {
			return err
		}
	}
	u.GoodsNums += model.GoodsNums
	u.Selected = model.Selected || u.Selected
	err := u.UpdateOneById()
	return err
}

func (u *Cart) DeleteOneByUserGoodsIds() error {
	res := gb.DB.Where("goods_id = ? AND user_id = ?", u.GoodsId, u.UserId).Delete(&Cart{})
	if res.RowsAffected == 0 {
		return ErrCartNotFound
	}
	if res.Error != nil {
		return ErrInternalWrong
	}
	return nil
}
