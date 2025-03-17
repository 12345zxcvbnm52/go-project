package orderform

type CreateCartItemForm struct {
	GoodsId  uint32 `json:"goods_id" form:"goods_id" uri:"" header:"" binding:"required"`
	GoodsNum int32  `json:"goods_num" form:"goods_num" uri:"" header:"" binding:"required"`
	UserId   uint32 `json:"user_id" form:"user_id" uri:"" header:"" binding:"required"`
}

type CreateOrderForm struct {
	Address      string `json:"address" form:"address" uri:"" header:"" binding:"required"`
	Message      string `json:"message" form:"message" uri:"" header:"" binding:"required"`
	PayWay       string `json:"pay_way" form:"pay_way" uri:"" header:"" binding:"required"`
	SignerMobile string `json:"signer_mobile" form:"signer_mobile" uri:"" header:"" binding:"required"`
	SignerName   string `json:"signer_name" form:"signer_name" uri:"" header:"" binding:"required"`
	UserId       uint32 `json:"user_id" form:"user_id" uri:"" header:"" binding:"required"`
}

type DelCartItemForm struct {
	GoodsId uint32 `json:"goods_id" form:"goods_id" uri:"" header:"" binding:""`
	Id      uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	UserId  uint32 `json:"user_id" form:"user_id" uri:"" header:"" binding:""`
}

type OrderFliterForm struct {
	PageSize int32  `json:"page_size" form:"page_size" uri:"" header:"" binding:""`
	PagesNum int32  `json:"pages_num" form:"pages_num" uri:"" header:"" binding:""`
	UserId   uint32 `json:"user_id" form:"user_id" uri:"user_id" header:"" binding:""`
}

type OrderInfoForm struct {
	Id     uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	UserId uint32 `json:"user_id" form:"user_id" uri:"" header:"" binding:""`
}

type OrderStatusForm struct {
	Id        uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	OrderSign string `json:"order_sign" form:"order_sign" uri:"" header:"" binding:""`
	Status    int32  `json:"status" form:"status" uri:"" header:"" binding:"required"`
	UserId    uint32 `json:"user_id" form:"user_id" uri:"" header:"" binding:"required"`
}

type UpdateCartItemForm struct {
	GoodsId  uint32 `json:"goods_id" form:"goods_id" uri:"" header:"" binding:""`
	GoodsNum int32  `json:"goods_num" form:"goods_num" uri:"" header:"" binding:""`
	Id       uint32 `json:"id" form:"id" uri:"id" header:"" binding:""`
	Selected bool   `json:"selected" form:"selected" uri:"" header:"" binding:""`
	UserId   uint32 `json:"user_id" form:"user_id" uri:"" header:"" binding:""`
}

type UserInfoForm struct {
	PageSize int32  `json:"page_size" form:"page_size" uri:"" header:"" binding:""`
	PagesNum int32  `json:"pages_num" form:"pages_num" uri:"" header:"" binding:""`
	UserId   uint32 `json:"user_id" form:"user_id" uri:"user_id" header:"" binding:""`
}
