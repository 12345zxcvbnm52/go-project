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
	AfterTime  time.Time
	BeforeTime time.Time
	PagesNum   int32
	PageSize   int32
}

func (u *Order) FindByOpt(opt *OrderFindOption) (*OrderResult, error) {
	res := &OrderResult{}
	LocDB := gb.DB.Model(&Order{})
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
	r := LocDB.Find(res.Data)
	if r.RowsAffected == 0 {
		return nil, ErrOrderNotFound
	}
	if r.Error != nil {
		return nil, ErrInternalWrong
	}
	return res, nil
}

func (u *Order) InsertOne(ctx context.Context) error {
	orderListener := OrderListener{Ctx: ctx}
	p, err := rocketmq.NewTransactionProducer(
		&orderListener,
		producer.WithNameServer([]string{fmt.Sprintf("%s:%d", gb.ServerConfig.RockMqConfig.Host, gb.ServerConfig.RockMqConfig.Port)}),
	)
	if err != nil {
		zap.S().Errorw("rocketmq生成producer失败", "msg", err.Error())
		return ErrInternalWrong
	}
	if err = p.Start(); err != nil {
		zap.S().Errorw("rocketmq中producer启动失败", "msg", err.Error())
		return ErrInternalWrong
	}

	jsonStr, _ := json.Marshal(u)
	if _, err = p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_reback", jsonStr)); err != nil {
		zap.S().Errorw("rockmq中producer发送事务消息失败", "msg", err.Error())
		return ErrInternalWrong
	}

	return orderListener.Err
}

// 隐式
func (u *Order) FindOne() error {
	if res := gb.DB.Model(&Order{}).Where("user_id = ?", u.UserId).First(u, u.ID); res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return ErrOrderNotFound
		}
		return ErrInternalWrong
	}
	return nil
}
