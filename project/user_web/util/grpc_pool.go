package util

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func init() {
	DefaultRpcConnOpt = RpcConnOptions{
		DialTimeout:        5 * time.Second,
		BackoffMaxTry:      3 * time.Second,
		KeepAliveCheck:     30 * time.Second,
		KeepAliveTimeout:   15 * time.Second,
		initWindowSize:     1 << 28,
		initConnWindowSize: 1 << 28,
		maxSendSize:        1 << 28,
		maxRecvSize:        1 << 28,
	}
}

// 没想到一个好的实现能够充分利用grpc的多路复用,就暂时不写了
// 没看明白grpc的心跳检测机制,可以考虑自己实现一个心跳检查

type Pooler interface {
	Value() (*conn, error)
	Close()
}

type Conner interface {
	Value() *grpc.ClientConn
	Close() error
}

type conn struct {
	//这个连接被引用的次数,达到最大值后不可再被使用
	//即利用grpc连接的多路复用性
	ref    int32
	client *grpc.ClientConn
	//对应Pool中的下标,从0开始
	index int32
	//是否可用(被删除),看看怎么通过心跳检查设置
	okFlag bool
	//指向拥有这个连接的池
	pool  *pool
	mutex *sync.RWMutex
}

func (c *conn) Value() *grpc.ClientConn {
	return c.client
}

func (c *conn) Close() {
	c.mutex.RLock()
	c.ref--
	c.mutex.RUnlock()
}

func (c *conn) close() {
	c.mutex.Lock()
	c.okFlag = false
	c.mutex.Unlock()
}

// 不允许直接修改value和opt,后续可以开api修改部分参数并开api同步到所有连接中
type pool struct {
	invoke func(*RpcConnOptions) (*grpc.ClientConn, error)
	opt    *RpcConnOptions
	size   int32
	okFlag bool
	conns  []*conn
	//用于负载均衡,采用简单的轮询即可
	balance int32
	//可以考虑用bitmap记录那几个conns是非okflag的
	//bitmap int32
}

func (p *pool) Value() (*conn, error) {
	if !p.okFlag {
		return nil, errors.New("连接池未初始化或已关闭")
	}
	for {
		i := (p.balance) % p.size
		//注意这里可以考虑复原balance如果balance超过了取模的值
		p.balance++
		if !p.conns[i].okFlag || p.conns[i] == nil {
			continue
		}
		p.conns[i].mutex.RLock()
		if p.conns[i].ref < p.opt.RpcPoolopt.ConnMaxRef {
			p.conns[i].ref++
			p.conns[i].mutex.RUnlock()
			return p.conns[i], nil
		}
		//这里可以考虑记录遇到一个连接被用完的情况,进而
		p.conns[i].mutex.RUnlock()
	}
}

func (p *pool) Close() {
	p.okFlag = false
	for i := 0; i < int(p.size); i++ {
		p.conns[i].close()
	}
}

func (p *pool) initPool(opt *RpcPoolOptions) {
	r := &RpcPoolOptions{}
	if opt.ConnMaxRef != 0 {
		r.ConnMaxRef = opt.ConnMaxRef
	} else {
		r.ConnMaxRef = DefaultRpcPoolOpt.ConnMaxRef
	}

	if opt.MaxIdle != 0 {
		r.MaxIdle = opt.MaxIdle
	} else {
		r.MaxIdle = DefaultRpcPoolOpt.MaxIdle
	}

	if opt.MinIdle != 0 {
		r.MinIdle = opt.MinIdle
	} else {
		r.MinIdle = DefaultRpcPoolOpt.MinIdle
	}

	if opt.PoolSize != 0 {
		r.PoolSize = opt.PoolSize
	} else {
		r.PoolSize = DefaultRpcPoolOpt.PoolSize
	}

	if opt.Reuse.Valid {
		r.Reuse = opt.Reuse
	} else {
		r.Reuse = DefaultRpcPoolOpt.Reuse
	}
	p.size = r.PoolSize
	p.opt.RpcPoolopt = r
	fmt.Println(p.opt.RpcPoolopt)
	p.conns = make([]*conn, r.PoolSize)
	go func() {
		for i := range p.size {
			c, err := p.invoke(p.opt)
			if err != nil {
				zap.S().Errorw("一个conn无法建立", "msg", err.Error(), "index", i)
			}
			p.conns[i] = new(conn)
			p.conns[i].mutex = new(sync.RWMutex)
			p.conns[i].okFlag = true
			p.conns[i].client = c
			p.conns[i].pool = p
			p.conns[i].index = i
		}
		p.okFlag = true
	}()
}

func NewDefaultGrpcPool() Pooler {
	pool := &pool{}
	pool.opt = new(RpcConnOptions)
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
	pool.initPool(&DefaultRpcPoolOpt)
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
	pool := &pool{opt: opt, invoke: DefaultDial}
	pool.initPool(opt.RpcPoolopt)
	return pool
}

type RpcPoolOptions struct {
	//池内维持的最小的空闲连接,保证缩容不会缩到没有
	MinIdle int32
	//池内允许的最大空闲连接,小于等于0则不开启,作为缩容的凭证之一
	MaxIdle int32
	//池的大小
	PoolSize   int32
	ConnMaxRef int32
	//是否开启pool的可重用功能
	Reuse sql.NullBool
}

type RpcConnOptions struct {
	RpcPoolopt *RpcPoolOptions
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
	//连接最大可重用性
}

var DefaultRpcConnOpt RpcConnOptions
var DefaultRpcPoolOpt RpcPoolOptions

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
