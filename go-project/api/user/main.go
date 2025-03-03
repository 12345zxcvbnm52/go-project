package main

import (
	"kenshop/api/user/internal/resourse"
)

// 用户服务
// @BasePath /
// @Description User management service API
// @Host NULL
// @Title User Service API
// @Version 1.0.0
func main() {
	resourse.UserServer.Execute()
}
