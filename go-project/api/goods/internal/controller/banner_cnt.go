package goodscontroller

import (
	goodsform "kenshop/api/goods/internal/form"
	"kenshop/pkg/common/httputil"
	proto "kenshop/proto/goods"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 轮播窗口服务
// @Description 获取轮播图列表
// @Produce application/json
// @Router /banners [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.BannerListRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
func (s *GoodsHttpServer) GetBannerList(c *gin.Context) {
	res, err := s.GoodsData.GetBannerListDB(s.Server.Ctx, &emptypb.Empty{})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 创建一个轮播图
// @Produce application/json
// @Router /banners [POST]
// @Success 200 {object} httputil.JsonResult{data=proto.BannerInfoRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Param index body integer true "(int32)轮播图序号"
// @Param image body string true "轮播图图片URL"
// @Param url body string true "轮播图跳转URL"
func (s *GoodsHttpServer) CreateBanner(c *gin.Context) {
	u := &goodsform.CreateBannerForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.CreateBannerDB(s.Server.Ctx, &proto.CreateBannerReq{
		Image: u.Image,
		Index: u.Index,
		Url:   u.Url,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Description 删除一个轮播图
// @Produce application/json
// @Router /banner/{id} [DELETE]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)轮播图ID"
// @Param index query integer false "(int32)轮播图索引"
func (s *GoodsHttpServer) DeleteBanner(c *gin.Context) {
	u := &goodsform.DelBannerForm{}

	if err := c.ShouldBindQuery(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.DeleteBannerDB(s.Server.Ctx, &proto.DelBannerReq{
		Id:    u.Id,
		Index: u.Index,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}

// @Accept application/json
// @Description 更新一个轮播图
// @Produce application/json
// @Router /banner/{id} [PATCH]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)轮播图ID"
// @Param index body integer false "(int32)轮播图序号"
// @Param image body string false "轮播图图片URL"
// @Param url body string false "轮播图跳转URL"
func (s *GoodsHttpServer) UpdateBanner(c *gin.Context) {
	u := &goodsform.UpdateBannerForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}
	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.GoodsData.UpdateBannerDB(s.Server.Ctx, &proto.UpdateBannerReq{
		Id:    u.Id,
		Image: u.Image,
		Index: u.Index,
		Url:   u.Url,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK, "", res, s.Server.UseAbort)
}
