package goodslogic

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌服务
func (s *GoodsService) GetBrandListLogic(ctx context.Context, in *proto.BrandFilterReq) (*proto.BrandListRes, error) {
	return s.GoodsData.GetBrandListDB(ctx, in)
}

func (s *GoodsService) CreateBrandLogic(ctx context.Context, in *proto.CreateBrandReq) (*proto.BrandInfoRes, error) {
	return s.GoodsData.CreateBrandDB(ctx, in)
}

func (s *GoodsService) DeleteBrandLogic(ctx context.Context, in *proto.DelBrandReq) (*emptypb.Empty, error) {
	return s.GoodsData.DeleteBrandDB(ctx, in)
}

func (s *GoodsService) UpdateBrandLogic(ctx context.Context, in *proto.UpdateBrandReq) (*emptypb.Empty, error) {
	return s.GoodsData.UpdateBrandDB(ctx, in)
}
