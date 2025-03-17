package resourse

import (
	"context"
	orderconfig "kenshop/api/order/internal/config"
	ordercontroller "kenshop/api/order/internal/controller"
	orderdata "kenshop/api/order/internal/data"
	"kenshop/pkg/config"
	"net"
	"os"
	"strconv"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var Ctx context.Context
var Logger *otelzap.Logger
var ConfLoader *config.Loader
var Conf *orderconfig.ServerConf
var OrderData *orderdata.GrpcOrderData
var OrderServer *ordercontroller.OrderHttpServer
var Pwd string

func init() {
	Ctx = context.Background()
	Pwd, _ = os.Getwd()
	InitLogger()
	InitConf()
	InitServer()
	InitOtel()
	InitValidate()
	ip, port, _ := net.SplitHostPort(OrderServer.Server.Host)
	Conf.Ip = ip
	Conf.Port, _ = strconv.Atoi(port)
	Logger.Sugar().Info("服务配置文件为: ", Conf)
}
