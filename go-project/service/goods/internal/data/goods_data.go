package goodsdata

import (
	"context"
	"errors"
	"kenshop/pkg/common/paginate"
	proto "kenshop/proto/goods"
	model "kenshop/service/goods/internal/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 获得商品列表
func (s *GormGoodsData) GetGoodListDB(ctx context.Context, in *proto.GoodsFilterReq) (*proto.GoodsListRes, error) {
	res := &proto.GoodsListRes{}
	data := []*model.Goods{}
	s.DB.Model(&model.Goods{}).Count(&res.Total)
	result := s.DB.Scopes(paginate.GormPaginate(int(in.PagesNum), int(in.PageSize))).Find(&data)
	if result.Error != nil {
		return nil, GormGoodsErrHandle(result.Error)
	}
	for _, v := range data {
		res.Data = append(res.Data, GoodsToGoodsInfoRes(v))
	}
	return res, nil
}

// 用于通过id数组得到所有商品信息,常用于从订单中获得所有商品信息,
func (s *GormGoodsData) GetGoodsListByIdDB(ctx context.Context, in *proto.GoodsIdsReq) (*proto.GoodsListRes, error) {
	goodsres := make([]*model.Goods, len(in.Ids))
	if res := s.DB.Model(&model.Goods{}).Where("id in (?)", in.Ids).Find(&goodsres); res.Error != nil {
		return nil, GormGoodsErrHandle(res.Error)
	}
	out := &proto.GoodsListRes{}

	out.Total = int64(len(goodsres))
	for _, v := range goodsres {
		out.Data = append(out.Data, GoodsToGoodsInfoRes(v))
	}
	return out, nil
}

func (s *GormGoodsData) GetGoodsDetailDB(ctx context.Context, in *proto.GoodsInfoReq) (*proto.GoodsDetailRes, error) {
	goods := &model.Goods{}
	goods.ID = in.Id

	if res := s.DB.Model(&model.Goods{}).Preload("GoodsDetail").Preload("Brand").Preload("Category").First(goods); res.Error != nil {
		return nil, GormGoodsErrHandle(res.Error)
	}

	return &proto.GoodsDetailRes{
		CategoryId:  goods.CategoryID,
		BrandId:     goods.BrandID,
		MarketPrice: goods.MarketPrice,
		Status:      goods.Status,
		SalePrice:   goods.SalePrice,
		ShipFree:    goods.ShipFree,
		FavorNum:    goods.FavorNum,
		ClickNum:    goods.ClickNum,
		SoldNum:     goods.SoldNum,
		FirstImage:  goods.FirstImage,
		IsHot:       goods.IsHot,
		IsNew:       goods.IsNew,
		Name:        goods.Name,
		GoodsSign:   goods.GoodsSign,
		Images:      goods.GoodsDetail.Images,
		DescImages:  goods.GoodsDetail.DescImages,
		GoodsBrief:  goods.GoodsDetail.GoodsBrief,
		Category: &proto.CategoryBriefInfoRes{
			Id:   goods.CategoryID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoRes{
			Id:   goods.Brand.ID,
			Name: goods.Brand.Name,
			Logo: goods.Brand.Logo,
		},
	}, nil
}

// TODO 可以在创建后把商品类型信息和品牌信息查询出来并返回
func (s *GormGoodsData) CreateGoodsDB(ctx context.Context, in *proto.CreateGoodsReq) (*proto.GoodsDetailRes, error) {
	goods := &model.Goods{
		CategoryID:  in.CategoryId,
		BrandID:     in.BrandId,
		MarketPrice: in.MarketPrice,
		Status:      model.StatusUnderReview,
		SalePrice:   in.SalePrice,
		ShipFree:    in.ShipFree,
		FavorNum:    0,
		ClickNum:    0,
		SoldNum:     0,
		FirstImage:  in.FirstImage,
		IsHot:       false,
		IsNew:       true,
		Name:        in.Name,
		GoodsSign:   "",
	}
	goodsDetail := &model.GoodsDetail{
		Images:     in.Images,
		DescImages: in.DescImages,
		GoodsBrief: in.GoodsBrief,
	}
	tx := s.DB.Begin()
	res := tx.Model(&model.Goods{}).Create(goods)
	if res.Error != nil {
		return nil, GormGoodsErrHandle(res.Error)
	}
	goodsDetail.GoodsId = goods.ID
	res = tx.Model(&model.GoodsDetail{}).Create(goodsDetail)
	if res.Error != nil {
		tx.Rollback()
		return nil, GormGoodsErrHandle(res.Error)
	}

	tx.Commit()
	return &proto.GoodsDetailRes{}, nil
}

func (s *GormGoodsData) DeleteGoodsDB(ctx context.Context, in *proto.DelGoodsReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormGoodsData) UpdeateGoodsDB(ctx context.Context, in *proto.UpdateGoodsReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}
