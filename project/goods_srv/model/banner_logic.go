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
			return ErrBannerNotFound
		}
		if res.Error != nil {
			return ErrInternalWrong
		}
	}
	//这个错误搁置转化为rpcErr
	if err = util.Unmarshal([]byte(s), u); err != nil {
		zap.S().Errorw("缓存数据格式错误", "msg", err.Error())
		return err
	}
	return nil
}

func (u *Banner) InsertOne() error {
	if res := gb.DB.Model(&Banner{}).Create(u); res.Error != nil {
		zap.S().Errorw("一件商品banner插入失败", "msg", res.Error.Error())
		if res.Error == gorm.ErrDuplicatedKey {
			return ErrDuplicatedBanner
		}
		return ErrInternalWrong
	}
	return nil
}

func (u *Banner) DeleteOneById() error {
	res := gb.DB.Model(&Banner{}).Delete(&Banner{}, u.ID)
	if res.Error != nil {
		zap.S().Errorw("一件商品banner删除失败", "msg", res.Error.Error())
		return ErrInternalWrong
	}
	if res.RowsAffected == 0 {
		return ErrBannerNotFound
	}
	return nil
}

func (u *Banner) UpdateOneById() error {
	res := gb.DB.Model(&Banner{}).Updates(u)
	if res.Error != nil {
		zap.S().Errorw("一件商品banner更新失败", "msg", res.Error.Error())
		return ErrInternalWrong
	}
	if res.RowsAffected == 0 {
		return ErrBannerNotFound
	}
	return nil
}
