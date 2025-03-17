package goodscontroller

import (
	proto "kenshop/proto/goods"
	goodslogic "kenshop/service/goods/internal/logic"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	protojson "google.golang.org/protobuf/encoding/protojson"
	gproto "google.golang.org/protobuf/proto"
)

// Goods服务中的Contoller层,用于对外暴露grpc接口
type GoodsServer struct {
	Service *goodslogic.GoodsService
	Logger  *otelzap.Logger
	proto.UnimplementedGoodsServer
}

var ProtoJson = protojson.MarshalOptions{
	EmitUnpopulated: true,
}

func MethodInfoRecord(data gproto.Message) string {
	r, err := ProtoJson.Marshal(data)
	if err != nil {
		return ""
	}
	return string(r)
}
