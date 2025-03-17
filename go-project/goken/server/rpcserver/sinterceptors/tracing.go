package sinterceptors

import (
	"context"
	"fmt"

	ktrace "kenshop/pkg/trace"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	gcodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryTracingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	var span trace.Span
	tr, spanCtx := ktrace.ExtractSpanFormCtx(ctx)
	name, attr := ktrace.ResolveGrpcInfo(info.FullMethod, ktrace.PeerAddrFromCtx(ctx))

	ctx, span = tr.Start(spanCtx, fmt.Sprintf("server-%s", name), trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attr...))

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

	tr, spanCtx := ktrace.ExtractSpanFormCtx(ss.Context())
	name, attr := ktrace.ResolveGrpcInfo(info.FullMethod, ktrace.PeerAddrFromCtx(ss.Context()))
	ctx, span := tr.Start(spanCtx, fmt.Sprintf("server-%s", name), trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attr...))

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

// wrapServerStream wraps the given grpc.ServerStream with the given context.
func wrapServerStream(ctx context.Context, ss grpc.ServerStream) *serverStream {
	return &serverStream{
		ServerStream: ss,
		ctx:          ctx,
	}
}
