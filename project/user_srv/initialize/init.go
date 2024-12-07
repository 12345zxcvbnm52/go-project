package initialize

import (
	gb "user_srv/global"
	"user_srv/model"

	"go.uber.org/zap"
)

func init() {
	InitLog()
	InitConfig()
	InitDB()
	InitConsul()
	gb.DB.AutoMigrate(&model.User{})
	zap.S().Infoln("ServerConfig is : ", gb.ServerConfig)
}
