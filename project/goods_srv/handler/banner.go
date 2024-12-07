package handler

import pb "goods_srv/proto"

type BannerServer struct {
	pb.UnimplementedGoodsServer
}

// // 轮播窗口服务
// func (s *GoodsServer) GetBannerList(context.Context, *emptypb.Empty) (*pb.BannerListRes, error)   {}
// func (s *GoodsServer) CreateBanner(context.Context, *pb.BannerInfoReq) (*pb.BannerInfoRes, error) {}
// func (s *GoodsServer) DeleteBanner(context.Context, *pb.DelBrandReq) (*emptypb.Empty, error)      {}
// func (s *GoodsServer) UpdateBanner(context.Context, *pb.BannerInfoReq) (*emptypb.Empty, error)    {}
