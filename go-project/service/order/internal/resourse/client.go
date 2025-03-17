package resourse

import (
	"fmt"
	"kenshop/goken/registry/discover"
	"kenshop/goken/registry/ways/consul"
	"kenshop/goken/server/rpcserver"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

func InitClient() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", Conf.Consul.Ip, Conf.Consul.Port)
	cli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	d := consul.MustNewConsulDiscover(cli)
	b := discover.MustNewBuilder(d)
	resolver.Register(b)
	InventoryClient = rpcserver.MustNewClient(Ctx, "inventory_srv", rpcserver.WithDiscover(d))
	GoodsClient = rpcserver.MustNewClient(Ctx, "goods_srv", rpcserver.WithDiscover(d))
}
