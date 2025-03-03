package cinterceptors

import (
	"context"
	"errors"
	"fmt"

	"io"

	ktrace "kenshop/pkg/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	gcodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	receiveEndEvent streamEventType = iota
	errorEvent
)

// UnaryTracingInterceptor returns a grpc.UnaryClientInterceptor for opentelemetry.
func UnaryTracingInterceptor(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	var span trace.Span
	ctx, span = getAndInjectMD(ctx, method, cc.Target())
	defer span.End()

	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {

		span.SetStatus(codes.Error, err.Error())
		//	span.SetAttributes(ktrace.StatusCodeAttr(s.Code()))
		return err
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// StreamTracingInterceptor returns a grpc.StreamClientInterceptor for opentelemetry.
func StreamTracingInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
	method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	var span trace.Span
	ctx, span = getAndInjectMD(ctx, method, cc.Target())
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			span.SetStatus(codes.Error, st.Message())
			span.SetAttributes(ktrace.StatusCodeAttr(st.Code()))
		} else {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
		return s, err
	}

	stream := wrapClientStream(ctx, s, desc)

	go func() {
		if err := <-stream.Finished; err != nil {
			s, ok := status.FromError(err)
			if ok {
				span.SetStatus(codes.Error, s.Message())
				span.SetAttributes(ktrace.StatusCodeAttr(s.Code()))
			} else {
				span.SetStatus(codes.Error, err.Error())
			}
		} else {
			span.SetAttributes(ktrace.StatusCodeAttr(gcodes.OK))
		}

		span.End()
	}()

	return stream, nil
}

type (
	streamEventType int

	streamEvent struct {
		Type streamEventType
		Err  error
	}

	clientStream struct {
		grpc.ClientStream
		Finished          chan error
		desc              *grpc.StreamDesc
		events            chan streamEvent
		eventsDone        chan struct{}
		receivedMessageID int
		sentMessageID     int
	}
)

func (w *clientStream) CloseSend() error {
	err := w.ClientStream.CloseSend()
	if err != nil {
		w.sendStreamEvent(errorEvent, err)
	}

	return err
}

func (w *clientStream) Header() (metadata.MD, error) {
	md, err := w.ClientStream.Header()
	if err != nil {
		w.sendStreamEvent(errorEvent, err)
	}

	return md, err
}

func (w *clientStream) RecvMsg(m any) error {
	err := w.ClientStream.RecvMsg(m)
	if err == nil && !w.desc.ServerStreams {
		w.sendStreamEvent(receiveEndEvent, nil)
	} else if errors.Is(err, io.EOF) {
		w.sendStreamEvent(receiveEndEvent, nil)
	} else if err != nil {
		w.sendStreamEvent(errorEvent, err)
	} else {
		w.receivedMessageID++
		ktrace.MessageReceived.Event(w.Context(), w.receivedMessageID, m)
	}

	return err
}

func (w *clientStream) SendMsg(m any) error {
	err := w.ClientStream.SendMsg(m)
	w.sentMessageID++
	ktrace.MessageSent.Event(w.Context(), w.sentMessageID, m)
	if err != nil {
		w.sendStreamEvent(errorEvent, err)
	}

	return err
}

func (w *clientStream) sendStreamEvent(eventType streamEventType, err error) {
	select {
	case <-w.eventsDone:
	case w.events <- streamEvent{Type: eventType, Err: err}:
	}
}

// 从ctx的metadata中提取传入的span,如果没有则从ctx中提取span,并将得到的span以metadata的形式注入到ctx中
func getAndInjectMD(ctx context.Context, method, target string) (context.Context, trace.Span) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	var traceName string
	//如果md中没有traceName则到传入的context里找
	if len(md.Get("tracer-name")) == 0 {
		traceName, ok = ctx.Value("tracer-name").(string)
		//如果context里也没有traceName则使用默认的traceName
		if !ok || traceName == "" {
			traceName = ktrace.TraceName
		}
		md.Set("tracer-name", traceName)
	} else {
		traceName = md.Get("tracer-name")[0]
	}

	var span trace.Span
	spanCtx := ktrace.ExtractMD(ctx, &md)
	sc := trace.SpanContextFromContext(spanCtx)
	//如果md里没有有效的spanCtx信息则从传入的ctx中找
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
	name, attr := ktrace.SpanInfo(method, target)
	cSpanCtx, cspan := tr.Start(spanCtx, fmt.Sprintf("client-%s", name),
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attr...),
	)

	ktrace.InjectMD(cSpanCtx, &md)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, cspan
}

// wrapClientStream wraps s with given ctx and desc.
func wrapClientStream(ctx context.Context, s grpc.ClientStream, desc *grpc.StreamDesc) *clientStream {
	events := make(chan streamEvent)
	eventsDone := make(chan struct{})
	finished := make(chan error)

	go func() {
		defer close(eventsDone)

		for {
			select {
			case event := <-events:
				switch event.Type {
				case receiveEndEvent:
					finished <- nil
					return
				case errorEvent:
					finished <- event.Err
					return
				}
			case <-ctx.Done():
				finished <- ctx.Err()
				return
			}
		}
	}()

	return &clientStream{
		ClientStream: s,
		desc:         desc,
		events:       events,
		eventsDone:   eventsDone,
		Finished:     finished,
	}
}
