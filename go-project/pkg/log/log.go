package log

import (
	"io"
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var gLogger *otelzap.Logger
var gSLogger *otelzap.SugaredLogger

func init() {
	gLogger = MustNewOtelLogger()
	gSLogger = gLogger.Sugar()
}

func MustNewOtelLogger(opts ...OptionFunc) *otelzap.Logger {
	logOpt := newDefaultOptions()
	for _, opt := range opts {
		opt(logOpt)
	}
	return mustNewOtelLogger(logOpt)
}

func mustNewOtelLogger(logOpt *Options) *otelzap.Logger {
	log := mustNewLogger(logOpt)
	return otelzap.New(log, logOpt.OtelZapOptions...)
}

func MustNewLogger(opts ...OptionFunc) *zap.Logger {
	logOpt := newDefaultOptions()
	for _, opt := range opts {
		opt(logOpt)
	}
	return mustNewLogger(logOpt)
}

func mustNewLogger(logOpt *Options) *zap.Logger {
	// 设置默认的日志输出位置
	//teeFlag 用于判断是否需要将错误日志和标准日志连接(tee)起来
	var teeFlag bool = false
	if len(logOpt.OutputPaths) == 0 {
		logOpt.OutputPaths = append(logOpt.OutputPaths, "stdout")
	}
	if len(logOpt.ErrorOutputPaths) == 0 {
		logOpt.ErrorOutputPaths = append(logOpt.ErrorOutputPaths, "stderr")
	} else {
		//若既有错误输出和标准输出则需要连接日志
		teeFlag = true
	}

	var encoderCfg zapcore.EncoderConfig
	var encoder zapcore.Encoder

	// 设置zap.logger的Mode
	if logOpt.Development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	// 无论是production还是development都用这两个配置
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	// 如果是文本且开启了彩色打印则在zap中开启彩色打印,后续也可以改为在有stdout作为outpath时才开启
	if logOpt.EnableColor && logOpt.Format == ConsoleFormat {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// 设置zap.logger的Format
	if logOpt.Format == JsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	var multiOut io.Writer
	var multiErrOut io.Writer

	for _, path := range logOpt.OutputPaths {
		var file io.Writer
		var err error = nil

		if path == "stdout" {
			file = os.Stdout
		} else if path == "stderr" {
			file = os.Stderr
		} else {
			file, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		}
		if err != nil {
			panic(err)
		}

		if multiOut == nil {
			multiOut = file
		} else {
			multiOut = io.MultiWriter(multiOut, file)
		}
	}

	for _, path := range logOpt.ErrorOutputPaths {
		var file io.Writer
		var err error = nil

		if path == "stdout" {
			file = os.Stdout
		} else if path == "stderr" {
			file = os.Stderr
		} else {
			file, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		}
		if err != nil {
			panic(err)
		}

		if multiErrOut == nil {
			multiErrOut = file
		} else {
			multiErrOut = io.MultiWriter(multiErrOut, file)
		}
	}

	wsyncer := zapcore.AddSync(multiOut)
	core := zapcore.NewCore(encoder, wsyncer, logOpt.Level)
	if teeFlag {
		werrSyncer := zapcore.AddSync(multiErrOut)
		core = zapcore.NewTee(core, zapcore.NewCore(encoder, werrSyncer, logOpt.ErrorLevel))
	}

	log := zap.New(core, logOpt.ZapOptions...)
	return log
}
