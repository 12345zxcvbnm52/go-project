package trace

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	gcodes "google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

const (
	// GRPCStatusCodeKey 是用于表示gRPC请求数字状态码的约定,
	GRPCStatusCodeKey = attribute.Key("rpc.grpc.status_code")
	// RPCNameKey 是传输或接收的消息的名称,
	RPCNameKey = attribute.Key("name")
	// RPCMessageTypeKey 是传输或接收的消息的类型,
	RPCMessageTypeKey = attribute.Key("message.type")
	// RPCMessageIDKey 是传输或接收的消息的标识符,
	RPCMessageIDKey = attribute.Key("message.id")
	// RPCMessageCompressedSizeKey 是传输或接收的消息的压缩大小,以字节为单位,
	RPCMessageCompressedSizeKey = attribute.Key("message.compressed_size")
	// RPCMessageUncompressedSizeKey 是传输或接收的消息的未压缩大小,以字节为单位,
	RPCMessageUncompressedSizeKey = attribute.Key("message.uncompressed_size")
)

// 常见RPC属性的otel自建的语义约定
var (
	// RPCSystemGRPC 是将 gRPC 作为远程系统的语义约定,
	RPCSystemGRPC = semconv.RPCSystemKey.String("grpc")
	// RPCNameMessage 是名为message的消息的语义约定,
	RPCNameMessage = RPCNameKey.String("message")
	// RPCMessageTypeSent 是发送的RPC消息类型的语义约定,
	RPCMessageTypeSent = RPCMessageTypeKey.String("SENT")
	// RPCMessageTypeReceived 是接收的RPC消息类型的语义约定,
	RPCMessageTypeReceived = RPCMessageTypeKey.String("RECEIVED")
)

// StatusCodeAttr 返回一个表示给定gRPC状态码的attribute.KeyValue,
func StatusCodeAttr(c gcodes.Code) attribute.KeyValue {
	return GRPCStatusCodeKey.Int64(int64(c))
}

// 定义了事件的类型
const messageEvent = "message"

var (
	// MessageSent 表示已发送的消息类型,
	MessageSent = messageType(RPCMessageTypeSent)
	// MessageReceived 表示已接收的消息类型,
	MessageReceived = messageType(RPCMessageTypeReceived)
)

// messageType 是基于 attribute.KeyValue 的一个类型别名,用来表示消息类型,
type messageType attribute.KeyValue

// Event方法将一个消息类型事件动态添加到与传入上下文关联的span上,
// 添加的信息包含消息的id,如果消息是proto.Message类型,还会添加消息的大小,
// AddEvent Event记录的是时间点的快照,表示在span的执行过程中某个特定时刻发生的事情,每个事件都有一个时间戳,表示它发生的准确时刻,
func (m messageType) Event(ctx context.Context, id int, message any) {
	span := trace.SpanFromContext(ctx)
	if p, ok := message.(proto.Message); ok {
		span.AddEvent(messageEvent, trace.WithAttributes(
			attribute.KeyValue(m),
			RPCMessageIDKey.Int(id),
			RPCMessageUncompressedSizeKey.Int(proto.Size(p)),
		))
	} else {
		span.AddEvent(messageEvent, trace.WithAttributes(
			attribute.KeyValue(m),
			RPCMessageIDKey.Int(id),
		))
	}
}
