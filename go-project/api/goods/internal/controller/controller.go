package goodscontroller

import (
	goodsdata "kenshop/api/goods/internal/data"
	"kenshop/goken/server/httpserver"
	"kenshop/pkg/log"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// 默认使用otelzap.Logger以及GrpcGoodsData
func MustNewGoodsHTTPServer(s *httpserver.Server, opts ...OptionFunc) *GoodsHttpServer {
	ss := &GoodsHttpServer{
		Server: s,
	}
	for _, opt := range opts {
		opt(ss)
	}
	if ss.Logger == nil {
		ss.Logger = log.MustNewOtelLogger()
	}
	if ss.GoodsData == nil {
		cli, err := s.GrpcCli.Dial()
		if err != nil {
			panic(err)
		}
		ss.GoodsData = goodsdata.MustNewGrpcGoodsData(cli)
	}
	return ss
}

// @BasePath /
// @Description Goods management service API
// @Host NULL
// @Title Goods Service API
// @Version 1.0.0
type GoodsHttpServer struct {
	Server    *httpserver.Server
	GoodsData goodsdata.GoodsDataService
	Logger    *otelzap.Logger
}

func (s *GoodsHttpServer) Execute() error {
	s.Server.Engine.GET("/goods", s.GetGoodList)
	s.Server.Engine.GET("/goods/ids", s.GetGoodsListById)
	s.Server.Engine.GET("/goods/detail/:id", s.GetGoodsDetail)
	s.Server.Engine.POST("/goods", s.CreateGoods)
	s.Server.Engine.DELETE("/goods/:id", s.DeleteGoods)
	s.Server.Engine.PATCH("/goods/:id", s.UpdeateGoods)
	s.Server.Engine.GET("/categories", s.GetCategoryList)
	s.Server.Engine.GET("/category/:id", s.GetCategoryInfo)
	s.Server.Engine.POST("/categories", s.CreateCategory)
	s.Server.Engine.DELETE("/category/:id", s.DeleteCategory)
	s.Server.Engine.PATCH("/category/:id", s.UpdateCategory)
	s.Server.Engine.GET("/brands", s.GetBrandList)
	s.Server.Engine.POST("/brand", s.CreateBrand)
	s.Server.Engine.DELETE("/brand/:id", s.DeleteBrand)
	s.Server.Engine.PATCH("/brand/:id", s.UpdateBrand)
	s.Server.Engine.GET("/banners", s.GetBannerList)
	s.Server.Engine.POST("/banners", s.CreateBanner)
	s.Server.Engine.DELETE("/banner/:id", s.DeleteBanner)
	s.Server.Engine.PATCH("/banner/:id", s.UpdateBanner)
	s.Server.Engine.GET("/categorybrands", s.GetCategoryBrandList)
	s.Server.Engine.POST("/categorybrands", s.CreateCategoryBrand)
	s.Server.Engine.DELETE("/categorybrand/:id", s.DeleteCategoryBrand)
	s.Server.Engine.PATCH("/categorybrand/:id", s.UpdateCategoryBrand)
	return s.Server.Serve()
}

type OptionFunc func(*GoodsHttpServer)

func WithLogger(l *otelzap.Logger) OptionFunc {
	return func(s *GoodsHttpServer) {
		s.Logger = l
	}
}

func WithGoodsDataService(s goodsdata.GoodsDataService) OptionFunc {
	return func(h *GoodsHttpServer) {
		h.GoodsData = s
	}
}
