package main

import (
	"context"
	"kenshop/goken/registry/ways/consul"
	"kenshop/goken/server/httpserver"
	opengintracing "kenshop/goken/server/httpserver/middlewares/tracing"
	"kenshop/goken/server/rpcserver"
	"kenshop/goken/server/rpcserver/cinterceptors"
	proto "kenshop/proto/test"

	_ "kenshop/docs"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	gswaggof "github.com/swaggo/files"
	gswaggo "github.com/swaggo/gin-swagger"
)

func f(c *gin.Context) (map[string]string, error) {
	r := make(map[string]string)
	r["username"] = "admin"
	r["password"] = "zxcvbnm52"
	r["id"] = "123456"
	return r, nil
}

// @BasePath /v1
// @Description This is a sample API
// @Host api.example.com
// @Title My API
// @Version 1.0.0
func main() {
	ctx := context.Background()
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.199.128:8500"
	cli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	d := consul.MustNewConsulDiscover(cli)

	server := httpserver.MustNewServer(ctx, "0.0.0.0:8080",
		httpserver.WithTracer("192.168.199.128:4318", opengintracing.WithTracerName("server-test")),
		httpserver.WithGrpcClient("ken",
			rpcserver.WithClientUnaryInterceptor(cinterceptors.UnaryTracingInterceptor),
			rpcserver.WithDiscover(d),
		),
		httpserver.WithJWTMiddleware("kensame"),
	)
	s := proto.RegisterMessagingHTTPServer(server)

	s.Server.Engine.GET("/test",
		server.Jwt.RefreshHandler,
		server.Jwt.JwtAuthHandler,
		server.Jwt.AuthorizationHandler,
		server.Tracer.WithSpanHandler("gin-span"),
		s.UpdateMessage,
	)

	s.Server.Engine.GET("/login", s.Server.Jwt.LoginHandler(f))
	s.Server.Engine.GET("/swagger/*any", gswaggo.WrapHandler(gswaggof.Handler))
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
