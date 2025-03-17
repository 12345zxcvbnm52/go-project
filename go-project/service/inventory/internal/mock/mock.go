package main

import (
	model "kenshop/service/inventory/internal/model"
	"kenshop/service/inventory/internal/resourse"
	"math/rand/v2"
)

func main() {
	for i := 1; i <= 1000; i++ {
		record := model.Inventory{
			GoodsId:  uint32(i),
			GoodsNum: rand.Int32N(1000),
		}
		// 插入数据
		result := resourse.InventoryData.DB.Create(&record)
		if result.Error != nil {
			panic(result.Error)
		}
	}

}
