package goodsdata

import (
	"context"
	"errors"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌服务
func (s *GormGoodsData) GetBrandListDB(ctx context.Context, in *proto.BrandFilterReq) (*proto.BrandListRes, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormGoodsData) CreateBrandDB(ctx context.Context, in *proto.CreateBrandReq) (*proto.BrandInfoRes, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormGoodsData) DeleteBrandDB(ctx context.Context, in *proto.DelBrandReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormGoodsData) UpdateBrandDB(ctx context.Context, in *proto.UpdateBrandReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}
