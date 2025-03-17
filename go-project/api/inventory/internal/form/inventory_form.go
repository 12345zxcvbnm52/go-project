package inventoryform

type CreateInventoryForm struct {
	GoodsId  uint32 `json:"goods_id" form:"goods_id" uri:"" header:"" binding:"required"`
	GoodsNum int32  `json:"goods_num" form:"goods_num" uri:"" header:"" binding:"required"`
}

type DecrData struct {
	GoodsId  uint32 `json:"goods_id" form:"goods_id" uri:"" header:"" binding:"required"`
	GoodsNum int32  `json:"goods_num" form:"goods_num" uri:"" header:"" binding:"required"`
}

type DecrStockForm struct {
	DecrData  []*DecrData `json:"decr_data" form:"decr_data" uri:"" header:"" binding:"required"`
	OrderSign string      `json:"order_sign" form:"order_sign" uri:"" header:"" binding:"required"`
}

type IncrData struct {
	GoodsId  uint32 `json:"goods_id" form:"goods_id" uri:"" header:"" binding:"required"`
	GoodsNum int32  `json:"goods_num" form:"goods_num" uri:"" header:"" binding:"required"`
}

type IncrStockForm struct {
	IncrData  []*IncrData `json:"incr_data" form:"incr_data" uri:"" header:"" binding:"required"`
	OrderSign string      `json:"order_sign" form:"order_sign" uri:"" header:"" binding:"required"`
}

type InventoryInfoForm struct {
	GoodsId uint32 `json:"goods_id" form:"goods_id" uri:"goods_id" header:"" binding:""`
}

type SetInventoryForm struct {
	GoodsId  uint32 `json:"goods_id" form:"goods_id" uri:"" header:"" binding:"required"`
	GoodsNum int32  `json:"goods_num" form:"goods_num" uri:"" header:"" binding:"required"`
}
