package global

import "github.com/gin-gonic/gin"

var Router *gin.Engine

// 这个是稳定的Router
var UseRter *gin.RouterGroup

// 用于测试的Router
var TestRter *gin.RouterGroup
