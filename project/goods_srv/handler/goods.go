package handler

import (
	"context"
	"goods_srv/model"
	pb "goods_srv/proto"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GoodsServer struct {
	pb.UnimplementedGoodsServer
}

func GoodsToGoodsInfoRes(goods *model.Goods) *pb.GoodsInfoRes {
	return &pb.GoodsInfoRes{
		Id:          goods.ID,
		CategyId:    goods.CategoryID,
		Name:        goods.Name,
		GoodsSign:   goods.GoodSign,
		ClickNum:    goods.ClickNum,
		SoldNum:     goods.SoldNum,
		FavorNum:    goods.FavorNum,
		MarketPrice: goods.MarketPrice,
		SalePrice:   goods.SalePrice,
		GoodsBrief:  goods.GoodsBrief,
		TransFree:   goods.TransFree,
		FirstImage:  goods.FirstImage,
		IsNew:       goods.IsNew,
		IsHot:       goods.IsHot,
		OnSale:      goods.OnSale,
		DescImages:  goods.DescImages,
		Images:      goods.Images,
		Categy: &pb.CategyBriefInfoRes{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &pb.BrandInfoRes{
			Id:   goods.Brand.ID,
			Name: goods.Brand.Name,
			Logo: goods.Brand.Logo,
		},
	}
}

func (s *GoodsServer) GetGoodList(ctx context.Context, req *pb.GoodsFilterReq) (*pb.GoodsListRes, error) {
	logic := &model.Goods{}
	res, err := logic.FindByOpt(&model.FindOption{
		PageSize: req.PageSize,
		IsHot:    req.IsHot,
		IsNew:    req.IsNew,
		MinPrice: req.MinPrice,
		MaxPrice: req.MaxPrice,
		CategyId: req.CategyId,
		OnTable:  req.OnTable,
		PagesNum: req.PagesNum,
		KeyWords: req.KeyWords,
		Brand:    req.Brand,
	})
	if err != nil {
		return nil, err
	}
	r := []*pb.GoodsInfoRes{}
	for _, v := range res.Data {
		r = append(r, GoodsToGoodsInfoRes(v))
	}
	return &pb.GoodsListRes{
		Total: res.Total,
		Data:  r,
	}, nil
}

// 用于通过id数组得到所有商品信息,常用于从订单中获得所有商品信息,
func (s *GoodsServer) GetGoodsListById(ctx context.Context, req *pb.GoodsIdsReq) (*pb.GoodsListRes, error) {
	logic := &model.Goods{}
	ans, err := logic.FindByIds(req.Id...)
	if err != nil {
		zap.S().Errorw("用户按id批量查询商品失败", "msg", err.Error())
		return nil, err
	}
	res := pb.GoodsListRes{}
	res.Total = ans.Total
	for _, v := range ans.Data {
		res.Data = append(res.Data, GoodsToGoodsInfoRes(v))
	}
	return &res, nil
}

// 增删改
func (s *GoodsServer) CreateGoods(ctx context.Context, req *pb.WriteGoodsInfoReq) (*pb.GoodsInfoRes, error) {
	goods := &model.Goods{
		CategoryID:  req.CategyId,
		BrandID:     req.BrandId,
		Name:        req.Name,
		GoodSign:    req.GoodsSign,
		GoodsBrief:  req.GoodsBrief,
		ClickNum:    0,
		FavorNum:    0,
		SoldNum:     0,
		FirstImage:  req.FirstImage,
		DescImages:  req.DescImages,
		Images:      req.Images,
		OnSale:      true,
		IsHot:       false,
		IsNew:       true,
		TransFree:   req.TransFree,
		MarketPrice: req.MarketPrice,
		SalePrice:   req.SalePrice,
	}
	if err := goods.InsertOne(); err != nil {
		return nil, err
	}
	return &pb.GoodsInfoRes{Id: goods.ID}, nil
}

func (s *GoodsServer) UpdeateGoods(ctx context.Context, req *pb.WriteGoodsInfoReq) (*emptypb.Empty, error) {
	goods := &model.Goods{
		CategoryID:  req.CategyId,
		BrandID:     req.BrandId,
		Name:        req.Name,
		GoodSign:    req.GoodsSign,
		GoodsBrief:  req.GoodsBrief,
		FirstImage:  req.FirstImage,
		DescImages:  req.DescImages,
		Images:      req.Images,
		OnSale:      req.OnSale,
		TransFree:   req.TransFree,
		MarketPrice: req.MarketPrice,
		SalePrice:   req.SalePrice,
	}
	if err := goods.UpdateOneByid(); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// func (s *GoodsServer) DeleteGoods(context.Context, *pb.DelGoodsReq) (*emptypb.Empty, error) {

// }

// func (s *GoodsServer) GetGoodsDetail(context.Context, *pb.GoodsInfoReq) (*pb.GoodsInfoRes, error) {

// }
