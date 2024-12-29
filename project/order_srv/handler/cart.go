package handler

import (
	"context"
	"order_srv/model"
	pb "order_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

func CartToCartItemRes(c *model.Cart) *pb.CartItemInfoRes {
	return &pb.CartItemInfoRes{
		GoodsId:  c.GoodsId,
		GoodsNum: c.GoodsNums,
		Id:       c.ID,
		Selected: c.Selected,
	}
}

func (us *OrderServer) GetUserCartItems(ctx context.Context, req *pb.UserInfoReq) (*pb.CartItemListRes, error) {
	logic := &model.Cart{}
	res, err := logic.FindByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	r := &pb.CartItemListRes{
		Total: res.Total,
	}
	for _, v := range res.Data {
		v.UserId = req.UserId
		r.Data = append(r.Data, CartToCartItemRes(v))
	}
	return r, nil
}

func (us *OrderServer) CreateCartItem(ctx context.Context, req *pb.CartItemInfoReq) (*pb.CartItemInfoRes, error) {
	u := &model.Cart{
		UserId:  req.UserId,
		GoodsId: req.GoodsId,
		//GoodsNums: req.GoodsNum,
		Selected: true,
	}
	if err := u.InsertOne(); err != nil {
		return nil, err
	}
	return CartToCartItemRes(u), nil
}

// 还要考虑到更新时如果goodsNum小于等于0了应该干嘛
func (us *OrderServer) UpdateCartItem(ctx context.Context, req *pb.WriteCartItemReq) (*emptypb.Empty, error) {
	u := &model.Cart{}
	if req.Id != 0 {
		u.ID = req.Id
		u.Selected = req.Selected
		u.GoodsNums = req.GoodsNum
		if err := u.UpdateOneById(); err != nil {
			return nil, err
		}
		return &emptypb.Empty{}, nil
	} else {
		u.Selected = req.Selected
		u.GoodsNums = req.GoodsNum
		u.GoodsId = req.GoodsId
		u.UserId = req.UserId
		return nil, u.UpdateOneByUserGoodsId()
	}
}

func (us *OrderServer) DeleteCartItem(ctx context.Context, req *pb.DelCartItemReq) (*emptypb.Empty, error) {
	u := &model.Cart{}
	u.GoodsId = req.GoodsId
	u.UserId = req.UserId
	if err := u.DeleteOneByUserGoodsIds(); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
