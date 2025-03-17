package cache

import (
	"context"
	"encoding/json"
	"errors"
	"kenshop/pkg/log"
	"runtime"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/sourcegraph/conc/pool"
)

type messageQueueExtractor struct {
	Key     string `json:"key"`
	Version int64  `json:"version"`
}

func (m *messageQueueExtractor) GetVersion() int64 {
	return m.Version
}

func (m *messageQueueExtractor) GetKey() string {
	return m.Key
}

var (
	ErrTurnConsumerFailed = errors.New("更换消费者失败")
	ErrExtractFailed      = errors.New("无法获取version和key信息")
)

type CacheExtracter interface {
	GetVersion() int64
	GetKey() string
}

func extractRocketMessage(data []byte) (*messageQueueExtractor, error) {
	r := &messageQueueExtractor{}
	if err := json.Unmarshal(data, r); err != nil {
		return nil, err
	}
	return r, nil
}

func WrapMessageQueueExtractor(c CacheExtracter) []byte {
	b, _ := json.Marshal(c)
	return b
}

func (c *MultiCache) RegisterRocketmq(pc rocketmq.PushConsumer, topic string) error {
	if c.RocketmqConsumer != nil {
		if err := c.RocketmqConsumer.Shutdown(); err != nil {
			log.Errorf("无法关闭欲更换的rocketmq consumer err = %v", err)
			return ErrTurnConsumerFailed
		}
		c.RocketmqConsumer = pc
	}
	pc.Subscribe(topic, consumer.MessageSelector{}, c.updateConsume)
	c.RocketmqConsumer = pc
	return nil
}

func (c *MultiCache) updateConsume(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	msgLen := len(me)
	if msgLen > runtime.NumCPU()*50 {
		msgLen = runtime.NumCPU() * 50
	}

	ep := pool.New().WithMaxGoroutines(msgLen).WithErrors()
	for _, msg := range me {
		ep.Go(
			func() error {
				r, err := extractRocketMessage(msg.Body)
				if err != nil {
					log.Errorf("[multi cache] 无法从rocketmq message中获取信息 err = %v", err)
					return err
				}
				return c.SetWithVersionMutex(ctx, r.Key, msg.Body, r.Version)
			},
		)
	}
	if err := ep.Wait(); err != nil {
		log.Error(err.Error())
		return consumer.ConsumeRetryLater, err
	}
	return consumer.ConsumeSuccess, nil
}

// 结合版本控制和分布式锁机制,注意value参数中应当也要有version字段,多出来的version参数则是避免再次unmarshal得到
func (c *MultiCache) SetWithVersionMutex(ctx context.Context, key string, value []byte, version int64) error {
	lock, err := c.redlock.GetRedLockAndLock(ctx, key)
	if err != nil {
		return err
	}

	//先基于版本控制检查,此时redis查不到数据可以直接更新(err为Nil),
	// 如果是内部原因不可用(其它错误)在之后也会检查,不用在此重复检查
	if res := c.distributedCache.Get(ctx, key); res.Err() == nil {
		data, _ := res.Bytes()
		if c.getVersion(ctx, data) >= version {
			//因为redis里存在的数据版本已经新于要更新的数据,故直接跳过即可
			return nil
		}
	}

	if err := c.distributedCache.Set(ctx, key, value, c.ExpireTime).Err(); err != nil {
		log.Errorf("[multi-cache] 从distributed cache中获取数据失败 err= %v", err)
		if err := c.redlock.UnlockRedLock(ctx, lock); err != nil {
			log.Errorf("[multi-cache] 释放悲观锁失败 err = %v", err)
		}
		return ErrBadMultiCache
	}

	c.localCache.Set(key, value)

	if err := c.redlock.UnlockRedLock(ctx, lock); err != nil {
		log.Errorf("[multi-cache] 释放悲观锁失败 err = %v", err)
	}
	return nil
}
