package userform

type CreateUserForm struct {
	Birth    int64  `json:"birth" form:"birth" uri:"" header:"" binding:"required"`
	Gender   string `json:"gender" form:"gender" uri:"" header:"" binding:""`
	Mobile   string `json:"mobile" form:"mobile" uri:"" header:"" binding:"required"`
	Password string `json:"password" form:"password" uri:"" header:"" binding:"required"`
	UserName string `json:"user_name" form:"user_name" uri:"" header:"" binding:"required"`
}

type DelUserForm struct {
	Id   uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	Name string `json:"name" form:"name" uri:"" header:"" binding:""`
}

type UpdateUserForm_0 struct {
	Birth    int64  `json:"birth" form:"birth" uri:"" header:"" binding:"required"`
	Gender   string `json:"gender" form:"gender" uri:"" header:"" binding:"required"`
	Id       uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	Mobile   string `json:"mobile" form:"mobile" uri:"" header:"" binding:"required"`
	Password string `json:"password" form:"password" uri:"" header:"" binding:"required"`
	Role     int32  `json:"role" form:"role" uri:"" header:"" binding:"required"`
	UserName string `json:"user_name" form:"user_name" uri:"" header:"" binding:"required"`
}

type UpdateUserForm_1 struct {
	Birth    int64  `json:"birth" form:"birth" uri:"" header:"" binding:""`
	Gender   string `json:"gender" form:"gender" uri:"" header:"" binding:""`
	Id       uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	Mobile   string `json:"mobile" form:"mobile" uri:"" header:"" binding:""`
	Password string `json:"password" form:"password" uri:"" header:"" binding:""`
	Role     int32  `json:"role" form:"role" uri:"" header:"" binding:""`
	UserName string `json:"user_name" form:"user_name" uri:"" header:"" binding:""`
}

type UserFliterForm struct {
	PageSize int32 `json:"page_size" form:"page_size" uri:"" header:"" binding:""`
	PagesNum int32 `json:"pages_num" form:"pages_num" uri:"" header:"" binding:""`
}

type UserIdForm struct {
	Id uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
}

type UserMobileForm struct {
	Mobile string `json:"mobile" form:"mobile" uri:"mobile" header:"" binding:"mobile"`
}

type UserPasswordForm struct {
	Password string `json:"password" form:"password" uri:"" header:"" binding:"required"`
	UserName string `json:"user_name" form:"user_name" uri:"" header:"" binding:"required"`
}
