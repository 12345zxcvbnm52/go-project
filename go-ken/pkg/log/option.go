package log

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	flagLevel             = "log.level"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagFormat            = "log.format"
	flagEnableColor       = "log.enable-color"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagDevelopment       = "log.development"
	flagName              = "log.name"

	consoleFormat = "console"
	jsonFormat    = "json"
)

// Cascadia Mono SemiLight
// Options contains configuration items related to log.
type Options struct {
	//标准日志的输出路径
	OutputPaths []string `json:"output-paths"       mapstructure:"output-paths"`
	//错误日志的输出路径
	ErrorOutputPaths []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	//启动的日志级别
	Level string `json:"level"              mapstructure:"level"`
	//打印的格式,这里只支持json格式和console格式
	Format string `json:"format"             mapstructure:"format"`
	//是否禁用显示调用者信息
	DisableCaller bool `json:"disable-caller"     mapstructure:"disable-caller"`
	//是否禁用显示错误日志的堆栈信息
	DisableStacktrace bool `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	//是否允许彩色日志
	EnableColor bool `json:"enable-color"       mapstructure:"enable-color"`
	//是否开启开发者模式,如果为false则开启为production模式
	Development bool `json:"development"        mapstructure:"development"`
	//是否为日志实例设置名称,以便用于区分不同日志的来源的Logger
	Name string `json:"name"               mapstructure:"name"`
	//以下两个选项用于控制分布式系统中的链路追踪
	//开启TraceID模式,通过TraceID可以关联不同服务的日志
	EnableTraceID bool `json:"enable-trace-id"     mapstructure:"enable-trace-id"`
	//是否启用链路追踪堆栈信息
	EnableTraceStack bool `json:"enable-trace-stack" mapstructure:"enable-trace-stack"`
}

// 创建一个默认的Logger Options,默认的options不支持分布式链路追踪
func NewOptions() *Options {
	return &Options{
		Level:             zapcore.InfoLevel.String(),
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            consoleFormat,
		EnableColor:       false,
		Development:       false,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}
}

// 验证options是否合规
func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	//检查level是否合法
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}
	//检查格式是否合法
	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

// 这里使用pflag实现通过命令行填充options参数,具体实现可以自定义
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

func (o *Options) MarshalJSON() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// TODO 将数据序列化为YAML
func (o *Options) MarshalYAML() string {
	data, _ := json.Marshal(o)
	return string(data)
}

// 生成全局的Logger
func (o *Options) Build() error {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	encodeLevel := zapcore.CapitalLevelEncoder
	if o.Format == consoleFormat && o.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zc := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       o.Development,
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: o.Format,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "timestamp",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeTime:     timeEncoder,
			EncodeDuration: milliSecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      o.OutputPaths,
		ErrorOutputPaths: o.ErrorOutputPaths,
	}
	logger, err := zc.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		return err
	}
	zap.RedirectStdLog(logger.Named(o.Name))
	zap.ReplaceGlobals(logger)

	return nil
}

type Option func(l *Logger)

// WithMinLevel 设置记录到span上的日志消息的最低日志级别
// 默认级别是警告级别及以上
func WithMinLevel(lvl zapcore.Level) Option {
	return func(l *Logger) {
		l.minLevel = lvl
	}
}

// WithMinLevel 设置记录到span上的错误日志消息的最低日志级别
// 默认级别是错误级别及以上
func WithErrorStatusLevel(lvl zapcore.Level) Option {
	return func(l *Logger) {
		l.errorStatusLevel = lvl
	}
}

func WithCaller(on bool) Option {
	return func(l *Logger) {
		l.caller = on
	}
}

func WithStackTrace(on bool) Option {
	return func(l *Logger) {
		l.stackTrace = on
	}
}

// WithTraceIDField 配置日志记录器在结构化日志消息中添加trace_id字段
// 此选项仅在后端不支持 OTLP(OpenTelemetry协议)且需要通过解析日志消息提取结构化信息时有用
func WithTraceIDField(on bool) Option {
	return func(l *Logger) {
		l.withTraceID = on
	}
}
