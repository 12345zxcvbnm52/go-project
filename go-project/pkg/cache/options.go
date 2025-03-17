package cache

import (
	"context"
	"kenshop/pkg/redlock"
	"time"

	"github.com/allegro/bigcache"
	"github.com/redis/go-redis/v9"
)

type OptionFunc func(*MultiCache)

func MustNewLocalCache(conf *bigcache.Config) *LocalCache {
	cache, err := bigcache.NewBigCache(*conf)
	if err != nil {
		panic(err)
	}
	return &LocalCache{
		BigCache: cache,
		conf:     conf,
	}
}

func MustNewDistributedCache(addrs []string, conf *redis.ClusterOptions) *DistributedCache {
	if len(addrs) <= 1 {
		panic("请直接使用单机redis")
	} else {
		cache := redis.NewClusterClient(conf)
		return &DistributedCache{
			ClusterClient: cache,
		}
	}
}

// 默认会将使用的分布式缓存作为分布式锁
func MustNewMultiCache(ctx context.Context, dc *DistributedCache, lc *LocalCache, opts ...OptionFunc) *MultiCache {
	c := &MultiCache{
		localCache:         lc,
		distributedCache:   dc,
		Ctx:                ctx,
		ExpireTime:         10 * time.Minute,
		UseVersionControll: true,
		EnableTracing:      false,
	}
	for _, opt := range opts {
		opt(c)
	}

	if lc.conf.LifeWindow >= c.ExpireTime {
		panic(ErrBadExpireTime)
	}

	if c.redlock == nil {
		c.redlock = redlock.MustNewRedLock(
			dc.Options().Addrs,
			redlock.WithCluster(true),
			redlock.WithMaxSleepTime(5*time.Second),
			redlock.WithPassword(dc.Options().Password),
			redlock.WithSleepTime(10*time.Millisecond),
		)
	}
	return c
}

func WithRedlock(r *redlock.RedLock) OptionFunc {
	return func(mc *MultiCache) {
		mc.redlock = r
	}
}

func WithExpireTime(t time.Duration) OptionFunc {
	return func(mc *MultiCache) {
		mc.ExpireTime = t
	}
}
