package resourse

import (
	"context"
	"kenshop/goken/server/rpcserver"
	"kenshop/pkg/config"
	userconfig "kenshop/service/user/internal/config"
	usercontroller "kenshop/service/user/internal/controller"
	userdata "kenshop/service/user/internal/data"
	"net"
	"os"
	"strconv"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var UserData *userdata.GormUserData
var Ctx context.Context
var ConfLoader *config.Loader
var Conf *userconfig.ServerConf
var Logger *otelzap.Logger
var Server *rpcserver.Server
var UserServer *usercontroller.UserServer
var Pwd string

func init() {
	Ctx = context.Background()
	Pwd, _ = os.Getwd()
	InitLogger()
	InitConf()
	InitDB()
	InitServer()
	ip, port, _ := net.SplitHostPort(Server.Host)
	Conf.Ip = ip
	Conf.Port, _ = strconv.Atoi(port)
	Logger.Sugar().Info("服务配置文件为: ", Conf)
}
