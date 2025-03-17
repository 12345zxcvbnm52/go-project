package log

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	Logger *otelzap.Logger
}

func (g GormLogger) LogMode(_ logger.LogLevel) logger.Interface {
	return g
}

func setGormPrefix(s string, i ...interface{}) string {
	s = "[gorm] " + s
	return fmt.Sprintf(s, i)
}

func (g GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	g.Logger.InfoContext(ctx, setGormPrefix(s, i))
}

func (g GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	g.Logger.WarnContext(ctx, setGormPrefix(s, i))
}

func (g GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	g.Logger.ErrorContext(ctx, setGormPrefix(s, i))
}

func (g GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	const traceStr = "Cost: %v, Rows: %v, SQL: %s"
	elapsed := time.Since(begin)
	sql, rows := fc()

	// 初始化字段
	var fields []zap.Field
	if err != nil {
		fields = append(fields, zap.String("error", err.Error()))
	}
	fields = append(fields, zap.Float64("cost_ms", float64(elapsed.Nanoseconds())/1e6))

	if rows == -1 {
		g.Logger.DebugContext(ctx, setGormPrefix(traceStr, "-", sql), fields...)
	} else {
		fields = append(fields, zap.Int64("rows", rows)) // 记录受影响行数
		g.Logger.DebugContext(ctx, setGormPrefix(traceStr, rows, sql), fields...)
	}
}
func MustNewGormLogger(l *otelzap.Logger) *GormLogger {
	return &GormLogger{
		Logger: l,
	}
}
