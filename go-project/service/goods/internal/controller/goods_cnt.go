package goodscontroller

import (
	"context"
	proto "kenshop/proto/goods"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 获得商品列表
func (s *GoodsServer) GetGoodList(ctx context.Context, in *proto.GoodsFilterReq) (*proto.GoodsListRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetGoodList调用,调用信息为: %s", info)
	res, err := s.Service.GetGoodListLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetGoodList失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

// 用于通过id数组得到所有商品信息,常用于从订单中获得所有商品信息,
func (s *GoodsServer) GetGoodsListById(ctx context.Context, in *proto.GoodsIdsReq) (*proto.GoodsListRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetGoodsListById调用,调用信息为: %s", info)
	res, err := s.Service.GetGoodsListByIdLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetGoodsListById失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) GetGoodsDetail(ctx context.Context, in *proto.GoodsInfoReq) (*proto.GoodsDetailRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetGoodsDetail调用,调用信息为: %s", info)
	res, err := s.Service.GetGoodsDetailLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetGoodsDetail失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, in *proto.CreateGoodsReq) (*proto.GoodsDetailRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次CreateGoods调用,调用信息为: %s", info)
	res, err := s.Service.CreateGoodsLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用CreateGoods失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, in *proto.DelGoodsReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次DeleteGoods调用,调用信息为: %s", info)
	res, err := s.Service.DeleteGoodsLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用DeleteGoods失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *GoodsServer) UpdeateGoods(ctx context.Context, in *proto.UpdateGoodsReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次UpdeateGoods调用,调用信息为: %s", info)
	res, err := s.Service.UpdeateGoodsLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用UpdeateGoods失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}
