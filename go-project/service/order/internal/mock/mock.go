package main

import (
	model "kenshop/service/order/internal/model"
	"kenshop/service/order/internal/resourse"
)

func main() {
	carts := []*model.Cart{}
	carts = append(carts, &model.Cart{GoodsId: 2, GoodsNums: 21, UserId: 1, Selected: true})
	carts = append(carts, &model.Cart{GoodsId: 1, GoodsNums: 11, UserId: 1, Selected: true})
	carts = append(carts, &model.Cart{GoodsId: 3, GoodsNums: 146, UserId: 1, Selected: true})
	carts = append(carts, &model.Cart{GoodsId: 7, GoodsNums: 21, UserId: 1, Selected: false})
	carts = append(carts, &model.Cart{GoodsId: 4, GoodsNums: 31, UserId: 1, Selected: true})
	res := resourse.OrderData.DB.Model(&model.Cart{}).Create(&carts)
	if res.Error != nil {
		panic(res.Error)
	}

}
