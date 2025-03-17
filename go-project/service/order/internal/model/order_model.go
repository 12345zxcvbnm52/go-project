package ordermodel

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint32         `gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 购物车
// 一对多关系,一个User的Cart可以有多个GoodsId
// 这里需要注意的业务逻辑是通过购物车内查询的商品架构应该是动态的(可被售者修改)
// 而订单中的价格则不应该被修改,查询购物车内商品价格时是有可能要跨微服务的
// 故最后并没有按范式中的一对多关系建两张关联表,而是通过冗余UserId解决
type Cart struct {
	Model
	//购物车所属的用户
	UserId uint32 `gorm:"type:uint;index;not null"`
	//购物车内有的商品
	GoodsId uint32 `gorm:"type:int;index;not null"`
	//商品要购买的数量
	GoodsNums int32 `gorm:"type:int;not null"`
	//商品是否选中
	Selected bool `gorm:"default:0"`
}

// 保证(除了订单取消和订单完成级别一样)过程前后和数字大小是一一对应即可
const (
	//订单创建成功,待支付
	StatusCreated = 1 + iota
	//已支付
	StatusPaid
	//未发货
	StatusUnhipped
	//待收货
	StatusUnRecived
	//已收货
	StatusRecived
	//订单退款中
	StatusRefunding
	//订单已取消
	StatusCancelled
	//订单已完成
	StatusFinished
)

// 订单,也是用户的一对多关系,
type Order struct {
	Model
	//订单所属的用户
	UserId uint32 `gorm:"type:uint;index;not null"`
	//订单号
	OrderSign string `gorm:"varchar(50);index;not null"`
	//支付方式
	PayWay string `gorm:"varchar(20);not null"`
	//订单的状态,如待支付,已支付
	Status int16 `gorm:"type:smallint;not null"`
	//来自支付宝,微信的交易号,便于后期查账
	TradeNum string `gorm:"type:varchar(100);not null"`
	//订单的交易金额
	Cost float32 `gorm:"not null"`
	//支付时间
	PayTime *time.Time
	//目的地址
	Address string `gorm:"type:varchar(100);not null"`
	//收件人的信息
	SignerName   string `gorm:"type:varchar(20);not null"`
	SignerMobile string `gorm:"type:varchar(11);not null"`
	//留言信息
	Message string `gorm:"type:varchar(60)"`
}

// 将完整的订单表拆分,以更小的冗余换取更大的冗余
// 否则需要跨微服务查找商品,同时也满足了订单内的价格一般是不变的
// 未采用外键
type OrderGoods struct {
	Model
	OrderId    uint32  `gorm:"type:uint;index;not null"`
	GoodsId    uint32  `gorm:"type:int;index;not null"`
	GoodsName  string  `gorm:"type:varchar(100);index;not null"`
	GoodsImage string  `gorm:"type:varchar(200);not null"`
	GoodsPrice float32 `gorm:"not null"`
	GoodsNum   int32   `gorm:"type:int;not null"`
}
