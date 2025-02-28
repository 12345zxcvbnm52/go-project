package main

import (
	"context"
	"kenshop/goken/server/httpserver"
	proto "kenshop/proto/test"
)

// PathsMap中key为对应的handler函数名
// @BasePath /v1
// @Description This is a sample API
// @Host api.example.com
// @Title My API
// @Version 1.0.0
func main() {
	s := proto.RegisterMessagingHTTPServer()
	server := httpserver.MustNewServer(context.Background(), "0.0.0.0:8080", s)
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
