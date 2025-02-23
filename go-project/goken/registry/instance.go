package registry

import (
	"context"
	"net/url"
)

//任意注册中心想嵌入到代码中只需实现以下接口

// 服务发现接口
type Discover interface {
	//创建服务监听器,而服务监听器用于从特定的注册中心持续监听Name为serviceName的服务的信息
	NewListener(ctx context.Context, serviceName string) (Listener, error)
	//返回解析器关联的注册中心的类型:consul,etcd等
	RegisrtyName() string
}

type Registor interface {
	Register(context.Context, *ServiceInstance) error
	Deregister(context.Context, string) error
}

type Listener interface {
	//第一次监听或者服务实例发生变化时调用会返回服务实例列表
	//其它情况下会阻塞直至context超时或服务实例发生变化
	ListenAndGet(context.Context) ([]*ServiceInstance, error)
	StopListen(context.Context) error
}

type ServiceInstance struct {
	//注册到注册中心的服务id
	ID string `json:"id"`

	//服务名称
	Name string `json:"name"`

	//服务版本
	Version string `json:"version"`

	//服务元数据
	Metadata map[string]string `json:"metadata"`

	//http://127.0.0.1:8000
	//grpc://127.0.0.1:9000
	//一般来说该切片只用当成一个string即可,若想在一台机器上既运行http服务也运行grpc服务即可作切片使用
	Endpoints []*url.URL `json:"endpoints"`
}
