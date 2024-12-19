package main

import (
	"context"
	"fmt"
	pb "user_srv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// u := &pb.WriteUserReq{
	// 	UserName: "ken",
	// 	Password: "zxcvbnm52",
	// 	Gender:   "boy",
	// 	Mobile:   "18174637586",
	// 	Role:     0,
	// 	Birth:    1633072800,
	// }
	cc, _ := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	c := pb.NewUserClient(cc)
	//c.CreateUser(context.Background(), u)
	flag, _ := c.CheckUserRole(context.Background(), &pb.UserPasswordReq{
		Id:       1,
		UserName: "ken",
		Password: "zxcvbnm52",
	})
	fmt.Println(flag)
}
