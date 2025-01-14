package registry

import "context"

type Registor interface {
	Register(context.Context, *ServiceInstance) error
	Deregister(context.Context, *ServiceInstance) error
}

type Descover interface {
	GetService(context.Context, string) ([]*ServiceInstance, error)
	Listen(context.Context, string) ([]*ServiceInstance, error)
}

type Listener interface {
	//第一次监听或者服务实例发生变化时调用会返回服务实例列表
	//其它情况下会阻塞直至context超时或服务实例发生变化
	ListenAndGet() ([]*ServiceInstance, error)
	StopListen() error
}

type ServiceInstance struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	//集群节点的所有地址
	// grpc://host:port 或者 https://host:port
	EndPoints []string          `json:"endpoints"`
	Metadata  map[string]string `json:"metadata"`
}
