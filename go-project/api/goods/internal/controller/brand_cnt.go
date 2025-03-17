package goodscontroller

import (
	goodsform "kenshop/api/goods/internal/form"
	"kenshop/pkg/common/httputil"
	proto "kenshop/proto/goods"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 品牌服务
// @Description 获取品牌列表
// @Produce application/json
// @Router /brands [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.BrandListRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param pagesNum query integer false "(int32)返回数据集的页号"
// @Param pageSize query integer false "(int32)返回数据集的页大小"
func (s *GoodsHttpServer) GetBrandList(c *gin.Context) {
	u := &goodsform.BrandFilterForm{}

	if err := c.ShouldBindQuery(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.GetBrandListDB(s.Server.Ctx, &proto.BrandFilterReq{
		PageSize: u.PageSize,
		PagesNum: u.PagesNum,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 创建一个品牌
// @Produce application/json
// @Router /brand [POST]
// @Success 200 {object} httputil.JsonResult{data=proto.BrandInfoRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Param name body string true "品牌名称"
// @Param logo body string true "品牌Logo"
func (s *GoodsHttpServer) CreateBrand(c *gin.Context) {
	u := &goodsform.CreateBrandForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.CreateBrandDB(s.Server.Ctx, &proto.CreateBrandReq{
		Logo: u.Logo,
		Name: u.Name,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Description 删除一个品牌
// @Produce application/json
// @Router /brand/{id} [DELETE]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)品牌ID"
// @Param name query string true "品牌名称"
func (s *GoodsHttpServer) DeleteBrand(c *gin.Context) {
	u := &goodsform.DelBrandForm{}

	if err := c.ShouldBindQuery(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.DeleteBrandDB(s.Server.Ctx, &proto.DelBrandReq{
		Id:   u.Id,
		Name: u.Name,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 更新一个品牌
// @Produce application/json
// @Router /brand/{id} [PATCH]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)品牌ID"
// @Param name body string false "品牌名称"
// @Param logo body string false "品牌Logo"
func (s *GoodsHttpServer) UpdateBrand(c *gin.Context) {
	u := &goodsform.UpdateBrandForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}
	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.UpdateBrandDB(s.Server.Ctx, &proto.UpdateBrandReq{
		Id:   u.Id,
		Logo: u.Logo,
		Name: u.Name,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}
