package goodslogic

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 轮播窗口服务
func (s *GoodsService) GetBannerListLogic(ctx context.Context, in *emptypb.Empty) (*proto.BannerListRes, error) {
	return s.GoodsData.GetBannerListDB(ctx, in)
}

func (s *GoodsService) CreateBannerLogic(ctx context.Context, in *proto.CreateBannerReq) (*proto.BannerInfoRes, error) {
	return s.GoodsData.CreateBannerDB(ctx, in)
}

func (s *GoodsService) DeleteBannerLogic(ctx context.Context, in *proto.DelBannerReq) (*emptypb.Empty, error) {
	return s.GoodsData.DeleteBannerDB(ctx, in)
}

func (s *GoodsService) UpdateBannerLogic(ctx context.Context, in *proto.UpdateBannerReq) (*emptypb.Empty, error) {
	return s.GoodsData.UpdateBannerDB(ctx, in)
}
