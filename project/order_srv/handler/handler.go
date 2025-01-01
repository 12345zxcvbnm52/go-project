package handler

import (
	"order_srv/model"
	pb "order_srv/proto"
)

type OrderServer struct {
	pb.UnimplementedOrderServer
}

func OrderGoodsToItemRes(c *model.OrderGoods) *pb.OrderItemRes {
	return &pb.OrderItemRes{
		Id:          c.ID,
		OrderId:     c.OrderId,
		GoodsId:     c.GoodsId,
		GoodsNum:    c.GoodsNum,
		GoodsName:   c.GoodsName,
		GoodsImages: c.GoodsImage,
		GoodsPrice:  c.GoodsPrice,
	}
}
