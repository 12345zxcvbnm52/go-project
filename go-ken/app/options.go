package app

import (
	"context"
	"goken/registry"
	"syscall"

	"net/url"
	"os"

	"goken/pkg/log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Option func(o *options)

type options struct {
	//id和name用于服务注册等
	id   string
	name string
	//用于收集终止的信号
	signals []os.Signal
	//用于服务注册
	registor registry.Registor
	//服务的主体部分
	rpcServer *grpc.Server
	//用于记录本服务的所有节点,默认只有一个节点且为本地节点
	endpoints []string
	//记录所需要的元数据
	metadata map[string]string
	//日志实例,注意这个日志实例不是具体服务的日志,这个只记录整个微服务的运行信息
	logger log.Log

	//钩子函数
	beforeRun  []func(context.Context, *App) error
	beforeStop []func(context.Context, *App) error
	afterRun   []func(context.Context, *App) error
	afterStop  []func(context.Context, *App) error
}

func newOptions(opts ...Option) *options {
	o := &options{}
	o.signals = []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	id, _ := uuid.NewUUID()
	//name需要注意,防止出现被卡掉
	o.id = id.String()
	logger := log.NewCompatLogger()
	o.logger = logger
	return o
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

// 自定义监听哪种信号
func WithSignals(signals []os.Signal) Option {
	return func(o *options) {
		o.signals = signals
	}
}

func WithRegistor(registor registry.Registor) Option {
	return func(o *options) {
		o.registor = registor
	}
}

func WithMetadata(md map[string]string) Option {
	return func(o *options) {
		o.metadata = md
	}
}

func WithEndpoints(endpoints ...*url.URL) Option {
	return func(o *options) {
		for _, v := range endpoints {
			o.endpoints = append(o.endpoints, v.String())
		}
	}
}

func WithLogger(logger log.Log) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// func RegistorTimeout(t time.Duration) Option {
// 	return func(o *options) { o.registrarTimeout = t }
// }

// // StopTimeout with app stop timeout.
// func StopTimeout(t time.Duration) Option {
// 	return func(o *options) { o.stopTimeout = t }
// }

func WithBeforeRun(fn ...func(context.Context, *App) error) Option {
	return func(o *options) {
		o.beforeRun = append(o.beforeRun, fn...)
	}
}

// BeforeStop run funcs before app stops
func WithBeforeStop(fn ...func(context.Context, *App) error) Option {
	return func(o *options) {
		o.beforeStop = append(o.beforeStop, fn...)
	}
}

// AfterStart run funcs after app starts
func WithAfterStart(fn ...func(context.Context, *App) error) Option {
	return func(o *options) {
		o.afterRun = append(o.afterRun, fn...)
	}
}

// AfterStop run funcs after app stops
func WithAfterStop(fn ...func(context.Context, *App) error) Option {
	return func(o *options) {
		o.afterStop = append(o.afterStop, fn...)
	}
}

// func WithContext(ctx context.Context) Option {
// 	return func(o *options) {
// 		o.ctx = ctx
// 	}
// }
