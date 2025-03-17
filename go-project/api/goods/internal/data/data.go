package goodsdata

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GoodsDataService interface {
	//获得商品列表
	GetGoodListDB(context.Context, *proto.GoodsFilterReq) (*proto.GoodsListRes, error)
	//用于通过id数组得到所有商品信息,常用于从订单中获得所有商品信息,
	GetGoodsListByIdDB(context.Context, *proto.GoodsIdsReq) (*proto.GoodsListRes, error)

	GetGoodsDetailDB(context.Context, *proto.GoodsInfoReq) (*proto.GoodsDetailRes, error)

	CreateGoodsDB(context.Context, *proto.CreateGoodsReq) (*proto.GoodsDetailRes, error)

	DeleteGoodsDB(context.Context, *proto.DelGoodsReq) (*emptypb.Empty, error)

	UpdeateGoodsDB(context.Context, *proto.UpdateGoodsReq) (*emptypb.Empty, error)
	//商品类型服务
	GetCategoryListDB(context.Context, *emptypb.Empty) (*proto.CategoryListRes, error)

	GetCategoryInfoDB(context.Context, *proto.SubCategoryReq) (*proto.SubCategoryListRes, error)

	CreateCategoryDB(context.Context, *proto.CreateCategoryReq) (*proto.CategoryInfoRes, error)

	DeleteCategoryDB(context.Context, *proto.DelCategoryReq) (*emptypb.Empty, error)

	UpdateCategoryDB(context.Context, *proto.UpdateCategoryReq) (*emptypb.Empty, error)
	//品牌服务
	GetBrandListDB(context.Context, *proto.BrandFilterReq) (*proto.BrandListRes, error)

	CreateBrandDB(context.Context, *proto.CreateBrandReq) (*proto.BrandInfoRes, error)

	DeleteBrandDB(context.Context, *proto.DelBrandReq) (*emptypb.Empty, error)

	UpdateBrandDB(context.Context, *proto.UpdateBrandReq) (*emptypb.Empty, error)
	//轮播窗口服务
	GetBannerListDB(context.Context, *emptypb.Empty) (*proto.BannerListRes, error)

	CreateBannerDB(context.Context, *proto.CreateBannerReq) (*proto.BannerInfoRes, error)

	DeleteBannerDB(context.Context, *proto.DelBannerReq) (*emptypb.Empty, error)

	UpdateBannerDB(context.Context, *proto.UpdateBannerReq) (*emptypb.Empty, error)
	//品牌分类服务
	GetCategoryBrandListDB(context.Context, *proto.CategoryBrandFilterReq) (*proto.CategoryBrandListRes, error)
	//通过一个类型获得所有有这个类型的品牌
	//rpc GetBrandListByCategory(CategoryInfoReq)returns(BrandListRes);
	CreateCategoryBrandDB(context.Context, *proto.CreateCategoryBrandReq) (*proto.CategoryBrandInfoRes, error)

	DeleteCategoryBrandDB(context.Context, *proto.DelCategoryBrandReq) (*emptypb.Empty, error)

	UpdateCategoryBrandDB(context.Context, *proto.UpdateCategoryBrandReq) (*emptypb.Empty, error)
}

func MustNewGrpcGoodsData(c *grpc.ClientConn) GoodsDataService {
	return &GrpcGoodsData{Cli: proto.NewGoodsClient(c)}
}

var _ GoodsDataService = (*GrpcGoodsData)(nil)

// Goods服务中的Data层,是数据操作的具体逻辑
type GrpcGoodsData struct {
	Cli proto.GoodsClient
}
