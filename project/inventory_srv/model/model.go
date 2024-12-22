package model

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

//一个完整的仓库应该考虑位置(距离筛选),数量等因素
//但是引进仓库系统会使业务变得复杂,故先搁置
// type Stock struct{
// 	Model
// 	StockName string
// 	Address string

// }

type Inventory struct {
	Model
	//库存中的商品id
	GoodsId int32 `gorm:"type:int;index;not null;unique"`
	//库存中的商品的数量
	GoodsNum int32 `gorm:"type:int;default:0"`
	//涉及分布式锁的乐观锁
	Version int32 `gorm:"type:int"`
}

var (
	StatusOK      int32 = 1
	StatusDecring int32 = 2
)

// 用于库存预扣减回滚的记录,即在微服务层记录一次库存的扣减,
// 如果过期就能根据Record恢复
type InvtRecord struct {
	//提出扣减的用户的id
	UserId int32
	//要扣减的商品id
	GoodsId int32
	//要扣减的数量
	GoodsNum int32
	//订单编号
	OrderId int32
	//当前库存状态,用int类型能实现幂等性
	Status int32
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
