package httpserver

import (
	"context"
	"fmt"

	//"github.com/penglongli/gin-metrics/ginmetrics"
	"net/http"
	"time"

	mws "kenshop/goken/server/httpserver/middlewares"
	//"kenshop/goken/server/httpserver/pprof"
	//"kenshop/goken/server/httpserver/validation"
	"kenshop/pkg/log"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

type JwtInfo struct {
	// 默认为"JWT"
	Realm string
	// 默认为empty
	Key string
	// jwt-token的有效时间,默认为七天
	Validity time.Duration
	// jwt-token的refresh最大间隔时间,默认为七天
	Refreshy time.Duration
}

// wrapper for gin.Engine
type Server struct {
	*gin.Engine

	//端口号,默认为8080
	port int

	//开发模式, 默认值debug
	mode string
	//是否开启健康检查接口,默认开启,如果开启会自动添加/health接口
	healthz bool

	//是否开启pprof接口, 默认开启, 如果开启会自动添加/debug/pprof接口
	enableProfiling bool

	//是否开启metrics接口,默认开启,如果开启会自动添加/metrics接口
	enableMetrics bool

	//中间件
	middlewares []string

	//默认值 zh
	locale string
	trans  ut.Translator
	//从
	server *http.Server

	serviceName string
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		port:            8080,
		mode:            "debug",
		healthz:         true,
		enableProfiling: true,
		Engine:          gin.New(),
		locale:          "zh",
		serviceName:     "goken",
	}

	for _, o := range opts {
		o(srv)
	}

	//srv.Use(mws.TracingHandler(srv.serviceName))
	for _, m := range srv.middlewares {
		mw, ok := mws.Middlewares[m]
		if !ok {
			log.Warnf("无法寻找使用到该中间件: %s", m)
		}
		srv.Use(mw())
	}

	return srv
}

func (s *Server) Translator() ut.Translator {
	return s.trans
}

// start rest server
func (s *Server) Start(ctx context.Context) error {
	//设置开发模式,打印路由信息格式
	gin.SetMode(s.mode)
	// gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	// 	log.Infof("%-6s %-s --> %s(%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	// }

	err := s.initTrans()
	if err != nil {
		log.Errorf("初始翻译器失败 %s", err.Error())
		return err
	}

	// //注册mobile验证码
	// validation.RegisterMobile(s.trans)

	// //根据配置初始化pprof路由
	// if s.enableProfiling {
	// 	pprof.Register(s.Engine)
	// }

	if s.enableMetrics {
		// get global Monitor object
		m := ginmetrics.GetMonitor()
		// +optional set metric path, default /debug/metrics
		m.SetMetricPath("/metrics")
		// +optional set slow time, default 5s
		// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
		// used to p95, p99
		m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
		m.Use(s)
	}

	log.Infof("服务正在启动中... 预计监听的端口为: %d", s.port)

	//gin.Run内部的逻辑抽取出来,便于优雅退出
	address := fmt.Sprintf(":%d", s.port)
	s.server = &http.Server{
		Addr:    address,
		Handler: s.Engine,
	}
	_ = s.SetTrustedProxies(nil)
	if err = s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Errorf("http服务启动失败 err= %s", err.Error())
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Infof("关闭http服务中...")
	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("http服务关闭失败 err= %s", err.Error())
		return err
	}
	log.Infoln("服务关闭成功")
	return nil
}
