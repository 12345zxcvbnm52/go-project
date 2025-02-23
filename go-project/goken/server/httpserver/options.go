package httpserver

import "github.com/gin-gonic/gin"

type ServerOption func(*Server)

func WithEnableProfiling(profiling bool) ServerOption {
	return func(s *Server) {
		s.enableProfiling = profiling
	}
}

func WithMode(mode string) ServerOption {
	return func(s *Server) {
		switch mode {
		case gin.ReleaseMode:
		case gin.DebugMode:
		case gin.TestMode:
		default:
			mode = gin.DebugMode
		}
		s.mode = mode
	}
}

func WithServiceName(srvName string) ServerOption {
	return func(s *Server) {
		s.serviceName = srvName
	}
}

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

func WithMiddlewares(middlewares []string) ServerOption {
	return func(s *Server) {
		s.middlewares = middlewares
	}
}

func WithHealthz(healthz bool) ServerOption {
	return func(s *Server) {
		s.healthz = healthz
	}
}

func WithTransLocale(locale string) ServerOption {
	return func(s *Server) {
		s.locale = locale
	}
}

func WithMetrics(enable bool) ServerOption {
	return func(o *Server) {
		o.enableMetrics = enable
	}
}
