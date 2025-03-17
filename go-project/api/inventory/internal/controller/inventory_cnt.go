package inventorycontroller

import (
	inventorydata "kenshop/api/inventory/internal/data"
	inventoryform "kenshop/api/inventory/internal/form"
	httpserver "kenshop/goken/server/httpserver"
	"kenshop/pkg/common/httputil"
	"kenshop/pkg/log"
	proto "kenshop/proto/inventory"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type InventoryHttpServer struct {
	Server        *httpserver.Server
	InventoryData inventorydata.InventoryDataService
	Logger        *otelzap.Logger
}

// 默认使用otelzap.Logger以及GrpcInventoryData
func MustNewInventoryHTTPServer(s *httpserver.Server, opts ...OptionFunc) *InventoryHttpServer {
	ss := &InventoryHttpServer{
		Server: s,
	}
	for _, opt := range opts {
		opt(ss)
	}
	if ss.Logger == nil {
		ss.Logger = log.MustNewOtelLogger()
	}
	if ss.InventoryData == nil {
		cli, err := s.GrpcCli.Dial()
		if err != nil {
			panic(err)
		}
		ss.InventoryData = inventorydata.MustNewGrpcInventoryData(cli)
	}
	return ss
}

// 后续可以考虑给每个库存添加id(同一个id则增加库存量,否则根据地址创建),address(地址)
// @Accept application/json
// @Description 创建一个商品库存
// @Produce application/json
// @Router /inventories [POST]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Param goodsId body integer true "(uint32)欲创建商品库存的id"
// @Param goodsNum body integer true "(int32)欲创建商品库存的数量"
func (s *InventoryHttpServer) CreateStock(c *gin.Context) {
	u := &inventoryform.CreateInventoryForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.InventoryData.CreateStockDB(s.Server.Ctx, &proto.CreateInventoryReq{
		GoodsId:  u.GoodsId,
		GoodsNum: u.GoodsNum,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 更新设置商品库存
// @Produce application/json
// @Router /inventory [PUT]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Param goodsId body integer true "(uint32)欲更新设置商品库存的id"
// @Param goodsNum body integer true "(int32)欲更新设置商品库存的数量"
func (s *InventoryHttpServer) SetStock(c *gin.Context) {
	u := &inventoryform.SetInventoryForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.InventoryData.SetStockDB(s.Server.Ctx, &proto.SetInventoryReq{
		GoodsId:  u.GoodsId,
		GoodsNum: u.GoodsNum,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Description 获取商品库存信息
// @Produce application/json
// @Router /inventory/{goods_id} [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.InventoryInfoRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param goodsId path integer true "(uint32)商品库存的id"
func (s *InventoryHttpServer) GetStockInfo(c *gin.Context) {
	u := &inventoryform.InventoryInfoForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.InventoryData.GetStockInfoDB(s.Server.Ctx, &proto.InventoryInfoReq{
		GoodsId: u.GoodsId,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 扣减商品库存
// @Produce application/json
// @Router /inventory/decr [POST]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param decrData body array true "([]*proto.UpdateInventoryReq)扣减库存的数据"
// @Param orderSign body string true "(string)订单签名"
func (s *InventoryHttpServer) DecrStock(c *gin.Context) {
	u := &inventoryform.DecrStockForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}
	decrData := []*proto.UpdateInventoryReq{}
	for _, v := range u.DecrData {
		decrData = append(decrData, &proto.UpdateInventoryReq{
			GoodsId:  v.GoodsId,
			GoodsNum: v.GoodsNum,
		})
	}
	res, err := s.InventoryData.DecrStockDB(s.Server.Ctx, &proto.DecrStockReq{
		DecrData:  decrData,
		OrderSign: u.OrderSign,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 增加商品库存
// @Produce application/json
// @Router /inventory/incr [POST]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param incrData body array true "([]*proto.UpdateInventoryReq)增加库存的数据"
// @Param orderSign body string true "(string)订单签名"
func (s *InventoryHttpServer) IncrStock(c *gin.Context) {
	u := &inventoryform.IncrStockForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}
	incrData := []*proto.UpdateInventoryReq{}
	for _, v := range u.IncrData {
		incrData = append(incrData, &proto.UpdateInventoryReq{
			GoodsId:  v.GoodsId,
			GoodsNum: v.GoodsNum,
		})
	}
	res, err := s.InventoryData.IncrStockDB(s.Server.Ctx, &proto.UpdateStockReq{
		IncrData:  incrData,
		OrderSign: u.OrderSign,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

func (s *InventoryHttpServer) Execute() error {
	s.Server.Engine.POST("/inventories", s.CreateStock)
	s.Server.Engine.PUT("/inventory", s.SetStock)
	s.Server.Engine.GET("/inventory/:goods_id", s.GetStockInfo)
	s.Server.Engine.POST("/inventory/decr", s.DecrStock)
	s.Server.Engine.POST("/inventory/incr", s.IncrStock)
	return s.Server.Serve()
}

type OptionFunc func(*InventoryHttpServer)

func WithLogger(l *otelzap.Logger) OptionFunc {
	return func(s *InventoryHttpServer) {
		s.Logger = l
	}
}

func WithInventoryDataService(s inventorydata.InventoryDataService) OptionFunc {
	return func(h *InventoryHttpServer) {
		h.InventoryData = s
	}
}
