package model

import (
	gb "goods_srv/global"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategyBrandResult struct {
	Total int64
	Data  []*CategoryWithBrand
}

type CategyBrandFindOption struct {
	PagesNum int32
	PageSize int32
}

func (u *CategoryWithBrand) FindByOpt(opt *CategyBrandFindOption) (*CategyBrandResult, error) {
	res := &CategyBrandResult{}
	LocDB := gb.DB.Model(&CategoryWithBrand{})
	LocDB.Count(&res.Total)
	r := LocDB.Preload("Category").Preload("Brand").Scopes(Paginate(int(opt.PagesNum), int(opt.PageSize))).Find(&res.Data)
	if r.Error != nil {
		zap.S().Errorw("按条件批量查找类型-品牌项失败", "msg", r.Error.Error())
		return nil, ErrInternalWrong
	}
	if r.RowsAffected == 0 {
		return nil, ErrCategyBrandNotFound
	}
	return res, nil
}

func (u *CategoryWithBrand) FindByCategyId(CategyId uint32) (*CategyBrandResult, error) {
	res := &CategyBrandResult{}
	LocDB := gb.DB.Model(&CategoryWithBrand{}).Where("category_id = ?", CategyId)
	LocDB.Count(&res.Total)
	r := LocDB.Preload("Brand").Find(&res.Data)
	if r.Error != nil {
		zap.S().Errorw("按根据类型查找类型-品牌项失败", "msg", r.Error.Error())
		return nil, ErrInternalWrong
	}
	if r.RowsAffected == 0 {
		return nil, ErrCategyBrandNotFound
	}
	return res, nil
}

func (u *CategoryWithBrand) FindOneById() error {
	if res := gb.DB.Model(&CategoryWithBrand{}).First(u, u.ID); res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return ErrCategyBrandNotFound
		}
		return ErrInternalWrong
	}
	return nil
}

//更新和插入可以先看看对应的Brand和Category是否存在,当然这个任务可以让上层做

func (u *CategoryWithBrand) InsertOne() error {
	if res := gb.DB.Model(&CategoryWithBrand{}).Create(u); res.Error != nil {
		zap.S().Errorw("一个类型-品牌项插入失败", "msg", res.Error.Error())
		if res.Error == gorm.ErrDuplicatedKey {
			return ErrDuplicatedCategyBrand
		}
		return ErrInternalWrong
	}
	return nil
}

func (u *CategoryWithBrand) DeleteOneById() error {
	res := gb.DB.Model(&CategoryWithBrand{}).Delete(&CategoryWithBrand{}, u.ID)
	if res.Error != nil {
		zap.S().Errorw("一个类型-品牌项删除失败", "msg", res.Error.Error())
		return ErrInternalWrong
	}
	if res.RowsAffected == 0 {
		return ErrCategyBrandNotFound
	}
	return nil
}

func (u *CategoryWithBrand) UpdateOneById() error {
	res := gb.DB.Model(&CategoryWithBrand{}).Updates(u)
	if res.Error != nil {
		zap.S().Errorw("一个类型-品牌项更新失败", "msg", res.Error.Error())
		return ErrInternalWrong
	}
	if res.RowsAffected == 0 {
		return ErrCategyBrandNotFound
	}
	return nil
}
