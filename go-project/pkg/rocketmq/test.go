package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.199.128:9876"}))
	if err != nil {
		panic(err)
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	defer p.Shutdown()
	res, err := p.SendSync(context.Background(), primitive.NewMessage("ken", []byte("hwqqwqwgo")))
	if err != nil {
		panic(err)
	}
	fmt.Println(res.String())
}
