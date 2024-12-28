package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
