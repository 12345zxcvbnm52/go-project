package initialize

import (
	"fmt"
	//gb "user_web/global"
	"user_web/global"
	"user_web/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter() {

	global.Router = gin.Default()
	global.UseRter = global.Router.Group(fmt.Sprintf("/v%s", global.ServerConfig.Version))
	global.UseRter.POST("/login", handler.UserLogin)
}
