package handler

import (
	"context"
	"encoding/json"
	gb "goods_srv/global"
	"goods_srv/model"
	pb "goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type CategoryServer struct {
	pb.UnimplementedGoodsServer
}

func CategoryToCategyInfoRes(categy *model.Category) *pb.CategyInfoRes {
	return &pb.CategyInfoRes{
		Id:             categy.ID,
		Name:           categy.Name,
		ParentCategyId: *categy.ParentCategoryID,
		Level:          categy.Level,
		OnTable:        categy.OnTab,
	}
}

// 商品类型服务
func (s *GoodsServer) GetCategyList(ctx context.Context, req *emptypb.Empty) (*pb.CategyListRes, error) {
	r := &pb.CategyListRes{}
	categys := []model.Category{}
	res := gb.DB.Where(&model.Category{Level: gb.TopLevel}).Preload("SubCategory.SubCategory").Find(&categys)
	if res.Error != nil {
		return nil, model.ErrInternalWrong
	}
	if res.RowsAffected == 0 {
		return nil, model.ErrCategoryNotFound
	}
	b, _ := json.Marshal(&categys)
	//http/2自带压缩,故可以直接发送长字符串
	r.JsonData = string(b)
	return r, nil
}

func (s *GoodsServer) GetCategyInfo(ctx context.Context, req *pb.SubCategyReq) (*pb.SubCategyListRes, error) {
	category := &model.Category{}
	if result := gb.DB.First(&category, req.Id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, model.ErrCategoryNotFound
		}
		return nil, model.ErrInternalWrong
	}
	res := pb.SubCategyListRes{}
	res.SelfInfo = &pb.CategyInfoRes{
		Id:      category.ID,
		Name:    category.Name,
		Level:   category.Level,
		OnTable: category.OnTab,
	}
	if category.ParentCategoryID != nil {
		res.SelfInfo.ParentCategyId = *category.ParentCategoryID
	}
	subCategy := []*model.Category{}
	if category.Level != gb.EndLevel {
		result := gb.DB.Model(&model.Category{}).Where("parent_category_id = ?", category.ID).Find(&subCategy)
		if result.Error != nil {
			return nil, model.ErrInternalWrong
		}
		if result.RowsAffected == 0 {
			return nil, model.ErrCategoryNotFound
		}
		for _, v := range subCategy {
			res.SubInfo = append(res.SubInfo, CategoryToCategyInfoRes(v))
		}
		res.Total = int64(len(subCategy))
	} else {
		res.SubInfo = nil
	}
	return &res, nil
}

func (s *GoodsServer) CreateCategy(ctx context.Context, req *pb.CategyInfoReq) (*pb.CategyInfoRes, error) {
	category := model.Category{}
	if req.Level != gb.TopLevel {
		category.ParentCategoryID = new(uint32)
		*category.ParentCategoryID = req.ParentCategyId
	}
	category.Level = req.Level
	category.Name = req.Name
	category.OnTab = req.OnTable
	if err := category.InsertOne(); err != nil {
		return nil, err
	}
	return &pb.CategyInfoRes{Id: category.ID}, nil
}

func (s *GoodsServer) DeleteCategy(ctx context.Context, req *pb.DelCategyReq) (*emptypb.Empty, error) {
	category := &model.Category{}
	category.ID = req.Id
	if err := category.DeleteOneById(); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// 这里逻辑可以考虑如果传入的数据没有Level就去查找Level
func (s *GoodsServer) UpdateCategy(ctx context.Context, req *pb.CategyInfoReq) (*emptypb.Empty, error) {
	category := &model.Category{}
	category.ID = req.Id
	category.Level = req.Level
	category.Name = req.Name
	category.OnTab = req.OnTable
	if category.Level != gb.TopLevel {
		category.ParentCategoryID = new(uint32)
		*category.ParentCategoryID = req.ParentCategyId
	}
	if err := category.UpdateOneById(); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
