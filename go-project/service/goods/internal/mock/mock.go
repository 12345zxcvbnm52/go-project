package main

import (
	"fmt"
	model "kenshop/service/goods/internal/model"
	"kenshop/service/goods/internal/resourse"
	"math/rand"

	"gorm.io/gorm"
)

func main() {
	Brand()
	Banner()
	Category()
	GoodsAndGoodsDetail()
	CategoryBrand()
}

func Brand() {
	for i := 1; i <= 500; i++ {
		brand := model.Brand{
			Name: fmt.Sprintf("Brand %d", i),
			Logo: fmt.Sprintf("https://example.com/logos/%d.png", rand.Intn(1000)),
		}
		res := resourse.GoodsData.DB.Model(&model.Brand{}).Create(&brand)
		if res.Error != nil {
			panic(res.Error)
		}
	}
}

func Banner() {
	// 生成 Banner 数据
	for i := 1; i <= 500; i++ {
		banner := model.Banner{
			Image: fmt.Sprintf("https://example.com/banners/%d.jpg", i),
			Url:   fmt.Sprintf("https://example.com/products/%d", i),
			Index: int32(rand.Intn(10) + 1), // 随机生成 1-10 的索引
		}
		res := resourse.GoodsData.DB.Model(&model.Banner{}).Create(&banner)
		if res.Error != nil {
			panic(res.Error)
		}
	}

}

func Category() {
	var level1Categories []model.Category
	for i := 1; i <= 50; i++ {
		category := model.Category{
			Name:    fmt.Sprintf("Level 1 Category %d", i),
			Level:   1,
			OnTable: rand.Intn(2) == 1, // 随机生成 true/false
		}

		res := resourse.GoodsData.DB.Model(&model.Category{}).Create(&category)
		if res.Error != nil {
			panic(res.Error)
		}
		level1Categories = append(level1Categories, category)
	}

	// 生成 100 个二级目录（Level 2），ParentCategoryID 为一级目录的 ID
	var level2Categories []model.Category
	for i := 51; i <= 150; i++ {
		parentID := level1Categories[rand.Intn(50)].ID // 随机选择一个一级目录作为父目录
		category := model.Category{
			Name:             fmt.Sprintf("Level 2 Category %d", i),
			Level:            2,
			OnTable:          rand.Intn(2) == 1,
			ParentCategoryID: &parentID, // 设置父目录 ID
		}

		res := resourse.GoodsData.DB.Model(&model.Category{}).Create(&category)
		if res.Error != nil {
			panic(res.Error)
		}
		level2Categories = append(level2Categories, category)
	}

	// 生成 350 个三级目录（Level 3），ParentCategoryID 为二级目录的 ID
	//var level3Categories []model.Category
	for i := 151; i <= 500; i++ {
		parentID := level2Categories[rand.Intn(100)].ID // 随机选择一个二级目录作为父目录
		category := model.Category{
			Name:             fmt.Sprintf("Level 3 Category %d", i),
			Level:            3,
			OnTable:          rand.Intn(2) == 1,
			ParentCategoryID: &parentID, // 设置父目录 ID
		}
		res := resourse.GoodsData.DB.Model(&model.Category{}).Create(&category)
		if res.Error != nil {
			panic(res.Error)
		}
	}
}

func CategoryBrand() {
	for i := 1; i <= 500; i++ {
		min := 350
		max := 500
		randomNumber := rand.Intn(max-min+1) + min
		categoryBrand := model.CategoryBrand{

			BrandID:    uint32(rand.Intn(500) + 1),
			CategoryID: uint32(randomNumber), // 随机生成 true/false
		}
		for {
			res := resourse.GoodsData.DB.Model(&model.CategoryBrand{}).Create(&categoryBrand)
			if res.Error != nil {
				if res.Error == gorm.ErrDuplicatedKey {
					categoryBrand.BrandID = uint32(rand.Intn(500) + 1)
					categoryBrand.CategoryID = uint32(rand.Intn(max-min+1) + min)
					continue
				} else {
					panic(res.Error)
				}
			}
			break
		}
	}

}

// 生成随机图片 URL
func generateRandomImageURL(i int) string {

	return fmt.Sprintf("https://example.com/images/%d.jpg", i)
}

// 生成随机 Goods 和 GoodsDetail 数据
func GoodsAndGoodsDetail() {
	for i := range 1000 {
		// 生成 Goods 数据
		goods := model.Goods{
			CategoryID:  uint32(rand.Intn(151) + 350),   // CategoryID 范围: 350-500
			BrandID:     uint32(rand.Intn(500) + 1),     // BrandID 范围: 1-500
			Status:      int32(rand.Intn(3)),            // 状态: 0, 1, 2
			ShipFree:    rand.Intn(2) == 1,              // 随机生成是否免运费
			IsHot:       rand.Intn(2) == 1,              // 随机生成是否热门
			IsNew:       rand.Intn(2) == 1,              // 随机生成是否新品
			Name:        fmt.Sprintf("商品-%d", i),        // 商品名称
			ClickNum:    int32(rand.Intn(10000)),        // 点击量: 0-9999
			SoldNum:     int32(rand.Intn(10000)),        // 售卖量: 0-9999
			FavorNum:    int32(rand.Intn(10000)),        // 收藏量: 0-9999
			MarketPrice: float32(rand.Intn(1000) + 100), // 市场价: 100-1099
			SalePrice:   float32(rand.Intn(1000) + 50),  // 售价: 50-1049
			FirstImage:  generateRandomImageURL(i),      // 封面图片 URL
		}

		// 生成 GoodsDetail 数据
		goodsDetail := model.GoodsDetail{
			GoodsBrief: fmt.Sprintf("商品简要描述-%d", i),                                                          // 商品简要描述
			Images:     model.GormList{generateRandomImageURL(i + 10000), generateRandomImageURL(i + 20000)}, // 外观预览图
			DescImages: model.GormList{generateRandomImageURL(i + 30000), generateRandomImageURL(i + 40000)}, // 详细信息图
		}
		res := resourse.GoodsData.DB.Model(&model.Goods{}).Create(&goods)
		if res.Error != nil {
			panic(res.Error)

		}
		goodsDetail.GoodsId = goods.ID
		res = resourse.GoodsData.DB.Model(&model.GoodsDetail{}).Create(&goodsDetail)
		if res.Error != nil {
			panic(res.Error)

		}
	}
}
