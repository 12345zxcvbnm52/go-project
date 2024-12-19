package initialize

import (
	gb "user_srv/global"
	"user_srv/model"
)

func init() {
	InitLog()
	InitConfig()
	InitDB()
	//InitConfig最后在main函数里调用
	//InitConsul()
	gb.DB.AutoMigrate(&model.User{})

}
