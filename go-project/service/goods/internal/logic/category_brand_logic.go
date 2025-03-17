package goodslogic

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌分类服务
func (s *GoodsService) GetCategoryBrandListLogic(ctx context.Context, in *proto.CategoryBrandFilterReq) (*proto.CategoryBrandListRes, error) {
	return s.GoodsData.GetCategoryBrandListDB(ctx, in)
}

// 通过一个类型获得所有有这个类型的品牌
// rpc GetBrandListByCategory(CategoryInfoReq)returns(BrandListRes);
func (s *GoodsService) CreateCategoryBrandLogic(ctx context.Context, in *proto.CreateCategoryBrandReq) (*proto.CategoryBrandInfoRes, error) {
	return s.GoodsData.CreateCategoryBrandDB(ctx, in)
}

func (s *GoodsService) DeleteCategoryBrandLogic(ctx context.Context, in *proto.DelCategoryBrandReq) (*emptypb.Empty, error) {
	return s.GoodsData.DeleteCategoryBrandDB(ctx, in)
}

func (s *GoodsService) UpdateCategoryBrandLogic(ctx context.Context, in *proto.UpdateCategoryBrandReq) (*emptypb.Empty, error) {
	return s.GoodsData.UpdateCategoryBrandDB(ctx, in)
}
