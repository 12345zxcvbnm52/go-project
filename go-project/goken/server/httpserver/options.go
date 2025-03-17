package httpserver

import (
	"kenshop/goken/registry"
	"kenshop/goken/server/httpserver/middlewares/jwt"
	opengintracing "kenshop/goken/server/httpserver/middlewares/tracing"
	"kenshop/goken/server/rpcserver"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServerOption func(*Server)

func WithServiceName(serviceName string) ServerOption {
	return func(s *Server) {
		s.Instance.Name = serviceName
	}
}

func WithServiceID(id string) ServerOption {
	return func(s *Server) {
		s.Instance.ID = id
	}
}

func WithServiceVersion(v string) ServerOption {
	return func(s *Server) {
		s.Instance.Version = v
	}
}

func WithMode(mode string) ServerOption {
	return func(s *Server) {
		s.Mode = mode
	}
}

func WithEnableProfiling(enable bool) ServerOption {
	return func(s *Server) {
		s.EnableProfiling = enable
	}
}

func WithEnableMetrics(enable bool) ServerOption {
	return func(s *Server) {
		s.EnableMetrics = enable
	}
}

func WithJWTMiddleware(key string) ServerOption {
	return func(s *Server) {
		s.Jwt = jwt.MustNewGinJWTMiddleware(key)
	}
}

func WithTracer(opts ...opengintracing.GinTracerOption) ServerOption {
	return func(s *Server) {
		s.Tracer = opengintracing.MustNewGinTracer(s.Ctx, opts...)
	}
}

func WithRegistor(reg registry.Registor) ServerOption {
	return func(s *Server) {
		s.Registor = reg
	}
}

func WithMiddlewares(middlewares map[string]gin.HandlerFunc) ServerOption {
	return func(s *Server) {
		for k, v := range middlewares {
			s.Middlewares[k] = v
		}
	}
}

func WithLocale(locale string) ServerOption {
	return func(s *Server) {
		s.Locale = locale
	}
}

func WithHTTPServer(server *http.Server) ServerOption {
	return func(s *Server) {
		s.Server = server
	}
}

func WithGrpcClient(target string, opts ...rpcserver.ClientOption) ServerOption {
	return func(s *Server) {
		s.GrpcCli = rpcserver.MustNewClient(s.Ctx, target, opts...)
	}
}
