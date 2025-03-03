package sinterceptors

import (
	"context"
	"fmt"

	ktrace "kenshop/pkg/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	gcodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryTracingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	var span trace.Span
	ctx, span = startSpan(ctx, info.FullMethod)
	defer span.End()
	resp, err := handler(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return resp, nil
}

// StreamTracingInterceptor returns a grpc.StreamServerInterceptor for opentelemetry.
func StreamTracingInterceptor(svr any, ss grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	ctx, span := startSpan(ss.Context(), info.FullMethod)
	defer span.End()

	if err := handler(svr, wrapServerStream(ctx, ss)); err != nil {
		s, ok := status.FromError(err)
		if ok {
			span.SetStatus(codes.Error, s.Message())
			span.SetAttributes(ktrace.StatusCodeAttr(s.Code()))
		} else {
			span.SetStatus(codes.Error, err.Error())
		}
		return err
	}

	span.SetAttributes(ktrace.StatusCodeAttr(gcodes.OK))
	return nil
}

// serverStream wraps around the embedded grpc.ServerStream,
// and intercepts the RecvMsg and SendMsg method call.
type serverStream struct {
	grpc.ServerStream
	ctx               context.Context
	receivedMessageID int
	sentMessageID     int
}

func (w *serverStream) Context() context.Context {
	return w.ctx
}

func (w *serverStream) RecvMsg(m any) error {
	err := w.ServerStream.RecvMsg(m)
	if err == nil {
		w.receivedMessageID++
		ktrace.MessageReceived.Event(w.Context(), w.receivedMessageID, m)
	}

	return err
}

func (w *serverStream) SendMsg(m any) error {
	err := w.ServerStream.SendMsg(m)
	w.sentMessageID++
	ktrace.MessageSent.Event(w.Context(), w.sentMessageID, m)

	return err
}

func startSpan(ctx context.Context, method string) (context.Context, trace.Span) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	var traceName string
	//同理先到md中找tracer-name,如果没有再从ctx中找
	if len(md.Get("tracer-name")) == 0 {
		traceName, ok = ctx.Value("tracer-name").(string)
		if !ok || traceName == "" {
			traceName = ktrace.TraceName
		}
		md.Set("tracer-name", traceName)
	} else {
		traceName = md.Get("tracer-name")[0]
	}

	var span trace.Span
	//从ExtractMD中提取得到带有spanContext的ctx
	spanCtx := ktrace.ExtractMD(ctx, &md)
	sc := trace.SpanContextFromContext(spanCtx)
	//如果md中的spanContext无效,则从传入的ctx获取
	if !sc.IsValid() {
		span = trace.SpanFromContext(ctx)
		sc = span.SpanContext()
		//如果传入ctx中仍然没有span信息则用传入的ctx自定义一个span
		if !sc.IsValid() {
			spanCtx = ctx
		} else {
			//把span中spanContext的信息注入到新的ctx中
			spanCtx = trace.ContextWithSpanContext(spanCtx, sc)
		}
	}

	tp := otel.GetTracerProvider()
	tr := tp.Tracer(traceName)
	name, attr := ktrace.SpanInfo(method, ktrace.PeerFromCtx(ctx))

	return tr.Start(spanCtx, fmt.Sprintf("server-%s", name), trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attr...))
}

// wrapServerStream wraps the given grpc.ServerStream with the given context.
func wrapServerStream(ctx context.Context, ss grpc.ServerStream) *serverStream {
	return &serverStream{
		ServerStream: ss,
		ctx:          ctx,
	}
}
