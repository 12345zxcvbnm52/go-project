package main

import (
	"fmt"
	gb "goods_web/global"
	_ "goods_web/initialize"
)

func main() {
	gb.Router.Run(fmt.Sprintf("%s:%d", gb.ServerConfig.Ip, gb.ServerConfig.Port))
}
