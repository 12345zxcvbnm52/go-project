package form

type OrderCreateForm struct {
	UserId       uint32 `form:"user_id" json:"user_id" binding:"required,min=1"`
	Address      string `form:"addr" json:"addr" binding:"required,max=100"`
	SignerName   string `form:"sname" json:"sname" binding:"required,max=20"`
	SignerMobile string `form:"smobile" json:"smobile" binding:"required,mobile"`
	Message      string `form:"msg" json:"msg" binding:"max=60"`
	PayWay       string `form:"pay_way" json:"pay_way" binding:"required,oneof=wechat alipay"`
}
