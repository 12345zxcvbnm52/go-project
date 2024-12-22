package handler

import (
	"context"
	gb "goods_srv/global"
	"goods_srv/model"
	pb "goods_srv/proto"

	"go.uber.org/zap"
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
	if err := banner.InsertOne(); err != nil {
		zap.S().Errorw("一件商品banner插入失败", "msg", err.Error())
		return nil, err
	}
	return &pb.BannerInfoRes{Id: banner.ID}, nil
}

func (s *GoodsServer) DeleteBanner(ctx context.Context, req *pb.DelBrandReq) (*emptypb.Empty, error) {
	banner := &model.Banner{}
	banner.ID = req.Id
	err := banner.DeleteOneById()
	if err != nil {
		zap.S().Errorw("一件商品banner删除失败", "msg", err.Error())
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateBanner(ctx context.Context, req *pb.BannerInfoReq) (*emptypb.Empty, error) {
	banner := &model.Banner{}
	banner.ID = req.Id
	banner.Url = req.Url
	banner.Image = req.Image
	banner.Index = req.Index
	if err := banner.UpdateOneById(); err != nil {
		zap.S().Errorw("一件商品banner更改失败", "msg", err.Error())
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
