package middlewares

import "github.com/gin-gonic/gin"

var Middlewares map[string]gin.HandlerFunc = defaultMiddlewares()

func defaultMiddlewares() (md map[string]gin.HandlerFunc) {
	md = make(map[string]gin.HandlerFunc)
	md["recovery"] = gin.Recovery()
	md["log"] = gin.Logger()
	return md
}

func CopyDefaultMiddlewares(w map[string]gin.HandlerFunc) {
	for k, v := range Middlewares {
		w[k] = v
	}
}
