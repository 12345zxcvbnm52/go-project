package handler

import (
	"context"
	gb "goods_srv/global"
	"goods_srv/model"
	pb "goods_srv/proto"
)

// 品牌分类服务
func (s *GoodsServer) GetCategyBrandList(ctx context.Context, req *pb.CategyBrandFilterReq) (*pb.CategyBrandListRes, error) {
	categyBrands := []model.CategoryWithBrand{}
	res := &pb.CategyBrandListRes{}
	result := gb.DB.Model(&model.CategoryWithBrand{}).Count(&res.Total)
	if result.Error != nil {
		return nil, result.Error
	}
	result = gb.DB.Preload("Category").Preload("Brand").Scopes(Paginate(int(req.PagesNum), int(req.PageSize))).Find(&categyBrands)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, categyBrand := range categyBrands {
		res.Data = append(res.Data, &pb.CategyBrandInfoRes{
			Id: categyBrand.ID,
			CategyInfo: &pb.CategyInfoRes{
				Id:             categyBrand.Category.ID,
				Name:           categyBrand.Category.Name,
				ParentCategyId: categyBrand.Category.ParentCategoryID,
				Level:          categyBrand.Category.Level,
				OnTable:        categyBrand.Category.OnTab,
			},
			BrandInfo: &pb.BrandInfoRes{
				Id:   categyBrand.Brand.ID,
				Logo: categyBrand.Brand.Logo,
				Name: categyBrand.Brand.Name,
			},
		})
	}
	return res, nil
}

// // 通过一个已经确定的类型获得所有有这个类型的品牌
// // 例如我想看看食品-面包里的几个品牌,有达利园,XXX等品牌是需要查到的
// func (s *GoodsServer) GetBrandListByCategy(context.Context, *pb.CategyInfoReq) (*pb.BrandListRes, error) {
// }
// func (s *GoodsServer) CreateCategyBrand(context.Context, *pb.CategyBrandInfoReq) (*pb.CategyBrandInfoRes, error) {
// }
// func (s *GoodsServer) DeleteCategyBrand(context.Context, *pb.CategyBrandInfoReq) (*emptypb.Empty, error) {
// }
// func (s *GoodsServer) UpdateCategyBrand(context.Context, *pb.CategyBrandInfoReq) (*emptypb.Empty, error) {
// }
