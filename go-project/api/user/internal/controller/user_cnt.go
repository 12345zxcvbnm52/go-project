package usercontroller

import (
	"context"
	userdata "kenshop/api/user/internal/data"
	userform "kenshop/api/user/internal/form"
	_ "kenshop/docs"
	"kenshop/goken/server/httpserver"
	"kenshop/pkg/common/httputil"
	"kenshop/pkg/log"
	proto "kenshop/proto/user"
	"net/http"

	"github.com/gin-gonic/gin"
	ginfile "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type OptionFunc func(*UserHttpServer)

// 默认使用otelzap.Logger以及GrpcUserData
func MustNewUserHTTPServer(s *httpserver.Server, opts ...OptionFunc) *UserHttpServer {
	ss := &UserHttpServer{
		Server: s,
	}
	for _, opt := range opts {
		opt(ss)
	}
	if ss.Logger == nil {
		ss.Logger = log.MustNewOtelLogger()
	}
	if ss.UserData == nil {
		cli, err := s.GrpcCli.Dial()
		if err != nil {
			panic(err)
		}
		ss.UserData = userdata.MustNewGrpcUserData(cli)
	}
	return ss
}

func (s *UserHttpServer) GetCtxFromGinCtx(c *gin.Context) context.Context {
	var ctx context.Context = s.Server.Ctx
	ctxAny, ok := c.Get(s.Server.Tracer.SpanGinCtxKey)
	if ok {
		ctx = ctxAny.(context.Context)
	}
	return ctx
}

type UserHttpServer struct {
	Server   *httpserver.Server
	UserData userdata.UserDataService
	Logger   *otelzap.Logger
}

func (s *UserHttpServer) WriteError(c *gin.Context, code int, msg gin.H) {
	if s.Server.UseAbort {
		c.Abort()
	}
	c.JSON(code, msg)
}

// 获取用户列表
// @Description 获取用户列表信息
// @Produce application/json
// @Router /users [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.UserListRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Param pagesNum query integer false "(int32)返回数据集的页号"
// @Param pageSize query integer false "(int32)返回数据集的页大小"

func (s *UserHttpServer) GetUserList(c *gin.Context) {
	u := &userform.UserFliterForm{}

	if err := c.ShouldBindQuery(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.UserData.GetUserListDB(s.GetCtxFromGinCtx(c), &proto.UserFliterReq{
		PageSize: u.PageSize,
		PagesNum: u.PagesNum,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": res.Total,
		"data":  res.Data,
	})
}

// 通过用户ID获取用户信息
// @Description 通过用户ID获取用户信息
// @Produce application/json
// @Router /user/id/{id} [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.UserInfoRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)用户ID"
func (s *UserHttpServer) GetUserById(c *gin.Context) {
	u := &userform.UserIdForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.UserData.GetUserByIdDB(s.GetCtxFromGinCtx(c), &proto.UserIdReq{
		Id: u.Id,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_data": res,
	})
}

// 通过手机号获取用户信息
// @Description 通过手机号获取用户信息
// @Produce application/json
// @Router /user/mobile/{mobile} [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.UserInfoRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param mobile path string true "用户手机号"
func (s *UserHttpServer) GetUserByMobile(c *gin.Context) {
	u := &userform.UserMobileForm{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.UserData.GetUserByMobileDB(s.GetCtxFromGinCtx(c), &proto.UserMobileReq{
		Mobile: u.Mobile,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_data": res,
	})
}

// 创建用户
// @Accept application/json
// @Description 创建用户
// @Produce application/json
// @Router /users [POST]
// @Success 200 {object} httputil.JsonResult{data=proto.CreateUserRes}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Param userName body string true "欲创建的用户名"
// @Param password body string true "欲创建的用户密码"
// @Param mobile body string true "欲创建的用户手机号"
// @Param gender body string false "欲创建的用户性别"
// @Param birth body integer true "(int32)欲创建的用户出生日"
func (s *UserHttpServer) CreateUser(c *gin.Context) {
	u := &userform.CreateUserForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.UserData.CreateUserDB(s.GetCtxFromGinCtx(c), &proto.CreateUserReq{
		Birth:    u.Birth,
		Gender:   u.Gender,
		Mobile:   u.Mobile,
		Password: u.Password,
		UserName: u.UserName,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_data": res,
	})
}

// 绝对更新用户(全量更新)
// @Accept application/json
// @Description 绝对更新用户(全量更新)
// @Produce application/json
// @Router /user/{id} [PUT]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)欲修改用户ID"
// @Param userName body string true "欲修改的用户名"
// @Param password body string true "欲修改的用户密码"
// @Param mobile body string true "欲修改的用户手机号"
// @Param gender body string true "欲修改的用户性别"
// @Param birth body integer true "(int32)欲修改的用户出生日"
// @Param role body integer true "(int32)欲修改的用户权限"
func (s *UserHttpServer) AbsUpdateUser(c *gin.Context) {
	u := &userform.UpdateUserForm_0{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	_, err := s.UserData.AbsUpdateUserDB(s.GetCtxFromGinCtx(c), &proto.UpdateUserReq{
		Birth:    u.Birth,
		Gender:   u.Gender,
		Id:       u.Id,
		Mobile:   u.Mobile,
		Password: u.Password,
		Role:     u.Role,
		UserName: u.UserName,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

// 局部更新用户(部分字段更新)
// @Accept application/json
// @Description 局部更新用户(部分更新)
// @Produce application/json
// @Router /user/{id} [PATCH]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Param id path integer true "(uint32)欲修改用户ID"
// @Param userName body string false "欲修改的用户名"
// @Param password body string false "欲修改的用户密码"
// @Param mobile body string false "欲修改的用户手机号"
// @Param gender body string false "欲修改的用户性别"
// @Param birth body integer false "(int32)欲修改的用户出生日"
// @Param role body integer false "(int32)欲修改的用户权限"
func (s *UserHttpServer) UpdateUser(c *gin.Context) {
	u := &userform.UpdateUserForm_1{}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	_, err := s.UserData.UpdateUserDB(s.GetCtxFromGinCtx(c), &proto.UpdateUserReq{
		Birth:    u.Birth,
		Gender:   u.Gender,
		Id:       u.Id,
		Mobile:   u.Mobile,
		Password: u.Password,
		Role:     u.Role,
		UserName: u.UserName,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

// 删除用户
// @Description 删除用户
// @Produce application/json
// @Router /user/{id} [DELETE]
// @Success 200 {object} httputil.JsonResult{data=nil}
// @Failure 401 {object} httputil.JsonResult{data=nil}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Param id path int true "欲删除的用户ID"
// @Param name query string true "欲删除的用户名"
func (s *UserHttpServer) DeleteUser(c *gin.Context) {
	u := &userform.DelUserForm{}

	if err := c.ShouldBindQuery(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	_, err := s.UserData.DeleteUserDB(s.GetCtxFromGinCtx(c), &proto.DelUserReq{
		Id:   u.Id,
		Name: u.Name,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

// 用户权限验证
// @Accept application/json
// @Produce multipart/form-data
// @Router /user/login [GET]
// @Success 200 {object} httputil.JsonResult{data=proto.UserCheckRes}
// @Failure 500 {object} httputil.JsonResult{data=nil}
// @Failure 400 {object} httputil.JsonResult{data=nil}
// @Failure 404 {object} httputil.JsonResult{data=nil}
// @Description 用户权限验证
// @Param userName formData string true "欲检查的用户名"
// @Param password formData string true "欲检查的用户密码"
func (s *UserHttpServer) CheckUserRole(c *gin.Context) {
	u := &userform.UserPasswordForm{}

	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}

	res, err := s.UserData.CheckUserRoleDB(s.GetCtxFromGinCtx(c), &proto.UserPasswordReq{
		Password: u.Password,
		UserName: u.UserName,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		c.Abort()
		return
	}

	if !res.Ok {
		s.WriteError(c, http.StatusUnauthorized, gin.H{
			"msg": "账户或密码错误",
		})
	} else {
		kv := []string{"username", u.UserName}
		token, expire, err := s.Server.Jwt.NewToken(kv...)
		if err != nil {
			s.WriteError(c, http.StatusInternalServerError, gin.H{
				"msg": "服务器内部错误,无法生成jwt token",
			})

		}
		c.JSON(http.StatusOK, gin.H{
			"msg":       "登入成功",
			"jwt-token": token,
			"expired":   expire,
		})
	}

}

func (s *UserHttpServer) Execute() error {
	s.Server.Engine.GET("/users", s.Server.Jwt.RefreshHandler, s.Server.Jwt.JwtAuthHandler, s.Server.Tracer.TraceHandler(), s.GetUserList)
	s.Server.Engine.GET("/user/id/:id", s.Server.Jwt.RefreshHandler, s.Server.Jwt.JwtAuthHandler, s.Server.Tracer.TraceHandler(), s.GetUserById)
	s.Server.Engine.GET("/user/mobile/:mobile", s.Server.Jwt.RefreshHandler, s.Server.Jwt.JwtAuthHandler, s.Server.Tracer.TraceHandler(), s.GetUserByMobile)
	s.Server.Engine.POST("/users", s.Server.Tracer.TraceHandler(), s.CreateUser)
	s.Server.Engine.PUT("/user/:id", s.Server.Jwt.RefreshHandler, s.Server.Jwt.JwtAuthHandler, s.Server.Tracer.TraceHandler(), s.AbsUpdateUser)
	s.Server.Engine.PATCH("/user/:id", s.Server.Jwt.RefreshHandler, s.Server.Jwt.JwtAuthHandler, s.Server.Tracer.TraceHandler(), s.UpdateUser)
	s.Server.Engine.DELETE("/user/:id", s.Server.Jwt.RefreshHandler, s.Server.Jwt.JwtAuthHandler, s.Server.Tracer.TraceHandler(), s.DeleteUser)
	s.Server.Engine.GET("/user/login", s.Server.Tracer.TraceHandler(), s.CheckUserRole)
	s.Server.Engine.GET("swagger/*any", ginswagger.WrapHandler(ginfile.Handler))
	if err := s.Server.Validator.Excute(); err != nil {
		return err
	}
	return s.Server.Serve()
}

func WithLogger(l *otelzap.Logger) OptionFunc {
	return func(s *UserHttpServer) {
		s.Logger = l
	}
}

func WithUserDataService(s userdata.UserDataService) OptionFunc {
	return func(h *UserHttpServer) {
		h.UserData = s
	}
}
