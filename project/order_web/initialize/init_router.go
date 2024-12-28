package initialize

import (
	"fmt"
	//gb "user_web/global"
	"order_web/global"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	global.Router = gin.Default()
	global.CartRter = global.Router.Group(fmt.Sprintf("/v%s/cart", global.ServerConfig.Version))
	global.OrderRter = global.Router.Group(fmt.Sprintf("/v%s/router", global.ServerConfig.Version))

}
