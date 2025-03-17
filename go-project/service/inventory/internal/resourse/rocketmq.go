package resourse

import (
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
)

func InitRocketmq() {
	var err error
	Consumer, err = rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{fmt.Sprintf("%s:%d", Conf.Rocketmq.Ip, Conf.Rocketmq.Port)}),
		consumer.WithGroupName(Conf.Rocketmq.ConsumerGroupName),
	)
	if err != nil {
		panic(err)
	}
	if err = Consumer.Start(); err != nil {
		panic(err)
	}
}
