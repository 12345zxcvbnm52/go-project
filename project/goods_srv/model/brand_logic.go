package model

import (
	gb "goods_srv/global"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BrandResult struct {
	Total int64
	Data  []*Brand
}

type BrandFindOption struct {
	PagesNum int32
	PageSize int32
}

func (u *Brand) FindByOpt(opt *BrandFindOption) (*BrandResult, error) {
	res := &BrandResult{}
	LocDB := gb.DB.Model(&Brand{})
	LocDB.Count(&res.Total)
	r := LocDB.Scopes(Paginate(int(opt.PagesNum), int(opt.PageSize))).Find(&res.Data)
	if r.Error != nil {
		zap.S().Errorw("按条件批量查找品牌失败", "msg", r.Error.Error())
		return nil, ErrInternalWrong
	}
	if r.RowsAffected == 0 {
		return nil, ErrBrandNotFound
	}
	return res, nil
}

func (u *Brand) FindOneById() error {
	if res := gb.DB.Model(&Brand{}).First(u, u.ID); res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return ErrBrandNotFound
		}
		return ErrInternalWrong
	}
	return nil
}

func (u *Brand) InsertOne() error {
	if res := gb.DB.Model(&Brand{}).Create(u); res.Error != nil {
		zap.S().Errorw("一个品牌插入失败", "msg", res.Error.Error())
		if res.Error == gorm.ErrDuplicatedKey {
			return ErrDuplicatedBrand
		}
		return ErrInternalWrong
	}
	return nil
}

func (u *Brand) DeleteOneById() error {
	res := gb.DB.Model(&Brand{}).Delete(&Brand{}, u.ID)
	if res.Error != nil {
		zap.S().Errorw("一个品牌删除失败", "msg", res.Error.Error())
		return ErrInternalWrong
	}
	if res.RowsAffected == 0 {
		return ErrBrandNotFound
	}
	return nil
}

func (u *Brand) UpdateOneById() error {
	res := gb.DB.Model(&Brand{}).Updates(u)
	if res.Error != nil {
		zap.S().Errorw("一个品牌更新失败", "msg", res.Error.Error())
		return ErrInternalWrong
	}
	if res.RowsAffected == 0 {
		return ErrBrandNotFound
	}
	return nil
}
