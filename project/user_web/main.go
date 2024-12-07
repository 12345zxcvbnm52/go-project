package main

import (
	"fmt"
	gb "user_web/global"
	_ "user_web/initialize"
)

func main() {
	gb.Router.Run(fmt.Sprintf("%s:%d", gb.ServerConfig.Ip, gb.ServerConfig.Port))
}
