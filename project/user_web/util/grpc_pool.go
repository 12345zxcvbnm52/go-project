package util

import (
	"context"
	"fmt"
	"time"
	gb "user_web/global"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// 没想到一个好的实现能够充分利用grpc的多路复用,就暂时不写了
type Pooler interface {
	Value() (*grpc.ClientConn, error)
}

type Pool struct{}

func (p *Pool) Value() (*grpc.ClientConn, error) {
	return DefaultDial(DefaultRpcConnOpt)
}

type RpcConnOptions struct {
	//用于指定客户端尝试建立与服务器的连接时的最大等待时间,
	DialTimeout time.Duration
	//最大尝试退避策略来尝试重连,
	BackoffMaxTry time.Duration
	//多少时间内连接未被访问过则尝试ping一次服务端以测试可用性
	KeepAliveCheck time.Duration
	//每一次ping的最大等待时间,超过则断开连接,
	KeepAliveTimeout time.Duration
	//以下这两个参数详见http/2的流窗口控制
	InitWindowSize     int32
	InitConnWindowSize int32
	//规定rpc每次最大的接发收帧的大小,默认给的比较低
	MaxSendSize int32
	MaxRecvSize int32
}

var DefaultRpcConnOpt RpcConnOptions = RpcConnOptions{
	DialTimeout:        5 * time.Second,
	BackoffMaxTry:      3 * time.Second,
	KeepAliveCheck:     10 * time.Second,
	KeepAliveTimeout:   3 * time.Second,
	InitWindowSize:     1 << 28,
	InitConnWindowSize: 1 << 28,
	MaxSendSize:        1 << 28,
	MaxRecvSize:        1 << 28,
}

func DefaultDial(opt RpcConnOptions) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opt.DialTimeout)
	defer cancel()

	connParams := grpc.ConnectParams{
		//里面可以配置退避重连时间
		Backoff: backoff.DefaultConfig,
		//最小超时时间
		MinConnectTimeout: opt.DialTimeout,
	}

	//"github.com/mbobakov/grpc-consul-resolver"已经帮忙封装了通过consul找到ip与port了
	return grpc.DialContext(ctx,
		fmt.Sprintf("consul://%s:%d/%s?healthy=true",
			gb.ServerConfig.ConsulConfig.ConsulIp,
			gb.ServerConfig.ConsulConfig.ConsulPort,
			gb.ServerConfig.ConsulConfig.Name,
		),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(connParams),
		grpc.WithInitialConnWindowSize(opt.InitConnWindowSize),
		grpc.WithInitialWindowSize(opt.InitWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(opt.MaxRecvSize))),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(int(opt.MaxSendSize))),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                opt.KeepAliveCheck,
			Timeout:             opt.KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	)
}
