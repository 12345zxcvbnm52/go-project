package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ErrInternalWrong    = status.Error(codes.Internal, "服务器内部未知错误")
	ErrCategoryNotFound = status.Error(codes.NotFound, "所选分类不存在")
	ErrDuplicatedCategy = status.Error(codes.AlreadyExists, "欲创建的分类已存在")
	ErrBannerNotFound   = status.Error(codes.NotFound, "所选滑窗未找到")
	ErrDuplicatedBanner = status.Error(codes.AlreadyExists, "欲创建的滑窗已存在")
	ErrGoodsNotFound    = status.Error(codes.NotFound, "不存在对应的商品")
	ErrDuplicatedGoods  = status.Error(codes.AlreadyExists, "欲创建的商品已存在")
	ErrBrandNotFound    = status.Error(codes.NotFound, "不存在对应的品牌")
	ErrDuplicatedBrand  = status.Error(codes.AlreadyExists, "欲创建的品牌已存在")
)

// gorm给出的分页函数的最佳实践
func Paginate(pageNum int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNum <= 0 {
			pageNum = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize < 0:
			pageSize = 10
		}
		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type Model struct {
	ID        int32          `gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 商品类型/目录
type Category struct {
	Model
	Name string `gorm:"type:varchar(20);not null" json:"name"`
	//第几级商品类型
	Level int32 `gorm:"type:int;not null;default 1" json:"level"`
	//是否可以在窗口上显示
	OnTab bool `gorm:"default:false;not null" json:"on_tab"`
	//自引用的从表外键
	ParentCategoryID int32 `json:"-"`
	//父层级商品类型,自引用的主表结构体字段
	ParentCategory *Category `json:"parent_category" gorm:"foreignKey:ParentCategoryID"`
	//装所有子商品分类,
	//一对多关系并且实现表的自引用,主表的从表切片
	SubCategory []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
}

// 商品品牌
type Brand struct {
	Model
	Name string `gorm:"type:varchar(50);not null" json:"name"`
	Logo string `gorm:"type:varchar(200);default:'';not null" json:"logo"`
}

// 多对多建立的连接表
// 一个品牌旗下有多个商品类型,一个商品类型也能来自多个品牌
type CategoryWithBrand struct {
	Model
	CategoryID int32    `gorm:"type:int;index:idx_category_brand,unique" json:"-"`
	Category   Category `json:"category"`
	BrandID    int32    `gorm:"type:int;index:idx_category_brand,unique"`
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

// 商品
type Goods struct {
	Model
	//这里自动生成绑定目录的外键,即商品必须要有目录
	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	//这里自动生成绑定品牌的外键,即商品必须要有品牌
	BrandID int32 `gorm:"type:int;not null"`
	Brand   Brand
	//是否上架了
	OnSale bool
	//运费是否免费
	TransFree bool
	//是否是热门产品
	IsHot bool
	//商品名称
	Name string `gorm:"type:varchar(100);not null"`
	//商品的编号
	GoodSign string `gorm:"type:varchar(50);not null"`
	//商品点击量
	ClickNum int32
	//售卖量
	SoldNum int32
	//收藏量
	FavorNum int32
	//原始价格
	MarketPrice float32 `gorm:"not null"`
	//实际价格,因打折等原因变动
	SalePrice float32 `gorm:"not null"`
	//商品简要评价
	GoodsBrief string `gorm:"type:varchar(100);not null"`
	//商品的预览图,用的是字符串切片,即一张图片用字符串存储
	Images GormList `gorm:"type:varchar(4000);not null"`
	//商品的边缘图
	DescImages GormList `gorm:"type:varchar(4000);not null"`
	//封面
	FirstImage string `gorm:"type:varchar(200);not null"`
	IsNew      bool   `gorm:"not null"`
}

type EsGoods struct {
	ID          int32   `json:"id"`
	CategoryID  int32   `json:"category_id"`
	BrandID     int32   `json:"brand_id"`
	OnSale      bool    `json:"on_sale"`
	TransFree   bool    `json:"trans_free"`
	IsHot       bool    `json:"is_hot"`
	Name        string  `json:"name"`
	GoodsSign   string  `json:"goods_sign"`
	ClickNum    int32   `json:"click_num"`
	SoldNum     int32   `json:"sold_num"`
	FavorNum    int32   `json:"favor_num"`
	MarketPrice float32 `json:"market_price"`
	SalePrice   float32 `json:"sale_price"`
	GoodsBrief  string  `json:"goods_brief"`
	IsNew       bool    `json:"is_new"`
}

func (u EsGoods) IndexName() string {
	return "goods"
}

func (u EsGoods) IndexMapping() string {
	return `
	{
		"mappings" : {
			"properties" : {
				"brand_id" : {
					"type" : "integer"
				},
				"goods_sign" : {
					"type" : "integer"
				},
				"category_id" : {
					"type" : "integer"
				},
				"click_num" : {
					"type" : "integer"
				},
				"favor_num" : {
					"type" : "integer"
				},
				"id" : {
					"type" : "integer"
				},
				"is_hot" : {
					"type" : "boolean"
				},
				"is_new" : {
					"type" : "boolean"
				},
				"market_price" : {
					"type" : "float"
				},
				"name" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"goods_brief" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"on_sale" : {
					"type" : "boolean"
				},
				"trans_free" : {
					"type" : "boolean"
				},
				"sale_price" : {
					"type" : "float"
				},
				"sold_num" : {
					"type" : "long"
				}
			}
		}
	}`
}
