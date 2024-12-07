package handler

import pb "goods_srv/proto"

type CategyBrandServer struct {
	pb.UnimplementedGoodsServer
}

// // 通过一个类型获得所有有这个类型的品牌
// func (s *GoodsServer) GetBrandListByCategy(context.Context, *pb.CategyInfoReq) (*pb.BrandListRes, error) {
// }
// func (s *GoodsServer) CreateCategyBrand(context.Context, *pb.CategyBrandInfoReq) (*pb.CategyBrandInfoRes, error) {
// }
// func (s *GoodsServer) DeleteCategyBrand(context.Context, *pb.CategyBrandInfoReq) (*emptypb.Empty, error) {
// }
// func (s *GoodsServer) UpdateCategyBrand(context.Context, *pb.CategyBrandInfoReq) (*emptypb.Empty, error) {
// }
