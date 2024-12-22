package model

import (
	"fmt"
	gb "goods_srv/global"
	"goods_srv/util"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (u *Banner) FindOneById() error {
	key := fmt.Sprintf("goods_%d", u.ID)
	s, err := gb.RedisConn.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
		} else {
			zap.S().Errorw("redis查询出现未检测的错误", "msg", err.Error())
		}
		res := gb.DB.Find(u)
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return res.Error
	}
	if err = util.Unmarshal([]byte(s), u); err != nil {
		zap.S().Errorw("缓存数据格式错误", "msg", err.Error())
		return err
	}
	return nil
}

func (u *Banner) InsertOne() error {
	result := gb.DB.Model(&Banner{}).Create(u)
	if result.Error == gorm.ErrDuplicatedKey {
		return gorm.ErrDuplicatedKey
	}
	return result.Error
}

func (u *Banner) DeleteOneById() error {
	result := gb.DB.Delete(&Banner{}, u.ID)
	return result.Error
}

func (u *Banner) UpdateOneById() error {
	result := gb.DB.Updates(u)
	return result.Error
}
