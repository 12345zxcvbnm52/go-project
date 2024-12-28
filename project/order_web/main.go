package main

import (
	"fmt"
	gb "order_web/global"
	_ "order_web/initialize"
)

func main() {
	gb.Router.Run(fmt.Sprintf("%s:%d", gb.ServerConfig.Ip, gb.ServerConfig.Port))
}
