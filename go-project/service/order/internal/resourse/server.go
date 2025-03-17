package resourse

import (
	"fmt"
	"kenshop/goken/registry/discover"
	"kenshop/goken/registry/registor"
	"kenshop/goken/registry/ways/consul"
	"kenshop/goken/server/rpcserver"
	"kenshop/goken/server/rpcserver/sinterceptors"
	ordercontroller "kenshop/service/order/internal/controller"
	orderlogic "kenshop/service/order/internal/logic"

	"github.com/dtm-labs/dtmdriver"
	"github.com/hashicorp/consul/api"
)

func InitServer() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", Conf.Consul.Ip, Conf.Consul.Port)
	cli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	regisor := registor.MustNewRegister(consul.MustNewConsulRegistor(cli))

	Server = rpcserver.MustNewServer(Ctx,
		rpcserver.WithServiceID(Conf.Id),
		rpcserver.WithServiceName(Conf.Name),
		rpcserver.WithUnaryInts(sinterceptors.UnaryTracingInterceptor),
		rpcserver.WithVersion(Conf.Version),
		rpcserver.WithRegistor(regisor),
		rpcserver.WithHost("192.168.199.128:0"),
	)

	kdtm := discover.MustNewKDtmDriver(Ctx, discover.MustNewBuilder(consul.MustNewConsulDiscover(cli)), discover.WithRegistor(regisor))
	dtmdriver.Register(kdtm)
	if err := dtmdriver.Use("dtm-driver-goken"); err != nil {
		panic(err)
	}

	OrderServer = &ordercontroller.OrderServer{}
	OrderServer.Logger = Logger
	OrderServer.Service = &orderlogic.OrderService{}
	OrderServer.Service.OrderData = OrderData
}
