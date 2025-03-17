package resourse

import (
	"context"
	"kenshop/goken/server/rpcserver"
	"kenshop/pkg/config"
	orderconfig "kenshop/service/order/internal/config"
	ordercontroller "kenshop/service/order/internal/controller"
	orderdata "kenshop/service/order/internal/data"
	"net"
	"os"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var Ctx context.Context
var OrderData *orderdata.GormOrderData
var ConfLoader *config.Loader
var Conf *orderconfig.ServerConf
var Logger *otelzap.Logger
var Server *rpcserver.Server
var OrderServer *ordercontroller.OrderServer
var InventoryClient *rpcserver.Client
var GoodsClient *rpcserver.Client
var Pwd string
var Producer rocketmq.Producer
var TransProducer rocketmq.TransactionProducer
var TimeoutConsumer rocketmq.PushConsumer

func init() {
	Ctx = context.Background()
	Pwd, _ = os.Getwd()
	InitLogger()
	InitConf()
	InitClient()
	InitDB()
	InitServer()
	InitOtel()

	ip, port, _ := net.SplitHostPort(Server.Host)
	Conf.Ip = ip
	Conf.Port, _ = strconv.Atoi(port)
	Logger.Sugar().Info("服务配置文件为: ", Conf)
}
