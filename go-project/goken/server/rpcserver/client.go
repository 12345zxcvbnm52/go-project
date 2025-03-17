package rpcserver

import (
	"context"
	"fmt"
	"kenshop/goken/registry"
	"net/url"

	discover "kenshop/goken/registry/discover"

	"google.golang.org/grpc"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	RoundRobin      string = "round_robin"
	LeastConnection string = "least_connection"
	PickFirst       string = "pick_first"
)

type ClientOption func(o *Client)

type Client struct {
	Ctx context.Context
	//要连接的端点
	Endpoint *url.URL
	//timeout  time.Duration
	//服务发现接口
	Discover   registry.Discover
	UnaryInts  []grpc.UnaryClientInterceptor
	StreamInts []grpc.StreamClientInterceptor
	GrpcOpts   []grpc.DialOption
	//用于grpc的负载均衡
	BalanceModel  string
	EnableTracing bool
	EnableMetrics bool
	Insecure      bool
	client        *grpc.ClientConn
}

func MustNewClient(ctx context.Context, target string, opts ...ClientOption) *Client {
	c := &Client{
		BalanceModel:  RoundRobin,
		EnableTracing: true,
		EnableMetrics: true,
		Ctx:           ctx,
		Insecure:      true,
	}

	u, err := url.Parse(target)
	if err != nil {
		target := fmt.Sprintf("grpc://%s", target)
		u, err = url.Parse(target)
		if err != nil {
			panic(err)
		}
	}
	c.Endpoint = u

	for _, o := range opts {
		o(c)
	}
	c.GrpcOpts = append(c.GrpcOpts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "`+c.BalanceModel+`"}`))
	c.GrpcOpts = append(c.GrpcOpts, grpc.WithChainUnaryInterceptor(c.UnaryInts...))
	c.GrpcOpts = append(c.GrpcOpts, grpc.WithChainStreamInterceptor(c.StreamInts...))

	if c.Discover != nil {
		c.Endpoint.Scheme = "discovery"
		c.Endpoint.Host = "127.0.0.1:65535"
		c.GrpcOpts = append(c.GrpcOpts, grpc.WithResolvers(
			discover.MustNewBuilder(c.Discover),
		))
	}

	if c.Insecure {
		c.GrpcOpts = append(c.GrpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}
	return c
}

func (c *Client) CtxWithMetadata(md metadata.MD) context.Context {
	return metadata.NewOutgoingContext(c.Ctx, md)
}

func WithEnableTracing(on bool) ClientOption {
	return func(o *Client) {
		o.EnableTracing = on
	}
}

func WithSercure(on bool) ClientOption {
	return func(o *Client) {
		o.Insecure = !on
	}
}

func WithEnableMetrics(on bool) ClientOption {
	return func(o *Client) {
		o.EnableTracing = on
	}
}

// 设置服务发现
func WithDiscover(d registry.Discover) ClientOption {
	return func(o *Client) {
		o.Discover = d
	}
}

// 设置拦截器
func WithClientUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *Client) {
		o.UnaryInts = append(o.UnaryInts, in...)
	}
}

// 设置stream拦截器
func WithClientStreamInterceptor(in ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *Client) {
		o.StreamInts = append(o.StreamInts, in...)
	}
}

// 设置grpc的dial选项
func WithDialOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *Client) {
		o.GrpcOpts = opts
	}
}

// 设置负载均衡器
func WithBalanceModel(model string) ClientOption {
	return func(o *Client) {
		switch model {
		case RoundRobin:
			o.BalanceModel = model
		case PickFirst:
			o.BalanceModel = model
		default:
			o.BalanceModel = RoundRobin
		}

	}
}

func (c *Client) Reset() {
	if c.client != nil {
		c.client.Close()
		c.client = nil
	}
}

func (c *Client) Dial() (*grpc.ClientConn, error) {
	var err error
	if c.client == nil {
		if c.Discover != nil {
			c.client, err = grpc.DialContext(c.Ctx, c.Endpoint.String(), c.GrpcOpts...)
		} else {
			c.client, err = grpc.DialContext(c.Ctx, c.Endpoint.Host, c.GrpcOpts...)
		}
		if err != nil {
			return nil, err
		}
	}
	return c.client, nil
}
