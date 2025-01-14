package app

import (
	"context"
	"goken/registry"
	"net/url"
	"os"
	"os/signal"
	"sync"

	"github.com/google/uuid"
)

type App struct {
	opts      *option
	endpoints []string
	mtx       sync.Mutex
	instance  *registry.ServiceInstance
}

func newDefaultOptions() *option {
	opts := &option{}
	opts.signals = []os.Signal{os.Kill, os.Interrupt}
	id, _ := uuid.NewUUID()
	opts.id = id.String()

	return opts
}

func New(endpoints []*url.URL, opts ...Option) *App {
	o := &App{}
	for _, v := range endpoints {
		o.endpoints = append(o.endpoints, v.String())
	}
	o.opts = newDefaultOptions()
	for _, v := range opts {
		v(o.opts)
	}
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
	copy(endpoints, a.endpoints)
	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		EndPoints: endpoints,
	}, nil
}
