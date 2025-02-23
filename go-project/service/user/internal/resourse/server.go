package resourse

import (
	"fmt"
	"kenshop/goken/registry"
	"kenshop/goken/registry/registor"
	"kenshop/goken/registry/ways/consul"
	"kenshop/goken/server/rpcserver"
	usercontroller "kenshop/service/user/internal/controller"
	userlogic "kenshop/service/user/internal/logic"

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
	Server = rpcserver.MustNewServer(Ctx,
		rpcserver.WithServiceInstance(ins),
		rpcserver.WithRegistor(regisor),
		rpcserver.WithHost("192.168.199.128:0"),
	)

	UserServer = &usercontroller.UserServer{}
	UserServer.Logger = Logger
	UserServer.Service = &userlogic.UserService{}
	UserServer.Service.UserData = UserData
}
