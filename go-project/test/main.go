package main

import (
	"context"
	ktrace "kenshop/pkg/trace"
	proto "kenshop/proto/test"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func testA(ctx context.Context, tracer trace.Tracer) {
	//span:=trace.SpanFromContext(ctx)
	_, cspan := tracer.Start(ctx, "chird pp")
	ktrace.MessageSent.Event(ctx, 222, &proto.ReqMessage{Req: "222"})
	defer cspan.End()
	time.Sleep(2 * time.Second)
}

func Span() {

	tp := ktrace.MustNewTracer(context.TODO(), ktrace.WithName("ken"))
	t, err := tp.NewTraceProvider("192.168.199.128:4318")
	if err != nil {
		panic(err)
	}
	tracer := t.Tracer("ken-tracer", trace.WithInstrumentationAttributes(attribute.Bool("isken", false)))
	ctx, span := tracer.Start(context.Background(), "pp", trace.WithAttributes(attribute.String("span key", "test")))
	testA(ctx, tracer)

	span.End()
	t.Shutdown(context.Background())
}

func main() {
	Span()
}
