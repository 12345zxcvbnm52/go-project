package global

import "github.com/gin-gonic/gin"

var Router *gin.Engine

// 这个是稳定的Router
var GoodsRter *gin.RouterGroup
var BannerRter *gin.RouterGroup
var CategoryRter *gin.RouterGroup
var BrandRter *gin.RouterGroup
var CategyWithBrandRter *gin.RouterGroup
