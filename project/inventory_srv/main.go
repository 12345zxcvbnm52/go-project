package main

import (
	"flag"
	"fmt"
	gb "inventory_srv/global"
	"inventory_srv/handler"
	_ "inventory_srv/initialize"
	initialize "inventory_srv/initialize"
	pb "inventory_srv/proto"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	//开一个口子能够自定义ip和port
	Ip := flag.String("ip", gb.ServerConfig.Ip, "bind ip")
	Port := flag.Int("port", gb.ServerConfig.Port, "bind port")
	flag.Parse()
	if *Ip != gb.ServerConfig.Ip {
		gb.ServerConfig.Ip = *Ip
	}
	if *Port != gb.ServerConfig.Port {
		gb.ServerConfig.Port = *Port
	}
	//在consul中注册的两个服务如果name相同而id不同,且配置了负载均衡,则会被算为同一类服务的不同实例
	//若id相同则服务会互相覆盖
	//负载均衡应当在客户端配置
	go func() {
		lis, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", *Ip, *Port))
		server := grpc.NewServer()
		InventoryServer := &handler.InventoryServer{}
		healthServer := healthpb.NewServer()
		pb.RegisterInventoryServer(server, InventoryServer)
		health.RegisterHealthServer(server, healthServer)
		server.Serve(lis)
	}()
	initialize.InitConsul()
	zap.S().Infoln("ServerConfig is : ", gb.ServerConfig)
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)
	<-sign
	gb.ConsulClient.Agent().ServiceDeregister(gb.ServerConfig.Name)

}