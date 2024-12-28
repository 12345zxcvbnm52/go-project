package global

import "github.com/gin-gonic/gin"

var Router *gin.Engine

// 这个是稳定的Router
var OrderRter *gin.RouterGroup
var CartRter *gin.RouterGroup
