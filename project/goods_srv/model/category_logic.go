package model

import gb "goods_srv/global"

func (u *Category) FindOneById() error {
	if err := gb.DB.Model(u).Find(&u); err != nil {
		return err.Error
	}
	return nil
}
