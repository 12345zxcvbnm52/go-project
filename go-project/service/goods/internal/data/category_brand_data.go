package goodsdata

import (
	"context"
	"errors"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌分类服务
func (s *GormGoodsData) GetCategoryBrandListDB(ctx context.Context, in *proto.CategoryBrandFilterReq) (*proto.CategoryBrandListRes, error) {
	return nil, errors.New("this method is not implemented")
}

// 通过一个类型获得所有有这个类型的品牌
// rpc GetBrandListByCategory(CategoryInfoReq)returns(BrandListRes);
func (s *GormGoodsData) CreateCategoryBrandDB(ctx context.Context, in *proto.CreateCategoryBrandReq) (*proto.CategoryBrandInfoRes, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormGoodsData) DeleteCategoryBrandDB(ctx context.Context, in *proto.DelCategoryBrandReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormGoodsData) UpdateCategoryBrandDB(ctx context.Context, in *proto.UpdateCategoryBrandReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}
