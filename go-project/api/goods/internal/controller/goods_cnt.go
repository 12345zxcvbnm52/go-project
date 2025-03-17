package goodscontroller

import (
	goodsform "kenshop/api/goods/internal/form"
	"kenshop/pkg/common/httputil"
	proto "kenshop/proto/goods"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 获得商品列表
// @Description 获得商品列表
// @Produce application/json
// @Router /goods [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.GoodsListRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param minPrice query integer false "(int32)商品最低价格"
// @Param maxPrice query integer false "(int32)商品最高价格"
// @Param isHot query bool false "是否热门商品"
// @Param isNew query bool false "是否新品"
// @Param status query integer false "(int32)商品状态"
// @Param categoryId query integer false "(uint32)类别目录ID,确定点击的目录是哪一级,会递归显示下一层类别目录"
// @Param pagesNum query integer false "(int32)返回数据集的页号"
// @Param pageSize query integer false "(int32)返回数据集的页大小"
// @Param keyWords query string false "搜索关键词"
// @Param brandId query integer false "(uint32)品牌ID"
// @Param id query integer false "(uint32)商品ID"
func (s *GoodsHttpServer) GetGoodList(c *gin.Context) {
	u := &goodsform.GoodsFilterForm{}

	if err := c.ShouldBindQuery(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.GetGoodListDB(s.Server.Ctx, &proto.GoodsFilterReq{
		BrandId:    u.BrandId,
		CategoryId: u.CategoryId,
		Id:         u.Id,
		IsHot:      u.IsHot,
		IsNew:      u.IsNew,
		KeyWords:   u.KeyWords,
		MaxPrice:   u.MaxPrice,
		MinPrice:   u.MinPrice,
		PageSize:   u.PageSize,
		PagesNum:   u.PagesNum,
		Status:     u.Status,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// 用于通过id数组得到所有商品信息,常用于从订单中获得所有商品信息,
// @Accept multipart/form-data
// @Description 用于通过id数组得到所有商品信息,常用于从订单中获得所有商品信息
// @Produce application/json
// @Router /goods/ids [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.GoodsListRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param ids formData []integer true "([]uint32)商品ID"
func (s *GoodsHttpServer) GetGoodsListById(c *gin.Context) {
	u := &goodsform.GoodsIdsForm{}

	if err := httputil.ShouldBindFormSlice(c, u); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			httputil.WriteValidateError(c, s.Server.Validator.Trans, verr, s.Server.UseAbort)
		}
		httputil.WriteResponse(c, http.StatusBadRequest, err.Error(), nil, s.Server.UseAbort)
		return
	}
	res, err := s.GoodsData.GetGoodsListByIdDB(s.Server.Ctx, &proto.GoodsIdsReq{
		Ids: u.Ids,
	})

	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept multipart/form-data
// @Description 通过商品id获取对应商品的详细信息
// @Produce application/json
// @Router /goods/detail/{id} [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.GoodsDetailRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)商品ID"
func (s *GoodsHttpServer) GetGoodsDetail(c *gin.Context) {
	u := &goodsform.GoodsInfoForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.GetGoodsDetailDB(s.Server.Ctx, &proto.GoodsInfoReq{
		Id: u.Id,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 创建一个商品
// @Produce application/json
// @Router /goods [POST]
// @Success 200 {object} httputil.JsonResult{data=proto.GoodsDetailRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Param categoryId body integer true "(uint32)商品类型id"
// @Param brandId body integer true "(uint32)商品品牌id"
// @Param name body string true "商品名"
// @Param marketPrice body number true "(float32)商品原始售价"
// @Param salePrice body number true "(float32)商品实际售价"
// @Param goodsBrief body string true "商品简要介绍"
// @Param shipFree body bool true "(int32)是否免运费"
// @Param images body model.GormList true "商品图片"
// @Param descImages body model.GormList true "商品功能图"
// @Param firstImage body string true "商品封面"
func (s *GoodsHttpServer) CreateGoods(c *gin.Context) {
	u := &goodsform.CreateGoodsForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.CreateGoodsDB(s.Server.Ctx, &proto.CreateGoodsReq{
		BrandId:     u.BrandId,
		CategoryId:  u.CategoryId,
		DescImages:  u.DescImages,
		FirstImage:  u.FirstImage,
		GoodsBrief:  u.GoodsBrief,
		Images:      u.Images,
		MarketPrice: u.MarketPrice,
		Name:        u.Name,
		SalePrice:   u.SalePrice,
		ShipFree:    u.ShipFree,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Description 删除一个商品
// @Produce application/json
// @Router /goods/{id} [DELETE]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)商品id"
func (s *GoodsHttpServer) DeleteGoods(c *gin.Context) {
	u := &goodsform.DelGoodsForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.DeleteGoodsDB(s.Server.Ctx, &proto.DelGoodsReq{
		Id: u.Id,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 局部更新一个商品
// @Produce application/json
// @Router /goods/{id} [PATCH]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)商品id"
// @Param categoryId body integer false "(uint32)商品类型id"
// @Param brandId body integer false "(uint32)商品品牌id"
// @Param name body string false "商品名"
// @Param marketPrice body number false "(float32)商品原始售价"
// @Param salePrice body number false "(float32)商品实际售价"
// @Param goodsBrief body string false "商品简要介绍"
// @Param shipFree body bool false "(int32)是否免运费"
// @Param images body model.GormList false "商品图片"
// @Param descImages body model.GormList false "商品功能图"
// @Param firstImage body string false "商品封面"
// @Param status body integer false "商品状态"
func (s *GoodsHttpServer) UpdeateGoods(c *gin.Context) {
	u := &goodsform.UpdateGoodsForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}
	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.UpdeateGoodsDB(s.Server.Ctx, &proto.UpdateGoodsReq{
		BrandId:     u.BrandId,
		CategoryId:  u.CategoryId,
		DescImages:  u.DescImages,
		FirstImage:  u.FirstImage,
		GoodsBrief:  u.GoodsBrief,
		Id:          u.Id,
		Images:      u.Images,
		MarketPrice: u.MarketPrice,
		Name:        u.Name,
		SalePrice:   u.SalePrice,
		ShipFree:    u.ShipFree,
		Status:      u.Status,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}
