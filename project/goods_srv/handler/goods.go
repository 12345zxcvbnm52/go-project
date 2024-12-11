package handler

import (
	pb "goods_srv/proto"
)

type GoodsServer struct {
	pb.UnimplementedGoodsServer
}

// func (s *GoodsServer) GetGoodList(ctx context.Context, req *pb.GoodsFilterReq) (*pb.GoodsListRes, error) {
// 	locDB := gb.DB.Model(&model.Goods{})
// 	if req.IsHot{}
// }

// // 用于通过id数组得到所有商品信息,常用于从订单中获得所有商品信息,
// func (s *GoodsServer) GetGoodsListById(context.Context, *pb.BatchGoodsByIdReq) (*pb.GoodsListRes, error) {
// }

// // 增删改
// func (s *GoodsServer) CreateGoods(context.Context, *pb.WriteGoodsInfoReq) (*pb.GoodsInfoRes, error) {

// }
// func (s *GoodsServer) DeleteGoods(context.Context, *pb.DelGoodsReq) (*emptypb.Empty, error) {

// }
// func (s *GoodsServer) UpdeateGoods(context.Context, *pb.WriteGoodsInfoReq) (*emptypb.Empty, error) {

// }
// func (s *GoodsServer) GetGoodsDetail(context.Context, *pb.GoodsInfoReq) (*pb.GoodsInfoRes, error) {

// }
