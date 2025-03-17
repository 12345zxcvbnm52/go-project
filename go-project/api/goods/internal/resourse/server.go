package resourse

import (
	"fmt"
	goodscontroller "kenshop/api/goods/internal/controller"
	"kenshop/goken/registry/ways/consul"
	"kenshop/goken/server/httpserver"
	opengintracing "kenshop/goken/server/httpserver/middlewares/tracing"
	"kenshop/goken/server/rpcserver"
	"kenshop/goken/server/rpcserver/cinterceptors"

	"github.com/hashicorp/consul/api"
)

func InitServer() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", Conf.Consul.Ip, Conf.Consul.Port)
	cli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	d := consul.MustNewConsulDiscover(cli)

	s := httpserver.MustNewServer(Ctx,
		fmt.Sprintf("%s:%d", Conf.Ip, Conf.Port),
		httpserver.WithGrpcClient(Conf.GoodsSrv.Name, rpcserver.WithDiscover(d), rpcserver.WithClientUnaryInterceptor(cinterceptors.UnaryTracingInterceptor)),
		httpserver.WithTracer(opengintracing.WithTracerName(Conf.Otel.TracerName)),
		httpserver.WithJWTMiddleware(Conf.Jwt.Key),
	)

	GoodsServer = goodscontroller.MustNewGoodsHTTPServer(s, goodscontroller.WithLogger(Logger))
}
