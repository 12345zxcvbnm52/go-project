package initialize

import (
	gb "order_web/global"
	"order_web/util"
)

func InitRpcPool() {
	util.DefaultRpcConnOpt.ConsulAddr = gb.ServerConfig.ConsulConfig.ConsulIp
	util.DefaultRpcConnOpt.ConsulPort = gb.ServerConfig.ConsulConfig.ConsulPort
	util.DefaultRpcConnOpt.ServerName = gb.ServerConfig.ConsulConfig.Name
	gb.RpcPool = util.NewDefaultGrpcPool()
}
