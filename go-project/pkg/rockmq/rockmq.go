package rockmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func MustNewProducer(namesrv []string, gropname string, opts ...producer.Option) rocketmq.Producer {
	opts = append(opts, producer.WithGroupName(gropname))
	opts = append(opts, producer.WithNameServer(namesrv))
	p, err := rocketmq.NewProducer(opts...)
	if err != nil {
		panic(err)
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	return p
}

func MustNewPullConsumer(namesrv []string, gropname string, opts ...consumer.Option) rocketmq.PullConsumer {
	opts = append(opts, consumer.WithGroupName(gropname))
	opts = append(opts, consumer.WithNameServer(namesrv))
	p, err := rocketmq.NewPullConsumer(opts...)
	if err != nil {
		panic(err)
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	return p
}

func MustNewPushConsumer(namesrv []string, gropname string, opts ...consumer.Option) rocketmq.PushConsumer {
	opts = append(opts, consumer.WithGroupName(gropname))
	opts = append(opts, consumer.WithNameServer(namesrv))
	p, err := rocketmq.NewPushConsumer(opts...)
	if err != nil {
		panic(err)
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	return p
}

type TransProducer struct {
	TransactionProducer rocketmq.TransactionProducer
	Listener            primitive.TransactionListener
}

func MustNewTransProducer(l primitive.TransactionListener,
	namesrv []string, gn string, opts ...producer.Option) *TransProducer {

	var err error
	t := &TransProducer{}
	t.Listener = l
	opts = append(opts, producer.WithNameServer(namesrv))
	opts = append(opts, producer.WithGroupName(gn))
	t.TransactionProducer, err = rocketmq.NewTransactionProducer(l, opts...)
	if err != nil {
		panic(err)
	}
	if err = t.TransactionProducer.Start(); err != nil {
		panic(err)
	}
	return t
}

// //便于事务消息之间的信息传递和错误传递
// type MessageQueuePeyload struct {
// 	//用于流程控制
// 	version int16
// 	//需要传递的信息
// 	Data  any
// 	//需要传递的错误
// 	Error error
// }

// type MessageQueueChan chan MessageQueuePeyload

func (t *TransProducer) SendMessageInTransaction(ctx context.Context, mq *primitive.Message) (*primitive.TransactionSendResult, error) {
	return t.TransactionProducer.SendMessageInTransaction(ctx, mq)
}

// func()
