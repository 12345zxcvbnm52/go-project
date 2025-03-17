package httputil

import (
	"kenshop/pkg/errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type JsonResult struct {
	Code    int         `json:"code" `
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func newJsonResult(code int, msg string, data interface{}) JsonResult {
	return JsonResult{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func (r JsonResult) WriteResponse(c *gin.Context, abort bool) {
	c.JSON(r.Code, gin.H{
		"data": r.Data,
		"msg":  r.Message,
		"code": r.Code,
	})
	if abort {
		c.Abort()
	}
}

func WriteResponse(c *gin.Context, code int, msg string, data interface{}, abort bool) {
	newJsonResult(code, msg, data).WriteResponse(c, abort)
}

func WriteRpcError(c *gin.Context, err error, abort bool) {
	e := errors.ExtractCodeErrorFromGRPC(err)
	if cerr := errors.ExtractCoderFromError(e); cerr != nil {
		newJsonResult(cerr.HTTPCode(), cerr.Message(), nil).WriteResponse(c, abort)
	} else {
		newJsonResult(http.StatusInternalServerError, "服务器内部错误,请稍后再试", nil).WriteResponse(c, abort)
	}
}

func WriteValidateError(c *gin.Context, trans ut.Translator, err error, abort bool) {
	verr, ok := err.(validator.ValidationErrors)
	if !ok {
		newJsonResult(http.StatusInternalServerError, err.Error(), nil).WriteResponse(c, abort)
	}
	newJsonResult(http.StatusBadRequest, "传入参数错误", translateErr(verr, trans)).WriteResponse(c, abort)
}

func translateErr(err validator.ValidationErrors, trans ut.Translator) map[string]string {
	// 移除默认字段检测时多出来的结构体名称.
	f := func(msg map[string]string) map[string]string {
		res := map[string]string{}
		for k, v := range msg {
			res[k[strings.Index(k, ".")+1:]] = v
		}
		return res
	}

	return f(err.Translate(trans))
}
