package initialize

import (
	"fmt"
	//gb "user_web/global"
	"goods_web/global"
	"goods_web/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	global.Router = gin.Default()
	global.GoodsRter = global.Router.Group(fmt.Sprintf("/v%s/goods", global.ServerConfig.Version))
	{
		global.GoodsRter.GET("/list", handler.GetGoodsList)
	}
	global.BannerRter = global.Router.Group(fmt.Sprintf("/v%s/banner", global.ServerConfig.Version))
	global.CategoryRter = global.Router.Group(fmt.Sprintf("/v%s/category", global.ServerConfig.Version))
	global.BrandRter = global.Router.Group(fmt.Sprintf("/v%s/brand", global.ServerConfig.Version))
	global.GoodsRter = global.Router.Group(fmt.Sprintf("/v%s/categy_brand", global.ServerConfig.Version))
}
