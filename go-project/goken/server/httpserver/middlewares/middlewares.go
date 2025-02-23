package middlewares

import "github.com/gin-gonic/gin"

var Middlewares map[string]func() gin.HandlerFunc = defaultMiddlewares()

func defaultMiddlewares() (md map[string]func() gin.HandlerFunc) {
	md = make(map[string]func() gin.HandlerFunc)
	md["recovery"] = gin.Recovery
	md["log"] = gin.Logger
	return md
}
