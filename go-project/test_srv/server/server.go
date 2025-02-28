package main

import (
	"context"
	"kenshop/goken/registry/registor"
	"kenshop/goken/registry/ways/consul"
	"kenshop/goken/server/rpcserver"
	"kenshop/goken/server/rpcserver/sinterceptors"
	"kenshop/pkg/trace"
	pb "kenshop/test_srv/proto"
	"net"

	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/otel"
)

type ConnServer struct {
	pb.UnimplementedTestServer
}

func (c *ConnServer) Conn(ctx context.Context, in *pb.Req) (*pb.Res, error) {

	return &pb.Res{Res: in.Name + " hello ken"}, nil
}

func main() {
	lis, _ := net.Listen("tcp", "192.168.199.128:22223")
	c := trace.MustNewTracer(context.Background(), trace.WithName("server-test"))
	tp, err := c.NewTraceProvider("192.168.199.128:4318")
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tp)
	cfg := api.DefaultConfig()

	cfg.Address = "192.168.199.128:8500"
	cli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	r := registor.MustNewRegister(
		consul.MustNewConsulRegistor(cli,
			consul.WithEnableHealthCheck(true),
			consul.WithDeregisterCriticalServiceAfter("30s"),
			consul.WithHealthcheckInterval("10s"),
		),
	)

	s := rpcserver.MustNewServer(context.Background(),
		rpcserver.WithHost("127.0.0.1:22223"),
		rpcserver.WithServiceName("ken"),
		rpcserver.WithListener(lis),
		rpcserver.WithRegistor(r),
		rpcserver.WithUnaryInts(sinterceptors.UnaryTracingInterceptor),
	)

	cs := &ConnServer{}
	pb.RegisterTestServer(s.Server, cs)

	if err := s.Serve(); err != nil {
		panic(err)
	}
}