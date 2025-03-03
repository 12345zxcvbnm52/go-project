package resourse

import (
	"fmt"
	usercontroller "kenshop/api/user/internal/controller"
	"kenshop/goken/registry/ways/consul"
	"kenshop/goken/server/httpserver"
	"kenshop/goken/server/rpcserver"

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
		httpserver.WithGrpcClient(Conf.UserSrv.Name, rpcserver.WithDiscover(d)),
		httpserver.WithTracer(fmt.Sprintf("%s:%d", Conf.Otel.Ip, Conf.Otel.Port)),
		httpserver.WithJWTMiddleware(Conf.Jwt.Key),
	)

	UserServer = usercontroller.MustNewUserHTTPServer(s, usercontroller.WithLogger(Logger))
}
