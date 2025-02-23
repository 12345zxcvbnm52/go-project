package log

import (
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ConsoleFormat = "console"
	JsonFormat    = "json"
)

type Options struct {
	OutputPaths      []string      `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths []string      `json:"error-output-paths" mapstructure:"error-output-paths"`
	Level            zapcore.Level `json:"level"              mapstructure:"level"`
	ErrorLevel       zapcore.Level `json:"error-level"              mapstructure:"error-level"`
	Format           string        `json:"format"             mapstructure:"format"`
	EnableColor      bool          `json:"enable-color"       mapstructure:"enable-color"`
	Development      bool          `json:"development"        mapstructure:"development"`
	ZapOptions       []zap.Option
	OtelZapOptions   []otelzap.Option
}

type OptionFunc func(o *Options)

func newDefaultOptions() *Options {
	opts := &Options{
		Level:            zapcore.InfoLevel,
		ErrorLevel:       zapcore.ErrorLevel,
		Format:           ConsoleFormat,
		EnableColor:      false,
		OutputPaths:      []string{},
		ErrorOutputPaths: []string{},
		Development:      false,
	}
	opts.ZapOptions = append(opts.ZapOptions, zap.AddCaller())
	opts.ZapOptions = append(opts.ZapOptions, zap.AddStacktrace(zap.WarnLevel))
	return opts
}

// 如果使用该函数且path.len>=1,则将不自动使用stdout,除非在path中添加stdout
// 如果path.len==0或不使用该函数则默认日志输出在stdout中
func WithOutputPaths(path ...string) OptionFunc {
	return func(o *Options) {
		o.OutputPaths = append(o.OutputPaths, path...)
	}
}

// 如果使用该函数且path.len>=1,则将不自动使用stderr,除非在path中添加stderr
// 如果path.len==0或不使用该函数则默认日志输出在stderr中
func WithErrOutPaths(path ...string) OptionFunc {
	return func(o *Options) {
		o.ErrorOutputPaths = append(o.ErrorOutputPaths, path...)
	}
}

// 默认的Level为Info级别
func WithLevel(lvl zapcore.Level) OptionFunc {
	return func(o *Options) {
		o.Level = lvl
	}
}

// 设置被认为错误的日志级别,默认的ErrorLevel为Info级别
func WithErrLevel(lvl zapcore.Level) OptionFunc {
	return func(o *Options) {
		o.ErrorLevel = lvl
	}
}

// 如果format格式错误或默认则使用console,可选"console"以及"json"
func WithFormat(format string) OptionFunc {
	return func(o *Options) {
		o.Format = format
	}
}

// 默认不开启彩色打印
func WithColor(on bool) OptionFunc {
	return func(o *Options) {
		o.EnableColor = on
	}
}

// 创建一个development级别的日志(日志将打印更多信息)
func WithDevelopmentLogging(on bool) OptionFunc {
	return func(o *Options) {
		o.Development = on
	}
}

// 允许创建时使用zap自带的Options
func WithZapOptions(opt ...zap.Option) OptionFunc {
	return func(o *Options) {
		o.ZapOptions = append(o.ZapOptions, opt...)
	}
}

// 允许创建时使用otelzap自带的Options
func WithOtelZapOptions(opt ...otelzap.Option) OptionFunc {
	return func(o *Options) {
		o.OtelZapOptions = append(o.OtelZapOptions, opt...)
	}
}
