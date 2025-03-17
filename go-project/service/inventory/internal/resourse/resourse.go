package resourse

import (
	"context"
	"kenshop/goken/server/rpcserver"
	"kenshop/pkg/config"
	inventoryconfig "kenshop/service/inventory/internal/config"
	inventorycontroller "kenshop/service/inventory/internal/controller"
	inventorydata "kenshop/service/inventory/internal/data"
	"net"
	"os"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var Ctx context.Context
var InventoryData *inventorydata.GormInventoryData
var ConfLoader *config.Loader
var Conf *inventoryconfig.ServerConf
var Logger *otelzap.Logger
var Server *rpcserver.Server
var InventoryServer *inventorycontroller.InventoryServer
var Pwd string
var Consumer rocketmq.PushConsumer

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
	rlog.SetLogLevel("warn")
	Logger.Sugar().Info("服务配置文件为: ", Conf)
}
