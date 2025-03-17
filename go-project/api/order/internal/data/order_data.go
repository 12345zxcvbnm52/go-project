package orderdata

import (
	"context"
	proto "kenshop/proto/order"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// OrderDataService是提供Order底层相关数据操作的接口
type OrderDataService interface {
	//获得用户购物车信息
	GetUserCartItemsDB(context.Context, *proto.UserInfoReq) (*proto.CartItemListRes, error)
	//为购物车添加商品
	CreateCartItemDB(context.Context, *proto.CreateCartItemReq) (*proto.CartItemInfoRes, error)
	//修改购物车的一条记录
	UpdateCartItemDB(context.Context, *proto.UpdateCartItemReq) (*emptypb.Empty, error)

	DeleteCartItemDB(context.Context, *proto.DelCartItemReq) (*emptypb.Empty, error)

	CreateOrderDB(context.Context, *proto.CreateOrderReq) (*proto.OrderInfoRes, error)

	GetOrderListDB(context.Context, *proto.OrderFliterReq) (*proto.OrderListRes, error)

	GetOrderInfoDB(context.Context, *proto.OrderInfoReq) (*proto.OrderDetailRes, error)

	UpdateOrderStatusDB(context.Context, *proto.OrderStatusReq) (*emptypb.Empty, error)
}

func MustNewGrpcOrderData(c *grpc.ClientConn) OrderDataService {
	return &GrpcOrderData{Cli: proto.NewOrderClient(c)}
}

var _ OrderDataService = (*GrpcOrderData)(nil)

// Order服务中的Data层,是数据操作的具体逻辑
type GrpcOrderData struct {
	Cli proto.OrderClient
}

// 获得用户购物车信息
func (d *GrpcOrderData) GetUserCartItemsDB(ctx context.Context, in *proto.UserInfoReq) (*proto.CartItemListRes, error) {
	return d.Cli.GetUserCartItems(ctx, in)
}

// 为购物车添加商品
func (d *GrpcOrderData) CreateCartItemDB(ctx context.Context, in *proto.CreateCartItemReq) (*proto.CartItemInfoRes, error) {
	return d.Cli.CreateCartItem(ctx, in)
}

// 修改购物车的一条记录
func (d *GrpcOrderData) UpdateCartItemDB(ctx context.Context, in *proto.UpdateCartItemReq) (*emptypb.Empty, error) {
	return d.Cli.UpdateCartItem(ctx, in)
}

func (d *GrpcOrderData) DeleteCartItemDB(ctx context.Context, in *proto.DelCartItemReq) (*emptypb.Empty, error) {
	return d.Cli.DeleteCartItem(ctx, in)
}

func (d *GrpcOrderData) CreateOrderDB(ctx context.Context, in *proto.CreateOrderReq) (*proto.OrderInfoRes, error) {
	return d.Cli.CreateOrder(ctx, in)
}

func (d *GrpcOrderData) GetOrderListDB(ctx context.Context, in *proto.OrderFliterReq) (*proto.OrderListRes, error) {
	return d.Cli.GetOrderList(ctx, in)
}

func (d *GrpcOrderData) GetOrderInfoDB(ctx context.Context, in *proto.OrderInfoReq) (*proto.OrderDetailRes, error) {
	return d.Cli.GetOrderInfo(ctx, in)
}

func (d *GrpcOrderData) UpdateOrderStatusDB(ctx context.Context, in *proto.OrderStatusReq) (*emptypb.Empty, error) {
	return d.Cli.UpdateOrderStatus(ctx, in)
}
