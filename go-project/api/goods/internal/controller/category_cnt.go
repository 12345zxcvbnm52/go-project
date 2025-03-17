package goodscontroller

import (
	goodsform "kenshop/api/goods/internal/form"
	"kenshop/pkg/common/httputil"
	proto "kenshop/proto/goods"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 商品类型服务
// @Description 获取商品目录列表
// @Produce application/json
// @Router /categories [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.CategoryListRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
func (s *GoodsHttpServer) GetCategoryList(c *gin.Context) {
	res, err := s.GoodsData.GetCategoryListDB(s.Server.Ctx, &emptypb.Empty{})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Description 获取商品目录及其子目录信息
// @Produce application/json
// @Router /category/{id} [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.SubCategoryListRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)商品目录ID"
// @Param level query integer false "(int32)目录层级"
func (s *GoodsHttpServer) GetCategoryInfo(c *gin.Context) {
	u := &goodsform.SubCategoryForm{}

	if err := c.ShouldBindQuery(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.GetCategoryInfoDB(s.Server.Ctx, &proto.SubCategoryReq{
		Id:    u.Id,
		Level: u.Level,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 创建一个商品目录
// @Produce application/json
// @Router /categories [POST]
// @Success 200 {object} httputil.JsonResult{data=proto.CategoryInfoRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Param name body string true "商品目录名称"
// @Param parentCategoryId body integer false "(uint32)父目录ID"
// @Param level body integer true "(int32)目录层级"
func (s *GoodsHttpServer) CreateCategory(c *gin.Context) {
	u := &goodsform.CreateCategoryForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.CreateCategoryDB(s.Server.Ctx, &proto.CreateCategoryReq{
		Level:            u.Level,
		Name:             u.Name,
		ParentCategoryId: u.ParentCategoryId,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Description 删除一个商品目录
// @Produce application/json
// @Router /category/{id} [DELETE]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)商品目录ID"
func (s *GoodsHttpServer) DeleteCategory(c *gin.Context) {
	u := &goodsform.DelCategoryForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.DeleteCategoryDB(s.Server.Ctx, &proto.DelCategoryReq{
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
// @Description 更新一个商品目录
// @Produce application/json
// @Router /category/{id} [PATCH]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)商品目录ID"
// @Param name body string false "商品目录名称"
// @Param parentCategoryId body integer false "(uint32)父目录ID"
// @Param level body integer false "(int32)目录层级"
func (s *GoodsHttpServer) UpdateCategory(c *gin.Context) {
	u := &goodsform.UpdateCategoryForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}
	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.UpdateCategoryDB(s.Server.Ctx, &proto.UpdateCategoryReq{
		Id:               u.Id,
		Level:            u.Level,
		Name:             u.Name,
		ParentCategoryId: u.ParentCategoryId,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}
