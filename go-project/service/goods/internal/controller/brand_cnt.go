package goodscontroller

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌服务
func (s *GoodsServer) GetBrandList(ctx context.Context, in *proto.BrandFilterReq) (*proto.BrandListRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetBrandList调用,调用信息为: %s", info)
	res, err := s.Service.GetBrandListLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetBrandList失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) CreateBrand(ctx context.Context, in *proto.CreateBrandReq) (*proto.BrandInfoRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次CreateBrand调用,调用信息为: %s", info)
	res, err := s.Service.CreateBrandLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用CreateBrand失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) DeleteBrand(ctx context.Context, in *proto.DelBrandReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次DeleteBrand调用,调用信息为: %s", info)
	res, err := s.Service.DeleteBrandLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用DeleteBrand失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) UpdateBrand(ctx context.Context, in *proto.UpdateBrandReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次UpdateBrand调用,调用信息为: %s", info)
	res, err := s.Service.UpdateBrandLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用UpdateBrand失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}
