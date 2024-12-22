package model

import (
	"fmt"
	gb "goods_srv/global"

	"go.uber.org/zap"
)

type FindOption struct {
	MinPrice int32
	MaxPrice int32
	IsHot    bool
	IsNew    bool
	OnTable  bool
	// 确定点击的目录是哪一级,会递归显示下一层(或者说下几层)目录
	//如果位最底层则忽视
	CategyId int32
	PagesNum int32
	PageSize int32
	KeyWords string
	Brand    int32
}

type Result struct {
	Data  []*Goods
	Total int64
}

func (u *Goods) FindByOpt(opt *FindOption) (*Result, error) {
	LocDB := gb.DB.Model(&Goods{})
	fmt.Println(opt)
	if opt.KeyWords != "" {
		LocDB = LocDB.Where("name LIKE ?", "%"+opt.KeyWords+"%")
	}
	if opt.MaxPrice > 0 {
		LocDB = LocDB.Where("sale_price <= ?", opt.MaxPrice)
	}
	if opt.MinPrice > 0 {
		LocDB = LocDB.Where("sale_price >= ?", opt.MinPrice)
	}
	if opt.Brand > 0 {
		LocDB = LocDB.Where("brand_id = ?", opt.Brand)
	}
	if opt.IsHot {
		LocDB = LocDB.Where("is_hot = true")
	}
	if opt.IsNew {
		LocDB = LocDB.Where("is_new = true")
	}
	if opt.OnTable {
		LocDB = LocDB.Where("on_tab = true")
	}

	//这里想要通过Category限制找到商品
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
		LocDB = LocDB.Where(fmt.Sprintf("category_id in (%s)", Query))
	}

	LocDB = LocDB.Scopes(Paginate(int(opt.PagesNum), int(opt.PageSize)))
	res := []*Goods{}
	var total int64 = 0
	LocDB.Count(&total)
	if err := LocDB.Find(&res).Error; err != nil {
		zap.S().Errorw("商品按条件查询失败", "msg", err.Error())
		return nil, err
	}
	return &Result{
		Data:  res,
		Total: total,
	}, nil
}

func (u *Goods) FindByIds(Id ...int32) (*Result, error) {
	res := []*Goods{}
	LocDB := gb.DB.Where("id in (?)", Id)
	r := LocDB.Find(&res)
	if r.Error != nil {
		zap.S().Errorw("商品按id批量查询失败", "msg", r.Error.Error())
		return nil, r.Error
	}
	return &Result{
		Data:  res,
		Total: r.RowsAffected,
	}, nil
}
