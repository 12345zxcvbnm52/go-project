package cache

import (
	"context"
	"errors"
	"kenshop/pkg/log"
	"kenshop/pkg/redlock"
	"time"

	"github.com/allegro/bigcache"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/redis/go-redis/v9"
)

//注意,考虑到操作本地缓存不会出现未定义的,无法预知的错误,该包将更新分布式缓存和本地缓存视为一个事务
//即不存在更新成功本地缓存却未更新分布式缓存这种情况,且该包保证分布式缓存默认超时时间一定不比本地缓存短

var (
	ErrRecordNotFound  = errors.New("缓存记录不存在")
	ErrBadExpireTime   = errors.New("分布式缓存的过期时间不得早于本地缓存")
	ErrDeleteKeyFailed = errors.New("key删除失败")
	ErrBadMultiCache   = errors.New("分布式锁暂不可用")
	ErrDirtyData       = errors.New("读入的数据为脏数据")
)

type LocalCache struct {
	*bigcache.BigCache
	conf *bigcache.Config
}

type DistributedCache struct {
	*redis.ClusterClient
}

type version struct {
	Version int64 `json:"version"`
}

var v *version = &version{}

type MultiCache struct {
	localCache       *LocalCache
	distributedCache *DistributedCache
	Ctx              context.Context
	EnableTracing    bool
	//该Expire只能为distributedCache设置,localCache一旦设置之后无法改变
	ExpireTime time.Duration
	redlock    *redlock.RedLock
	//使用版本控制将带来更高的一致性,但是需要Set系列的函数传入的data是可以得到version字段的
	//暂时只支持json
	UseVersionControll bool
	RocketmqConsumer   rocketmq.PushConsumer
}

func (c *MultiCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.localCache.Get(key)
	if err != nil {
		log.Warnf("[multi-cache] %s未命中本地缓存", key)
	} else {
		return data, nil
	}

	res := c.distributedCache.Get(ctx, key)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			log.Warnf("[multi-cache] %s未命中分布式缓存", key)
			return nil, ErrRecordNotFound
		}
		log.Errorf("[multi-cache] 从distributed cache中获取数据失败 err= %v", err)
		return nil, ErrBadMultiCache
	}

	//因为整个multiCache都只支持存[]byte,故不用担心这里出错
	data, _ = res.Bytes()
	//尝试获取锁并更新到本地缓存中
	lock, err := c.redlock.GetRedLockAndLockFast(ctx, key)
	//如果拿不到锁就直接返回
	if err != nil {
		return data, nil
	}
	//先获取一次本地缓存,如果获取成功了就一定说明在自己到分布式缓存获取数据时已有进程修改过cache
	oldData, err := c.localCache.Get(key)
	//此时就直接把新的值返回即可,也没必要检查锁是否释放成功
	if err == nil {
		c.redlock.UnlockRedLock(ctx, lock)
		return oldData, nil
	}
	//如果此间没有人更新就回写到本地缓存
	c.localCache.Set(key, data)

	//没必要为了更新本地缓存的内出现锁失效问题进行rollback,如果真出了问题就删除本地缓存
	if err := c.redlock.UnlockRedLock(ctx, lock); err != nil {
		c.localCache.Delete(key)
	}
	return data, nil
}

func (c *MultiCache) GetWithMutex(ctx context.Context, key string) ([]byte, error) {
	lock, err := c.redlock.GetRedLockAndLock(ctx, key)
	if err != nil {
		return nil, err
	}

	data, err := c.localCache.Get(key)
	if err != nil {
		log.Warnf("[multi-cache] %s未命中本地缓存", key)
	} else {
		//直接释放即可,出错也不影响
		c.redlock.UnlockRedLock(ctx, lock)
		return data, nil
	}

	res := c.distributedCache.Get(ctx, key)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			log.Warnf("[multi-cache] %s未命中分布式缓存", key)
			return nil, ErrRecordNotFound
		}
		log.Errorf("[multi-cache] 从distributed cache中获取数据失败 err= %v", err)
		return nil, ErrBadMultiCache
	}

	//因为整个multiCache都只支持存[]byte,故不用担心这里出错
	data, _ = res.Bytes()
	c.localCache.Set(key, data)
	c.redlock.UnlockRedLock(ctx, lock)
	return data, nil
}

// 本包不建议在希望最终一致性的情况下直接使用该Set函数
func (c *MultiCache) Set(ctx context.Context, key string, value []byte) error {
	if err := c.distributedCache.Set(ctx, key, value, c.ExpireTime); err != nil {
		log.Errorf("[multi-cache] 从distributed cache中获取数据失败 err= %v", err)
		return ErrBadMultiCache
	}
	c.localCache.Set(key, value)
	return nil
}

func (c *MultiCache) SetWithMutex(ctx context.Context, key string, value []byte) error {
	lock, err := c.redlock.GetRedLockAndLock(ctx, key)
	if err != nil {
		return err
	}

	if err := c.distributedCache.Set(ctx, key, value, c.ExpireTime).Err(); err != nil {
		log.Errorf("[multi-cache] 从distributed cache中获取数据失败 err= %v", err)
		if err := c.redlock.UnlockRedLock(ctx, lock); err != nil {
			log.Errorf("[multi-cache] 释放悲观锁失败 err = %v", err)
		}
		return ErrBadMultiCache
	}

	//这里出错的唯一原因是需要存储的对象大小超过了整个分区所能承受的最大限度,
	//考虑到很难出现这个错误,暂时可以忽视
	c.localCache.Set(key, value)

	//在不使用版本控制的情况下(无法考虑哪个更新信息是最先的,一致性稍弱)锁出错就删掉本地缓存即可
	//这样如果更新成功了不会有大的影响,如果更新失败也可以通过错误把失败逻辑交给使用者手上
	if err := c.redlock.UnlockRedLock(ctx, lock); err != nil {
		c.localCache.Delete(key)
		return err
	}
	return nil
}

// 本包不建议在希望最终一致性的情况下直接使用该Delete函数
func (c *MultiCache) Delete(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		//无论是否存在这个key,删除都没关系(幂等的),故不处理error
		c.localCache.Delete(key)
	}

	if err := c.distributedCache.Del(ctx, keys...).Err(); err != nil {
		log.Errorf("[multi-cache] 从distributed cache中删除数据失败 err= %v", err)
		return ErrBadMultiCache
	}
	return nil
}

func (c *MultiCache) DeleteWithMutex(ctx context.Context, key string) error {
	lock, err := c.redlock.GetRedLockAndLock(ctx, key)
	if err != nil {
		return err
	}

	//这里直接删除本地缓存,可以解决后续因为锁释放失败导致不清楚多级缓存间的删除关系问题
	c.localCache.Delete(key)

	if err := c.distributedCache.Del(ctx, key).Err(); err != nil {
		log.Errorf("[multi-cache] 从distributed cache中删除数据失败 err= %v", err)
		if err := c.redlock.UnlockRedLock(ctx, lock); err != nil {
			log.Errorf("[multi-cache] 释放悲观锁失败 err = %v", err)
		}
		return ErrBadMultiCache
	}

	//因为之前已经删除了本地缓存,哪怕最终的Delete操作是否有效都不影响(没有效则可以通过分布式缓存重新获取数据)
	if err := c.redlock.UnlockRedLock(ctx, lock); err != nil {
		log.Errorf("[multi-cache] 释放悲观锁失败 err = %v", err)
		return err
	}
	return nil
}
