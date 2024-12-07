package form

// 这个对应proto里的UserPasswordReq
type PasswordLogin struct {
	Password string `json:"password" form:"password" binding:"required"`
	Id       uint32 `json:"id" form:"id"`
	UserName string `json:"username" form:"username" binding:"required,mobile"`
}
