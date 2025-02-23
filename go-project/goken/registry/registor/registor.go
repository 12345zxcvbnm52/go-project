package registor

import (
	"context"
	"kenshop/goken/registry"
)

// 最后权衡:goken内部的注册器名为Register
type Register struct {
	//要注册到的服务注册中心地址
	registor registry.Registor
}

type RegisterOption func(*Register)

// r 可以是goken内部提供的几种注册方式(通过不同的注册中心,也可以自己实现)
func MustNewRegister(registor registry.Registor, opts ...RegisterOption) *Register {
	r := &Register{
		registor: registor,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Register) Register(ctx context.Context, ins *registry.ServiceInstance) error {
	return r.registor.Register(ctx, ins)
}

func (r *Register) Deregister(ctx context.Context, insID string) error {
	return r.registor.Deregister(ctx, insID)
}
