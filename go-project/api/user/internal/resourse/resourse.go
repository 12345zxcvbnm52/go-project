package resourse

import (
	"context"
	userconfig "kenshop/api/user/internal/config"
	usercontroller "kenshop/api/user/internal/controller"
	userdata "kenshop/api/user/internal/data"
	"kenshop/pkg/config"
	"net"
	"os"
	"strconv"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var Ctx context.Context
var Logger *otelzap.Logger
var ConfLoader *config.Loader
var Conf *userconfig.ServerConf
var UserData *userdata.GrpcUserData
var UserServer *usercontroller.UserHttpServer
var Pwd string

func init() {
	Ctx = context.Background()
	Pwd, _ = os.Getwd()
	InitLogger()
	InitConf()
	InitServer()
	InitOtel()
	InitValidate()
	ip, port, _ := net.SplitHostPort(UserServer.Server.Host)
	Conf.Ip = ip
	Conf.Port, _ = strconv.Atoi(port)
	Logger.Sugar().Info("服务配置文件为: ", Conf)
}
