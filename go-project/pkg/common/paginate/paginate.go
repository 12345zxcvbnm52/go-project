package paginate

import "gorm.io/gorm"

// gorm给出的分页函数的最佳实践
func GormPaginate(pagesNum int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagesNum < 0 {
			pagesNum = 0
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (pagesNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
