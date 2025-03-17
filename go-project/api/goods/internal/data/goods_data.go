package goodsdata

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 获得商品列表
func (d *GrpcGoodsData) GetGoodListDB(ctx context.Context, in *proto.GoodsFilterReq) (*proto.GoodsListRes, error) {
	return d.Cli.GetGoodList(ctx, in)
}

// 用于通过id数组得到所有商品信息,常用于从订单中获得所有商品信息,
func (d *GrpcGoodsData) GetGoodsListByIdDB(ctx context.Context, in *proto.GoodsIdsReq) (*proto.GoodsListRes, error) {
	return d.Cli.GetGoodsListById(ctx, in)
}

func (d *GrpcGoodsData) GetGoodsDetailDB(ctx context.Context, in *proto.GoodsInfoReq) (*proto.GoodsDetailRes, error) {
	return d.Cli.GetGoodsDetail(ctx, in)
}

func (d *GrpcGoodsData) CreateGoodsDB(ctx context.Context, in *proto.CreateGoodsReq) (*proto.GoodsDetailRes, error) {
	return d.Cli.CreateGoods(ctx, in)
}

func (d *GrpcGoodsData) DeleteGoodsDB(ctx context.Context, in *proto.DelGoodsReq) (*emptypb.Empty, error) {
	return d.Cli.DeleteGoods(ctx, in)
}

func (d *GrpcGoodsData) UpdeateGoodsDB(ctx context.Context, in *proto.UpdateGoodsReq) (*emptypb.Empty, error) {
	return d.Cli.UpdeateGoods(ctx, in)
}
