package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
)

func NewJaegerTraceProvider(ctx context.Context, endpoint string) (*traceSDK.TracerProvider, error) {
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	tp := traceSDK.NewTracerProvider(
		traceSDK.WithBatcher(exp, traceSDK.WithBatchTimeout(time.Second)),
		traceSDK.WithSampler(traceSDK.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}
