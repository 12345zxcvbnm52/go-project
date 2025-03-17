package trace

import (
	"context"
	"errors"
	"net/http"

	"kenshop/pkg/common/hostgen"
	"kenshop/pkg/log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

const KTraceName = "goken"

// TraceIdKey is the trace id header.
// https://www.w3.org/TR/trace-context/#trace-id
// May change it to trace-id afterward.
var TraceIdKey = http.CanonicalHeaderKey("x-trace-id")

var kPropagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

type Tracer struct {
	Ctx context.Context
	// 记录traceProvider的serviceName
	ServiceName string

	// Sampler 记录基准采样率,派生的span采样率不会低于该百分比,至多百分百采样
	Sampler float64

	ExtraResources []attribute.KeyValue

	// OtlpHeaders是OTLP gRPC或HTTP传输的headers,例如:
	//  uptrace-dsn: 'http://project2_secret_token@localhost:14317/2'
	OtlpHeaders map[string]string

	// OtlpHttpPath是OTLP HTTP传输的路径,例如:/v1/traces
	OtlpHttpPath string

	// OtlpHttpSecure 指定 OTLP HTTP传输是否使用安全的HTTPS协议,
	OtlpHttpSecure bool

	// Global 表示是否将生成的TraceProvider加入到全局中
	Global bool
}

type TracerOption func(*Tracer)

func init() {
	otel.SetTextMapPropagator(kPropagator)
	te, err := stdouttrace.New()
	if err != nil {
		panic(err)
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(te))
	otel.SetTracerProvider(tp)
}

// NewTracer 创建一个默认的 Tracer
func MustNewTracer(ctx context.Context, opts ...TracerOption) *Tracer {
	t := &Tracer{
		Ctx:            ctx,
		Sampler:        1.0, // 默认采样率为 100%
		OtlpHeaders:    make(map[string]string),
		OtlpHttpPath:   "/v1/traces", // 默认路径
		OtlpHttpSecure: false,        // 默认使用不安全的HTTP
		Global:         true,         // 默认不加入全局
		ServiceName:    "unknown_service",
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (c *Tracer) NewTraceProvider(host string) (*sdktrace.TracerProvider, error) {
	if ok := hostgen.ValidListenHost(host); !ok {
		return nil, errors.New("无效的host")
	}
	var HttpOpt []otlptracehttp.Option = []otlptracehttp.Option{otlptracehttp.WithEndpoint(host)}
	if c.OtlpHttpSecure {
	} else {
		HttpOpt = append(HttpOpt, otlptracehttp.WithInsecure())
	}

	exp, err := otlptracehttp.New(c.Ctx, HttpOpt...)
	if err != nil {
		return nil, err
	}

	r, err := resource.New(c.Ctx,
		resource.WithAttributes(semconv.ServiceName(c.ServiceName)),
		resource.WithAttributes(c.ExtraResources...))
	if err != nil {
		return nil, err
	}

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(c.Sampler))),
		sdktrace.WithResource(r),
		sdktrace.WithBatcher(exp),
	}

	tp := sdktrace.NewTracerProvider(opts...)
	if c.Global {
		otel.SetTracerProvider(tp)
	}

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Errorf("[otel]tracing 内部出错, err= %v", err)
	}))

	return tp, nil
}

// 注册对应的traceProvider
func RegistorTP(ctx context.Context, host string, opts ...TracerOption) error {
	opts = append(opts, WithGlobal(true))
	tr := MustNewTracer(ctx, opts...)
	_, err := tr.NewTraceProvider(host)
	if err != nil {
		return nil
	}

	return nil
}

// 注意,Extract不会把span信息以context的形式抽取出,而是把span的spanContext信息抽取到ctx中
func Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return kPropagator.Extract(ctx, carrier)
}

// 注意,Inject不会把span信息以context的形式注入到ctx中,而是把span的spanContext信息注入到carrier中
func Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	kPropagator.Inject(ctx, carrier)
}

// 从ctx中获取记录的tracer和spanCtx以便直接生成新span
func ExtractSpanFormCtx(ctx context.Context, opts ...trace.TracerOption) (trace.Tracer, context.Context) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	var traceName string
	//同理先到md中找tracer-name,如果没有再从ctx中找
	if len(md.Get("tracer-name")) == 0 {
		traceName, ok = ctx.Value("tracer-name").(string)
		if !ok || traceName == "" {
			traceName = KTraceName
		}
		md.Set("tracer-name", traceName)
	} else {
		traceName = md.Get("tracer-name")[0]
	}

	//从ExtractMD中提取得到带有spanContext的ctx
	ctx = ExtractMD(ctx, md)
	sc := trace.SpanContextFromContext(ctx)
	//如果md中的spanContext无效,则从传入的ctx获取
	if !sc.IsValid() {
		return otel.Tracer(traceName, opts...), ctx
	}
	ctx = trace.ContextWithSpanContext(ctx, sc)
	return otel.Tracer(traceName, opts...), ctx
}

// 将ctx中的span信息注入到metadata中
func NewSpanOutgoingContext(ctx context.Context, span trace.Span) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	InjectMD(ctx, md)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}

func (c *Tracer) AddResources(attrs ...attribute.KeyValue) {
	c.ExtraResources = append(c.ExtraResources, attrs...)
}

// WithName 设置tracing的服务名称
func WithName(name string) TracerOption {
	return func(o *Tracer) {
		o.ServiceName = name
	}
}

// WithSampler 设置采样率，值为 [0.0, 1.0] 之间
func WithSampler(sampler float64) TracerOption {
	return func(o *Tracer) {
		if sampler < 0.0 {
			sampler = 0.0
		} else if sampler > 1.0 {
			sampler = 1.0
		}
		o.Sampler = sampler
	}
}

// WithOtlpHeaders 设置 OTLP 传输使用的 headers
func WithOtlpHeaders(headers map[string]string) TracerOption {
	return func(o *Tracer) {
		o.OtlpHeaders = headers
	}
}

// WithOtlpHttpPath 设置 OTLP HTTP 传输的路径
func WithOtlpHttpPath(path string) TracerOption {
	return func(o *Tracer) {
		o.OtlpHttpPath = path
	}
}

// WithOtlpHttpSecure 设置 OTLP HTTP 传输是否使用 HTTPS
func WithOtlpHttpSecure(secure bool) TracerOption {
	return func(o *Tracer) {
		o.OtlpHttpSecure = secure
	}
}

// WithGlobal 设置是否将生成的 TraceProvider 加入到全局中
func WithGlobal(global bool) TracerOption {
	return func(o *Tracer) {
		o.Global = global
	}
}

// WithGlobal 添加额外的字段信息
func WithExtraResource(kv []attribute.KeyValue) TracerOption {
	return func(o *Tracer) {
		o.ExtraResources = append(o.ExtraResources, kv...)
	}
}
