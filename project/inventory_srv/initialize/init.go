package initialize

import (
	gb "inventory_srv/global"
	"inventory_srv/model"
)

//gb "inventory/global"
//"inventory_srv/model"

func init() {
	InitLog()
	InitConfig()
	InitDB()
	InitRedis()
	//InitConfig最后在main函数里调用
	//InitConsul()
	gb.DB.AutoMigrate(&model.Inventory{})
}
