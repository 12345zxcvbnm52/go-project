package goodsdata

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 商品类型服务
func (d *GrpcGoodsData) GetCategoryListDB(ctx context.Context, in *emptypb.Empty) (*proto.CategoryListRes, error) {
	return d.Cli.GetCategoryList(ctx, in)
}

func (d *GrpcGoodsData) GetCategoryInfoDB(ctx context.Context, in *proto.SubCategoryReq) (*proto.SubCategoryListRes, error) {
	return d.Cli.GetCategoryInfo(ctx, in)
}

func (d *GrpcGoodsData) CreateCategoryDB(ctx context.Context, in *proto.CreateCategoryReq) (*proto.CategoryInfoRes, error) {
	return d.Cli.CreateCategory(ctx, in)
}

func (d *GrpcGoodsData) DeleteCategoryDB(ctx context.Context, in *proto.DelCategoryReq) (*emptypb.Empty, error) {
	return d.Cli.DeleteCategory(ctx, in)
}

func (d *GrpcGoodsData) UpdateCategoryDB(ctx context.Context, in *proto.UpdateCategoryReq) (*emptypb.Empty, error) {
	return d.Cli.UpdateCategory(ctx, in)
}
