package main

import (
	"kenshop/api/user/internal/resourse"
	proto "kenshop/proto/user"
)

func main() {
	proto.RegisterUserHTTPServer(resourse.Server.Server, resourse.UserServer)
	resourse.Server.Serve()
}
