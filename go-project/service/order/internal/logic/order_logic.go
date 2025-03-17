package orderlogic

import (
	"context"
	"kenshop/pkg/errors"
	proto "kenshop/proto/order"
	orderdata "kenshop/service/order/internal/data"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Order服务中的Service层,编写具体的服务逻辑
type OrderService struct {
	OrderData orderdata.OrderDataService
}

// 获得用户购物车信息
func (s *OrderService) GetUserCartItemsLogic(ctx context.Context, in *proto.UserInfoReq) (*proto.CartItemListRes, error) {
	return s.OrderData.GetUserCartItemsDB(ctx, in)
}

// 为购物车添加商品
func (s *OrderService) CreateCartItemLogic(ctx context.Context, in *proto.CreateCartItemReq) (*proto.CartItemInfoRes, error) {
	return s.OrderData.CreateCartItemDB(ctx, in)
}

// 修改购物车的一条记录
func (s *OrderService) UpdateCartItemLogic(ctx context.Context, in *proto.UpdateCartItemReq) (*emptypb.Empty, error) {
	return s.OrderData.UpdateCartItemDB(ctx, in)
}

func (s *OrderService) DeleteCartItemLogic(ctx context.Context, in *proto.DelCartItemReq) (*emptypb.Empty, error) {
	return s.OrderData.DeleteCartItemDB(ctx, in)
}

func (s *OrderService) CreateOrderLogic(ctx context.Context, in *proto.CreateOrderReq) (*proto.OrderInfoRes, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, errors.WithCoder(err, errors.CodeInternalError, "")
	}
	in.OrderSign = strings.ReplaceAll(uuid.String(), "-", "")
	return s.OrderData.CreateOrderDB(ctx, in)
}

func (s *OrderService) GetOrderListLogic(ctx context.Context, in *proto.OrderFliterReq) (*proto.OrderListRes, error) {
	return s.OrderData.GetOrderListDB(ctx, in)
}

func (s *OrderService) GetOrderInfoLogic(ctx context.Context, in *proto.OrderInfoReq) (*proto.OrderDetailRes, error) {
	return s.OrderData.GetOrderInfoDB(ctx, in)
}

func (s *OrderService) UpdateOrderStatusLogic(ctx context.Context, in *proto.OrderStatusReq) (*emptypb.Empty, error) {
	return s.OrderData.UpdateOrderStatusDB(ctx, in)
}
