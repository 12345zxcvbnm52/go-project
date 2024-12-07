package util

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"sync"
// 	"sync/atomic"
// 	"time"
// 	gb "user_web/global"

// 	"go.uber.org/zap"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/backoff"
// 	"google.golang.org/grpc/credentials/insecure"
// 	"google.golang.org/grpc/keepalive"
// )

// var ErrClosed = errors.New("pool is closed")

// // 用于描述一个grpc.Conn,用于封装grpc.Conn
// type Connecter interface {
// 	//返回完整的grpc连接
// 	Value() *grpc.ClientConn
// 	Close() error
// }

// type conn struct {
// 	//对应的连接对象
// 	client *grpc.ClientConn
// 	//指向被收录的连接池
// 	pool *connPool
// 	//确定该连接是否已被删除
// 	deleted bool
// }

// func (c *conn) Value() *grpc.ClientConn {
// 	return c.client
// }

// func (c *conn) Close() error {
// 	c.pool.decrOneRef()
// 	if c.deleted {
// 		return c.reset()
// 	}
// 	return nil
// }

// func (c *conn) reset() error {
// 	client := c.client
// 	c.client = nil
// 	c.deleted = false
// 	if client != nil {
// 		return client.Close()
// 	}
// 	return nil
// }

// func (p *connPool) newConn(client *grpc.ClientConn) *conn {
// 	return &conn{
// 		client:  client,
// 		pool:    p,
// 		deleted: false,
// 	}
// }

// // 用于抽象一个池
// type ConnPooler interface {
// 	Get() (Connecter, error)
// 	Close() error
// 	//返回连接池的各项状态
// 	Status() string
// }

// // 池内的index,current,ref,closed应当逻辑上都是原子
// type connPool struct {
// 	//用于获取可用连接的最大下标
// 	index int32

// 	//当前连接池中被使用的连接数
// 	ref int32

// 	//连接池选项
// 	opt ConnPoolOptions

// 	//用于创建连接的服务器地址ip
// 	name string

// 	//所有已创建的物理连接
// 	conns []*conn

// 	//用于创建连接的服务器地址ip
// 	address string

// 	//用于创建连接的服务器地址port
// 	port int

// 	//当调用Close方法,设置为true
// 	closed bool

// 	// 控制对原子变量current的并发读写
// 	sync.RWMutex
// }

// func NewPool(addr string, port int, name string, opt ConnPoolOptions) (ConnPooler, error) {
// 	if addr == "" {
// 		return nil, errors.New("无效的服务器地址")
// 	}
// 	if opt.Dial == nil {
// 		return nil, errors.New("无效的连接创建函数")
// 	}
// 	if opt.MaxIdle <= 0 || opt.MaxSize <= 0 || opt.MaxIdle > opt.MaxSize {
// 		return nil, errors.New("错误的opt参数,请保证最大连接数,最大空闲数是合理的")
// 	}
// 	p := &connPool{
// 		index:   0,
// 		current: opt.MaxIdle,
// 		ref:     0,
// 		opt:     opt,
// 		conns:   make([]*conn, opt.MaxSize),
// 		address: addr,
// 		closed:  false,
// 	}
// 	for i := range opt.MaxIdle {
// 		c, err := p.opt.Dial(addr, port, name)
// 		if err != nil {
// 			p.Close()
// 			return nil, err
// 		}
// 		p.conns[i] = p.newConn(c)
// 		zap.S().Infof("第%d号连接池节点建立成功\n", i)
// 	}
// 	zap.S().Infof("新建连接池成功,当前状态为%s\n", p.Status())
// 	return p, nil
// }

// func (p *connPool) incrOneRef() error {
// 	p.RWMutex.Lock()
// 	defer p.RWMutex.Unlock()
// 	{
// 		newRef := atomic.AddInt32(&p.ref, 1)

// 		if newRef == p.opt.MaxSize {
// 			zap.S().Infoln("此时连接池使用达到上限")
// 		}

// 	}
// 	return
// }

// func (p *connPool) decrOneRef() {
// 	newRef := atomic.AddInt32(&p.ref, -1)
// }

// // 用于创建连接池的参数
// type ConnPoolOptions struct {
// 	//用于创建和配置连接的函数
// 	Dial func(addr string, port int, name string) (*grpc.ClientConn, error)
// 	//连接池中空闲连接的最大数量
// 	MaxIdle int32
// 	//最大连接数,当此值为零时,连接池中的连接数量没有限制
// 	MaxSize int32
// 	//如果为true且连接池已经到达最大连接数,则开始重用池内连接,
// 	//若为false,则哪怕池满也会new一个连接返回
// 	Reuse bool
// }

// var DefaultOptions = ConnPoolOptions{
// 	Dial:    Dial,
// 	MaxIdle: 8,
// 	MaxSize: 64,
// 	Reuse:   true,
// }

// func Dial(addr string, port int, name string) (*grpc.ClientConn, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(gb.ConnPoolConfig.DialTimeout)*time.Second)
// 	defer cancel()

// 	connParams := grpc.ConnectParams{
// 		//里面可以配置退避重连时间
// 		Backoff: backoff.DefaultConfig,
// 		//最小超时时间
// 		MinConnectTimeout: time.Duration(gb.ConnPoolConfig.DialTimeout) * time.Second,
// 	}

// 	//	"github.com/mbobakov/grpc-consul-resolver"已经帮忙封装了通过consul找到ip与port了
// 	return grpc.DialContext(ctx,
// 		fmt.Sprintf("consul://%s:%d/%s?healthy=true",
// 			addr,
// 			port,
// 			name,
// 		),
// 		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithConnectParams(connParams),
// 		grpc.WithKeepaliveParams(keepalive.ClientParameters{
// 			Time:    time.Duration(gb.ConnPoolConfig.KeepAliveCheck) * time.Second,
// 			Timeout: time.Duration(gb.ConnPoolConfig.KeepAliveTimeout) * time.Second,
// 		}))
// }
