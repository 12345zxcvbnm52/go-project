package resourse

import (
	"context"
	userconfig "kenshop/api/user/internal/config"
	usercontroller "kenshop/api/user/internal/controller"
	userdata "kenshop/api/user/internal/data"
	"kenshop/goken/server/httpserver"
	"kenshop/pkg/config"
	"net"
	"os"
	"strconv"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var UserData *userdata.GrpcUserData
var Ctx context.Context
var ConfLoader *config.Loader
var Conf *userconfig.ServerConf
var Logger *otelzap.Logger
var Server *httpserver.Server
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
