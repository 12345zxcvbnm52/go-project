package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func main() {
	secureMid := secure.New(secure.Options{
		SSLRedirect: true,
		SSLHost:     "192.168.199.128:8080",
	})
	secureFunc := func() gin.HandlerFunc {
		return func(ctx *gin.Context) {
			if err := secureMid.Process(ctx.Writer, ctx.Request); err != nil {
				panic(err)
			}
			if status := ctx.Writer.Status(); status > 300 && status < 400 {
				panic(errors.New("err"))
			}
		}
	}()
	//app:=secureMid.Handler()
	r := gin.Default()
	r.Use(secureFunc)
	r.GET("/a", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "httpsOK"})
	})
	r.RunTLS("192.168.199.128:8080", "./cert.pem", "./private_key.pem")
}
