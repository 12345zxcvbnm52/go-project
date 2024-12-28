package model

import (
	gb "goods_srv/global"

	"gorm.io/gorm"
)

func (u *Brand) FindOneById() error {
	if res := gb.DB.Model(&Brand{}).First(u, u.ID); res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return ErrBrandNotFound
		} else {
			return ErrInternalWrong
		}
	}
	return nil
}
