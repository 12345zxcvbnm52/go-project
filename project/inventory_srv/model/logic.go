package model

import (
	gb "inventory_srv/global"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ErrNotFound        = status.Error(codes.NotFound, "未找到对应的库存信息")
	ErrInternalWrong   = status.Error(codes.Internal, "服务器内部未知的错误,请稍后尝试")
	ErrDuplicated      = status.Error(codes.AlreadyExists, "欲创建的库存已存在")
	ErrInvalidArgument = status.Error(codes.InvalidArgument, "修改库存失败,错误的参数")
	ErrBadRequest      = status.Error(codes.Aborted, "因未知原因操作失败")
	ErrLackInventory   = status.Error(codes.ResourceExhausted, "缺少库存")
)

type Result struct {
	Data  []*Inventory
	Total int64
}

// type FindOption struct {
// 	KeyWords string
// 	Age      int32
// 	Gender   string
// 	PagesNum int32
// 	PageSize int32
// }

func (u *Inventory) InsertOne() error {
	if res := gb.DB.Create(u); res.Error != nil {
		if res.Error == gorm.ErrDuplicatedKey {
			return ErrDuplicated
		}
		return ErrInternalWrong
	}
	return nil
}

func (u *Inventory) UpdateOneByGoodsId() error {
	res := gb.DB.Model(u).Update("goods_num", u.GoodsNum)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return ErrInternalWrong
	}
	return nil
}

func (u *Inventory) FindOneByGoodsId() error {
	if res := gb.DB.Model(u).Where("goods_id = ?", u.GoodsId).First(u); res.Error != nil {
		if res.RowsAffected == 0 {
			return ErrNotFound
		}
		return ErrInternalWrong
	}
	return nil
}

func (u *Inventory) FindByGoodsIds(Ids ...uint32) (*Result, error) {
	invt := []*Inventory{}
	if len(Ids) == 0 {
		return nil, ErrNotFound
	}
	res := gb.DB.Where("goods_id in (?)", Ids).Find(&invt)
	if res.Error != nil {
		return nil, ErrInternalWrong
	}
	return &Result{
		Data:  invt,
		Total: res.RowsAffected,
	}, nil
}
