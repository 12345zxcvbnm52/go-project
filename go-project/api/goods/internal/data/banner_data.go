package goodsdata

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 轮播窗口服务
func (d *GrpcGoodsData) GetBannerListDB(ctx context.Context, in *emptypb.Empty) (*proto.BannerListRes, error) {
	return d.Cli.GetBannerList(ctx, in)
}

func (d *GrpcGoodsData) CreateBannerDB(ctx context.Context, in *proto.CreateBannerReq) (*proto.BannerInfoRes, error) {
	return d.Cli.CreateBanner(ctx, in)
}

func (d *GrpcGoodsData) DeleteBannerDB(ctx context.Context, in *proto.DelBannerReq) (*emptypb.Empty, error) {
	return d.Cli.DeleteBanner(ctx, in)
}

func (d *GrpcGoodsData) UpdateBannerDB(ctx context.Context, in *proto.UpdateBannerReq) (*emptypb.Empty, error) {
	return d.Cli.UpdateBanner(ctx, in)
}
