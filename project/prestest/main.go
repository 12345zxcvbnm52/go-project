package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	c := gin.Default()
	c.LoadHTMLFiles("./test.html")
	c.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "test.html", gin.H{})
	})
	c.Run()
}
