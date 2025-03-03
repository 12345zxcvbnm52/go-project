package form

type UserFliterForm struct {
	PageSize int32 `json:"page_size" form:"page_size" binding:""`
	PagesNum int32 `json:"pages_num" form:"pages_num" binding:""`
}

type UserIdForm struct {
	Id uint32 `json:"id" form:"id" binding:"required"`
}

type UserMobileForm struct {
	Mobile string `json:"mobile" form:"mobile" binding:"required"`
}

type CreateUserForm struct {
	Birth    int64  `json:"birth" form:"birth" binding:"required"`
	Gender   string `json:"gender" form:"gender" binding:""`
	Mobile   string `json:"mobile" form:"mobile" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	UserName string `json:"user_name" form:"user_name" binding:""`
}

type UpdateUserForm struct {
	Birth    int64  `json:"birth" form:"birth" binding:""`
	Gender   string `json:"gender" form:"gender" binding:""`
	Id       uint32 `json:"id" form:"id" binding:"required"`
	Mobile   string `json:"mobile" form:"mobile" binding:""`
	Password string `json:"password" form:"password" binding:""`
	Role     int32  `json:"role" form:"role" binding:""`
	UserName string `json:"user_name" form:"user_name" binding:""`
}

type DelUserForm struct {
	Id   uint32 `json:"id" form:"id" binding:"required"`
	Name string `json:"name" form:"name" binding:"required"`
}

type UserPasswordForm struct {
	Id       uint32 `json:"id" form:"id" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	UserName string `json:"user_name" form:"user_name" binding:"required"`
}
