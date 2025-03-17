package inventorycontroller

import (
	"context"
	proto "kenshop/proto/inventory"
	inventorylogic "kenshop/service/inventory/internal/logic"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	protojson "google.golang.org/protobuf/encoding/protojson"
	gproto "google.golang.org/protobuf/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// Inventory服务中的Contoller层,用于对外暴露grpc接口
type InventoryServer struct {
	Service *inventorylogic.InventoryService
	Logger  *otelzap.Logger
	proto.UnimplementedInventoryServer
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

// 后续可以考虑给每个库存添加id(同一个id则增加库存量,否则根据地址创建),address(地址)
func (s *InventoryServer) CreateStock(ctx context.Context, in *proto.CreateInventoryReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次CreateStock调用,调用信息为: %s", info)
	res, err := s.Service.CreateStockLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用CreateStock失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *InventoryServer) SetStock(ctx context.Context, in *proto.SetInventoryReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次SetStock调用,调用信息为: %s", info)
	res, err := s.Service.SetStockLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用SetStock失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *InventoryServer) GetStockInfo(ctx context.Context, in *proto.InventoryInfoReq) (*proto.InventoryInfoRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetStockInfo调用,调用信息为: %s", info)
	res, err := s.Service.GetStockInfoLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetStockInfo失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *InventoryServer) DecrStock(ctx context.Context, in *proto.UpdateStockReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次DecrStock调用,调用信息为: %s", info)
	res, err := s.Service.DecrStockLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用DecrStock失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *InventoryServer) IncrStock(ctx context.Context, in *proto.UpdateStockReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次IncrStock调用,调用信息为: %s", info)
	res, err := s.Service.IncrStockLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用IncrStock失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *InventoryServer) RebackStock(ctx context.Context, in *proto.RebackStockReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次RebackStock调用,调用信息为: %s", info)
	res, err := s.Service.RebackStockLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用RebackStock失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}
