package goodslogic

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 商品类型服务
func (s *GoodsService) GetCategoryListLogic(ctx context.Context, in *emptypb.Empty) (*proto.CategoryListRes, error) {
	return s.GoodsData.GetCategoryListDB(ctx, in)
}

func (s *GoodsService) GetCategoryInfoLogic(ctx context.Context, in *proto.SubCategoryReq) (*proto.SubCategoryListRes, error) {
	return s.GoodsData.GetCategoryInfoDB(ctx, in)
}

func (s *GoodsService) CreateCategoryLogic(ctx context.Context, in *proto.CreateCategoryReq) (*proto.CategoryInfoRes, error) {
	return s.GoodsData.CreateCategoryDB(ctx, in)
}

func (s *GoodsService) DeleteCategoryLogic(ctx context.Context, in *proto.DelCategoryReq) (*emptypb.Empty, error) {
	return s.GoodsData.DeleteCategoryDB(ctx, in)
}

func (s *GoodsService) UpdateCategoryLogic(ctx context.Context, in *proto.UpdateCategoryReq) (*emptypb.Empty, error) {
	return s.GoodsData.UpdateCategoryDB(ctx, in)
}
