package main

import (
	proto "kenshop/proto/inventory"
	"kenshop/service/inventory/internal/resourse"
)

func main() {
	proto.RegisterInventoryServer(resourse.Server.Server, resourse.InventoryServer)
	resourse.Server.Serve()
}
