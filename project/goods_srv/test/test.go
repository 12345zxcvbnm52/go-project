package main

import (
	"context"
	"fmt"
	"goods_srv/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	cc, err := grpc.NewClient(fmt.Sprintf("consul://%s:%d/%s?healthy=true",
		"192.168.199.128",
		8500,
		"goods_srv",
	),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	c := proto.NewGoodsClient(cc)
	// r, err := c.GetSubCategy(context.Background(), &proto.SubCategyReq{Id: 1000})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(r.SelfInfo)
	d, err := c.GetBannerList(context.Background(), &emptypb.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(d)
}

// var DB *gorm.DB

// type User struct {
// 	Id   int    `gorm:"primarykey"`
// 	Item []Item `gorm:"many2many:user_item"`
// }

// type Item struct {
// 	Id int `gorm:"primarykey"`
// 	//Userw int
// 	User []User `gorm:"many2many:user_item"`
// }

// func main() {
// 	l := log.New(os.Stdout, "", log.LstdFlags)
// 	log := logger.New(l, logger.Config{
// 		LogLevel: logger.Info,
// 		Colorful: true,
// 	})
// 	DB, _ = gorm.Open(
// 		mysql.Open("ken:123@tcp(127.0.0.1:3306)/test01?charset=utf8mb4&parseTime=True&loc=Local"),
// 		&gorm.Config{
// 			NamingStrategy: schema.NamingStrategy{
// 				SingularTable: false,
// 			},
// 			Logger: log,
// 		},
// 	)
// 	DB.AutoMigrate(&User{}, &Item{})
// }
