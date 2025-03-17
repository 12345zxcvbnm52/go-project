package goodsdata

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌服务
func (d *GrpcGoodsData) GetBrandListDB(ctx context.Context, in *proto.BrandFilterReq) (*proto.BrandListRes, error) {
	return d.Cli.GetBrandList(ctx, in)
}

func (d *GrpcGoodsData) CreateBrandDB(ctx context.Context, in *proto.CreateBrandReq) (*proto.BrandInfoRes, error) {
	return d.Cli.CreateBrand(ctx, in)
}

func (d *GrpcGoodsData) DeleteBrandDB(ctx context.Context, in *proto.DelBrandReq) (*emptypb.Empty, error) {
	return d.Cli.DeleteBrand(ctx, in)
}

func (d *GrpcGoodsData) UpdateBrandDB(ctx context.Context, in *proto.UpdateBrandReq) (*emptypb.Empty, error) {
	return d.Cli.UpdateBrand(ctx, in)
}
