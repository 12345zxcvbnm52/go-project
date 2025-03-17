package goodsdata

import (
	"context"
	"encoding/json"
	"errors"
	proto "kenshop/proto/goods"
	model "kenshop/service/goods/internal/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 商品类型服务
func (s *GormGoodsData) GetCategoryListDB(ctx context.Context, in *emptypb.Empty) (*proto.CategoryListRes, error) {
	r := &proto.CategoryListRes{}
	categorys := []*model.Category{}
	res := s.DB.Where("level = ?", model.TopLevel).Preload("SubCategory.SubCategory").Find(&categorys)
	if res.Error != nil {
		return nil, GormCategoryErrHandle(res.Error)
	}
	b, _ := json.Marshal(&categorys)
	//http/2自带压缩,故可以直接发送长字符串
	r.JsonData = string(b)
	return r, nil
}

func (s *GormGoodsData) GetCategoryInfoDB(ctx context.Context, in *proto.SubCategoryReq) (*proto.SubCategoryListRes, error) {
	category := &model.Category{}
	if result := s.DB.First(&category, in.Id); result.Error != nil {
		return nil, GormCategoryErrHandle(result.Error)
	}
	res := proto.SubCategoryListRes{}
	res.SelfInfo = &proto.CategoryInfoRes{
		Id:      category.ID,
		Name:    category.Name,
		Level:   category.Level,
		OnTable: category.OnTable,
	}
	if category.ParentCategoryID != nil {
		res.SelfInfo.ParentCategoryId = *category.ParentCategoryID
	}
	subCategy := []*model.Category{}
	if category.Level != model.EndLevel {
		result := s.DB.Model(&model.Category{}).Where("parent_category_id = ?", category.ID).Find(&subCategy)
		if result.Error != nil {
			return nil, GormCategoryErrHandle(result.Error)
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

func (s *GormGoodsData) CreateCategoryDB(ctx context.Context, in *proto.CreateCategoryReq) (*proto.CategoryInfoRes, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormGoodsData) DeleteCategoryDB(ctx context.Context, in *proto.DelCategoryReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormGoodsData) UpdateCategoryDB(ctx context.Context, in *proto.UpdateCategoryReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}
