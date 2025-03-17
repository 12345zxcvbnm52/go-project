package ordercontroller

import (
	"fmt"
	orderdata "kenshop/api/order/internal/data"
	orderform "kenshop/api/order/internal/form"
	httpserver "kenshop/goken/server/httpserver"
	"kenshop/pkg/common/httputil"
	log "kenshop/pkg/log"
	proto "kenshop/proto/order"
	"net/http"

	gin "github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type OrderHttpServer struct {
	Server    *httpserver.Server
	OrderData orderdata.OrderDataService
	Logger    *otelzap.Logger
}

// 默认使用otelzap.Logger以及GrpcOrderData
func MustNewOrderHTTPServer(s *httpserver.Server, opts ...OptionFunc) *OrderHttpServer {
	ss := &OrderHttpServer{
		Server: s,
	}
	for _, opt := range opts {
		opt(ss)
	}
	if ss.Logger == nil {
		ss.Logger = log.MustNewOtelLogger()
	}
	if ss.OrderData == nil {
		cli, err := s.GrpcCli.Dial()
		if err != nil {
			panic(err)
		}
		ss.OrderData = orderdata.MustNewGrpcOrderData(cli)
	}
	return ss
}

// 获得用户购物车信息
// @Description 获取用户购物车信息
// @Produce application/json
// @Router /cart/{user_id} [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.CartItemListRes}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param userId path integer true "(uint32)购物车所属的用户的用户ID"
// @Param pagesNum query integer flase "(int32)"显示购物车内商品的页号
// @Param pageSize query integer false "(int32)显示购物车内商品一页的数量"
func (s *OrderHttpServer) GetUserCartItems(c *gin.Context) {
	u := &orderform.UserInfoForm{}

	if err := c.ShouldBindQuery(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}
	fmt.Println(u.UserId)
	res, err := s.OrderData.GetUserCartItemsDB(s.Server.Ctx, &proto.UserInfoReq{
		PageSize: u.PageSize,
		PagesNum: u.PagesNum,
		UserId:   u.UserId,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// 为购物车添加商品
// @Accept application/json
// @Description 为购物车添加商品
// @Produce application/json
// @Router /carts [POST]
// @Success 200 {object} httputil.JsonResult{data=proto.CartItemInfoRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Param userId body integer true "(uint32)购物车所属的用户的用户ID"
// @Param goodsId body integer true "(uint32)欲添加到购物车的商品ID"
// @Param goodsNum body integer true "(uint32)欲添加到购物车的商品数量"
func (s *OrderHttpServer) CreateCartItem(c *gin.Context) {
	u := &orderform.CreateCartItemForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.OrderData.CreateCartItemDB(s.Server.Ctx, &proto.CreateCartItemReq{
		GoodsId:  u.GoodsId,
		GoodsNum: u.GoodsNum,
		UserId:   u.UserId,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// 修改购物车的一条记录
// @Accept application/json
// @Description 修改购物车的一条记录
// @Produce application/json
// @Router /cart/{id} [PUT]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)购物车记录的ID"
// @Param userId body integer false "(uint32)购物车所属的用户的用户ID"
// @Param goodsId body integer false "(uint32)欲修改的购物车商品的ID"
// @Param goodsNum body integer false "(int32)修改后的商品数量"
// @Param selected body bool false "(uint32)购物车所属的用户的用户ID"
func (s *OrderHttpServer) UpdateCartItem(c *gin.Context) {
	u := &orderform.UpdateCartItemForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}
	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.OrderData.UpdateCartItemDB(s.Server.Ctx, &proto.UpdateCartItemReq{
		GoodsId:  u.GoodsId,
		GoodsNum: u.GoodsNum,
		Id:       u.Id,
		Selected: u.Selected,
		UserId:   u.UserId,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 删除购物车的一条记录
// @Produce application/json
// @Router /cart/{id} [DELETE]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)欲删除的购物车记录ID"
// @Param userId body integer false "(uint32)购物车所属的用户的用户ID"
// @Param goodsId body integer false "(uint32)欲删除的购物车商品的ID"
func (s *OrderHttpServer) DeleteCartItem(c *gin.Context) {
	u := &orderform.DelCartItemForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}
	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.OrderData.DeleteCartItemDB(s.Server.Ctx, &proto.DelCartItemReq{
		GoodsId: u.GoodsId,
		Id:      u.Id,
		UserId:  u.UserId,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 创建订单
// @Produce application/json
// @Router /orders [POST]
// @Success 200 {object} httputil.JsonResult{data=proto.OrderInfoRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Param userId body integer true "(uint32)订单所属的用户的用户ID"
// @Param address body string true "订单的收货地址"
// @Param signerName body string true "收货人"
// @Param signerMobile body string true "收货人电话"
// @Param message body string true "订单额外信息"
// @Param payWay body string true "支付方式"
func (s *OrderHttpServer) CreateOrder(c *gin.Context) {
	u := &orderform.CreateOrderForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.OrderData.CreateOrderDB(s.Server.Ctx, &proto.CreateOrderReq{
		Address:      u.Address,
		Message:      u.Message,
		PayWay:       u.PayWay,
		SignerMobile: u.SignerMobile,
		SignerName:   u.SignerName,
		UserId:       u.UserId,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Description 获取订单列表
// @Produce application/json
// @Router /orders/{user_id} [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.OrderListRes}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Param userId path integer true "(uint32)订单所属的用户的用户ID"
// @Param pagesNum query integer flase "(int32)"显示购物车内商品的页号
// @Param pageSize query integer false "(int32)显示购物车内商品一页的数量"
func (s *OrderHttpServer) GetOrderList(c *gin.Context) {
	u := &orderform.OrderFliterForm{}

	if err := c.ShouldBindQuery(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.OrderData.GetOrderListDB(s.Server.Ctx, &proto.OrderFliterReq{
		PageSize: u.PageSize,
		PagesNum: u.PagesNum,
		UserId:   u.UserId,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Description 获取订单详情
// @Produce application/json
// @Router /order/{id} [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.OrderDetailRes}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)订单ID"
// @Param userId query integer true "(uint32)订单所属的用户的用户ID"
func (s *OrderHttpServer) GetOrderInfo(c *gin.Context) {
	u := &orderform.OrderInfoForm{}

	if err := c.ShouldBindQuery(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.OrderData.GetOrderInfoDB(s.Server.Ctx, &proto.OrderInfoReq{
		Id:     u.Id,
		UserId: u.UserId,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 更新订单状态
// @Produce application/json
// @Router /order/status/{id} [PUT]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)订单ID"
// @Param status body integer true "(int32)订单的新状态"
// @Param userId body integer true "(uint32)订单所属的用户的用户ID"
// @Param orderSign body string false "(int32)订单号"
func (s *OrderHttpServer) UpdateOrderStatus(c *gin.Context) {
	u := &orderform.OrderStatusForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}
	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.OrderData.UpdateOrderStatusDB(s.Server.Ctx, &proto.OrderStatusReq{
		Id:        u.Id,
		OrderSign: u.OrderSign,
		Status:    u.Status,
		UserId:    u.UserId,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

func (s *OrderHttpServer) Execute() error {
	s.Server.Engine.GET("/cart/:user_id", s.GetUserCartItems)
	s.Server.Engine.POST("/carts", s.CreateCartItem)
	s.Server.Engine.PUT("/cart/:id", s.UpdateCartItem)
	s.Server.Engine.DELETE("/cart/:id", s.DeleteCartItem)
	s.Server.Engine.POST("/orders", s.CreateOrder)
	s.Server.Engine.GET("/orders/:user_id", s.GetOrderList)
	s.Server.Engine.GET("/order/:id", s.GetOrderInfo)
	s.Server.Engine.PUT("/order/status/:id", s.UpdateOrderStatus)
	return s.Server.Serve()
}

type OptionFunc func(*OrderHttpServer)

func WithLogger(l *otelzap.Logger) OptionFunc {
	return func(s *OrderHttpServer) {
		s.Logger = l
	}
}

func WithOrderDataService(s orderdata.OrderDataService) OptionFunc {
	return func(h *OrderHttpServer) {
		h.OrderData = s
	}
}
