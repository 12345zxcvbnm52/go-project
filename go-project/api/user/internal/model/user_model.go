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

type User struct {
	Model
	Mobile string `gorm:"type:varchar(12);unique;not null;index:idx_mobile"`
	//在sql中存储的password应当是加密的
	Password string     `gorm:"type:varchar(100);not null;"`
	UserName string     `gorm:"type:varchar(20)"`
	Birth    *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"type:varchar(4);default:'boy'"`
	//为0时权限为普通用户,用户等级随Role增大而增大
	Role int32 `gorm:"default:0"`
}

func (u *User) TableName() string {
	return "user_srv"
}

type UserList struct {
	Total int64   `json:"total,omitempty"` //总数
	Data  []*User `json:"data"`            //数据
}
