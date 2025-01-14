package main

import (
	"flag"
	"fmt"
	"net"
	gb "order_srv/global"
	"order_srv/handler"
	initialize "order_srv/initialize"
	otgrpc "order_srv/otgrpc"
	pb "order_srv/proto"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
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
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", gb.ServerConfig.Jaeger.Host, gb.ServerConfig.Jaeger.Port),
		},
		ServiceName: "kensame",
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	lis, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", *Ip, *Port))
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	OrderServer := &handler.OrderServer{}
	healthServer := healthpb.NewServer()
	pb.RegisterOrderServer(server, OrderServer)
	health.RegisterHealthServer(server, healthServer)

	go func() {
		server.Serve(lis)
	}()
	initialize.InitConsul()
	zap.S().Infoln("ServerConfig is : ", gb.ServerConfig)

	//启动阶段出错直接就panic
	p, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{fmt.Sprintf("%s:%d", gb.ServerConfig.RockMq.Host, gb.ServerConfig.RockMq.Port)}),
		consumer.WithGroupName("comsumer-"+gb.ServerConfig.Name),
	)
	if err != nil {
		panic(err)
	}
	if err = p.Subscribe(gb.ServerConfig.RockMq.TimeoutTopic, consumer.MessageSelector{}, handler.OrderTimeout); err != nil {
		panic(err)
	}
	if err = p.Start(); err != nil {
		panic(err)
	}

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)
	<-sign
	p.Shutdown()
	gb.ConsulClient.Agent().ServiceDeregister(gb.ServerConfig.Name)
}
