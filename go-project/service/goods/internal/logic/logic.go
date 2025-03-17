package goodslogic

import goodsdata "kenshop/service/goods/internal/data"

// Goods服务中的Service层,编写具体的服务逻辑
type GoodsService struct {
	GoodsData goodsdata.GoodsDataService
}
