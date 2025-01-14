package app

import (
	"context"
	"goken/registry"
	"os"
	"os/signal"
	"sync"
)

type App struct {
	//这两个参数用于控制程序运行
	ctx      context.Context
	cancel   context.CancelFunc
	opts     *options
	mtx      sync.Mutex
	instance *registry.ServiceInstance
}

func NewApp(opts ...Option) *App {
	o := &App{}
	o.opts = newOptions(opts...)
	for _, v := range opts {
		v(o.opts)
	}
	o.ctx, o.cancel = context.WithCancel(context.Background())
	return o
}

func (a *App) Run() error {
	var err error
	a.mtx.Lock()
	a.instance, err = a.ServiceBuild()
	a.mtx.Unlock()
	if err != nil {
		panic(err)
	}

	if a.opts.rpcServer != nil {
		a.opts.rpcServer.Serve()
	}

	if a.opts.registor != nil {
		if err = a.opts.registor.Register(context.Background(), a.instance); err != nil {
			panic(err)
		}
	} else {

	}
	closeCh := make(chan os.Signal, 1)
	signal.Notify(closeCh, a.opts.signals...)
	<-closeCh
	return nil
}

func (a *App) Stop() error {
	a.mtx.Lock()
	if a.opts.registor != nil && a.instance != nil {
		if err := a.opts.registor.Deregister(context.Background(), a.instance); err != nil {
			panic(err)
		}
	}
	a.mtx.Unlock()
	return nil
}

func (a *App) ServiceBuild() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 1)
	copy(endpoints, a.opts.endpoints)
	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		EndPoints: endpoints,
	}, nil
}
