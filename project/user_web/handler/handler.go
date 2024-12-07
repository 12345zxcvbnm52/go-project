package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"user_web/form"
	gb "user_web/global"
	"user_web/middlewares"
	pb "user_web/proto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

//巨坑:在gin中如果方法是get,无论指定的Content-Type是什么,ShouldBind都会采取Form表单读取
//巨坑:如果导入的gin下的validator没有加/v10,则标签检测产生的错误永远不能转化为validator.ValidationError
//而binding.Validator.Engine().(*Validator.Validate)也永远无法转化成功

// 移除默认字段检测时多出来的结构体名称.
func RemoveStructPrefix(msg map[string]string) map[string]string {
	res := map[string]string{}
	for k, v := range msg {
		res[k[strings.Index(k, ".")+1:]] = v
	}
	return res
}

func GrpcErrorToHttp(err error, c *gin.Context) {
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"msg": e.Message(),
			})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "服务器内部错误",
			})

		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "传入参数错误",
			})
		case codes.Unavailable:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "用户服务暂不可用",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "其他错误",
			})
		}
	}
}

func UserLogin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := grpc.DialContext(ctx,
		fmt.Sprintf("consul://%s:%d/%s?healthy=true",
			"192.168.199.128",
			8500,
			"user_srv",
		),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	cc := pb.NewUserClient(client)

	u := &form.PasswordLogin{}
	//待改
	if err := c.ShouldBindJSON(u); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": RemoveStructPrefix(errs.Translate(gb.Translator)),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
		}
		c.Abort()
		return
	}

	flag, err := cc.CheckUserRole(ctx, &pb.UserPasswordReq{
		Password: u.Password,
		UserName: u.UserName,
		Id:       u.Id,
	})

	if err != nil {
		zap.S().Errorw("微服务调用失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"错误信息": "服务器内部出错,请稍后尝试",
		})
		c.Abort()
		return
	}
	if flag.Ok {
		j := middlewares.NewJwtAuth()
		claim := middlewares.CustomClaims{
			ID:       u.Id,
			UserName: u.UserName,
			RegisteredClaims: jwt.RegisteredClaims{
				NotBefore: jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
				Issuer:    gb.ServerConfig.Name,
			},
		}
		//这里的错误可以记录,
		//得到的str可选择记录到redis里
		str, err := j.CreateToken(claim)
		if err != nil {
			zap.S().Errorw("未知原因导致无法创建token", "msg", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":     "登入成功",
			"x-token": str,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "登入失败,用户名或密码错误",
		})
	}
}
