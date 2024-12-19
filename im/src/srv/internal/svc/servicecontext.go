package svc

import (
	"srv/internal/config"
	"srv/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.MysqlConf)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUsersModel(sqlConn, c.Cache),
	}
}
