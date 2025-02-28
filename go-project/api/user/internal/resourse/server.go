package resourse

import (
	"fmt"
	usercontroller "kenshop/api/user/internal/controller"
	"kenshop/goken/registry"
	"kenshop/goken/registry/registor"
	"kenshop/goken/registry/ways/consul"
	"kenshop/goken/server/httpserver"

	"github.com/hashicorp/consul/api"
)

func InitServer() {
	ins := &registry.ServiceInstance{
		ID:      Conf.Id,
		Name:    Conf.Name,
		Version: Conf.Version,
	}

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", Conf.Consul.Ip, Conf.Consul.Port)
	cli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	regisor := registor.MustNewRegister(consul.MustNewConsulRegistor(cli))
	Server = httpserver.MustNewServer(Ctx,
		httpserver.WithServiceInstance(ins),
		httpserver.WithRegistor(regisor),
		httpserver.WithHost("192.168.199.128:0"),
	)

	UserServer = &usercontroller.UserServer{}
	UserServer.Logger = Logger
	UserServer.Service = &userlogic.UserService{}
	UserServer.Service.UserData = UserData
}
