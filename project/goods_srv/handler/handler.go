package handler

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrBannerNotExists = errors.New("选择的轮播图不存在")
)

// gorm给出的分页函数的最佳实践
func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize < 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
