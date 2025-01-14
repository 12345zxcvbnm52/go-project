package initialize

import (
	"fmt"
	//gb "user_web/gb"
	gb "order_web/global"
	"order_web/handler"
	"order_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gb.Router = gin.Default()
	gb.CartRter = gb.Router.Group(fmt.Sprintf("/v%s/cart", gb.ServerConfig.Version))
	gb.OrderRter = gb.Router.Group(fmt.Sprintf("/v%s/order", gb.ServerConfig.Version), middlewares.TraceMarking())
	{
		gb.OrderRter.POST("insert", handler.CreateOrder)
	}
}
