package main

import (
	"context"
	"kenshop/goken/registry/registor"
	"kenshop/goken/registry/ways/consul"
	"kenshop/goken/server/rpcserver"
	pb "kenshop/test_srv/proto"
	"net"

	"github.com/hashicorp/consul/api"
)

type ConnServer struct {
	pb.UnimplementedTestServer
}

func (c *ConnServer) Conn(ctx context.Context, in *pb.Req) (*pb.Res, error) {
	return &pb.Res{Res: in.Name + " hello ken"}, nil
}

func main() {
	lis, _ := net.Listen("tcp", "192.168.199.128:22222")

	cfg := api.DefaultConfig()

	cfg.Address = "192.168.199.128:8500"
	cli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	r := registor.MustNewRegister(
		consul.MustNewConsulRegistor(cli,
			consul.WithEnableHealthCheck(true),
			consul.WithDeregisterCriticalServiceAfter("30s"),
			consul.WithHealthcheckInterval("10s"),
		),
	)

	s := rpcserver.MustNewServer(context.Background(), rpcserver.WithHost("127.0.0.1:22222"), rpcserver.WithListener(lis), rpcserver.WithRegistor(r))
	cs := &ConnServer{}
	pb.RegisterTestServer(s.Server, cs)

	if err := s.Serve(); err != nil {
		panic(err)
	}
}

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"math/rand"

// 	"github.com/gin-gonic/gin"
// 	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/attribute"
// 	jaeger "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
// 	"go.opentelemetry.io/otel/propagation"
// 	"go.opentelemetry.io/otel/sdk/resource"
// 	"go.opentelemetry.io/otel/sdk/trace"
// 	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
// 	traceo "go.opentelemetry.io/otel/trace"
// )

// func initTracer() {
// 	otel.SetTextMapPropagator(
// 		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
// 	)
// }

// func Server() {
// 	r := gin.Default()
// 	r.Use(otelgin.Middleware("ken"))
// 	r.Run()
// }

// func main() {
// 	exporter, err := jaeger.New(context.Background(),
// 		jaeger.WithEndpoint("127.0.0.1:4318"),
// 		jaeger.WithInsecure(),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	res, err := resource.New(context.Background(), resource.WithAttributes(
// 		semconv.ServiceName("goken"),
// 	))

// 	res = resource.NewWithAttributes(
// 		semconv.SchemaURL,
// 		semconv.ServiceNameKey.String("goken"),
// 		attribute.String("name", "gin"),
// 		attribute.Bool("env", true),
// 	)

// 	if err != nil {
// 		panic(err)
// 	}

// 	tp := trace.NewTracerProvider(
// 		trace.WithResource(res),
// 		trace.WithSampler(trace.AlwaysSample()),
// 		trace.WithBatcher(exporter, trace.WithBatchTimeout(time.Second)),
// 	)

// 	otel.SetTracerProvider(tp)
// 	st := tp.Shutdown

// 	tracer := otel.Tracer("test-tracer")
// 	baseAttrs := []attribute.KeyValue{
// 		attribute.Int("num", 2),
// 	}

// 	c, span := tracer.Start(context.Background(), "parent-span", traceo.WithAttributes(baseAttrs...))
// 	defer span.End()
// 	// 使用for循环创建多个子span，方便查看效果
// 	for i := range 10 { // Go1.22+
// 		// 传入父ctx，开启子span
// 		_, iSpan := tracer.Start(c, fmt.Sprintf("span-%d", i))
// 		// 随机sleep，模拟子span中耗时的操作
// 		time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
// 		// 子span结束
// 		iSpan.End()
// 	}
// 	fmt.Println("done!")
// 	st(c)
// }
