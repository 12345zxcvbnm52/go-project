package discover

import (
	"context"

	"strings"
	"time"

	"kenshop/goken/registry"
	errors "kenshop/pkg/errors"

	"google.golang.org/grpc/resolver"
)

// 为客户端Dial时address的前缀(协议段)
const name = "discovery"

type BuilderOption func(o *builder)

type builder struct {
	discover registry.Discover
	timeout  time.Duration
	insecure bool
}

type DiscoverModel string

const (
	ConsulModel DiscoverModel = "consul"
	EtcdModel   DiscoverModel = "etcd"
)

// NewBuilder创建一个用于registry解析程序的构建器
func MustNewBuilder(discover registry.Discover, opts ...BuilderOption) resolver.Builder {
	builder := &builder{
		discover: discover,
		timeout:  time.Second * 10,
		insecure: false,
	}
	for _, o := range opts {
		o(builder)
	}
	return builder
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var err error
	var listener registry.Listener
	done := make(chan struct{}, 1)
	//创建最顶层的ctx,用于后续流程控制
	ctx, cancel := context.WithCancel(context.Background())
	//这里采取了观察者模式,这个协程等待服务推送,只有有服务推送这里listener才能被创建
	go func() {
		//异步创建Listener,该接口由具体的注册中心提供
		listener, err = b.discover.NewListener(ctx, strings.TrimPrefix(target.URL.Path, "/"))
		close(done)
	}()
	select {
	//等待listener被建立成功
	case <-done:
	case <-time.After(b.timeout):
		err = errors.New("discover创建超时")
	}
	if err != nil {
		cancel()
		return nil, err
	}
	resolver := &Resolver{
		listener: listener,
		cc:       cc,
		ctx:      ctx,
		cancel:   cancel,
		insecure: b.insecure,
	}

	go resolver.listen()
	return resolver, nil
}

// Scheme return scheme of discovery
func (*builder) Scheme() string {
	return name
}

func WithTimeout(timeout time.Duration) BuilderOption {
	return func(b *builder) {
		b.timeout = timeout
	}
}

func WithInsecure(insecure bool) BuilderOption {
	return func(b *builder) {
		b.insecure = insecure
	}
}
