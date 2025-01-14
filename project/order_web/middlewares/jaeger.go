package middlewares

import (
	"fmt"
	"net/http"
	gb "order_web/global"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

func TraceMarking() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: fmt.Sprintf("%s:%d", gb.ServerConfig.Jaeger.Host, gb.ServerConfig.Jaeger.Port),
			},
			ServiceName: "kensame",
		}
		//后续应当修改日志的写入地方
		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			zap.S().Errorw("链路追踪器生成失败", "msg", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "服务器内部错误",
			})
			ctx.Abort()
			return
		}
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()
		startSpan := tracer.StartSpan(ctx.Request.URL.Path)
		defer startSpan.Finish()
		ctx.Set("tracer", tracer)
		ctx.Set("web-span", startSpan)
		ctx.Next()
	}
}
