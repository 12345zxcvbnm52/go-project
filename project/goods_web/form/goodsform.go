package form

type GoodsWriteForm struct {
	Name        string   `form:"name" json:"name" binding:"required,min=1,max=100"`
	GoodsSign   string   `form:"goods_sign" json:"goods_sign" binding:"required,min=2,lt=20"`
	Stocks      int32    `form:"stocks" json:"stocks"`
	CategoryID  int32    `form:"category_id" json:"category_id" binding:"required"`
	MarketPrice float32  `form:"market_price" json:"market_price" binding:"required,min=0"`
	SalePrice   float32  `form:"sale_price" json:"sale_price" binding:"required,min=0"`
	GoodsBrief  string   `form:"goods_brief" json:"goods_brief" binding:"required,min=3"`
	Images      []string `form:"images" json:"images" binding:"required,min=1"`
	DescImages  []string `form:"desc_images" json:"desc_images" binding:"required,min=1"`
	TransFree   bool     `form:"trans_free" json:"trans_free" binding:"required"`
	BrandID     int32    `form:"brand_id" json:"brand_id" binding:"required"`
	FirstImage  string   `form:"first_image" json:"first_image" binding:"required,url"`
}

type GoodsStatusForm struct {
	IsNew  *bool `form:"is_new" json:"is_new" binding:"required"`
	IsHot  *bool `form:"is_hot" json:"is_hot" binding:"required"`
	OnSale *bool `form:"on_sale" json:"on_sale" binding:"required"`
}
