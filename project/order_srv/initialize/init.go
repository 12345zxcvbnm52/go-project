package initialize

import (
	gb "order_srv/global"
	"order_srv/model"
)

func init() {
	InitLog()
	InitConfig()
	InitDB()
	InitRedis()
	//InitConfig最后在main函数里调用
	//InitConsul()
	gb.DB.AutoMigrate(&model.Cart{}, &model.Order{}, &model.OrderGoods{})
}
