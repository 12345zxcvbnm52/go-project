package goodscontroller

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌分类服务
func (s *GoodsServer) GetCategoryBrandList(ctx context.Context, in *proto.CategoryBrandFilterReq) (*proto.CategoryBrandListRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetCategoryBrandList调用,调用信息为: %s", info)
	res, err := s.Service.GetCategoryBrandListLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetCategoryBrandList失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

// 通过一个类型获得所有有这个类型的品牌
// rpc GetBrandListByCategory(CategoryInfoReq)returns(BrandListRes);
func (s *GoodsServer) CreateCategoryBrand(ctx context.Context, in *proto.CreateCategoryBrandReq) (*proto.CategoryBrandInfoRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次CreateCategoryBrand调用,调用信息为: %s", info)
	res, err := s.Service.CreateCategoryBrandLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用CreateCategoryBrand失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) DeleteCategoryBrand(ctx context.Context, in *proto.DelCategoryBrandReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次DeleteCategoryBrand调用,调用信息为: %s", info)
	res, err := s.Service.DeleteCategoryBrandLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用DeleteCategoryBrand失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) UpdateCategoryBrand(ctx context.Context, in *proto.UpdateCategoryBrandReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次UpdateCategoryBrand调用,调用信息为: %s", info)
	res, err := s.Service.UpdateCategoryBrandLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用UpdateCategoryBrand失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}
