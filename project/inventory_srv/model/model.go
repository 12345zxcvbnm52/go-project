package model

import (
	"database/sql/driver"
	"encoding/json"
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

const (
	StatusWaitingReback int32 = 1
	StatusAlreadyReback int32 = 2
)

type Inventory struct {
	Model
	//库存中的商品id
	GoodsId uint32 `gorm:"type:int;index;not null;unique"`
	//库存中的商品的数量
	GoodsNum int32 `gorm:"type:int;default:0"`
	//涉及分布式锁的乐观锁
	Version int32 `gorm:"type:int;default:0"`
}

// 用于记录一件商品的扣除回滚记录
type GoodsRecord struct {
	GoodsId  uint32 `json:"id"`
	GoodsNum int32  `json:"num"`
}

type GoodsRecordList []GoodsRecord

func (g GoodsRecordList) Value() (driver.Value, error) {
	return json.Marshal(g)
}
func (g *GoodsRecordList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

// 用于库存扣减回滚的记录,即在微服务层记录一次库存的扣减,
// 如果事务出现问题能根据Record恢复
type InvtRecord struct {
	//订单编号
	OrderSign string `gorm:"column:order_sign;type:varchar(50);index:idx_order_sign,unique"`
	//当前库存状态,用int类型能实现幂等性
	Status      int32           `gorm:"type:int"`
	GoodsRecord GoodsRecordList `gorm:"type:varchar(256)"`
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
