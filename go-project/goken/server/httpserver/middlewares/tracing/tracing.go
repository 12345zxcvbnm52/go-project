package opengintracing

import (
	"context"
	"errors"
	ktrace "kenshop/pkg/trace"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type GinTracer struct {
	Ctx            context.Context
	SpanContextKey string
	AbortOnErrors  bool
}

var (
	ErrSpanNotFound = errors.New("span was not found in context")
)

// NewSpan会返回一个Handler,在其中启动一个新的Span并注入到请求上下文中,
// 它会测量所有后续处理程序的执行时间,
func (g *GinTracer) NewSpan(tracer trace.Tracer, spanName string, opts ...trace.SpanStartOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, span := tracer.Start(g.Ctx, spanName, opts...)
		ctx.Set(g.SpanContextKey, c)
		defer span.End()
		ctx.Next()
	}
}

// SpanFromHeaders返回一个HandlerFunc,它会从HTTP头中以TextMap(键值对)格式提取父Span数据,
// 并使用Derive启动一个与父Span相关联的新Span并记录后续时间
func (g *GinTracer) SpanFromHeaders(tracer trace.Tracer, spanName string, opts ...trace.SpanStartOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		spanContext := ktrace.Extract(g.Ctx, propagation.HeaderCarrier(ctx.Request.Header))

		cspanContext, cspan := tracer.Start(spanContext, spanName, opts...)
		ctx.Set(g.SpanContextKey, cspanContext)
		defer cspan.End()
		ctx.Next()
	}
}

// 该函数会从gin.context中提取父Span数据,并使用Derive启动一个与Span相关联的新Span,
func (g *GinTracer) SpanFromContext(tracer trace.Tracer, spanName string, opts ...trace.SpanStartOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		spanContextAny, _ := ctx.Get(g.SpanContextKey)
		spanContext, typeOk := spanContextAny.(context.Context)
		if spanContext != nil && typeOk {
			opts = append(opts, trace.WithLinks(trace.LinkFromContext(spanContext)))
		} else {
			if g.AbortOnErrors {
				_ = ctx.AbortWithError(http.StatusInternalServerError, ErrSpanNotFound)
			}
			return
		}

		cSpanContext, cspan := tracer.Start(spanContext, spanName, opts...)
		ctx.Set(g.SpanContextKey, cSpanContext)
		defer cspan.End()

		ctx.Next()
	}
}

// 该函数将Span的元信息注入到请求头中,适合跟踪链式请求(如客户端->服务1->服务2),
func (g *GinTracer) InjectToHeaders(tracer trace.Tracer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		spanContextAny, _ := ctx.Get(g.SpanContextKey)
		spanContext, typeOk := spanContextAny.(context.Context)
		if spanContext != nil && typeOk {
			ktrace.Inject(spanContext, propagation.HeaderCarrier(ctx.Request.Header))
		} else {
			if g.AbortOnErrors {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrSpanNotFound)
			}
			return
		}

	}
}

// GetSpan从上下文中提取Span
func (g *GinTracer) GetSpanContext(ctx *gin.Context) (context.Context, bool) {
	spanContextAny, _ := ctx.Get(g.SpanContextKey)
	spanContext, ok := spanContextAny.(context.Context)
	return spanContext, ok
}
