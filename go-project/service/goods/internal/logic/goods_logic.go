package goodslogic

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 获得商品列表
func (s *GoodsService) GetGoodListLogic(ctx context.Context, in *proto.GoodsFilterReq) (*proto.GoodsListRes, error) {
	return s.GoodsData.GetGoodListDB(ctx, in)
}

// 用于通过id数组得到所有商品信息,常用于从订单中获得所有商品信息,
func (s *GoodsService) GetGoodsListByIdLogic(ctx context.Context, in *proto.GoodsIdsReq) (*proto.GoodsListRes, error) {
	return s.GoodsData.GetGoodsListByIdDB(ctx, in)
}

func (s *GoodsService) GetGoodsDetailLogic(ctx context.Context, in *proto.GoodsInfoReq) (*proto.GoodsDetailRes, error) {
	return s.GoodsData.GetGoodsDetailDB(ctx, in)
}

func (s *GoodsService) CreateGoodsLogic(ctx context.Context, in *proto.CreateGoodsReq) (*proto.GoodsDetailRes, error) {
	return s.GoodsData.CreateGoodsDB(ctx, in)
}

func (s *GoodsService) DeleteGoodsLogic(ctx context.Context, in *proto.DelGoodsReq) (*emptypb.Empty, error) {
	return s.GoodsData.DeleteGoodsDB(ctx, in)
}

func (s *GoodsService) UpdeateGoodsLogic(ctx context.Context, in *proto.UpdateGoodsReq) (*emptypb.Empty, error) {
	return s.GoodsData.UpdeateGoodsDB(ctx, in)
}
