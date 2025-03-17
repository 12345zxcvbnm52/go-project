package main

import (
	"context"
	"fmt"
	"kenshop/pkg/cache"
	"kenshop/pkg/log"
	"kenshop/pkg/redlock"
	"kenshop/pkg/rockmq"
	"sync"
	"time"

	"github.com/allegro/bigcache"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/redis/go-redis/v9"
)

type OrderListener struct {
	Ctx context.Context
}

func (ol *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	return primitive.CommitMessageState
}

func (s *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Println(msg.Body)
	msg.WithProperty("error", "badrequest")
	return primitive.RollbackMessageState
}

func BigCache() {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		panic(err)
	}
	fmt.Println(cache.Capacity())
	w := []byte{'w', '2'}
	cache.Set("ken", w)
	fmt.Println(cache.Get("ken"))
}

// func Sentinel() {
// 	if err := sentinel.InitDefault(); err != nil {
// 		panic(err)
// 	}
// 	_, err := flow.LoadRules([]*flow.Rule{
// 		{
// 			Resource:               "ken",
// 			TokenCalculateStrategy: flow.WarmUp,
// 			ControlBehavior:        flow.Reject,
// 			Threshold:              10000,
// 			StatIntervalInMs:       1000,
// 			WarmUpPeriodSec:        60,
// 		},
// 	})
// }

func Cache() {
	d := bigcache.DefaultConfig(5 * time.Minute)

	addr := []string{
		"192.168.199.128:6380",
		"192.168.199.128:6381",
		"192.168.199.128:6382",
		"192.168.199.128:6383",
		"192.168.199.128:6384",
		"192.168.199.128:6379",
	}
	opts := &redis.ClusterOptions{}
	opts.Password = "123"
	opts.Addrs = addr

	cache := cache.MustNewMultiCache(context.Background(),
		cache.MustNewDistributedCache(addr, opts),
		cache.MustNewLocalCache(&d),
	)
	fmt.Println("begin")
	t1 := time.Now()
	f := sync.WaitGroup{}
	f.Add(500)
	i := 0
	for range 500 {
		go func() {
			l := redlock.MustNewRedLock(addr, redlock.WithCluster(true), redlock.WithPassword("123"))
			lock, err := l.GetRedLockAndLock(context.TODO(), "ken")
			if err != nil {
				log.Error(err.Error())
			} else {
				i = i + 1
				//time.Sleep(30 * time.Millisecond)
			}
			if err := l.UnlockRedLock(context.TODO(), lock); err != nil {
				log.Error(err.Error())
			}
			f.Done()
		}()
	}

	f.Wait()
	t2 := time.Now()
	fmt.Println(t2.Sub(t1))
	fmt.Println(i)
	cache.Ctx.Deadline()
	// cache.SetWithMutex(context.Background(), "ken", []byte("1235s"))
	// data, err := cache.Get(context.Background(), "ken")
	// fmt.Println(string(data), err)
	// fmt.Println(cache.Stats())
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Version  int64  `json:"version"`
	Key      string `json:"key"`
}

func (m *User) GetVersion() int64 {
	return m.Version
}

func (m *User) GetKey() string {
	return m.Key
}

func RocketCache() {
	d := bigcache.DefaultConfig(5 * time.Minute)

	addr := []string{
		"192.168.199.128:6380",
		"192.168.199.128:6381",
		"192.168.199.128:6382",
		"192.168.199.128:6383",
		"192.168.199.128:6384",
		"192.168.199.128:6379",
	}
	opts := &redis.ClusterOptions{}
	opts.Password = "123"
	opts.Addrs = addr

	ca := cache.MustNewMultiCache(context.Background(),
		cache.MustNewDistributedCache(addr, opts),
		cache.MustNewLocalCache(&d),
	)

	// w1 := &User{
	// 	Password: "12345678",
	// 	Username: "kensame",
	// 	Version:  2,
	// 	Key:      "ken",
	// }
	// w2 := &User{
	// 	Password: "zxcvbnm52",
	// 	Username: "kensame42",
	// 	Version:  1,
	// 	Key:      "ken",
	// }
	//req1 := cache.WrapMessageQueueExtractor(w1)
	//req2 := cache.WrapMessageQueueExtractor(w2)
	rlog.SetLogLevel("warn")
	//pd := rockmq.MustNewProducer([]string{"192.168.199.128:9876"}, "name1")
	cs := rockmq.MustNewPushConsumer([]string{"192.168.199.128:9876"}, "name2")
	if err := ca.RegisterRocketmq(cs, "cache"); err != nil {
		panic(err)
	}

	// msg1 := primitive.NewMessage("cache", req1)
	// pd.SendSync(context.Background(), msg1)
	// time.Sleep(500 * time.Millisecond)
	// msg2 := primitive.NewMessage("cache", req2)
	// pd.SendSync(context.Background(), msg2)
	// time.Sleep(8 * time.Second)
	// b, err := ca.GetWithMutex(context.Background(), "ken")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(b))
	time.Sleep(10 * time.Second)
}

func main() {
	Cache()
	//RocketCache()
}
