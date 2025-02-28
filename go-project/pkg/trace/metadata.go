package trace

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

var _ propagation.TextMapCarrier = (*metadataSupplier)(nil)

type metadataSupplier struct {
	metadata *metadata.MD
}

func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func (s *metadataSupplier) Set(key, value string) {
	s.metadata.Set(key, value)
}

func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(*s.metadata))
	for key := range *s.metadata {
		out = append(out, key)
	}

	return out
}

// spanContext注入到metadata中
func InjectMD(ctx context.Context, metadata *metadata.MD) {
	Inject(ctx, &metadataSupplier{metadata: metadata})
}

// 从metadata中提取spanContext并将其注入到context中
func ExtractMD(ctx context.Context, metadata *metadata.MD) context.Context {
	return Extract(ctx, &metadataSupplier{metadata: metadata})
}
