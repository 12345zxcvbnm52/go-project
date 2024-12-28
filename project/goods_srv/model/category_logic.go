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
	res := gb.DB.Model(&Category{}).Create(u)
	if res.Error == gorm.ErrDuplicatedKey {
		return ErrDuplicatedCategy
	}
	if res.Error != nil {
		return ErrInternalWrong
	}
	return nil
}

func (u *Category) UpdateOneById() error {
	res := gb.DB.Model(&Category{}).Updates(u)
	if res.Error == gorm.ErrRecordNotFound {
		return ErrCategoryNotFound
	}
	if res.Error != nil {
		return ErrInternalWrong
	}
	return nil
}
