package usercontroller

import (
	"kenshop/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func RpcErrorHandle(c *gin.Context, err error) {
	if e, ok := err.(errors.Coder); ok {
		c.JSON(e.HTTPCode(), e.Message())
	} else {
		c.JSON(500, "服务器内部错误,请稍后再试")
	}
}

func (s *UserHttpServer) ValidatorErrorHandle(c *gin.Context, err error) {
	verr, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"msg": s.Server.TranslateErr(verr),
	})
}
