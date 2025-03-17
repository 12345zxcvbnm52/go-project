package goodscontroller

import (
	goodsform "kenshop/api/goods/internal/form"
	"kenshop/pkg/common/httputil"
	proto "kenshop/proto/goods"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 品牌分类服务
// @Description 获取品牌与目录关联列表
// @Produce application/json
// @Router /categorybrands [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.CategoryBrandListRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param pagesNum query integer false "(int32)返回数据集的页号"
// @Param pageSize query integer false "(int32)返回数据集的页大小"
func (s *GoodsHttpServer) GetCategoryBrandList(c *gin.Context) {
	u := &goodsform.CategoryBrandFilterForm{}

	if err := c.ShouldBindQuery(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.GetCategoryBrandListDB(s.Server.Ctx, &proto.CategoryBrandFilterReq{
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

// 通过一个类型获得所有有这个类型的品牌
// rpc GetBrandListByCategory(CategoryInfoReq)returns(BrandListRes);
// @Accept application/json
// @Description 创建一个品牌与目录关联
// @Produce application/json
// @Router /categorybrands [POST]
// @Success 200 {object} httputil.JsonResult{data=proto.CategoryBrandInfoRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Param categoryId body integer true "(uint32)目录ID"
// @Param brandId body integer true "(uint32)品牌ID"
func (s *GoodsHttpServer) CreateCategoryBrand(c *gin.Context) {
	u := &goodsform.CreateCategoryBrandForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.CreateCategoryBrandDB(s.Server.Ctx, &proto.CreateCategoryBrandReq{
		BrandId:    u.BrandId,
		CategoryId: u.CategoryId,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Description 删除一个品牌与目录关联
// @Produce application/json
// @Router /categorybrand/{id} [DELETE]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)关联ID"
func (s *GoodsHttpServer) DeleteCategoryBrand(c *gin.Context) {
	u := &goodsform.DelCategoryBrandForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.DeleteCategoryBrandDB(s.Server.Ctx, &proto.DelCategoryBrandReq{
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
// @Description 更新一个品牌与目录关联
// @Produce application/json
// @Router /categorybrand/{id} [PATCH]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)关联ID"
// @Param categoryId body integer false "(uint32)目录ID"
// @Param brandId body integer false "(uint32)品牌ID"
func (s *GoodsHttpServer) UpdateCategoryBrand(c *gin.Context) {
	u := &goodsform.UpdateCategoryBrandForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}
	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.UpdateCategoryBrandDB(s.Server.Ctx, &proto.UpdateCategoryBrandReq{
		BrandId:    u.BrandId,
		CategoryId: u.CategoryId,
		Id:         u.Id,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}
