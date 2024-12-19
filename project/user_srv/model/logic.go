package model

import (
	"fmt"
	gb "user_srv/global"
	"user_srv/util"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 这些错误应当只在内部使用

func SelectAll(db *gorm.DB) *gorm.DB {
	return db.Model(&User{}).Select("Mobile", "UserName", "Password", "Gender", "Birth", "Role")
}

func (u *User) FindOne() error {
	key := fmt.Sprintf("user_%d", u.ID)
	s, err := gb.RedisConn.Get(key).Result()
	if err != nil {
		if err == redis.ErrEmptyKey {

		} else {
			zap.S().Errorw("redis连接池出现未知问题", "msg", err.Error())
		}
	} else {
		util.Unmarshal([]byte(s), u)
		return nil
	}

	res := gb.DB.Find(u)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}

func (u *User) FindOneByMobile() error {
	res := gb.DB.Where("Mobile=?", u.Mobile).Find(u)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}

func (u *User) DeleteById() error {
	res := gb.DB.Delete(&User{}, u.ID)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}

func (u *User) InsertOne() error {
	res := gb.DB.Create(u)
	return res.Error
}

// 严格限制只更新一个
func (u *User) UpdateOne() error {
	return gb.DB.Updates(u).Error
}
