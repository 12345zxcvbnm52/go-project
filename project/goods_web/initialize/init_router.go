package initialize

import (
	"fmt"
	//gb "user_web/gb"
	gb "goods_web/global"
	"goods_web/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gb.Router = gin.Default()
	gb.GoodsRter = gb.Router.Group(fmt.Sprintf("/v%s/goods", gb.ServerConfig.Version))
	{
		gb.GoodsRter.GET("/list", handler.GetGoodsList)
		gb.GoodsRter.POST("/newgoods", handler.CreateGoods)
	}
	gb.BannerRter = gb.Router.Group(fmt.Sprintf("/v%s/banner", gb.ServerConfig.Version))
	gb.CategoryRter = gb.Router.Group(fmt.Sprintf("/v%s/category", gb.ServerConfig.Version))
	gb.BrandRter = gb.Router.Group(fmt.Sprintf("/v%s/brand", gb.ServerConfig.Version))
	gb.GoodsRter = gb.Router.Group(fmt.Sprintf("/v%s/categy_brand", gb.ServerConfig.Version))
}
