package app

import (
	"goken/registry"
	"os"

	"google.golang.org/grpc"
)

type Option func(o *option)

type option struct {
	id        string
	name      string
	signals   []os.Signal
	registor  registry.Registor
	rpcServer *grpc.Server
}

func WithID(id string) Option {
	return func(o *option) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *option) {
		o.name = name
	}
}

// 自定义监听哪种信号
func WithSignals(signals []os.Signal) Option {
	return func(o *option) {
		o.signals = signals
	}
}

func WithRegistor(registor registry.Registor) Option {
	return func(o *option) {
		o.registor = registor
	}
}
