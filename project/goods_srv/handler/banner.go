package handler

import (
	"context"
	gb "goods_srv/global"
	"goods_srv/model"
	pb "goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 轮播窗口服务
func (s *GoodsServer) GetBannerList(ctx context.Context, emp *emptypb.Empty) (*pb.BannerListRes, error) {
	res := &pb.BannerListRes{}
	banners := []model.Banner{}
	result := gb.DB.Model(&model.Banner{}).Find(&banners)
	if result.Error != nil {
		return nil, result.Error
	}
	res.Total = result.RowsAffected
	for _, banner := range banners {
		res.Data = append(res.Data, &pb.BannerInfoRes{
			Id:    banner.ID,
			Image: banner.Image,
			Index: banner.Index,
			Url:   banner.Url,
		})
	}
	return res, nil
}

func (s *GoodsServer) CreateBanner(ctx context.Context, req *pb.BannerInfoReq) (*pb.BannerInfoRes, error) {
	banner := model.Banner{
		Url:   req.Url,
		Index: req.Index,
		Image: req.Image,
	}
	result := gb.DB.Model(&model.Banner{}).Create(&banner)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pb.BannerInfoRes{Id: banner.ID}, nil
}

//嘶,或许需要通过req的各项属性来查询更新?
// func (s *GoodsServer) DeleteBanner(ctx context.Context, req *pb.DelBrandReq) (*emptypb.Empty, error) {
// 	result := gb.DB.Delete(&model.Banner{}, req.Id)
// 	if result.RowsAffected == 0 {
// 		return nil, ErrBannerNotExists
// 	}
// 	return &emptypb.Empty{}, nil
// }

// 嘶,或许需要通过req的各项属性来查询更新?
func (s *GoodsServer) UpdateBanner(ctx context.Context, req *pb.BannerInfoReq) (*emptypb.Empty, error) {
	return nil, nil
}
