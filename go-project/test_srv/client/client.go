package main

import (
	"fmt"
	"kenshop/goken/registry/discover"
	"kenshop/goken/registry/ways/consul"
	"kenshop/goken/server/rpcserver/cinterceptors"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DefaultDial() (*grpc.ClientConn, error) {
	// connParams := grpc.ConnectParams{
	// 	//里面可以配置退避重连时间
	// 	Backoff: backoff.DefaultConfig,
	// 	//最小超时时间
	// 	MinConnectTimeout: 4 * time.Second,
	// }

	cfg := api.DefaultConfig()

	cfg.Address = "192.168.199.128:8500"
	cli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	d := discover.MustNewBuilder(consul.MustNewConsulDiscover(cli, consul.WithTTL(8*time.Second)))
	//理论来讲NewClient的地址应该是consul的地址,我搞错逻辑了
	return grpc.NewClient(
		fmt.Sprintf("discovery://%s:%d/%s",
			"192.168.199.128",
			22222,
			"ken",
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(d),
		grpc.WithUnaryInterceptor(cinterceptors.UnaryTracingInterceptor),
		// grpc.WithConnectParams(connParams),
		// grpc.WithKeepaliveParams(keepalive.ClientParameters{
		// 	Time:                3 * time.Second,
		// 	Timeout:             3 * time.Second,
		// 	PermitWithoutStream: true,
		// }),
	)
}

func test() {

}

func main() {
	test()

}
