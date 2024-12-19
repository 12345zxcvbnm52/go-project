package initialize

import (
	"fmt"
	gb "user_srv/global"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	// 这个包在内部的init函数内修改了grpc.NewClient函数的resolver解析器
)

func RegisterServer() {
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", gb.ServerConfig.Ip, gb.ServerConfig.Port),
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "1m",
	}
	server := &api.AgentServiceRegistration{
		ID:      gb.ServerConfig.Id,
		Name:    gb.ServerConfig.Name,
		Tags:    gb.ServerConfig.Tags,
		Port:    gb.ServerConfig.Port,
		Address: gb.ServerConfig.Ip,
		Check:   check,
	}

	err := gb.ConsulClient.Agent().ServiceRegister(server)
	if err != nil {
		zap.S().Errorw("consul服务注册失败", "msg", err.Error())
		panic(err)
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
	RegisterServer()
}
