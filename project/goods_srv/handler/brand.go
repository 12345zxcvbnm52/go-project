package handler

import (
	"context"
	gb "goods_srv/global"
	"goods_srv/proto"
	pb "goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌服务
func (s *GoodsServer) GetBrandList(context.Context, *pb.BrandFilterReq) (*pb.BrandListRes, error) {
	res := &proto.BrandListRes{}
	gb.DB.Find()
}
func (s *GoodsServer) CreateBrand(context.Context, *pb.BrandInfoReq) (*pb.BrandInfoRes, error) {}
func (s *GoodsServer) DeleteBrand(context.Context, *pb.DelBrandReq) (*emptypb.Empty, error)    {}
func (s *GoodsServer) UpdateBrand(context.Context, *pb.BrandInfoReq) (*emptypb.Empty, error)   {}
