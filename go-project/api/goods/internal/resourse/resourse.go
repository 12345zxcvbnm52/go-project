package resourse

import (
	"context"
	goodsconfig "kenshop/api/goods/internal/config"
	goodscontroller "kenshop/api/goods/internal/controller"
	goodsdata "kenshop/api/goods/internal/data"
	"kenshop/pkg/config"
	"net"
	"os"
	"strconv"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var Ctx context.Context
var Logger *otelzap.Logger
var ConfLoader *config.Loader
var Conf *goodsconfig.ServerConf
var GoodsData *goodsdata.GrpcGoodsData
var GoodsServer *goodscontroller.GoodsHttpServer
var Pwd string

func init() {
	Ctx = context.Background()
	Pwd, _ = os.Getwd()
	InitLogger()
	InitConf()
	InitServer()
	InitOtel()
	InitValidate()
	ip, port, _ := net.SplitHostPort(GoodsServer.Server.Host)
	Conf.Ip = ip
	Conf.Port, _ = strconv.Atoi(port)
	Logger.Sugar().Info("服务配置文件为: ", Conf)
}
