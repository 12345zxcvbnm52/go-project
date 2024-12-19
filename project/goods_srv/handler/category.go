package handler

import (
	"context"
	"encoding/json"
	gb "goods_srv/global"
	"goods_srv/model"
	pb "goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryServer struct {
	pb.UnimplementedGoodsServer
}

var ErrCategoryNotFound error = status.Errorf(codes.NotFound, "商品分类不存在")

// 商品类型服务
func (s *GoodsServer) GetAllCategyList(ctx context.Context, req *emptypb.Empty) (*pb.CategyListRes, error) {
	res := &pb.CategyListRes{}
	categys := []model.Category{}
	gb.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categys)
	b, _ := json.Marshal(&categys)
	//http/2自带压缩,故可以直接发送长字符串
	res.JsonData = string(b)
	return res, nil
}

func (s *GoodsServer) GetSubCategy(ctx context.Context, req *pb.SubCategyReq) (*pb.SubCategyListRes, error) {
	category := &model.Category{}
	result := gb.DB.First(&category, req.Id)
	if result.RowsAffected == 0 {
		return nil, ErrCategoryNotFound
	}
	res := pb.SubCategyListRes{}
	res.SelfInfo = &pb.CategyInfoRes{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		OnTable:        category.OnTab,
		ParentCategyId: category.ParentCategoryID,
	}
	subCategy := []*pb.CategyInfoRes{}
	switch category.Level {
	case 1:
		result = gb.DB.Where(&model.Category{Level: 1}).Preload("SubCategory").Find(&subCategy)
		res.SubInfo = append(res.SubInfo, subCategy...)
	case 2:
		result = gb.DB.Where(&model.Category{Level: 2}).Preload("SubCategory.SubCategory").Find(&subCategy)
		res.SubInfo = append(res.SubInfo, subCategy...)
	case 3:
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (s *GoodsServer) CreateCategy(ctx context.Context, req *pb.CategyInfoReq) (*pb.CategyInfoRes, error) {
	category := model.Category{}
	if req.Level != 1 {
		category.ParentCategoryID = req.ParentCategyId
	}
	category.Level = req.Level
	category.Name = req.Name
	category.OnTab = req.OnTable
	result := gb.DB.Model(&model.Category{}).Create(&category)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pb.CategyInfoRes{Id: category.ID}, nil
}

// func (s *GoodsServer) DeleteCategy(context.Context, *pb.DelCategyReq) (*emptypb.Empty, error)       {}
// func (s *GoodsServer) UpdateCategy(context.Context, *pb.CategyInfoReq) (*emptypb.Empty, error)      {}
