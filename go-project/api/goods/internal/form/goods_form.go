package goodsform

type UpdateGoodsForm struct {
	BrandId     uint32   `json:"brand_id" form:"brand_id" uri:"" header:"" binding:""`
	CategoryId  uint32   `json:"category_id" form:"category_id" uri:"" header:"" binding:""`
	DescImages  []string `json:"desc_images" form:"desc_images" uri:"" header:"" binding:""`
	FirstImage  string   `json:"first_image" form:"first_image" uri:"" header:"" binding:""`
	GoodsBrief  string   `json:"goods_brief" form:"goods_brief" uri:"" header:"" binding:""`
	Id          uint32   `json:"id" form:"id" uri:"id" header:"" binding:""`
	Images      []string `json:"images" form:"images" uri:"" header:"" binding:""`
	MarketPrice float32  `json:"market_price" form:"market_price" uri:"" header:"" binding:""`
	Name        string   `json:"name" form:"name" uri:"" header:"" binding:""`
	SalePrice   float32  `json:"sale_price" form:"sale_price" uri:"" header:"" binding:""`
	ShipFree    bool     `json:"ship_free" form:"ship_free" uri:"" header:"" binding:""`
	Status      int32    `json:"status" form:"status" uri:"" header:"" binding:""`
}

type DelGoodsForm struct {
	Id uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
}

type GoodsFilterForm struct {
	BrandId    uint32 `json:"brand_id" form:"brand_id" uri:"" header:"" binding:""`
	CategoryId uint32 `json:"category_id" form:"category_id" uri:"" header:"" binding:""`
	Id         uint32 `json:"id" form:"id" uri:"" header:"" binding:""`
	IsHot      bool   `json:"is_hot" form:"is_hot" uri:"" header:"" binding:""`
	IsNew      bool   `json:"is_new" form:"is_new" uri:"" header:"" binding:""`
	KeyWords   string `json:"key_words" form:"key_words" uri:"" header:"" binding:""`
	MaxPrice   int32  `json:"max_price" form:"max_price" uri:"" header:"" binding:""`
	MinPrice   int32  `json:"min_price" form:"min_price" uri:"" header:"" binding:""`
	PageSize   int32  `json:"page_size" form:"page_size" uri:"" header:"" binding:""`
	PagesNum   int32  `json:"pages_num" form:"pages_num" uri:"" header:"" binding:""`
	Status     int32  `json:"status" form:"status" uri:"" header:"" binding:""`
}

type GoodsIdsForm struct {
	Ids []uint32 `json:"ids" form:"ids" uri:"" header:"" binding:"required"`
}

type GoodsInfoForm struct {
	Id uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
}

type CreateGoodsForm struct {
	BrandId     uint32   `json:"brand_id" form:"brand_id" uri:"" header:"" binding:"required"`
	CategoryId  uint32   `json:"category_id" form:"category_id" uri:"" header:"" binding:"required"`
	DescImages  []string `json:"desc_images" form:"desc_images" uri:"" header:"" binding:"required"`
	FirstImage  string   `json:"first_image" form:"first_image" uri:"" header:"" binding:"required"`
	GoodsBrief  string   `json:"goods_brief" form:"goods_brief" uri:"" header:"" binding:"required"`
	Images      []string `json:"images" form:"images" uri:"" header:"" binding:"required"`
	MarketPrice float32  `json:"market_price" form:"market_price" uri:"" header:"" binding:"required"`
	Name        string   `json:"name" form:"name" uri:"" header:"" binding:"required"`
	SalePrice   float32  `json:"sale_price" form:"sale_price" uri:"" header:"" binding:"required"`
	ShipFree    bool     `json:"ship_free" form:"ship_free" uri:"" header:"" binding:"required"`
}
