package main

import (
	proto "kenshop/proto/goods"
	"kenshop/service/goods/internal/resourse"
)

func main() {
	proto.RegisterGoodsServer(resourse.Server.Server, resourse.GoodsServer)
	resourse.Server.Serve()
}
