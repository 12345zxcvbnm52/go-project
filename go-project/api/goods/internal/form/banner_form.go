package goodsform

type CreateBannerForm struct {
	Image string `json:"image" form:"image" uri:"" header:"" binding:"required"`
	Index int32  `json:"index" form:"index" uri:"" header:"" binding:"required"`
	Url   string `json:"url" form:"url" uri:"" header:"" binding:"required"`
}

type DelBannerForm struct {
	Id    uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	Index int32  `json:"index" form:"index" uri:"" header:"" binding:""`
}

type UpdateBannerForm struct {
	Id    uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	Image string `json:"image" form:"image" uri:"" header:"" binding:""`
	Index int32  `json:"index" form:"index" uri:"" header:"" binding:""`
	Url   string `json:"url" form:"url" uri:"" header:"" binding:""`
}
