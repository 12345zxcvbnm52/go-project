package model

import (
	gb "goods_srv/global"

	"gorm.io/gorm"
)

func (u *Category) FindOneById() error {
	res := gb.DB.Model(u).Find(&u)
	if res.Error != nil {
		return ErrInternalWrong
	}
	if res.RowsAffected == 0 {
		return ErrCategoryNotFound
	}
	return nil
}

func (u *Category) InsertOne() error {
	LocDB := gb.DB.Model(&Category{})
	if u.Level == 0 {
		LocDB = LocDB.Omit("Level")
	}
	res := LocDB.Create(u)
	if res.Error == gorm.ErrDuplicatedKey {
		return ErrDuplicatedCategy
	}
	if res.Error != nil {
		return ErrInternalWrong
	}
	return nil
}

func (u *Category) UpdateOneById() error {
	LocDB := gb.DB.Model(&Category{})
	if u.ParentCategoryID == nil {
		LocDB = LocDB.Omit("parent_category_id")
	}
	res := LocDB.Where("id = ?", u.ID).Updates(u)
	if res.Error == gorm.ErrRecordNotFound {
		return ErrCategoryNotFound
	}
	if res.Error != nil {
		return ErrInternalWrong
	}
	return nil
}

func (u *Category) DeleteOneById() error {
	if u.Level == gb.EndLevel {
		res := gb.DB.Model(&Category{}).Delete(u)
		if res.Error != nil {
			return ErrInternalWrong
		}
		if res.RowsAffected == 0 {
			return ErrCategoryNotFound
		}
	} else {
		var num int64
		res := gb.DB.Model(&Category{}).Where("parent_category_id = ?", u.ID).Count(&num)
		if res.Error != nil {
			return ErrInternalWrong
		}
		if num > 0 {
			return ErrCategoryRefered
		}
		res = gb.DB.Model(&Category{}).Delete(u)
		if res.Error != nil {
			return ErrInternalWrong
		}
		if res.RowsAffected == 0 {
			return ErrCategoryNotFound
		}
	}
	return nil
}
