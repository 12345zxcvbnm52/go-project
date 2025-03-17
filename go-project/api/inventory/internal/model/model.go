package inventorymodel

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

const (
	StatusWaitingReback int32 = 1
	StatusAlreadyReback int32 = 2
)

type Inventory struct {
	Model
	//库存中的商品id
	GoodsId uint32 `gorm:"index;not null;unique"`
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
	OrderSign string `gorm:"column:order_sign;type:varchar(50)"`
	//当前库存状态,用int类型能实现幂等性
	Status      int32           `gorm:"type:int"`
	GoodsRecord GoodsRecordList `gorm:"type:varchar(4000)"`
}
