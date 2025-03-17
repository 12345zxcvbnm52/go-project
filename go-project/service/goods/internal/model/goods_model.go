package goodsmodel

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

// json标签是考虑到返回数据时的格式,一般不返回的数据则json标签为'-'忽视
// 商品目录标签
type Category struct {
	Model
	Name string `gorm:"type:varchar(20);not null" json:"name"`
	//第几级商品类型
	Level int32 `gorm:"not null" json:"level"`
	//是否可以在窗口上显示
	OnTable bool `gorm:"default:false;not null" json:"on_tab"`
	//自引用的从表外键,注意,这里的类型必须是指针,否则在gorm中无法创建
	ParentCategoryID *uint32 `gorm:"column:parent_category_id" json:"-"`
	//父层级商品类型,自引用的主表结构体字段
	ParentCategory *Category `json:"parent_category" gorm:"foreignKey:ParentCategoryID"`
	//装所有子商品分类,一对多关系并且实现表的自引用,主表的从表切片
	SubCategory []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
}

const (
	TopLevel = iota + 1
	SecondLevel
	EndLevel
)

// 商品品牌
type Brand struct {
	Model
	Name string `gorm:"type:varchar(50);not null;unique" json:"name"`
	Logo string `gorm:"type:varchar(200);default:'';not null" json:"logo"`
}

// 多对多建立的连接表
// 一个品牌旗下有多个商品类型,一个商品类型也能来自多个品牌
type CategoryBrand struct {
	Model
	CategoryID uint32   `gorm:"type:int;index:idx_category_brand,unique" json:"-"`
	Category   Category `json:"category"`
	BrandID    uint32   `gorm:"type:int;index:idx_category_brand,unique"`
	Brand      Brand    `json:"brand"`
}

// 滑动窗口
type Banner struct {
	Model
	//放在窗口上的预览图片
	Image string `gorm:"type:varchar(200);not null" json:"image"`
	//点击窗口跳转到对应商品售处
	Url string `gorm:"type:varchar(200);not null" json:"url"`
	//该预览窗口的下标位置
	Index int32 `gorm:"type:int;default:1;not null" json:"index"`
}

type GormList []string

// 自定义gorm类型需要实现的两个接口
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

const (
	StatusUnderReview = iota
	StatusOutShelves
	StatusOnShelves
)

// 商品
type Goods struct {
	//model里有主键ID(uint),updated_at,created_at,deleted_at
	Model
	//这里自动生成绑定目录的外键,即商品必须要有目录
	CategoryID uint32 `gorm:"not null"`
	Category   Category
	//这里自动生成绑定品牌的外键,即商品必须要有品牌
	BrandID uint32 `gorm:"not null"`
	Brand   Brand
	//商品的条码,这个由服务自动生成,测试时不要开unique
	GoodsSign string `gorm:"type:varchar(100)"`
	//状态,如审核中,下架,上架中几种状态
	Status int32 `gorm:"default:0"`
	//是否免运费
	ShipFree bool `gorm:"default:false"`
	//是否是热门产品
	IsHot bool `gorm:"default:false"`
	//标识是否是新发布的商品,可以直接在代码中
	IsNew bool `gorm:"default:true"`
	//商品名称,这里冗余字段,避免联表查询减低效率
	Name string `gorm:"type:varchar(100);not null"`
	//商品点击量
	ClickNum int32 `gorm:"default:0"`
	//售卖量
	SoldNum int32 `gorm:"default:0"`
	//收藏量
	FavorNum int32 `gorm:"default:0"`
	//原始价格
	MarketPrice float32 `gorm:"not null"`
	//实际价格,因打折等原因变动
	SalePrice float32 `gorm:"not null"`
	//封面
	FirstImage  string `gorm:"type:varchar(200);not null"`
	GoodsDetail GoodsDetail
}

type GoodsDetail struct {
	Model
	GoodsId uint32 `gorm:"not null"`
	//商品简要评价
	GoodsBrief string `gorm:"type:varchar(200)"`
	//商品的页内外观预览图,用的是字符串切片,即一张图片用字符串存储
	Images GormList `gorm:"type:varchar(4000);not null"`
	//商品的页内品详细信息图,如尺码表,功能示意图
	DescImages GormList `gorm:"type:varchar(4000);not null"`
}
