package httpserver

import (

	//"github.com/penglongli/gin-metrics/ginmetrics"

	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"kenshop/goken/registry"
	mws "kenshop/goken/server/httpserver/middlewares"
	"kenshop/goken/server/httpserver/middlewares/jwt"
	ginotel "kenshop/goken/server/httpserver/middlewares/tracing"
	"kenshop/goken/server/httpserver/validate"

	//"kenshop/goken/server/httpserver/pprof"

	"kenshop/pkg/common/hostgen"
	errors "kenshop/pkg/error"
	"kenshop/pkg/log"

	"github.com/gin-gonic/gin"
)

var ErrNilHttpRegistor = errors.New("该http服务不存在注册器")

// 辅助protoc-gen-gin
type HandlerFuncMux interface {
	RegisterHandlerFunc(method string, path string, handler gin.HandlerFunc)
	Execute(gin.IRouter)
}

type Server struct {
	Ctx      context.Context
	Engine   *gin.Engine
	Host     string
	Mode     string
	InSecure bool

	Jwt       *jwt.GinJWTMiddleware
	Registor  registry.Registor
	Validator *validate.Validator
	Tracer    *ginotel.GinTracer
	Mux       HandlerFuncMux
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
}

func MustNewServer(ctx context.Context, host string, handlerMux HandlerFuncMux, opts ...ServerOption) *Server {
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
		Mux:             handlerMux,
		Server:          &http.Server{},
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

	s.Server.Addr = u.Host
	s.Server.Handler = s.Engine
	s.Instance.Endpoints = append(s.Instance.Endpoints, u)

	s.Engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{})
	})

	for _, m := range s.Middlewares {
		s.Engine.Use(m)
	}

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
	return s.Registor.Deregister(ctx, s.Instance.ID)
}

func (s *Server) Serve() error {
	//设置开发模式,打印路由信息格式
	gin.SetMode(s.Mode)
	//运行前前打印配置信息
	log.Infof("服务启动中,服务信息为: msg= %+v", s.Instance)

	if err := s.Validator.Excute(); err != nil {
		return err
	}

	//如果注册器为空就不进行注册而不是返回错误,
	if err := s.Register(s.Ctx, s.Instance); err != nil && err != ErrNilHttpRegistor {
		return err
	}

	s.Mux.Execute(s.Engine)
	//监听终止信号,优雅退出

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)
	ech := make(chan error, 1)
	go func() {
		if err := s.Server.ListenAndServe(); err != nil {
			//同理若注册器为空就不进行注销
			if e := s.Deregister(s.Ctx); e != nil && e != ErrNilHttpRegistor {
				log.Errorf("服务注销失败, err= %v", e)
			}
			ech <- err
		}
	}()
	select {
	case <-sign:
		close(sign)
		if e1 := s.Server.Shutdown(s.Ctx); e1 != nil {
			log.Errorf("服务注销失败, err= %v", e1)
			return e1
		}
		if e2 := s.Deregister(s.Ctx); e2 != nil && e2 != ErrNilHttpRegistor {
			log.Errorf("服务注销失败, err= %v", e2)
			return e2
		}
		log.Info("服务正常注销")
		return nil
	case err := <-ech:
		close(ech)
		return err
	}
}
