package handler

import (
	gb "goods_web/global"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
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

func RpcErrorHandle(c *gin.Context, err error) {
	zap.S().Errorw("微服务调用失败", "msg", err.Error())
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
				"msg": e.Message(),
			})
		case codes.Unavailable:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "商品服务暂不可用",
			})
		case codes.AlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"msg": e.Message(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "未知错误",
			})
		}
	}
}

func ValidatorErrorHandle(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		zap.S().Errorw("Validate认证失败", "msg", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"msg": RemoveStructPrefix(errs.Translate(gb.Translator)),
	})
	return
}
