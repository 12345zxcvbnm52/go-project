package model

import (
	"errors"
	"fmt"
	gb "user_srv/global"
	"user_srv/util"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ErrNotFound      = status.Error(codes.NotFound, "未找到对应的用户")
	ErrInternalWrong = status.Error(codes.Internal, "服务器内部未知的错误,请稍后尝试")
	ErrDuplicated    = status.Error(codes.AlreadyExists, "欲创建的用户已存在")
	ErrBadAuth       = status.Error(codes.InvalidArgument, "错误的账户或密码")
)

type Result struct {
	Data  []*User
	Total int64
}

type FindOption struct {
	KeyWords string
	Age      int32
	Gender   string
	PagesNum int32
	PageSize int32
}

// 考虑到只允许一台设备登入一个用户,那么用户的增删改就可以直接单机而不用开事务和锁
//哪怕允许多人操作下,用户的增删改流量也不会大,可以直接锁表锁缓存

func (u *User) FindByOpt(opt *FindOption) (*Result, error) {
	LocDB := gb.DB.Model(u)
	if opt.KeyWords != "" {
		LocDB = LocDB.Where("username like ?", "%"+opt.KeyWords+"%")
	}
	if opt.Gender == "boy" || opt.Gender == "girl" {
		LocDB = LocDB.Where("gender = ?", opt.Gender)
	}
	if opt.PageSize > 0 {
		LocDB = LocDB.Scopes(Paginate(int(opt.PagesNum), int(opt.PageSize)))
	}
	var total int64
	LocDB.Count(&total)
	res := []*User{}
	if err := LocDB.Find(&res); err != nil {
		zap.S().Errorw("通过条件获取用户失败", "msg", err, "条件为", opt)
		return nil, ErrInternalWrong
	}
	return &Result{
		Total: total,
		Data:  res,
	}, nil
}

func (u *User) FindByIds(ids ...int32) (*Result, error) {
	res := []*User{}
	if len(ids) == 0 {
		return nil, ErrNotFound
	}
	LocDB := gb.DB.Where("id in (?)", ids)
	r := LocDB.Find(&res)
	if r.Error != nil {
		zap.S().Errorw("商品按id批量查询失败", "msg", r.Error.Error(), "id号为", ids)
		return nil, ErrInternalWrong
	}
	return &Result{
		Data:  res,
		Total: r.RowsAffected,
	}, nil
}

func (u *User) FindOneById() error {
	key := fmt.Sprintf("user_%d", u.ID)
	//不好说生产上会不会出现在go-mysql-transfer监控下仍然出现缓存不一致的情况
	//暂时也不知道怎么写来通过transfer主动同步一次缓存,就先搁置吧
	//可以考虑复杂环境下换canal
	s, err := gb.RedisConn.Get(key).Result()
	if err != nil {
		if err != redis.Nil {
			zap.S().Errorw("redis查找出现未检测的问题", "msg", err.Error())
		}
		res := gb.DB.Take(u)
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				return ErrNotFound
			} else {
				zap.S().Errorw("gorm查找出现未检测的问题", "msg", err.Error())
				return ErrInternalWrong
			}
		}
		return nil
	} else {
		util.Unmarshal([]byte(s), u)
		return nil
	}
}

// 可以考虑在缓存中多设计一个mobile前缀key,暂时不想这么多
func (u *User) FindOneByMobile() error {
	res := gb.DB.Where("Mobile=?", u.Mobile).Find(u)
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	if res.Error != nil {
		zap.S().Errorw("出现未检测的问题", "msg", res.Error.Error())
		return ErrInternalWrong
	}
	return nil
}

// 删除用户的逻辑可以先考虑考虑,
func (u *User) DeleteById() error {
	res := gb.DB.Delete(&User{}, u.ID)
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return res.Error
}

// web端传入的数据是保证完整的,故无需检测
func (u *User) InsertOne() error {
	res := gb.DB.Create(u)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return ErrDuplicated
		} else {
			zap.S().Errorw("UserCreate失败", "msg", res.Error.Error())
			return ErrInternalWrong
		}
	}
	return nil
}

// 严格限制只更新一个
func (u *User) UpdateOneById() error {
	res := gb.DB.Updates(u)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return ErrInternalWrong
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
