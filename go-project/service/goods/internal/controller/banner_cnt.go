package goodscontroller

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 轮播窗口服务
func (s *GoodsServer) GetBannerList(ctx context.Context, in *emptypb.Empty) (*proto.BannerListRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetBannerList调用,调用信息为: %s", info)
	res, err := s.Service.GetBannerListLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetBannerList失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) CreateBanner(ctx context.Context, in *proto.CreateBannerReq) (*proto.BannerInfoRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次CreateBanner调用,调用信息为: %s", info)
	res, err := s.Service.CreateBannerLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用CreateBanner失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) DeleteBanner(ctx context.Context, in *proto.DelBannerReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次DeleteBanner调用,调用信息为: %s", info)
	res, err := s.Service.DeleteBannerLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用DeleteBanner失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) UpdateBanner(ctx context.Context, in *proto.UpdateBannerReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次UpdateBanner调用,调用信息为: %s", info)
	res, err := s.Service.UpdateBannerLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用UpdateBanner失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}
