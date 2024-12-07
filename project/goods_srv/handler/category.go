package handler

import pb "goods_srv/proto"

type CategoryServer struct {
	pb.UnimplementedGoodsServer
}

// // 商品类型服务
// func (s *GoodsServer) GetAllCategyList(context.Context, *emptypb.Empty) (*pb.CategyListRes, error)  {}
// func (s *GoodsServer) GetSubCategy(context.Context, *pb.SubCategyReq) (*pb.SubCategyListRes, error) {}
// func (s *GoodsServer) CreateCategy(context.Context, *pb.CategyInfoReq) (*pb.CategyInfoRes, error)   {}
// func (s *GoodsServer) DeleteCategy(context.Context, *pb.DelCategyReq) (*emptypb.Empty, error)       {}
// func (s *GoodsServer) UpdateCategy(context.Context, *pb.CategyInfoReq) (*emptypb.Empty, error)      {}
