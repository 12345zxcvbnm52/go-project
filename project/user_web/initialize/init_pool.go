package initialize

import (
	gb "user_web/global"
	"user_web/util"
)

func InitRpcPool() {
	util.DefaultRpcConnOpt.ConsulAddr = gb.ServerConfig.ConsulConfig.ConsulIp
	util.DefaultRpcConnOpt.ConsulPort = gb.ServerConfig.ConsulConfig.ConsulPort
	util.DefaultRpcConnOpt.ServerName = gb.ServerConfig.ConsulConfig.Name
	util.DefaultRpcPoolOpt.ConnMaxRef = 64
	util.DefaultRpcPoolOpt.MaxIdle = 1
	util.DefaultRpcPoolOpt.MinIdle = 1
	util.DefaultRpcPoolOpt.PoolSize = 5
	util.DefaultRpcPoolOpt.Reuse.Valid = true
	util.DefaultRpcPoolOpt.Reuse.Bool = true
	gb.RpcPool = util.NewDefaultGrpcPool()
}
