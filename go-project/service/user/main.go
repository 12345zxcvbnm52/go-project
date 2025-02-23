package main

import (
	proto "kenshop/proto/user"
	"kenshop/service/user/internal/resourse"
)

func main() {
	proto.RegisterUserServer(resourse.Server.Server, resourse.UserServer)
	resourse.Server.Serve()
}
