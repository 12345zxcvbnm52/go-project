package initialize

import (
	"fmt"
	gb "order_web/global"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

func FindOrderServer() {
	filter := fmt.Sprintf("Service==`%s`", gb.ServerConfig.ConsulConfig.Name)
	serverData, err := gb.ConsulClient.Agent().Services()
	gb.ConsulClient.Agent().ServicesWithFilter(filter)
	if err != nil {
		zap.S().Errorw("web获得Order服务信息失败", "msg", err.Error())
		panic(err)
	}
	for _, v := range serverData {
		gb.ServerConfig.OrderServerConfig.OrderServerIp = v.Address
		gb.ServerConfig.OrderServerConfig.OrderServerPort = v.Port
		gb.ServerConfig.OrderServerConfig.OrderServerId = v.ID
		gb.ServerConfig.OrderServerConfig.OrderServerTags = v.Tags
	}
}

func InitConsul() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d",
		gb.ServerConfig.ConsulConfig.ConsulIp,
		gb.ServerConfig.ConsulConfig.ConsulPort,
	)
	var err error
	gb.ConsulClient, err = api.NewClient(cfg)
	if err != nil {
		zap.S().Errorw("连接到consul服务端失败", "msg", err.Error())
		panic(err)
	}
	FindOrderServer()
}
