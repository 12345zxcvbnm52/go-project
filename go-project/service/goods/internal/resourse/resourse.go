package resourse

import (
	"context"
	"kenshop/goken/server/rpcserver"
	"kenshop/pkg/config"
	goodsconfig "kenshop/service/goods/internal/config"
	goodscontroller "kenshop/service/goods/internal/controller"
	goodsdata "kenshop/service/goods/internal/data"
	"net"
	"os"
	"strconv"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var Ctx context.Context
var GoodsData *goodsdata.GormGoodsData
var ConfLoader *config.Loader
var Conf *goodsconfig.ServerConf
var Logger *otelzap.Logger
var Server *rpcserver.Server
var GoodsServer *goodscontroller.GoodsServer
var Pwd string

func init() {
	Ctx = context.Background()
	Pwd, _ = os.Getwd()
	InitLogger()
	InitConf()
	InitDB()
	InitServer()
	InitOtel()
	ip, port, _ := net.SplitHostPort(Server.Host)
	Conf.Ip = ip
	Conf.Port, _ = strconv.Atoi(port)
	Logger.Sugar().Info("服务配置文件为: ", Conf)
}
