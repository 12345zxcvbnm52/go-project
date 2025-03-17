package goodsdata

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌分类服务
func (d *GrpcGoodsData) GetCategoryBrandListDB(ctx context.Context, in *proto.CategoryBrandFilterReq) (*proto.CategoryBrandListRes, error) {
	return d.Cli.GetCategoryBrandList(ctx, in)
}

// 通过一个类型获得所有有这个类型的品牌
// rpc GetBrandListByCategory(CategoryInfoReq)returns(BrandListRes);
func (d *GrpcGoodsData) CreateCategoryBrandDB(ctx context.Context, in *proto.CreateCategoryBrandReq) (*proto.CategoryBrandInfoRes, error) {
	return d.Cli.CreateCategoryBrand(ctx, in)
}

func (d *GrpcGoodsData) DeleteCategoryBrandDB(ctx context.Context, in *proto.DelCategoryBrandReq) (*emptypb.Empty, error) {
	return d.Cli.DeleteCategoryBrand(ctx, in)
}

func (d *GrpcGoodsData) UpdateCategoryBrandDB(ctx context.Context, in *proto.UpdateCategoryBrandReq) (*emptypb.Empty, error) {
	return d.Cli.UpdateCategoryBrand(ctx, in)
}
