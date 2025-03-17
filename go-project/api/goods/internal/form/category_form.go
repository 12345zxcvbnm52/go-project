package goodsform

type DelCategoryForm struct {
	Id uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
}

type SubCategoryForm struct {
	Id    uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	Level int32  `json:"level" form:"level" uri:"" header:"" binding:""`
}

type CreateCategoryForm struct {
	Level            int32  `json:"level" form:"level" uri:"" header:"" binding:"required"`
	Name             string `json:"name" form:"name" uri:"" header:"" binding:"required"`
	ParentCategoryId uint32 `json:"parent_category_id" form:"parent_category_id" uri:"" header:"" binding:""`
}

type UpdateCategoryForm struct {
	Id               uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	Level            int32  `json:"level" form:"level" uri:"" header:"" binding:""`
	Name             string `json:"name" form:"name" uri:"" header:"" binding:""`
	ParentCategoryId uint32 `json:"parent_category_id" form:"parent_category_id" uri:"" header:"" binding:""`
}
