package util

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func init() {
	DefaultRpcConnOpt = RpcConnOptions{
		DialTimeout:        5 * time.Second,
		BackoffMaxTry:      3 * time.Second,
		KeepAliveCheck:     10 * time.Second,
		KeepAliveTimeout:   3 * time.Second,
		initWindowSize:     1 << 28,
		initConnWindowSize: 1 << 28,
		maxSendSize:        1 << 28,
		maxRecvSize:        1 << 28,
	}
}

// 没想到一个好的实现能够充分利用grpc的多路复用,就暂时不写了
type Pooler interface {
	Value() (*grpc.ClientConn, error)
}

// 不允许直接修改value和opt,后续可以开api修改部分参数并开api同步到所有连接中
type Pool struct {
	invoke func(*RpcConnOptions) (*grpc.ClientConn, error)
	opt    *RpcConnOptions
}

func (p *Pool) Value() (*grpc.ClientConn, error) {
	return p.invoke(p.opt)
}

func NewDefaultGrpcPool() Pooler {
	pool := &Pool{}
	pool.opt = new(RpcConnOptions)
	fmt.Println(DefaultRpcConnOpt)
	pool.opt.maxRecvSize = DefaultRpcConnOpt.maxRecvSize
	pool.opt.maxSendSize = DefaultRpcConnOpt.maxSendSize
	pool.opt.initConnWindowSize = DefaultRpcConnOpt.initConnWindowSize
	pool.opt.initWindowSize = DefaultRpcConnOpt.initWindowSize
	pool.opt.BackoffMaxTry = DefaultRpcConnOpt.BackoffMaxTry
	pool.opt.DialTimeout = DefaultRpcConnOpt.DialTimeout
	pool.opt.KeepAliveTimeout = DefaultRpcConnOpt.KeepAliveCheck
	pool.opt.KeepAliveCheck = DefaultRpcConnOpt.KeepAliveCheck
	pool.opt.ConsulAddr = DefaultRpcConnOpt.ConsulAddr
	pool.opt.ConsulPort = DefaultRpcConnOpt.ConsulPort
	pool.opt.ServerName = DefaultRpcConnOpt.ServerName
	pool.invoke = DefaultDial
	return pool
}

func NewGrpcPool(opt *RpcConnOptions) Pooler {
	if opt.ConsulAddr == "" {
		opt.ConsulAddr = DefaultRpcConnOpt.ConsulAddr
	}
	if opt.ConsulPort == 0 {
		opt.ConsulPort = DefaultRpcConnOpt.ConsulPort
	}
	if opt.ServerName == "" {
		opt.ServerName = DefaultRpcConnOpt.ServerName
	}
	if opt.BackoffMaxTry == 0 {
		opt.BackoffMaxTry = DefaultRpcConnOpt.BackoffMaxTry
	}
	if opt.DialTimeout == 0 {
		opt.DialTimeout = DefaultRpcConnOpt.DialTimeout
	}
	opt.maxRecvSize = DefaultRpcConnOpt.maxRecvSize
	opt.maxSendSize = DefaultRpcConnOpt.maxSendSize
	opt.initConnWindowSize = DefaultRpcConnOpt.initConnWindowSize
	opt.initWindowSize = DefaultRpcConnOpt.initWindowSize
	if opt.KeepAliveCheck == 0 {
		opt.KeepAliveCheck = DefaultRpcConnOpt.KeepAliveCheck
	}
	if opt.KeepAliveTimeout == 0 {
		opt.KeepAliveTimeout = DefaultRpcConnOpt.KeepAliveCheck
	}
	pool := &Pool{opt: opt, invoke: DefaultDial}
	return pool
}

type RpcConnOptions struct {
	//预留的口子,允许外部修改这三个参数改变对应的连接对象
	ConsulAddr string
	ConsulPort int
	ServerName string
	//用于指定客户端尝试建立与服务器的连接时的最大等待时间,
	DialTimeout time.Duration
	//最大尝试退避策略来尝试重连,
	BackoffMaxTry time.Duration
	//多少时间内连接未被访问过则尝试ping一次服务端以测试可用性
	KeepAliveCheck time.Duration
	//每一次ping的最大等待时间,超过则断开连接,
	KeepAliveTimeout time.Duration
	//以下这两个参数详见http/2的流窗口控制
	initWindowSize     int32
	initConnWindowSize int32
	//规定rpc每次最大的接发收帧的大小,默认给的比较低
	maxSendSize int32
	maxRecvSize int32
}

var DefaultRpcConnOpt RpcConnOptions

func DefaultDial(opt *RpcConnOptions) (*grpc.ClientConn, error) {
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
			opt.ConsulAddr,
			opt.ConsulPort,
			opt.ServerName,
		),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(connParams),
		grpc.WithInitialConnWindowSize(opt.initConnWindowSize),
		grpc.WithInitialWindowSize(opt.initWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(opt.maxRecvSize))),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(int(opt.maxSendSize))),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                opt.KeepAliveCheck,
			Timeout:             opt.KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	)
}
