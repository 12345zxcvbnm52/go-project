package initialize

import (
	"fmt"
	gb "user_web/global"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

func FindUserServer() {
	filter := fmt.Sprintf("Service==`%s`", gb.ServerConfig.ConsulConfig.Name)
	serverData, err := gb.ConsulClient.Agent().Services()
	gb.ConsulClient.Agent().ServicesWithFilter(filter)
	if err != nil {
		zap.S().Errorw("web获得user服务信息失败", "msg", err.Error())
		panic(err)
	}
	for _, v := range serverData {
		gb.ServerConfig.UserServerConfig.UserServerIp = v.Address
		gb.ServerConfig.UserServerConfig.UserServerPort = v.Port
		gb.ServerConfig.UserServerConfig.UserServerId = v.ID
		gb.ServerConfig.UserServerConfig.UserServerTags = v.Tags
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
	FindUserServer()
}
