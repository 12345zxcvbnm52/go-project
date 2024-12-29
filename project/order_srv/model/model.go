package model

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ErrBadRequest      = status.Error(codes.Aborted, "因未知原因操作失败")
	ErrInternalWrong   = status.Error(codes.Internal, "服务器内部未知的错误,请稍后尝试")
	ErrInvalidArgument = status.Error(codes.InvalidArgument, "修改库存失败,错误的参数")

	ErrOrderNotFound   = status.Error(codes.NotFound, "未找到对应的订单信息")
	ErrOrderDuplicated = status.Error(codes.AlreadyExists, "欲创建的订单已存在")

	ErrCartNotFound = status.Error(codes.NotFound, "购物车无对应的商品")
	//这个Err似乎只能在内部使用?
	ErrCartNoItems          = status.Error(codes.NotFound, "购物车内是空的")
	ErrCartDuplicated       = status.Error(codes.AlreadyExists, "欲创建的已存在")
	ErrCartNoSelected       = status.Error(codes.InvalidArgument, " 购物车内没有选中的商品")
	ErrOrderGoodsNotFound   = status.Error(codes.NotFound, "未找到对应的订单内商品信息")
	ErrOrderGoodsDuplicated = status.Error(codes.AlreadyExists, "欲创建的订单内商品信息已存在")
	ErrOrderFailedCreate    = status.Error(codes.Aborted, "订单因未知原因创建失败")
	ErrBadGoodsClient       = status.Error(codes.Internal, "商品服务连接失败")
	ErrBadInventoryClient   = status.Error(codes.Internal, "库存服务连接失败")

	ErrBadRockMq = status.Error(codes.Internal, "消息队列操作失败")
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

const (
	//待支付
	StatusUnPay int16 = 1
	//已支付
	StatusPaid int16 = 2
	//待发货
	StatusUnDispached int16 = 3
	//已发货
	StatusDisPached int16 = 4
	//待收货
	StatusUnRecv int16 = 5
	StatusRecv   int16 = 6
	//订单已取消
	StatusCancelled int16 = 7
	//订单退款中
	StatusRefunding int16 = 8
	//订单已完成
	StatusFinished int16 = 9
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
	Cost float32 `gormL:"not null"`
	//支付时间
	PayTime time.Time
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

// gorm给出的分页函数的最佳实践
func Paginate(pagesNum int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagesNum <= 0 {
			pagesNum = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize < 0:
			pageSize = 10
		}
		offset := (pagesNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
