package model

import (
	"errors"
	"fmt"
	gb "user_srv/global"
	"user_srv/util"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 这些错误应当只在内部使用
var (
	ErrLackParam      = errors.New("User缺少指定参数")
	ErrUserNotFind    = errors.New("没查找到指定User")
	ErrWrongParam     = errors.New("错误的参数")
	ErrDuplicatedData = errors.New("插入重复的unique键或主键")
)

func SelectAll(db *gorm.DB) *gorm.DB {
	return db.Model(&User{}).Select("Mobile", "UserName", "Password", "Gender", "Birth", "Role")
}

func (u *User) FindOne() error {
	if u.ID == 0 {
		return ErrLackParam
	}
	key := fmt.Sprintf("user_srv:%d", u.ID)
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
		return ErrUserNotFind
	}
	return res.Error
}

func (u *User) FindOneByMobile() error {
	if u.Mobile == "" {
		return ErrLackParam
	}
	res := gb.DB.Where("Mobile=?", u.Mobile).Find(u)
	if res.RowsAffected == 0 {
		return ErrUserNotFind
	}
	return res.Error
}

func (u *User) DeleteById() error {
	if u.ID == 0 {
		return ErrLackParam
	}
	res := gb.DB.Delete(&User{}, u.ID)
	if res.RowsAffected == 0 {
		return ErrUserNotFind
	}
	return res.Error
}

func (u *User) InsertOne() error {
	if u.ID != 0 {
		return ErrWrongParam
	}
	res := gb.DB.Create(u)
	if res.Error == gorm.ErrDuplicatedKey {
		return ErrDuplicatedData
	}
	return res.Error
}

// 严格限制只更新一个
func (u *User) UpdateOne() error {
	if u.ID == 0 {
		return ErrLackParam
	}
	return gb.DB.Updates(u).Error
}
