package handler

import (
	"context"
	"goods_srv/model"
	pb "goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌分类服务
func (s *GoodsServer) GetCategyBrandList(ctx context.Context, req *pb.CategyBrandFilterReq) (*pb.CategyBrandListRes, error) {
	logic := model.CategoryWithBrand{}
	res, err := logic.FindByOpt(&model.CategyBrandFindOption{PagesNum: req.PagesNum, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	r := &pb.CategyBrandListRes{}
	r.Total = res.Total
	for _, categyBrand := range res.Data {
		r.Data = append(r.Data, &pb.CategyBrandInfoRes{
			Id:         categyBrand.ID,
			CategyInfo: CategoryToCategyInfoRes(&categyBrand.Category),
			BrandInfo:  BrandToBrandInfoRes(&categyBrand.Brand),
		})
	}
	return r, nil
}

// 通过一个已经确定的类型获得所有有这个类型的品牌
// 例如我想看看食品-面包里的几个品牌,有达利园,XXX等品牌是需要查到的
func (s *GoodsServer) GetBrandListByCategy(ctx context.Context, req *pb.CategyInfoReq) (*pb.BrandListRes, error) {
	logic := model.CategoryWithBrand{}
	res, err := logic.FindByCategyId(req.Id)
	if err != nil {
		return nil, err
	}
	r := &pb.BrandListRes{}
	r.Total = res.Total
	for _, categyBrand := range res.Data {
		r.Data = append(r.Data, BrandToBrandInfoRes(&categyBrand.Brand))
	}
	return r, nil
}

func (s *GoodsServer) CreateCategyBrand(ctx context.Context, req *pb.CategyBrandInfoReq) (*pb.CategyBrandInfoRes, error) {
	categyBrand := &model.CategoryWithBrand{
		CategoryID: req.CategyId,
		BrandID:    req.BrandId,
	}
	if err := categyBrand.InsertOne(); err != nil {
		return nil, err
	}
	return &pb.CategyBrandInfoRes{Id: categyBrand.ID}, nil
}

func (s *GoodsServer) DeleteCategyBrand(ctx context.Context, req *pb.CategyBrandInfoReq) (*emptypb.Empty, error) {
	categyBrand := &model.CategoryWithBrand{
		CategoryID: req.CategyId,
		BrandID:    req.BrandId,
	}
	categyBrand.ID = req.Id
	if err := categyBrand.DeleteOneById(); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategyBrand(ctx context.Context, req *pb.CategyBrandInfoReq) (*emptypb.Empty, error) {
	categyBrand := &model.CategoryWithBrand{
		CategoryID: req.CategyId,
		BrandID:    req.BrandId,
	}
	categyBrand.ID = req.Id
	if err := categyBrand.UpdateOneById(); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
