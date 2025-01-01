package model

import (
	"context"
	"fmt"
	gb "goods_srv/global"
	"strconv"

	es "github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FindOption struct {
	MinPrice int32
	MaxPrice int32
	IsHot    bool
	IsNew    bool
	OnTable  bool
	// 确定点击的目录是哪一级,会递归显示下一层(或者说下几层)目录
	//如果位最底层则忽视
	CategyId uint32
	PagesNum int32
	PageSize int32
	KeyWords string
	BrandId  uint32
}

type Result struct {
	Data  []*Goods
	Total int64
}

// 使用es找到商品的id再通过redis-mysql查询详细数据
// 只把需要作为查询条件的字段存储在es中
func (u *Goods) FindByOpt(opt *FindOption) (*Result, error) {
	esQuery := es.NewBoolQuery()
	if opt.KeyWords != "" {
		esQuery = esQuery.Must(es.NewMultiMatchQuery(opt.KeyWords, "name", "goods_brief"))
	}
	if opt.MaxPrice > 0 {
		esQuery = esQuery.Filter(es.NewRangeQuery("sale_price").Lte(opt.MaxPrice))
	}
	if opt.MinPrice > 0 {
		esQuery = esQuery.Filter(es.NewRangeQuery("sale_price").Gte(opt.MinPrice))

	}
	if opt.BrandId > 0 {
		esQuery = esQuery.Filter(es.NewTermQuery("brand_id", opt.BrandId))

	}
	if opt.IsHot {
		esQuery = esQuery.Filter(es.NewTermQuery("is_hot", opt.IsHot))
	}
	if opt.IsNew {
		esQuery = esQuery.Filter(es.NewTermQuery("is_new", opt.IsNew))
	}
	if opt.OnTable {
		esQuery = esQuery.Filter(es.NewTermQuery("on_tab", opt.OnTable))
	}

	//这里想要通过Category限制找到商品

	if opt.CategyId > 0 {
		categy := &Category{}
		categy.ID = opt.CategyId
		if err := categy.FindOneById(); err == nil {
			var Query string
			switch categy.Level {
			case gb.TopLevel:
				subQuery := fmt.Sprintf("select id from categories where parent_category_id = %d", categy.ID)
				Query = fmt.Sprintf("select id from categories where parent_category_id in (%s)", subQuery)
			case gb.SecondLevel:
				Query = fmt.Sprintf("select id from categories where parent_category_id = %d", categy.ID)
			case gb.EndLevel:
				Query = fmt.Sprintf("%d", categy.ID)
			default:
			}

			type Result struct {
				ID int32 `json:"id"`
			}
			results := []Result{}
			if res := gb.DB.Raw(Query).Scan(&results); res.Error == nil {
				fmt.Println(results)
				ids := []interface{}{}
				for _, v := range results {
					ids = append(ids, v.ID)
				}
				esQuery = esQuery.Filter(es.NewTermsQuery("category_id", ids...))
			}

		}
	}

	if opt.PagesNum <= 0 {
		opt.PagesNum = 1
	}
	switch {
	case opt.PageSize > 100:
		opt.PageSize = 100
	case opt.PageSize <= 0:
		opt.PageSize = 10
	}
	logic := EsGoods{}
	result, err := gb.EsConn.Search().Index(logic.IndexName()).Query(esQuery).
		From(int(opt.PagesNum)).Size(int(opt.PageSize)).Do(context.Background())
	if err != nil {
		zap.S().Errorw("es按条件查询失败", "msg", err.Error())
		return nil, ErrInternalWrong
	}
	ids := []uint32{}
	for _, v := range result.Hits.Hits {
		i, _ := strconv.Atoi(v.Id)
		ids = append(ids, uint32(i))
	}
	res, err := u.FindByIds(ids...)
	if err != nil {
		return nil, err
	}
	return &Result{
		Data:  res.Data,
		Total: result.TotalHits(),
	}, nil
}

func (u *Goods) FindByIds(Id ...uint32) (*Result, error) {
	data := []*Goods{}
	if len(Id) == 0 {
		return nil, ErrGoodsNotFound
	}
	//这里有个大问题,传入切片前一定要检查是否为空,否则会全表检查
	res := gb.DB.Model(&Goods{}).Preload("Category").Preload("Brand").Find(&data, Id)
	if res.RowsAffected == 0 {
		return nil, ErrGoodsNotFound
	}
	if res.Error != nil {
		zap.S().Errorw("商品按id批量查询失败", "msg", res.Error.Error())
		return nil, ErrInternalWrong
	}
	return &Result{
		Data:  data,
		Total: res.RowsAffected,
	}, nil
}

func (u *Goods) InsertOne() error {
	u.Category.ID = u.CategoryID
	if err := u.Category.FindOneById(); err != nil {
		return err
	}
	u.Brand.ID = u.BrandID
	if err := u.Brand.FindOneById(); err != nil {
		return err
	}
	tx := gb.DB.Begin()
	if res := tx.Model(&Goods{}).Create(u); res.Error != nil {
		var err error
		if res.Error == gorm.ErrDuplicatedKey {
			err = ErrDuplicatedGoods
		} else {
			err = ErrInternalWrong
		}
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *Goods) FindOneById() error {
	if res := gb.DB.Model(&Goods{}).Preload("Brand").Preload("Category").First(u); res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return ErrGoodsNotFound
		}
		return ErrInternalWrong
	}
	return nil
}

func (u *Goods) UpdateOneById() error {
	return nil
}

func (u *Goods) AfterCreate(tx *gorm.DB) error {
	d := EsGoods{
		ID:          u.ID,
		BrandID:     u.BrandID,
		CategoryID:  u.CategoryID,
		IsHot:       u.IsHot,
		IsNew:       u.IsNew,
		GoodsSign:   u.GoodSign,
		GoodsBrief:  u.GoodsBrief,
		TransFree:   u.TransFree,
		OnSale:      u.OnSale,
		Name:        u.Name,
		ClickNum:    u.ClickNum,
		SoldNum:     u.SoldNum,
		FavorNum:    u.FavorNum,
		SalePrice:   u.SalePrice,
		MarketPrice: u.MarketPrice,
	}
	_, err := gb.EsConn.Index().Index(d.IndexName()).BodyJson(d).Id(strconv.Itoa(int(d.ID))).Do(context.Background())
	return err
}

// 钩子函数解决不了数据的同步
func (u *Goods) AfterUpdate(tx *gorm.DB) error {
	d := EsGoods{
		ID:          u.ID,
		BrandID:     u.BrandID,
		CategoryID:  u.CategoryID,
		IsHot:       u.IsHot,
		IsNew:       u.IsNew,
		GoodsSign:   u.GoodSign,
		GoodsBrief:  u.GoodsBrief,
		TransFree:   u.TransFree,
		OnSale:      u.OnSale,
		Name:        u.Name,
		ClickNum:    u.ClickNum,
		SoldNum:     u.SoldNum,
		FavorNum:    u.FavorNum,
		SalePrice:   u.SalePrice,
		MarketPrice: u.MarketPrice,
	}
	_, err := gb.EsConn.Update().Index(d.IndexName()).Doc(d).Id(strconv.Itoa(int(d.ID))).Do(context.Background())
	return err
}
