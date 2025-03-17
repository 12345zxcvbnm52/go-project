package main

import "kenshop/api/goods/internal/resourse"

// 用户服务
// @BasePath /
// @Description User management service API
// @Host NULL
// @Title User Service API
// @Version 1.0.0
func main() {
	resourse.GoodsServer.Execute()
}
