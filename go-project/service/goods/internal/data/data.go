package goodsdata

import (
	"context"
	"database/sql"
	proto "kenshop/proto/goods"
	model "kenshop/service/goods/internal/model"

	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// GoodsDataService是提供Goods底层相关数据操作的接口
type GoodsDataService interface {
	//获得商品列表
	GetGoodListDB(ctx context.Context, in *proto.GoodsFilterReq) (*proto.GoodsListRes, error)
	//用于通过id数组得到所有商品信息,常用于从订单中获得所有商品信息,
	GetGoodsListByIdDB(ctx context.Context, in *proto.GoodsIdsReq) (*proto.GoodsListRes, error)

	GetGoodsDetailDB(ctx context.Context, in *proto.GoodsInfoReq) (*proto.GoodsDetailRes, error)

	CreateGoodsDB(ctx context.Context, in *proto.CreateGoodsReq) (*proto.GoodsDetailRes, error)

	DeleteGoodsDB(ctx context.Context, in *proto.DelGoodsReq) (*emptypb.Empty, error)

	UpdeateGoodsDB(ctx context.Context, in *proto.UpdateGoodsReq) (*emptypb.Empty, error)
	//商品类型服务
	GetCategoryListDB(ctx context.Context, in *emptypb.Empty) (*proto.CategoryListRes, error)

	GetCategoryInfoDB(ctx context.Context, in *proto.SubCategoryReq) (*proto.SubCategoryListRes, error)

	CreateCategoryDB(ctx context.Context, in *proto.CreateCategoryReq) (*proto.CategoryInfoRes, error)

	DeleteCategoryDB(ctx context.Context, in *proto.DelCategoryReq) (*emptypb.Empty, error)

	UpdateCategoryDB(ctx context.Context, in *proto.UpdateCategoryReq) (*emptypb.Empty, error)
	//品牌服务
	GetBrandListDB(ctx context.Context, in *proto.BrandFilterReq) (*proto.BrandListRes, error)

	CreateBrandDB(ctx context.Context, in *proto.CreateBrandReq) (*proto.BrandInfoRes, error)

	DeleteBrandDB(ctx context.Context, in *proto.DelBrandReq) (*emptypb.Empty, error)

	UpdateBrandDB(ctx context.Context, in *proto.UpdateBrandReq) (*emptypb.Empty, error)
	//轮播窗口服务
	GetBannerListDB(ctx context.Context, in *emptypb.Empty) (*proto.BannerListRes, error)

	CreateBannerDB(ctx context.Context, in *proto.CreateBannerReq) (*proto.BannerInfoRes, error)

	DeleteBannerDB(ctx context.Context, in *proto.DelBannerReq) (*emptypb.Empty, error)

	UpdateBannerDB(ctx context.Context, in *proto.UpdateBannerReq) (*emptypb.Empty, error)
	//品牌分类服务
	GetCategoryBrandListDB(ctx context.Context, in *proto.CategoryBrandFilterReq) (*proto.CategoryBrandListRes, error)
	//通过一个类型获得所有有这个类型的品牌
	//rpc GetBrandListByCategory(CategoryInfoReq)returns(BrandListRes);
	CreateCategoryBrandDB(ctx context.Context, in *proto.CreateCategoryBrandReq) (*proto.CategoryBrandInfoRes, error)

	DeleteCategoryBrandDB(ctx context.Context, in *proto.DelCategoryBrandReq) (*emptypb.Empty, error)

	UpdateCategoryBrandDB(ctx context.Context, in *proto.UpdateCategoryBrandReq) (*emptypb.Empty, error)
}

// Goods服务中的Data层,是数据操作的具体逻辑
type GormGoodsData struct {
	DB              *gorm.DB
	OpenTransaction bool
}

func (s *GormGoodsData) WithTransaction(opts ...*sql.TxOptions) *GormGoodsData {
	return &GormGoodsData{
		OpenTransaction: true,
		DB:              s.DB.Begin(opts...),
	}
}

func (s *GormGoodsData) Rollback() {
	if !s.OpenTransaction {
		s.DB.Rollback()
	}
}

func MustNewGormGoodsData(db *gorm.DB) *GormGoodsData {
	return &GormGoodsData{DB: db}
}

var _ GoodsDataService = (*GormGoodsData)(nil)

func GoodsToGoodsInfoRes(goods *model.Goods) *proto.GoodsInfoRes {
	return &proto.GoodsInfoRes{
		Id:          goods.ID,
		CategoryId:  goods.CategoryID,
		BrandId:     goods.BrandID,
		Name:        goods.Name,
		GoodsSign:   goods.GoodsSign,
		ClickNum:    goods.ClickNum,
		SoldNum:     goods.SoldNum,
		FavorNum:    goods.FavorNum,
		MarketPrice: goods.MarketPrice,
		SalePrice:   goods.SalePrice,
		ShipFree:    goods.ShipFree,
		FirstImage:  goods.FirstImage,
		IsNew:       goods.IsNew,
		IsHot:       goods.IsHot,
		Status:      goods.Status,
	}
}

func CategoryToCategyInfoRes(categy *model.Category) *proto.CategoryInfoRes {
	return &proto.CategoryInfoRes{
		Id:               categy.ID,
		Name:             categy.Name,
		ParentCategoryId: *categy.ParentCategoryID,
		Level:            categy.Level,
		OnTable:          categy.OnTable,
	}
}
