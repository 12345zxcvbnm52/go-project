package handler

import (
	pb "goods_srv/proto"

	"gorm.io/gorm"
)

type GoodsServer struct {
	pb.UnimplementedGoodsServer
}

// gorm给出的分页函数的最佳实践
func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize < 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// func (s *GoodsServer) GetGoodList(context.Context, *pb.GoodsFilterReq) (*pb.GoodsListRes, error) {

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
