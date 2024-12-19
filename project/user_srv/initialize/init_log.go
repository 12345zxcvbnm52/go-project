package initialize

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//用于zap注册的io.Writer不要close
	logErr, err := os.OpenFile("./log/error.log", os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	// logWarn, err := os.OpenFile("./log/warning.log", os.O_RDWR|os.O_TRUNC, 0666)
	// if err != nil {
	// 	panic(err)
	// }

	logMsg, err := os.OpenFile("./log/msg.log", os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	// logErr.Truncate(0)
	// logWarn.Truncate(0)
	// logMsg.Truncate(0)
	encoderCore := zapcore.NewJSONEncoder(encoderConfig)
	//注意,日志等级是有包含关系的,即error级日志会被warn级记录
	coreErr := zapcore.NewCore(encoderCore, zapcore.AddSync(logErr), zap.ErrorLevel)
	//coreWarn := zapcore.NewCore(encoderCore, zapcore.AddSync(logWarn), zap.WarnLevel)
	coreMsg := zapcore.NewCore(encoderCore, zapcore.AddSync(logMsg), zap.InfoLevel)
	//单独将error日志和warning日志记录到特殊文件中
	core := zapcore.NewTee(coreErr, coreMsg)

	log := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(log)
}
