package goodsdata

import (
	"context"
	"errors"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 轮播窗口服务
func (s *GormGoodsData) GetBannerListDB(ctx context.Context, in *emptypb.Empty) (*proto.BannerListRes, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormGoodsData) CreateBannerDB(ctx context.Context, in *proto.CreateBannerReq) (*proto.BannerInfoRes, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormGoodsData) DeleteBannerDB(ctx context.Context, in *proto.DelBannerReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormGoodsData) UpdateBannerDB(ctx context.Context, in *proto.UpdateBannerReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}
