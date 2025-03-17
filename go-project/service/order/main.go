package main

import (
	proto "kenshop/proto/order"
	"kenshop/service/order/internal/resourse"
)

func main() {
	proto.RegisterOrderServer(resourse.Server.Server, resourse.OrderServer)
	resourse.Server.Serve()
}
