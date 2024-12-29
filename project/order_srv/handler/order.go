package handler

import (
	"context"
	"fmt"
	"math/rand/v2"
	"order_srv/model"
	pb "order_srv/proto"
	"time"
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
	res.OrderSign = GenOrderSign(req.UserId)
	if err := res.InsertOne(ctx); err != nil {
		return nil, err
	}
	return OrderToOrderInfoRes(res), nil
}

func (us *OrderServer) GetOrderList(ctx context.Context, req *pb.OrderFliterReq) (*pb.OrderListRes, error) {
	logic := &model.Order{}
	res, err := logic.FindByOpt(&model.OrderFindOption{PagesNum: req.PagesNum, PageSize: req.PageSize})
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
	if err := u.FindOne(); err != nil {
		return nil, err
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
	if r, err := logic.FindByOrderId(); err != nil {
		return nil, err
	} else {
		for _, v := range r.Data {
			res.Items = append(res.Items, OrderGoodsToItemRes(v))
		}
	}
	return res, nil
}
