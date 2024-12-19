package initialize

import (
	gb "goods_srv/global"
	"goods_srv/model"
)

func init() {
	InitLog()
	InitConfig()
	InitDB()
	//InitConsul()
	gb.DB.AutoMigrate(&model.Category{}, &model.Goods{}, &model.Brand{}, &model.CategoryWithBrand{}, &model.Banner{})

}
