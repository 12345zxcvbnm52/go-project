package main

import (
	_ "user_web/initialize"

	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver"
)

func main() {
	// conn, err := util.Dial()
	// if err != nil {
	// 	panic(err)
	// }
	// c := proto.NewUserClient(conn)
	// a, err := c.CheckUserRole(context.Background(), &proto.UserPasswordReq{Password: "zxcvbnm52"})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(a.Ok)
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		var s string
		ctx.BindJSON(s)

	})
	r.Run(":8080")
}
