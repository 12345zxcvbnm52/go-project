package goodsform

type CreateBrandForm struct {
	Logo string `json:"logo" form:"logo" uri:"" header:"" binding:"required"`
	Name string `json:"name" form:"name" uri:"" header:"" binding:"required"`
}

type BrandFilterForm struct {
	PageSize int32 `json:"page_size" form:"page_size" uri:"" header:"" binding:""`
	PagesNum int32 `json:"pages_num" form:"pages_num" uri:"" header:"" binding:""`
}

type DelBrandForm struct {
	Id   uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	Name string `json:"name" form:"name" uri:"" header:"" binding:""`
}

type UpdateBrandForm struct {
	Id   uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	Logo string `json:"logo" form:"logo" uri:"" header:"" binding:""`
	Name string `json:"name" form:"name" uri:"" header:"" binding:""`
}
