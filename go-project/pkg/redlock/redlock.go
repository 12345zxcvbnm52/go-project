package redlock

import (
	"context"
	"fmt"
	"kenshop/pkg/errors"
	"kenshop/pkg/log"
	"time"

	"github.com/go-redsync/redsync/v4"
	redpool "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	redis "github.com/redis/go-redis/v9"
)

// 不确定要不要支持redis集群
type RedLock struct {
	addr     []string
	password string
	sync     *redsync.Redsync
	//用于检测addr,password是否改变过
	UseCluster bool
	//如果未抢到锁则休眠的时间,若为0则不休眠直接退出(不尝试自旋)
	SleepTime time.Duration
	//最大允许睡眠的时间
	MaxSleepTime time.Duration
}

var (
	ErrLockFailed   = errors.New("获取分布式锁失败")
	ErrLockTimeout  = errors.New("获取分布式锁超时")
	ErrUnlockFailed = errors.New("释放分布式锁失败")
)

type OptionFunc func(r *RedLock)

func MustNewRedLock(addr []string, opts ...OptionFunc) *RedLock {
	r := &RedLock{
		UseCluster:   false,
		SleepTime:    25 * time.Millisecond,
		MaxSleepTime: 8 * time.Second,
	}
	for _, opt := range opts {
		opt(r)
	}

	if r.UseCluster && len(addr) <= 1 {
		panic("提供的节点数不支持组成为集群")
	}
	r.addr = addr
	return r
}

func (r *RedLock) newPool(cluster bool) []redpool.Pool {
	lent := len(r.addr)
	if cluster {
		lent = 1
	}
	var pool []redpool.Pool = make([]redpool.Pool, lent)

	if cluster {
		cs := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    r.addr,
			Password: r.password,
		})
		pool[0] = goredis.NewPool(cs)

	} else {
		for i, v := range r.addr {
			cs := redis.NewClient(&redis.Options{
				Addr:     v,
				Password: r.password,
			})
			pool[i] = goredis.NewPool(cs)
		}
	}
	return pool
}

func (r *RedLock) NewMutex(name string, opts ...redsync.Option) *redsync.Mutex {
	if r.sync == nil {
		p := r.newPool(r.UseCluster)
		r.sync = redsync.New(p...)
	}
	return r.sync.NewMutex(name, opts...)
}

// 这个函数一旦得不到锁就立马返回
func (c *RedLock) GetRedLockAndLockFast(_ context.Context, key string) (*redsync.Mutex, error) {
	lockKey := fmt.Sprintf("%s-lock", key)
	lock := c.NewMutex(lockKey)
	if err := lock.Lock(); err != nil {
		log.Errorf("[redlock] 获取分布式锁失败 err = %v", err)
		return nil, ErrLockFailed
	}
	return lock, nil
}

func (c *RedLock) GetRedLockAndLock(_ context.Context, key string) (*redsync.Mutex, error) {
	lockKey := fmt.Sprintf("%s-lock", key)
	lock := c.NewMutex(lockKey)

	// 初次尝试获取锁
	if err := lock.Lock(); err == nil {
		return lock, nil
	}

	// 退避重试
	sleepTime := c.SleepTime
	timeout := time.After(c.MaxSleepTime)
	totalTimeSleep := time.Duration(0)

	for {
		select {
		case <-timeout:
			return nil, ErrLockTimeout
		default:
			remainingTime := c.MaxSleepTime - totalTimeSleep
			if remainingTime <= c.SleepTime {
				return nil, ErrLockTimeout
			}

			time.Sleep(sleepTime)
			totalTimeSleep += sleepTime

			if err := lock.Lock(); err == nil {
				return lock, nil
			}

			sleepTime *= 2
		}
	}
}

func (c *RedLock) UnlockRedLock(ctx context.Context, lock *redsync.Mutex) error {
	if _, err := lock.UnlockContext(ctx); err != nil {
		log.Fatalf("[redlock] 释放分布式锁失败 err = %v", err)
		return ErrUnlockFailed
	}
	return nil
}

func WithPassword(password string) OptionFunc {
	return func(r *RedLock) {
		r.password = password
	}
}

func WithCluster(useCluster bool) OptionFunc {
	return func(r *RedLock) {
		r.UseCluster = useCluster
	}
}

func WithSleepTime(st time.Duration) OptionFunc {
	return func(r *RedLock) {
		r.SleepTime = st
	}
}

func WithMaxSleepTime(mst time.Duration) OptionFunc {
	return func(r *RedLock) {
		r.MaxSleepTime = mst
	}
}
