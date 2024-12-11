package handler

import (
	"context"
	"errors"
	gb "goods_srv/global"
	"goods_srv/model"
	pb "goods_srv/proto"
)

var (
	ErrBrandAlreadExists = errors.New("选择的品牌已存在")
)

// 品牌服务
func (s *GoodsServer) GetBrandList(ctx context.Context, req *pb.BrandFilterReq) (*pb.BrandListRes, error) {
	res := &pb.BrandListRes{}
	brands := []model.Brand{}

	result := gb.DB.Scopes(Paginate(int(req.PagesNum), int(req.PageSize))).Find(brands)
	if result.Error != nil {
		return nil, result.Error
	}
	gb.DB.Model(&model.Brand{}).Count(&res.Total)
	for _, brand := range brands {
		res.Data = append(res.Data, &pb.BrandInfoRes{
			Id:   brand.ID,
			Logo: brand.Logo,
			Name: brand.Name,
		})
	}
	return res, nil
}

func (s *GoodsServer) CreateBrand(ctx context.Context, req *pb.BrandInfoReq) (*pb.BrandInfoRes, error) {
	result := gb.DB.Where("Name=?", req.Name).Find(&model.Brand{})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return nil, ErrBrandAlreadExists
	}
	gb.DB.Model(&model.Brand{}).Create(&model.Brand{
		Name: req.Name,
		Logo: req.Logo,
	})
	return &pb.BrandInfoRes{Id: req.Id}, nil
}

// func (s *GoodsServer) DeleteBrand(context.Context, *pb.DelBrandReq) (*emptypb.Empty, error)    {}
// func (s *GoodsServer) UpdateBrand(context.Context, *pb.BrandInfoReq) (*emptypb.Empty, error)   {}
