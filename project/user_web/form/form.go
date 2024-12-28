package form

// 这个对应proto里的UserPasswordReq
type UserLoginForm struct {
	Password string `json:"password" form:"password" binding:"required"`
	Id       uint32 `json:"id,omitempty" form:"id"`
	UserName string `json:"username" form:"username" binding:"required,mobile"`
}

type UserWriteForm struct {
	UserName string `json:"userName,omitempty" binding:"username"`
	Password string `json:"password" binding:"required,password"`
	Mobile   string `json:"mobile" binding:"required,mobile"`
	Id       uint32 `json:"id,omitempty"`
	Gender   string `json:"gender,omitempty"`
	Birth    int64  `json:"birth,omitempty"`
	Role     int32  `json:"Role,omitempty"`
}

type UserDeleteForm struct {
	ID uint32 `json:"id" binding:"required"`
}
