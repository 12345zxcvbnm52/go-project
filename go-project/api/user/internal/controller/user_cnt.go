package usercontroller

import (
	"fmt"
	userdata "kenshop/api/user/internal/data"
	"kenshop/api/user/internal/form"
	"kenshop/goken/server/httpserver"
	"kenshop/pkg/log"
	proto "kenshop/proto/user"
	"net/http"

	"github.com/gin-gonic/gin"
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

type UserHttpServer struct {
	Server   *httpserver.Server
	UserData userdata.UserDataService
	Logger   *otelzap.Logger
}

// 获取用户列表
// @Failure 404 {object} map[string]interface{}
// @Produce application/json
// @Router /users [GET]
// @Success 200 {object} map[string]interface{}
// @Param pagesNum query int false "返回数据集的页号"
// @Param pageSize query int false "返回数据集的页大小"
func (s *UserHttpServer) GetUserList(c *gin.Context) {
	u := &form.UserFliterForm{}
	if err := c.ShouldBindQuery(u); err != nil {
		s.ValidatorErrorHandle(c, err)
		return
	}
	fmt.Println(u)
	res, err := s.UserData.GetUserListDB(s.Server.Ctx, &proto.UserFliterReq{
		PageSize: u.PageSize,
		PagesNum: u.PagesNum,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		RpcErrorHandle(c, err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": res.Total,
		"data":  res.Data,
	})
}

// 通过用户ID获取用户信息
// @Failure 404 {object} map[string]interface{}
// @Produce application/json
// @Router /user/id/{id} [GET]
// @Success 200 {object} map[string]interface{}
// @Param id path int true "用户ID"
func (s *UserHttpServer) GetUserById(c *gin.Context) {
	u := &form.UserIdForm{}
	if err := c.ShouldBind(u); err != nil {
		s.ValidatorErrorHandle(c, err)
		return
	}

	res, err := s.UserData.GetUserByIdDB(s.Server.Ctx, &proto.UserIdReq{
		Id: u.Id,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		RpcErrorHandle(c, err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_data": res,
	})
}

// 通过手机号获取用户信息
// @Failure 404 {object} map[string]interface{}
// @Produce application/json
// @Router /user/model/{mobile} [GET]
// @Success 200 {object} map[string]interface{}
// @Param mobile path string true "用户手机号"
func (s *UserHttpServer) GetUserByMobile(c *gin.Context) {
	u := &form.UserMobileForm{}
	if err := c.ShouldBind(u); err != nil {
		s.ValidatorErrorHandle(c, err)
		return
	}

	res, err := s.UserData.GetUserByMobileDB(s.Server.Ctx, &proto.UserMobileReq{
		Mobile: u.Mobile,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		RpcErrorHandle(c, err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_data": res,
	})
}

// 创建用户
// @Accept application/json
// @Failure 400 {object} map[string]interface{}
// @Produce application/json
// @Router /users [POST]
// @Success 201 {object} map[string]interface{}
// @Param username body string true "欲创建的用户名"
// @Param password body string true "欲创建的用户密码"
// @Param mobile body string true "欲创建的用户手机号"
// @Param gender body string false "欲创建的用户性别"
// @Param birth body int true "欲创建的用户出生日"
func (s *UserHttpServer) CreateUser(c *gin.Context) {
	u := &form.CreateUserForm{}
	if err := c.ShouldBind(u); err != nil {
		s.ValidatorErrorHandle(c, err)
		return
	}

	res, err := s.UserData.CreateUserDB(s.Server.Ctx, &proto.CreateUserReq{
		Birth:    u.Birth,
		Gender:   u.Gender,
		Mobile:   u.Mobile,
		Password: u.Password,
		UserName: u.UserName,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		RpcErrorHandle(c, err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_data": res,
	})
}

// 绝对更新用户(全量更新)
// @Accept application/json
// @Failure 400 {object} map[string]interface{}
// @Produce application/json
// @Router /user/{id} [PUT]
// @Success 204 {object} map[string]interface{}
// @Param id path int true "欲修改用户ID"
// @Param username body string true "欲修改的用户名"
// @Param password body string true "欲修改的用户密码"
// @Param mobile body string true "欲修改的用户手机号"
// @Param gender body string true "欲修改的用户性别"
// @Param birth body int true "欲修改的用户出生日"
// @Param role body int true "欲修改的用户权限"
func (s *UserHttpServer) AbsUpdateUser(c *gin.Context) {
	u := &form.UpdateUserForm{}
	if err := c.ShouldBind(u); err != nil {
		s.ValidatorErrorHandle(c, err)
		return
	}

	_, err := s.UserData.AbsUpdateUserDB(s.Server.Ctx, &proto.UpdateUserReq{
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
		RpcErrorHandle(c, err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

// 局部更新用户(部分字段更新)
// @Accept application/json
// @Failure 400 {object} map[string]interface{}
// @Produce application/json
// @Router /user/{id} [PATCH]
// @Success 204 {object} map[string]interface{}
// @Param id path int true "欲修改用户ID"
// @Param username body string false "欲修改的用户名"
// @Param password body string false "欲修改的用户密码"
// @Param mobile body string false "欲修改的用户手机号"
// @Param gender body string false "欲修改的用户性别"
// @Param birth body int false "欲修改的用户出生日"
// @Param role body int false "欲修改的用户权限"
func (s *UserHttpServer) UpdateUser(c *gin.Context) {
	u := &form.UpdateUserForm{}
	if err := c.ShouldBind(u); err != nil {
		s.ValidatorErrorHandle(c, err)
		return
	}

	_, err := s.UserData.UpdateUserDB(s.Server.Ctx, &proto.UpdateUserReq{
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
		RpcErrorHandle(c, err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

// 删除用户
// @Failure 404 {object} map[string]interface{}
// @Produce application/json
// @Router /user/{id} [DELETE]
// @Success 204 {object} map[string]interface{}
// @Param id path int true "欲删除的用户ID"
// @Param name query string true "欲删除的用户名"
func (s *UserHttpServer) DeleteUser(c *gin.Context) {
	u := &form.DelUserForm{}
	if err := c.ShouldBind(u); err != nil {
		s.ValidatorErrorHandle(c, err)
		return
	}

	_, err := s.UserData.DeleteUserDB(s.Server.Ctx, &proto.DelUserReq{
		Id:   u.Id,
		Name: u.Name,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		RpcErrorHandle(c, err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

// 用户权限验证
// @Accept application/json
// @Failure 403 {object} map[string]interface{}
// @Produce application/json
// @Router /user/{id}/check [POST]
// @Success 200 {object} map[string]interface{}
// @Param id path int true "用户ID"
// @Param userName body string true "欲检查的用户名"
// @Param password body string true "欲检查的用户密码"
func (s *UserHttpServer) CheckUserRole(c *gin.Context) {
	u := &form.UserPasswordForm{}
	if err := c.ShouldBind(u); err != nil {
		s.ValidatorErrorHandle(c, err)
		return
	}

	res, err := s.UserData.CheckUserRoleDB(s.Server.Ctx, &proto.UserPasswordReq{
		Id:       u.Id,
		Password: u.Password,
		UserName: u.UserName,
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		RpcErrorHandle(c, err)
		c.Abort()
		return
	}

	if !res.Ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "账户或密码错误",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "登入成功",
		})
	}

}

func (s *UserHttpServer) Execute() error {
	s.Server.Engine.GET("/users", s.GetUserList)
	s.Server.Engine.GET("/user/id/:id", s.GetUserById)
	s.Server.Engine.GET("/user/model/:mobile", s.GetUserByMobile)
	s.Server.Engine.POST("/users", s.CreateUser)
	s.Server.Engine.PUT("/user/:id", s.AbsUpdateUser)
	s.Server.Engine.PATCH("/user/:id", s.UpdateUser)
	s.Server.Engine.DELETE("/user/:id", s.DeleteUser)
	s.Server.Engine.POST("/user/:id/check", s.CheckUserRole)
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
