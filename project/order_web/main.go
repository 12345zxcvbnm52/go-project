package main

import (
	"fmt"
	gb "order_web/global"
	_ "order_web/initialize"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
)

func test() {
	if err:=sentinel.InitDefault();err!=nil{
		panic(err)
	}
	_,err:=flow.LoadRules([]*flow.Rule{
		{Resource:},
	})
}

func main() {
	gb.Router.Run(fmt.Sprintf("%s:%d", gb.ServerConfig.Ip, gb.ServerConfig.Port))
}
