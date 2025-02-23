package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func send() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.199.128:9876"}))
	if err != nil {
		panic(err)
	}
	p.Start()
	msg := primitive.NewMessage("hello", []byte("this is a message"))
	msg.WithDelayTimeLevel(3)
	res, err := p.SendSync(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.String())
	p.Shutdown()
}

func receive() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.199.128:9876"}),
		consumer.WithGroupName("ken"),
	)
	c.Subscribe("hello", consumer.MessageSelector{}, func(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range me {
			fmt.Println(me[i])
		}
		return consumer.ConsumeSuccess, nil
	})
	c.Start()
	time.Sleep(1 * time.Hour)
	c.Shutdown()
}

func main() {
	send()
	receive()
}
