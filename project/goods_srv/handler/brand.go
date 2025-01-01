package handler

import (
	"context"
	"goods_srv/model"
	pb "goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

func BrandToBrandInfoRes(brand *model.Brand) *pb.BrandInfoRes {
	return &pb.BrandInfoRes{
		Name: brand.Name,
		Logo: brand.Logo,
		Id:   brand.ID,
	}
}

// 品牌服务
func (s *GoodsServer) GetBrandList(ctx context.Context, req *pb.BrandFilterReq) (*pb.BrandListRes, error) {
	logic := &model.Brand{}
	res, err := logic.FindByOpt(&model.BrandFindOption{PagesNum: req.PagesNum, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	r := &pb.BrandListRes{}
	for _, brand := range res.Data {
		r.Data = append(r.Data, BrandToBrandInfoRes(brand))
	}
	r.Total = res.Total
	return r, nil
}

func (s *GoodsServer) CreateBrand(ctx context.Context, req *pb.BrandInfoReq) (*pb.BrandInfoRes, error) {
	brand := &model.Brand{
		Name: req.Name,
		Logo: req.Logo,
	}
	if err := brand.InsertOne(); err != nil {
		return nil, err
	}
	return &pb.BrandInfoRes{Id: brand.ID}, nil
}

func (s *GoodsServer) DeleteBrand(ctx context.Context, req *pb.DelBrandReq) (*emptypb.Empty, error) {
	brand := &model.Brand{}
	brand.ID = req.Id
	if err := brand.DeleteOneById(); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateBrand(ctx context.Context, req *pb.BrandInfoReq) (*emptypb.Empty, error) {
	brand := &model.Brand{}
	brand.ID = req.Id
	brand.Logo = req.Logo
	brand.Name = req.Name
	if err := brand.UpdateOneById(); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
