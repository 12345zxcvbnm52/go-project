package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	gb "order_srv/global"
	"order_srv/model"
	pb "order_srv/proto"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func OrderToOrderInfoRes(c *model.Order) *pb.OrderInfoRes {
	return &pb.OrderInfoRes{
		Id:           c.ID,
		UserId:       c.UserId,
		OrderSign:    c.OrderSign,
		Status:       int32(c.Status),
		PayWay:       c.PayWay,
		SignerName:   c.SignerName,
		SignerMobile: c.SignerMobile,
		Address:      c.Address,
		Cost:         c.Cost,
		Message:      c.Message,
	}
}

const randGen = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenOrderSign(userId uint32) string {
	now := time.Now()
	s := []byte(now.Format(time.DateTime) + fmt.Sprintf("%d", userId) + "@")
	if 32-len(s) < 0 {
		fmt.Println("err")
	}
	for range 32 - len(s) {
		s = append(s, []byte(randGen)[rand.IntN(62)])
	}
	return string(s)
}

// 新建订单需要计算商品的金额,不能由前端传入价格数据(防止爬虫等),
// 顾需要1.根据购物车商品跨微服务调用查询金额,2.跨微服务扣减库存,3生成订单,填充本服务的几个表
// 其中详细:根据购物车表获得商品的信息并写入OrderGoods表,并把信息同步到Order表中
func (us *OrderServer) CreateOrder(ctx context.Context, req *pb.OrderInfoReq) (*pb.OrderInfoRes, error) {

	res := &model.Order{}
	res.Address = req.Address
	res.Message = req.Message
	res.SignerMobile = req.SignerMobile
	res.SignerName = req.SignerName
	res.UserId = req.UserId
	res.PayWay = req.PayWay
	res.Status = model.StatusUnPay
	res.OrderSign = GenOrderSign(req.UserId)
	if err := res.InsertOne(ctx); err != nil {
		return nil, err
	}
	return OrderToOrderInfoRes(res), nil
}

func (us *OrderServer) GetOrderList(ctx context.Context, req *pb.OrderFliterReq) (*pb.OrderListRes, error) {
	logic := &model.Order{}
	res, err := logic.FindByOpt(&model.OrderFindOption{PagesNum: req.PagesNum, PageSize: req.PageSize, UserId: req.UserId})
	if err != nil {
		return nil, err
	}
	r := &pb.OrderListRes{
		Total: res.Total,
	}
	for _, v := range res.Data {
		r.Data = append(r.Data, OrderToOrderInfoRes(v))
	}
	return r, nil
}

func (us *OrderServer) GetOrderInfo(ctx context.Context, req *pb.OrderInfoReq) (*pb.OrderDetailRes, error) {
	u := &model.Order{}
	u.ID = req.Id
	u.UserId = req.UserId
	//先看看这个订单是否属于这个用户
	if err := u.FindOneById(); err != nil {
		return nil, err
	}
	if u.UserId != req.UserId {
		return nil, model.ErrOrderNotFound
	}
	res := &pb.OrderDetailRes{
		Id:           u.ID,
		UserId:       u.UserId,
		OrderSign:    u.OrderSign,
		Status:       int32(u.Status),
		PayWay:       u.PayWay,
		SignerName:   u.SignerName,
		SignerMobile: u.SignerMobile,
		Address:      u.Address,
		Cost:         u.Cost,
		Message:      u.Message,
	}
	//在逻辑上是能保证查找的订单是一定有相应的商品的
	logic := &model.OrderGoods{}
	logic.OrderId = u.ID
	if r, err := logic.FindByOrderId(u.ID); err != nil {
		return nil, err
	} else {
		for _, v := range r.Data {
			res.Items = append(res.Items, OrderGoodsToItemRes(v))
		}
	}
	return res, nil
}

func (us *OrderServer) UpdateOrderStatus(ctx context.Context, req *pb.OrderStatusReq) (*emptypb.Empty, error) {
	u := &model.Order{}
	u.ID = req.Id
	u.Status = int16(req.Status)
	if err := u.UpdateById(); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		order := &model.Order{}
		json.Unmarshal(msgs[i].Body, order)

		zap.S().Infoln(order)

		if err := order.FindOneByOrderSign(); err != nil {
			if err == model.ErrOrderNotFound {
				return consumer.ConsumeSuccess, nil
			}
			return consumer.ConsumeRetryLater, err
		}
		//如果状态小于等于订单未支付就归还
		if order.Status <= model.StatusUnPay {
			tx := gb.DB.Begin()
			order.Status = model.StatusCancelled
			if err := order.UpdateById(tx); err != nil {
				zap.S().Errorw("订单状态更新失败", "msg", err.Error())
				tx.Rollback()
				return consumer.ConsumeRetryLater, err
			}

			orderGoodsLogic := &model.OrderGoods{}
			res, err := orderGoodsLogic.FindByOrderId(order.ID, tx)
			if err != nil {
				if err == model.ErrOrderGoodsNotFound {
					return consumer.ConsumeSuccess, nil
				}
				tx.Rollback()
				return consumer.ConsumeRetryLater, err
			}
			//这里暂时只考虑使用gorm软删除下的购物车恢复功能
			goodsIds := []uint32{}
			for _, v := range res.Data {
				goodsIds = append(goodsIds, v.GoodsId)
			}
			r := tx.Exec("update carts set deleted_at = null where goods_id in (?)", goodsIds)
			if r.Error != nil {
				tx.Rollback()
				zap.S().Errorw("购物车恢复失败", "msg", r.Error.Error())
				return consumer.ConsumeRetryLater, r.Error
			}
			if r.RowsAffected == 0 {
				tx.Rollback()
				zap.S().Errorw("购物车恢复失败", "msg", "不存在要恢复的数据")
				return consumer.ConsumeSuccess, nil
			}
			p, err := rocketmq.NewProducer(
				producer.WithNameServer([]string{fmt.Sprintf("%s:%d", gb.ServerConfig.RockMq.Host, gb.ServerConfig.RockMq.Port)}),
				producer.WithGroupName(order.OrderSign+"-"+gb.ServerConfig.RockMq.RebackTopic),
			)
			if err != nil {
				tx.Rollback()
				zap.S().Errorw("rocketmq生成producer失败", "msg", err.Error())
				return consumer.ConsumeRetryLater, err
			}
			if err = p.Start(); err != nil {
				tx.Rollback()
				zap.S().Errorw("rocketmq中producer启动失败", "msg", err.Error())
				return consumer.ConsumeRetryLater, err
			}
			_, err = p.SendSync(context.Background(), primitive.NewMessage(gb.ServerConfig.RockMq.RebackTopic, msgs[i].Body))
			if err != nil {
				tx.Rollback()
				zap.S().Errorw("rockmq中producer发送事务消息失败", "msg", err.Error())
				return consumer.ConsumeRetryLater, err
			}
			tx.Commit()
		}
	}
	return consumer.ConsumeSuccess, nil
}
