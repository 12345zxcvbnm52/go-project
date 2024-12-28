package global

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ConsulClient *api.Client

func DefaultDial(ServerName string) (*grpc.ClientConn, error) {
	RecvSize := 1 << 24
	SendSize := 1 << 24
	return grpc.NewClient(
		fmt.Sprintf("consul://%s:%d/%s?healthy=true",
			ServerConfig.ConsulConfig.ConsulIp,
			ServerConfig.ConsulConfig.ConsulPort,
			ServerName,
		),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(RecvSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(SendSize)),
	)
}
