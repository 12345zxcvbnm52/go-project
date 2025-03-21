package otelkgin

import (
	"context"
	"errors"
	"fmt"
	ktrace "kenshop/pkg/trace"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type GinTracer struct {
	Ctx context.Context
	// SpanGinCtxKey是gin.Context中找到Span的键
	SpanGinCtxKey string
	UseAbort      bool
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
		UseAbort:      true,
		TracerName:    "goken",
	}

	for _, opt := range opts {
		opt(gt)
	}
	return gt
}

func (g *GinTracer) addDefaultAttributes(ctx *gin.Context, span trace.Span) trace.Span {
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
	return span
}

// 在TraceHandler的基础上允许自定义span名称而非方法+路径名
func (g *GinTracer) TraceHandlerWithSpanName(spanName string, opts ...trace.SpanStartOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tracer := otel.Tracer(g.TracerName)
		c, span := tracer.Start(ctx.Request.Context(), spanName, opts...)
		ctx.Set(g.SpanGinCtxKey, c)
		ctx.Set("tracer-name", g.TracerName)
		defer func() {
			g.addDefaultAttributes(ctx, span)
			span.End()
		}()
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
			g.addDefaultAttributes(ctx, span)
			span.End()
		}()
		ctx.Next()
	}
}

// TraceHandlerFormSpan返回一个HandlerFunc,它会从传入的http报文中或gin.Context中以TextMap(键值对)格式提取父Span数据,
// 并使用Derive启动一个与父Span相关联的新Span并记录后续时间
func (g *GinTracer) TraceHandlerFormCarrier(opts ...trace.SpanStartOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := g.SpanFromGinContext(ctx)
		if span.SpanContext().IsValid() {
			spanCtx := trace.ContextWithSpan(ctx.Request.Context(), span)
			tracer := otel.GetTracerProvider().Tracer(g.TracerName)
			cspanContext, cspan := tracer.Start(spanCtx, fmt.Sprintf("%s-%s", ctx.Request.Method, ctx.Request.URL.Path), opts...)
			ctx.Set(g.SpanGinCtxKey, cspanContext)
			defer func() {
				g.addDefaultAttributes(ctx, span)
				cspan.End()
			}()
		}

		ctx.Next()
	}
}

// TraceHandlerFormSpan返回一个HandlerFunc,它会从传入的http报文中或gin.Context中以TextMap(键值对)格式提取父Span数据,
// 并使用Derive启动一个与父Span相关联的新Span并记录后续时间,这个函数允许自定义的spanName
func (g *GinTracer) TraceHandlerFormCarrierWithSpanName(spanName string, opts ...trace.SpanStartOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := g.SpanFromGinContext(ctx)
		if span.SpanContext().IsValid() {
			spanCtx := trace.ContextWithSpan(ctx.Request.Context(), span)
			tracer := otel.GetTracerProvider().Tracer(g.TracerName)
			cspanContext, cspan := tracer.Start(spanCtx, spanName, opts...)
			ctx.Set(g.SpanGinCtxKey, cspanContext)
			defer func() {
				g.addDefaultAttributes(ctx, span)
				cspan.End()
			}()
		}
		ctx.Next()
	}
}

// 该函数会从gin.context中提取父Span数据,其中会先从gin.context的键值对中寻找span信息
// 再者会从gin.Request.Header中寻找,否则将会返回一个noop span
func (g *GinTracer) SpanFromGinContext(ctx *gin.Context) trace.Span {
	spanContextAny, _ := ctx.Get(g.SpanGinCtxKey)
	spanContext, typeOk := spanContextAny.(context.Context)
	if spanContext != nil && typeOk {
		if trace.SpanContextFromContext(spanContext).IsValid() {
			return trace.SpanFromContext(spanContext)
		}
		spanContext = ktrace.Extract(g.Ctx, propagation.HeaderCarrier(ctx.Request.Header))
		if trace.SpanContextFromContext(spanContext).IsValid() {
			return trace.SpanFromContext(spanContext)
		}
	}
	return noop.Span{}
}

// 该函数将Span的元信息注入到请求头中,适合跟踪链式请求(如客户端->服务1->服务2),
func (g *GinTracer) InjectToGinHeaders(ctx *gin.Context, spanCtx context.Context) {
	ktrace.Inject(spanCtx, propagation.HeaderCarrier(ctx.Request.Header))
}

// 该函数将Span的元信息注入到请求头中,适合跟踪链式请求(如客户端->服务1->服务2),
func (g *GinTracer) InjectToGinContext(ctx *gin.Context, spanCtx context.Context) {
	ctx.Set(g.SpanGinCtxKey, spanCtx)
}

func WithUseAbort(UseAbort bool) GinTracerOption {
	return func(g *GinTracer) {
		g.UseAbort = UseAbort
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
