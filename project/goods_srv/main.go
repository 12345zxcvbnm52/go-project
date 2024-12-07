package main

import (
	"flag"
	"fmt"
	gb "goods_srv/global"
	"goods_srv/handler"
	_ "goods_srv/initialize"
	pb "goods_srv/proto"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	Ip := flag.String("ip", gb.ServerConfig.Ip, "bind ip")
	Port := flag.Int("port", gb.ServerConfig.Port, "bind port")
	flag.Parse()
	//在consul中注册的两个服务如果name相同而id不同,且配置了负载均衡,则会被算为同一类服务的不同实例
	//若id相同则服务会互相覆盖
	//负载均衡应当在客户端配置
	go func() {
		lis, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", *Ip, *Port))
		server := grpc.NewServer()
		userServer := &handler.GoodsServer{}
		healthServer := healthpb.NewServer()
		pb.RegisterGoodsServer(server, userServer)
		health.RegisterHealthServer(server, healthServer)
		server.Serve(lis)
	}()
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)
	<-sign
	gb.ConsulClient.Agent().ServiceDeregister(gb.ServerConfig.Name)

}
