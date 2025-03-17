package ordercontroller

import (
	"context"
	proto "kenshop/proto/order"

	orderlogic "kenshop/service/order/internal/logic"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	protojson "google.golang.org/protobuf/encoding/protojson"
	gproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Order服务中的Contoller层,用于对外暴露grpc接口
type OrderServer struct {
	Service *orderlogic.OrderService
	Logger  *otelzap.Logger
	proto.UnimplementedOrderServer
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

// 获得用户购物车信息
func (s *OrderServer) GetUserCartItems(ctx context.Context, in *proto.UserInfoReq) (*proto.CartItemListRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetUserCartItems调用,调用信息为: %s", info)
	res, err := s.Service.GetUserCartItemsLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetUserCartItems失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

// 为购物车添加商品
func (s *OrderServer) CreateCartItem(ctx context.Context, in *proto.CreateCartItemReq) (*proto.CartItemInfoRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次CreateCartItem调用,调用信息为: %s", info)
	res, err := s.Service.CreateCartItemLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用CreateCartItem失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

// 修改购物车的一条记录
func (s *OrderServer) UpdateCartItem(ctx context.Context, in *proto.UpdateCartItemReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次UpdateCartItem调用,调用信息为: %s", info)
	res, err := s.Service.UpdateCartItemLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用UpdateCartItem失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *OrderServer) DeleteCartItem(ctx context.Context, in *proto.DelCartItemReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次DeleteCartItem调用,调用信息为: %s", info)
	res, err := s.Service.DeleteCartItemLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用DeleteCartItem失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *OrderServer) CreateOrder(ctx context.Context, in *proto.CreateOrderReq) (*proto.OrderInfoRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次CreateOrder调用,调用信息为: %s", info)
	res, err := s.Service.CreateOrderLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用CreateOrder失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *OrderServer) GetOrderList(ctx context.Context, in *proto.OrderFliterReq) (*proto.OrderListRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetOrderList调用,调用信息为: %s", info)
	res, err := s.Service.GetOrderListLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetOrderList失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *OrderServer) GetOrderInfo(ctx context.Context, in *proto.OrderInfoReq) (*proto.OrderDetailRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetOrderInfo调用,调用信息为: %s", info)
	res, err := s.Service.GetOrderInfoLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetOrderInfo失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, in *proto.OrderStatusReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次UpdateOrderStatus调用,调用信息为: %s", info)
	res, err := s.Service.UpdateOrderStatusLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用UpdateOrderStatus失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}
