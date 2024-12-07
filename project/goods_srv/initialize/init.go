package initialize

import (
	gb "goods_srv/global"

	"go.uber.org/zap"
)

func init() {
	InitLog()
	InitConfig()
	InitDB()
	InitConsul()
	//gb.DB.AutoMigrate(&model.Category{}, &model.Goods{}, &model.Brand{}, &model.CategoryWithBrand{}, &model.Banner{})
	zap.S().Info("ServerConfig is : ", gb.ServerConfig)
}
