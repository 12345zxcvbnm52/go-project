package initialize

import (
	gb "goods_web/global"
	"goods_web/util"
)

func InitRpcPool() {
	util.DefaultRpcConnOpt.ConsulAddr = gb.ServerConfig.ConsulConfig.ConsulIp
	util.DefaultRpcConnOpt.ConsulPort = gb.ServerConfig.ConsulConfig.ConsulPort
	util.DefaultRpcConnOpt.ServerName = gb.ServerConfig.ConsulConfig.Name
	gb.RpcPool = util.NewDefaultGrpcPool()

}
