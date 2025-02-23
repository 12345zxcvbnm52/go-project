package log

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap/zapcore"
)

func Logger() *otelzap.Logger {
	return gLogger
}

func L() *otelzap.Logger {
	return gLogger
}

func Sugar() *otelzap.SugaredLogger {
	return gSLogger
}

func S() *otelzap.SugaredLogger {
	return gSLogger
}

func Flush() {
	_ = gLogger.Logger.Sync()
}

func Ctx(ctx context.Context) otelzap.LoggerWithCtx {
	return gLogger.Ctx(ctx)

}

func Debug(msg string, fields ...zapcore.Field) {
	gLogger.Debug(msg, fields...)
}

func Debugf(format string, v ...interface{}) {
	gSLogger.Debugf(format, v...)
}

func Debugw(format string, v ...interface{}) {
	gSLogger.Debugw(format, v...)
}

func Debugln(args ...interface{}) {
	gSLogger.Debugln(args...)
}

func DebugContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	gLogger.DebugContext(ctx, msg, fields...)
}

func DebugfContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.DebugfContext(ctx, format, v...)
}

func DebugwContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.DebugwContext(ctx, format, v...)
}

func Info(msg string, fields ...zapcore.Field) {
	gLogger.Info(msg, fields...)
}

func Infof(format string, v ...interface{}) {
	gSLogger.Infof(format, v...)
}

func Infow(format string, v ...interface{}) {
	gSLogger.Infow(format, v...)
}

func Infoln(args ...interface{}) {
	gSLogger.Infoln(args...)
}

func InfoContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	gLogger.InfoContext(ctx, msg, fields...)
}

func InfofContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.InfofContext(ctx, format, v...)
}

func InfowContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.InfowContext(ctx, format, v...)
}

func Warn(msg string, fields ...zapcore.Field) {
	gLogger.Warn(msg, fields...)
}

func Warnf(format string, v ...interface{}) {
	gSLogger.Warnf(format, v...)
}

func Warnw(format string, v ...interface{}) {
	gSLogger.Warnw(format, v...)
}

func Warnln(args ...interface{}) {
	gSLogger.Warnln(args...)
}

func WarnContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	gLogger.WarnContext(ctx, msg, fields...)
}

func WarnfContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.WarnfContext(ctx, format, v...)
}

func WarnwContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.ErrorwContext(ctx, format, v...)
}

func Error(msg string, fields ...zapcore.Field) {
	gLogger.Error(msg, fields...)
}

func Errorf(format string, v ...interface{}) {
	gSLogger.Errorf(format, v...)
}

func Errorw(format string, v ...interface{}) {
	gSLogger.Errorw(format, v...)
}

func Errorln(args ...interface{}) {
	gSLogger.Errorln(args...)
}

func ErrorContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	gLogger.ErrorContext(ctx, msg, fields...)
}

func ErrorfContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.ErrorfContext(ctx, format, v...)
}

func ErrorwContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.ErrorwContext(ctx, format, v...)
}

func DPanic(msg string, fields ...zapcore.Field) {
	gLogger.DPanic(msg, fields...)
}

func DPanicf(format string, v ...interface{}) {
	gSLogger.DPanicf(format, v...)
}

func DPanicw(format string, v ...interface{}) {
	gSLogger.DPanicw(format, v...)
}

func DPanicln(args ...interface{}) {
	gSLogger.DPanicln(args...)
}

func DPanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	gLogger.DPanicContext(ctx, msg, fields...)
}

func DPanicfContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.DPanicfContext(ctx, format, v...)
}

func DPanicwContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.DPanicwContext(ctx, format, v...)
}

func Panic(msg string, fields ...zapcore.Field) {
	gLogger.Panic(msg, fields...)
}

func Panicf(format string, v ...interface{}) {
	gSLogger.Panicf(format, v...)
}

func Panicw(format string, v ...interface{}) {
	gSLogger.Panicw(format, v...)
}

func Panicln(args ...interface{}) {
	gSLogger.Panicln(args...)
}

func PanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	gLogger.PanicContext(ctx, msg, fields...)
}

func PanicfContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.PanicfContext(ctx, format, v...)
}

func PanicwContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.PanicwContext(ctx, format, v...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	gLogger.Fatal(msg, fields...)
}

func Fatalf(format string, v ...interface{}) {
	gSLogger.Fatalf(format, v...)
}

func Fatalw(format string, v ...interface{}) {
	gSLogger.Fatalw(format, v...)
}

func Fatalln(args ...interface{}) {
	gSLogger.Fatalln(args...)
}

func FatalContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	gLogger.FatalContext(ctx, msg, fields...)
}

func FatalfContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.FatalfContext(ctx, format, v...)
}

func FatalwContext(ctx context.Context, format string, v ...interface{}) {
	gSLogger.FatalwContext(ctx, format, v...)
}
