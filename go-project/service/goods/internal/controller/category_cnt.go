package goodscontroller

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 商品类型服务
func (s *GoodsServer) GetCategoryList(ctx context.Context, in *emptypb.Empty) (*proto.CategoryListRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetCategoryList调用,调用信息为: %s", info)
	res, err := s.Service.GetCategoryListLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetCategoryList失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) GetCategoryInfo(ctx context.Context, in *proto.SubCategoryReq) (*proto.SubCategoryListRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetCategoryInfo调用,调用信息为: %s", info)
	res, err := s.Service.GetCategoryInfoLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetCategoryInfo失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) CreateCategory(ctx context.Context, in *proto.CreateCategoryReq) (*proto.CategoryInfoRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次CreateCategory调用,调用信息为: %s", info)
	res, err := s.Service.CreateCategoryLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用CreateCategory失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, in *proto.DelCategoryReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次DeleteCategory调用,调用信息为: %s", info)
	res, err := s.Service.DeleteCategoryLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用DeleteCategory失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, in *proto.UpdateCategoryReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次UpdateCategory调用,调用信息为: %s", info)
	res, err := s.Service.UpdateCategoryLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用UpdateCategory失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}
