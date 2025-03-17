package resourse

import (
	"context"
	inventoryconfig "kenshop/api/inventory/internal/config"
	inventorycontroller "kenshop/api/inventory/internal/controller"
	inventorydata "kenshop/api/inventory/internal/data"
	"kenshop/pkg/config"
	"net"
	"os"
	"strconv"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var Ctx context.Context
var Logger *otelzap.Logger
var ConfLoader *config.Loader
var Conf *inventoryconfig.ServerConf
var InventoryData *inventorydata.GrpcInventoryData
var InventoryServer *inventorycontroller.InventoryHttpServer
var Pwd string

func init() {
	Ctx = context.Background()
	Pwd, _ = os.Getwd()
	InitLogger()
	InitConf()
	InitServer()
	InitOtel()
	InitValidate()
	ip, port, _ := net.SplitHostPort(InventoryServer.Server.Host)
	Conf.Ip = ip
	Conf.Port, _ = strconv.Atoi(port)
	Logger.Sugar().Info("服务配置文件为: ", Conf)
}
