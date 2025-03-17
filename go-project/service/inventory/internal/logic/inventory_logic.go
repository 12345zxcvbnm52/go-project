package inventorylogic

import (
	"context"
	proto "kenshop/proto/inventory"
	inventorydata "kenshop/service/inventory/internal/data"

	"google.golang.org/protobuf/types/known/emptypb"
)

// Inventory服务中的Service层,编写具体的服务逻辑
type InventoryService struct {
	InventoryData inventorydata.InventoryDataService
}

// 后续可以考虑给每个库存添加id(同一个id则增加库存量,否则根据地址创建),address(地址)
func (s *InventoryService) CreateStockLogic(ctx context.Context, in *proto.CreateInventoryReq) (*emptypb.Empty, error) {
	return s.InventoryData.CreateStockDB(ctx, in)
}

func (s *InventoryService) SetStockLogic(ctx context.Context, in *proto.SetInventoryReq) (*emptypb.Empty, error) {
	return s.InventoryData.SetStockDB(ctx, in)
}

func (s *InventoryService) GetStockInfoLogic(ctx context.Context, in *proto.InventoryInfoReq) (*proto.InventoryInfoRes, error) {
	return s.InventoryData.GetStockInfoDB(ctx, in)
}

func (s *InventoryService) DecrStockLogic(ctx context.Context, in *proto.UpdateStockReq) (*emptypb.Empty, error) {
	return s.InventoryData.DecrStockDB(ctx, in)
}

func (s *InventoryService) IncrStockLogic(ctx context.Context, in *proto.UpdateStockReq) (*emptypb.Empty, error) {
	return s.InventoryData.IncrStockDB(ctx, in)
}

func (s *InventoryService) RebackStockLogic(ctx context.Context, in *proto.RebackStockReq) (*emptypb.Empty, error) {
	return s.InventoryData.RebackStockDB(ctx, in)
}
