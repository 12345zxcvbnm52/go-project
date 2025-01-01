package model

import (
	"context"
	"encoding/json"
	"fmt"
	gb "order_srv/global"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderResult struct {
	Data  []*Order
	Total int64
}

type OrderFindOption struct {
	UserId     uint32
	AfterTime  time.Time
	BeforeTime time.Time
	PagesNum   int32
	PageSize   int32
}

func (u *Order) FindByOpt(opt *OrderFindOption) (*OrderResult, error) {
	res := &OrderResult{}
	LocDB := gb.DB.Model(&Order{}).Where("user_id = ?", opt.UserId)
	if !opt.AfterTime.IsZero() {
		LocDB = LocDB.Where("created_at >= ?", opt.AfterTime)
	}
	if !opt.BeforeTime.IsZero() {
		LocDB = LocDB.Where("created_at <= ?", opt.BeforeTime)
	}
	LocDB.Count(&res.Total)
	if opt.PageSize > 0 {
		LocDB = LocDB.Scopes(Paginate(int(opt.PagesNum), int(opt.PageSize)))
	}
	r := LocDB.Find(&res.Data)
	if r.RowsAffected == 0 {
		return nil, ErrOrderNotFound
	}
	if r.Error != nil {
		zap.S().Errorw("订单通过条件查询失败", "msg", r.Error.Error())
		return nil, ErrInternalWrong
	}
	return res, nil
}

// 基于可靠消息的最终事务一致性的订单生成
func (u *Order) InsertOne(ctx context.Context) error {
	orderListener := OrderListener{Ctx: ctx}
	p, err := rocketmq.NewTransactionProducer(
		&orderListener,
		producer.WithNameServer([]string{fmt.Sprintf("%s:%d", gb.ServerConfig.RockMq.Host, gb.ServerConfig.RockMq.Port)}),
		producer.WithGroupName("Transaction-"+u.OrderSign),
	)
	if err != nil {
		zap.S().Errorw("rocketmq生成producer失败", "msg", err.Error())
		return ErrInternalWrong
	}
	if err = p.Start(); err != nil {
		zap.S().Errorw("rocketmq中producer启动失败", "msg", err.Error())
		return ErrInternalWrong
	}
	defer p.Shutdown()
	jsonStr, _ := json.Marshal(u)
	//发送库存归还半消息到rocketMQ中,并等待本地事务完成或出错
	if _, err = p.SendMessageInTransaction(
		context.Background(),
		primitive.NewMessage(gb.ServerConfig.RockMq.RebackTopic, jsonStr),
	); err != nil {
		zap.S().Errorw("rockmq中producer发送事务消息失败", "msg", err.Error())
		return ErrInternalWrong
	}
	u.Cost = orderListener.Cost
	u.ID = orderListener.OrderId
	return orderListener.Err
}

// 后续可能还会扩展根据订单id找
func (u *Order) FindOneByOrderSign() error {
	if res := gb.DB.Model(&Order{}).Where("order_sign = ?", u.OrderSign).First(u); res.Error != nil {
		zap.S().Errorw("订单通过订单号查询失败", "msg", res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return ErrOrderNotFound
		}
		return ErrInternalWrong
	}
	return nil
}

func (u *Order) FindOneById() error {
	if res := gb.DB.Model(&Order{}).First(u, u.ID); res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			zap.S().Errorw("订单通过Id查询失败", "msg", res.Error.Error())
			return ErrOrderNotFound
		}
		return ErrInternalWrong
	}
	return nil
}

func (u *Order) UpdateById(tx ...*gorm.DB) error {
	LocDB := gb.DB
	if len(tx) >= 1 {
		LocDB = tx[0]
	}
	if res := LocDB.Model(&Order{}).
		Where("status < ? and order_sign = ?", u.Status, u.OrderSign).Update("status", u.Status); res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound || res.RowsAffected == 0 {
			return ErrOrderNotFound
		}
		zap.S().Errorw("订单更新失败", "msg", res.Error.Error())
		return ErrInternalWrong
	}
	return nil
}
