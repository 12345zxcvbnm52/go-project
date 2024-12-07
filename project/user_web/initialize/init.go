package initialize

import (
	gb "user_web/global"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
)

func init() {
	//log和config的初始一定要写在最前面
	//consul和router也具有先后顺序
	InitLog()
	InitConfig()
	InitConsul()
	InitRouter()
	err := InitTranslator("zh")
	if err != nil {
		zap.S().Errorw("无法获取中文翻译器,尝试获取英文翻译器", "msg", err.Error())
		nerr := InitTranslator("en")
		if nerr != nil {
			zap.S().Errorw("无法获取任意翻译器")
			panic(nerr)
		}
	}
	InitValidator()

	zap.S().Infoln("ServerConfig is : ", gb.ServerConfig)
	zap.S().Infoln("ConnPoolConfig is : ", gb.ConnPoolConfig)
}
