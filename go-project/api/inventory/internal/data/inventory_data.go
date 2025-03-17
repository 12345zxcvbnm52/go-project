package inventorydata

import (
	"context"
	proto "kenshop/proto/inventory"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// InventoryDataService是提供Inventory底层相关数据操作的接口
type InventoryDataService interface {
	//后续可以考虑给每个库存添加id(同一个id则增加库存量,否则根据地址创建),address(地址)
	CreateStockDB(context.Context, *proto.CreateInventoryReq) (*emptypb.Empty, error)

	SetStockDB(context.Context, *proto.SetInventoryReq) (*emptypb.Empty, error)

	GetStockInfoDB(context.Context, *proto.InventoryInfoReq) (*proto.InventoryInfoRes, error)

	DecrStockDB(context.Context, *proto.UpdateStockReq) (*emptypb.Empty, error)

	IncrStockDB(context.Context, *proto.UpdateStockReq) (*emptypb.Empty, error)
}

func MustNewGrpcInventoryData(c *grpc.ClientConn) InventoryDataService {
	return &GrpcInventoryData{Cli: proto.NewInventoryClient(c)}
}

var _ InventoryDataService = (*GrpcInventoryData)(nil)

// Inventory服务中的Data层,是数据操作的具体逻辑
type GrpcInventoryData struct {
	Cli proto.InventoryClient
}

// 后续可以考虑给每个库存添加id(同一个id则增加库存量,否则根据地址创建),address(地址)
func (d *GrpcInventoryData) CreateStockDB(ctx context.Context, in *proto.CreateInventoryReq) (*emptypb.Empty, error) {
	return d.Cli.CreateStock(ctx, in)
}

func (d *GrpcInventoryData) SetStockDB(ctx context.Context, in *proto.SetInventoryReq) (*emptypb.Empty, error) {
	return d.Cli.SetStock(ctx, in)
}

func (d *GrpcInventoryData) GetStockInfoDB(ctx context.Context, in *proto.InventoryInfoReq) (*proto.InventoryInfoRes, error) {
	return d.Cli.GetStockInfo(ctx, in)
}

func (d *GrpcInventoryData) DecrStockDB(ctx context.Context, in *proto.UpdateStockReq) (*emptypb.Empty, error) {
	return d.Cli.DecrStock(ctx, in)
}

func (d *GrpcInventoryData) IncrStockDB(ctx context.Context, in *proto.UpdateStockReq) (*emptypb.Empty, error) {
	return d.Cli.IncrStock(ctx, in)
}
