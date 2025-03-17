package opengintracing

import (
	"context"
	"errors"
	"fmt"
	ktrace "kenshop/pkg/trace"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type GinTracer struct {
	Ctx context.Context
	// SpanGinCtxKey是gin.Context中找到Span的键
	SpanGinCtxKey string
	AbortOnErrors bool
	TracerName    string
}

var (
	ErrSpanNotFound = errors.New("span was not found in context")
)

type GinTracerOption func(*GinTracer)

func MustNewGinTracer(ctx context.Context, opts ...GinTracerOption) *GinTracer {
	gt := &GinTracer{
		Ctx:           ctx,
		SpanGinCtxKey: "gin-goken",
		AbortOnErrors: true,
		TracerName:    "goken",
	}

	for _, opt := range opts {
		opt(gt)
	}
	return gt
}

// 在TraceHandler的基础上允许自定义span名称而非方法+路径名
func (g *GinTracer) TraceHandlerWithSpanName(spanName string, opts ...trace.SpanStartOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tracer := otel.Tracer(g.TracerName)
		c, span := tracer.Start(g.Ctx, spanName, opts...)
		ctx.Set(g.SpanGinCtxKey, c)
		ctx.Set("tracer-name", g.TracerName)
		defer span.End()
		ctx.Next()
	}
}

// NewSpan会返回一个Handler,在其中启动一个新的Span并注入到请求上下文中,
// 它会测量所有后续处理程序的执行时间,
func (g *GinTracer) TraceHandler(opts ...trace.SpanStartOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tracer := otel.Tracer(g.TracerName)
		c, span := tracer.Start(g.Ctx, fmt.Sprintf("%s-%s", ctx.Request.Method, ctx.Request.URL.Path), opts...)
		c = ktrace.NewSpanOutgoingContext(c, span)
		ctx.Set("tracer-name", g.TracerName)
		ctx.Set(g.SpanGinCtxKey, c)
		defer func() {
			span.SetAttributes(
				// 记录 HTTP 方法
				semconv.HTTPMethodKey.String(ctx.Request.Method),
				// 记录完整URL
				semconv.HTTPURLKey.String(ctx.Request.URL.String()),
				// 记录客户端Host
				semconv.HTTPHostKey.String(ctx.Request.Host),
				// 记录响应状态码
				semconv.HTTPStatusCodeKey.Int(ctx.Writer.Status()),
				attribute.String("http.referer", ctx.Request.Referer()),
			)
			span.End()
		}()
		ctx.Next()
	}
}

// SpanFromHeaders返回一个HandlerFunc,它会从HTTP头中以TextMap(键值对)格式提取父Span数据,
// 并使用Derive启动一个与父Span相关联的新Span并记录后续时间
func (g *GinTracer) SpanFromHeaders(tracer trace.Tracer, spanName string, opts ...trace.SpanStartOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		spanContext := ktrace.Extract(g.Ctx, propagation.HeaderCarrier(ctx.Request.Header))

		cspanContext, cspan := tracer.Start(spanContext, spanName, opts...)
		ctx.Set(g.SpanGinCtxKey, cspanContext)
		defer cspan.End()
		ctx.Next()
	}
}

// 该函数会从gin.context中提取父Span数据,并使用Derive启动一个与Span相关联的新Span,
func (g *GinTracer) SpanFromContext(tracer trace.Tracer, spanName string, opts ...trace.SpanStartOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		spanContextAny, _ := ctx.Get(g.SpanGinCtxKey)
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
		ctx.Set(g.SpanGinCtxKey, cSpanContext)
		defer cspan.End()

		ctx.Next()
	}
}

// 该函数将Span的元信息注入到请求头中,适合跟踪链式请求(如客户端->服务1->服务2),
func (g *GinTracer) InjectToHeaders(tracer trace.Tracer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		spanContextAny, _ := ctx.Get(g.SpanGinCtxKey)
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
	spanContextAny, _ := ctx.Get(g.SpanGinCtxKey)
	spanContext, ok := spanContextAny.(context.Context)
	return spanContext, ok
}

func WithAbortOnErrors(abortOnErrors bool) GinTracerOption {
	return func(g *GinTracer) {
		g.AbortOnErrors = abortOnErrors
	}
}

func WithSpanContextKey(spanContextKey string) GinTracerOption {
	return func(g *GinTracer) {
		g.SpanGinCtxKey = spanContextKey
	}
}

func WithTracerName(tracerName string) GinTracerOption {
	return func(g *GinTracer) {
		g.TracerName = tracerName
	}
}
