package httpserver

import (

	//"github.com/penglongli/gin-metrics/ginmetrics"

	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"kenshop/goken/registry"
	mws "kenshop/goken/server/httpserver/middlewares"
	"kenshop/goken/server/httpserver/middlewares/jwt"
	otelkgin "kenshop/goken/server/httpserver/middlewares/otel"
	"kenshop/goken/server/httpserver/validate"
	"kenshop/goken/server/rpcserver"

	//"kenshop/goken/server/httpserver/pprof"

	"kenshop/pkg/common/hostgen"
	errors "kenshop/pkg/errors"
	"kenshop/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
)

var ErrNilHttpRegistor = errors.New("该http服务不存在注册器")

type Server struct {
	Ctx      context.Context
	Engine   *gin.Engine
	Host     string
	Mode     string
	InSecure bool
	UseAbort bool

	Jwt       *jwt.GinJWTMiddleware
	Registor  registry.Registor
	Validator *validate.Validator
	Tracer    *otelkgin.GinTracer
	GrpcCli   *rpcserver.Client
	//是否开启pprof接口, 默认开启, 如果开启会自动添加/debug/pprof接口
	EnableProfiling bool
	//是否开启metrics接口,默认开启,如果开启会自动添加/metrics接口
	EnableMetrics bool

	//全局的中间件,在这里面不要添加非全局用的中间件
	Middlewares map[string]gin.HandlerFunc

	//翻译器
	Locale string

	Instance *registry.ServiceInstance

	Server *http.Server
	closed bool
}

func MustNewServer(ctx context.Context, host string, opts ...ServerOption) *Server {
	s := &Server{
		Ctx:             ctx,
		Host:            host,
		Mode:            "debug",
		EnableProfiling: false,
		EnableMetrics:   false,
		Engine:          gin.New(),
		Locale:          "zh",
		Middlewares:     make(map[string]gin.HandlerFunc),
		InSecure:        true,
		Server:          &http.Server{},
		UseAbort:        true,
	}

	s.Instance = &registry.ServiceInstance{
		Name: host,
		ID:   host,
	}
	for _, o := range opts {
		o(s)
	}
	if len(s.Middlewares) == 0 {
		mws.CopyDefaultMiddlewares(s.Middlewares)
	}
	gin.SetMode(s.Mode)

	if ok := hostgen.ValidListenHost(s.Host); !ok {
		panic(errors.New("无效的监听地址"))
	}
	s.Server.Addr = s.Host
	s.Server.Handler = s.Engine
	//无论如何都开启/health路径便于后续服务注册,健康检查
	s.Engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{})
	})

	//只有注册器存在才构造instance
	if s.Registor != nil {
		u, err := url.Parse(host)
		if err != nil {
			if s.InSecure {
				host = fmt.Sprintf("http://%s", host)
			} else {
				host = fmt.Sprintf("https://%s", host)
			}
			u, err = url.Parse(host)
			if err != nil {
				panic(err)
			}
		}
		s.Instance.Endpoints = append(s.Instance.Endpoints, u)
	}

	for _, m := range s.Middlewares {
		s.Engine.Use(m)
	}

	if s.Jwt == nil {
		log.Warn("[httpserver] Server的Jwt为nil,可能将导致错误")
	}
	if s.Tracer == nil {
		log.Warn("[httpserver] Server的Tracer为nil,可能将导致错误")
	}
	if s.GrpcCli == nil {
		log.Warn("[httpserver] Server的GrpcCli为nil,可能将导致错误")
	}

	var err error
	s.Validator, err = validate.NewValidator(s.Locale)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *Server) Register(ctx context.Context, ins *registry.ServiceInstance) error {
	if s.Registor == nil {
		return ErrNilHttpRegistor
	}
	return s.Registor.Register(ctx, ins)
}

// Deregister会注销Server内Instance存储的服务Id
func (s *Server) Deregister(ctx context.Context) error {
	if s.Registor == nil {
		return ErrNilHttpRegistor
	}
	if s.closed == false {
		s.closed = true
		return s.Registor.Deregister(ctx, s.Instance.ID)
	}
	return nil
}

func (s *Server) Serve() error {
	// 设置开发模式,打印路由信息格式
	gin.SetMode(s.Mode)
	g := &run.Group{}

	// 运行前打印配置信息
	log.Infof("[httpserver] 服务启动中,监听信息为: host = %s,服务信息为: msg = %+v", s.Host, s.Instance)

	// 执行 Validator 校验
	if err := s.Validator.Excute(); err != nil {
		return err
	}

	// 如果注册器为空就不进行注册,而不是返回错误
	if err := s.Register(s.Ctx, s.Instance); err != nil && err != ErrNilHttpRegistor {
		return err
	}

	// 监听终止信号,优雅退出
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)

	// 确保Deregister只执行一次
	var deregisterOnce sync.Once
	deregisterFunc := func() {
		deregisterOnce.Do(func() {
			if e := s.Deregister(s.Ctx); e != nil && e != ErrNilHttpRegistor {
				log.Errorf("[httpserver] 服务注销失败, err= %v", e)
			} else {
				log.Info("[httpserver] 服务正常注销")
			}
		})
	}

	// 启动 HTTP 服务器
	g.Add(
		func() error {
			if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				return err
			}
			return nil
		},
		func(err error) {
			if e1 := s.Server.Shutdown(s.Ctx); e1 != nil {
				log.Errorf("[httpserver] http根服务注销失败, err= %v", e1)
			}
			deregisterFunc()
		},
	)

	// 监听终止信号,优雅退出
	g.Add(
		func() error {
			<-sign
			if e1 := s.Server.Shutdown(s.Ctx); e1 != nil {
				log.Errorf("[httpserver] 服务注销失败, err= %v", e1)
				return e1
			}
			deregisterFunc()
			return nil
		},
		func(err error) {
			// 发送终止信号,而不是 close(sign)
			select {
			case sign <- syscall.SIGINT:
			default:
			}
		},
	)

	return g.Run()
}
