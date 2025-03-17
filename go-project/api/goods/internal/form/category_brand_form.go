package goodsform

type CategoryBrandFilterForm struct {
	PageSize int32 `json:"page_size" form:"page_size" uri:"" header:"" binding:""`
	PagesNum int32 `json:"pages_num" form:"pages_num" uri:"" header:"" binding:""`
}

type CreateCategoryBrandForm struct {
	BrandId    uint32 `json:"brand_id" form:"brand_id" uri:"" header:"" binding:"required"`
	CategoryId uint32 `json:"category_id" form:"category_id" uri:"" header:"" binding:"required"`
}

type DelCategoryBrandForm struct {
	Id uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
}

type UpdateCategoryBrandForm struct {
	BrandId    uint32 `json:"brand_id" form:"brand_id" uri:"" header:"" binding:""`
	CategoryId uint32 `json:"category_id" form:"category_id" uri:"" header:"" binding:""`
	Id         uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
}
