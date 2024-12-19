package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	MysqlConf string
	Cache     cache.CacheConf

	Jwt struct {
		AccessSecret string
		AccessExpire int64
	}
}
