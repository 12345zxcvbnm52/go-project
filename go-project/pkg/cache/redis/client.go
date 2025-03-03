package redis

// import (
// 	"context"
// 	"fmt"
// 	"sync/atomic"
// 	"time"

// 	"kenshop/pkg/errors"

// 	uuid "github.com/google/uuid"
// 	"github.com/redis/go-redis/extra/redisotel/v9"
// 	redis "github.com/redis/go-redis/v9"

// 	"kenshop/pkg/log"
// )

// // ErrRedisIsDown is returned when we can't communicate with redis.
// var ErrRedisIsDown = errors.New("storage: Redis is either down or ws not configured")

// var disableRedis atomic.Value

// type RedisOptions redis.UniversalOptions

// // DisableRedis very handy when testsing it allows to dynamically enable/disable talking with redisW.
// func DisableRedis(ok bool) {
// 	if ok {
// 		redisUp.Store(false)
// 		disableRedis.Store(true)

// 		return
// 	}
// 	redisUp.Store(true)
// 	disableRedis.Store(false)
// }

// func shouldConnect() bool {
// 	if v := disableRedis.Load(); v != nil {
// 		return !v.(bool)
// 	}

// 	return true
// }

// // Connected returns true if we are connected to redis.
// func Connected() bool {
// 	if v := redisUp.Load(); v != nil {
// 		return v.(bool)
// 	}

// 	return false
// }

// // 根据传入的cache获得不同的Pool,但是注意该函数不会初始化连接池
// func getClientPool(cache bool) redis.UniversalClient {
// 	if cache {
// 		v := singleCachePool.Load()
// 		if v != nil {
// 			return v.(redis.UniversalClient)
// 		}

// 		return nil
// 	}
// 	if v := singlePool.Load(); v != nil {
// 		return v.(redis.UniversalClient)
// 	}
// 	return nil
// }

// // 建立到redis的Singleton Client,注意此时并没有建立连接,只有第一次调用函数时才开始建立连接
// func buildRedisClient(cache bool, opts ...OptionFunc) bool {
// 	if getClientPool(cache) == nil {
// 		log.Info("开始建立到Redis的Singleton Client")
// 		if cache {
// 			singleCachePool.Store(NewRedisClusterPool(true, opts...))
// 			return true
// 		}
// 		singlePool.Store(NewRedisClusterPool(true, opts...))
// 		return true
// 	}

// 	return true
// }

// // RedisCluster is a storage manager that uses the redis database.
// type RedisClient struct {
// 	//设置一个redisCluster实例的前缀key,用于标识
// 	KeyPrefix string
// 	//似乎是用hash存储key而非直接用字符串?
// 	HashKeys bool

// }

// func redisPoolHealthCheck(ctx context.Context, cache bool) error {
// 	c := getClientPool(cache)
// 	testKey := "health-test-" + uuid.Must(uuid.NewV7()).String()
// 	if err := c.Set(ctx, testKey, "health-test", time.Second).Err(); err != nil {
// 		return errors.New(fmt.Sprintf("尝试连接Redis并设置key失败 err=%s", err.Error()))
// 	}
// 	if _, err := c.Get(ctx, testKey).Result(); err != nil {
// 		return errors.New(fmt.Sprintf("尝试连接Redis并获取key失败 err=%s", err.Error()))
// 	}
// 	return nil
// }

// // ConnectToRedis会启动一个协程定期尝试连接到redis
// func ConnectToRedis(ctx context.Context, opts ...OptionFunc) {
// 	//通过计时器定期
// 	tick := time.NewTicker(time.Second)
// 	defer tick.Stop()
// 	//两个空实例用于初始化两个不同的连接池

// 	var ok bool = true

// 	buildRedisClient(true, opts...)
// 	buildRedisClient(false, opts...)

// 	if err := redisPoolHealthCheck(ctx, true); err != nil {
// 		redisUp.Store(false)
// 		log.Errorf("cache redis pool ", err.Error())
// 	}
// 	if err := redisPoolHealthCheck(ctx, false); err != nil {
// 		redisUp.Store(false)
// 		log.Errorf("nomal redis pool ", err.Error())
// 	}
// 	redisUp.Store(ok)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-tick.C:
// 			if !shouldConnect() {
// 				continue
// 			}
// 			buildRedisClient(true, opts...)
// 			buildRedisClient(false, opts...)

// 			if err := redisPoolHealthCheck(ctx, true); err != nil {
// 				redisUp.Store(false)
// 				log.Errorf("cache redis pool ", err.Error())
// 				continue
// 			}
// 			if err := redisPoolHealthCheck(ctx, false); err != nil {
// 				redisUp.Store(false)
// 				log.Errorf("nomal redis pool ", err.Error())
// 				continue
// 			}
// 			redisUp.Store(true)
// 		}
// 	}
// }

// // NewRedisClusterPool 创建一个Redis集群连接池
// func MustNewRedisClient(enableOtel bool, opts ...OptionFunc) redis.UniversalClient {
// 	config := MustNewRedisOptions(opts...)

// 	var client redis.UniversalClient

// 	if config.MasterName != "" {
// 		if len(config.Addrs) > 1 {
// 		} else {
// 			log.Info("[redis] 新建一个 singleton sentinel client")
// 			client = redis.NewFailoverClient(config.failover())
// 		}
// 	} else if len(config.Addrs) > 1 {
// 		log.Info("[redis] 新建一个 cluster client")
// 		client = redis.NewClusterClient(config.cluster())

// 	} else {
// 		log.Info("[redis] 新建一个 single-node client")
// 		client = redis.NewClient(config.simple())
// 	}

// 	if enableOtel {
// 		err := redisotel.InstrumentTracing(client)
// 		if err != nil {
// 			log.Errorf("Error instrumenting redis tracing: %s", err.Error())
// 		}
// 	}
// 	return client
// }

// // RedisOpts is the overridden type of redis.UniversalOptions. simple() and cluster() functions are not public in redis
// // library.
// // Therefore, they are redefined in here to use in creation of new redis cluster logic.
// // We don't want to use redis.NewUniversalClient() logic.
// func (o *RedisOptions) cluster() *redis.ClusterOptions {
// 	if len(o.Addrs) == 0 {
// 		o.Addrs = []string{"127.0.0.1:6379"}
// 	}

// 	return &redis.ClusterOptions{
// 		Addrs:      o.Addrs,
// 		ClientName: o.ClientName,
// 		Dialer:     o.Dialer,
// 		OnConnect:  o.OnConnect,

// 		Protocol: o.Protocol,
// 		Username: o.Username,
// 		Password: o.Password,

// 		MaxRedirects:   o.MaxRedirects,
// 		ReadOnly:       o.ReadOnly,
// 		RouteByLatency: o.RouteByLatency,
// 		RouteRandomly:  o.RouteRandomly,

// 		MaxRetries:      o.MaxRetries,
// 		MinRetryBackoff: o.MinRetryBackoff,
// 		MaxRetryBackoff: o.MaxRetryBackoff,

// 		DialTimeout:           o.DialTimeout,
// 		ReadTimeout:           o.ReadTimeout,
// 		WriteTimeout:          o.WriteTimeout,
// 		ContextTimeoutEnabled: o.ContextTimeoutEnabled,

// 		PoolFIFO: o.PoolFIFO,

// 		PoolSize:        o.PoolSize,
// 		PoolTimeout:     o.PoolTimeout,
// 		MinIdleConns:    o.MinIdleConns,
// 		MaxIdleConns:    o.MaxIdleConns,
// 		MaxActiveConns:  o.MaxActiveConns,
// 		ConnMaxIdleTime: o.ConnMaxIdleTime,
// 		ConnMaxLifetime: o.ConnMaxLifetime,

// 		TLSConfig: o.TLSConfig,

// 		DisableIndentity: o.DisableIndentity,
// 		IdentitySuffix:   o.IdentitySuffix,
// 		UnstableResp3:    o.UnstableResp3,
// 	}
// }

// func (o *RedisOptions) simple() *redis.Options {
// 	addr := "127.0.0.1:6379"
// 	if len(o.Addrs) > 0 {
// 		addr = o.Addrs[0]
// 	}

// 	return &redis.Options{
// 		Addr:       addr,
// 		ClientName: o.ClientName,
// 		Dialer:     o.Dialer,
// 		OnConnect:  o.OnConnect,

// 		DB:       o.DB,
// 		Protocol: o.Protocol,
// 		Username: o.Username,
// 		Password: o.Password,

// 		MaxRetries:      o.MaxRetries,
// 		MinRetryBackoff: o.MinRetryBackoff,
// 		MaxRetryBackoff: o.MaxRetryBackoff,

// 		DialTimeout:           o.DialTimeout,
// 		ReadTimeout:           o.ReadTimeout,
// 		WriteTimeout:          o.WriteTimeout,
// 		ContextTimeoutEnabled: o.ContextTimeoutEnabled,

// 		PoolFIFO:        o.PoolFIFO,
// 		PoolSize:        o.PoolSize,
// 		PoolTimeout:     o.PoolTimeout,
// 		MinIdleConns:    o.MinIdleConns,
// 		MaxIdleConns:    o.MaxIdleConns,
// 		MaxActiveConns:  o.MaxActiveConns,
// 		ConnMaxIdleTime: o.ConnMaxIdleTime,
// 		ConnMaxLifetime: o.ConnMaxLifetime,

// 		TLSConfig: o.TLSConfig,

// 		DisableIndentity: o.DisableIndentity,
// 		IdentitySuffix:   o.IdentitySuffix,
// 		UnstableResp3:    o.UnstableResp3,
// 	}
// }

// func (o *RedisOptions) failover() *redis.FailoverOptions {
// 	if len(o.Addrs) == 0 {
// 		o.Addrs = []string{"127.0.0.1:16379"}
// 	}

// 	return &redis.FailoverOptions{
// 		SentinelAddrs: o.Addrs,
// 		MasterName:    o.MasterName,
// 		ClientName:    o.ClientName,

// 		Dialer:    o.Dialer,
// 		OnConnect: o.OnConnect,

// 		DB:               o.DB,
// 		Protocol:         o.Protocol,
// 		Username:         o.Username,
// 		Password:         o.Password,
// 		SentinelUsername: o.SentinelUsername,
// 		SentinelPassword: o.SentinelPassword,

// 		MaxRetries:      o.MaxRetries,
// 		MinRetryBackoff: o.MinRetryBackoff,
// 		MaxRetryBackoff: o.MaxRetryBackoff,

// 		DialTimeout:           o.DialTimeout,
// 		ReadTimeout:           o.ReadTimeout,
// 		WriteTimeout:          o.WriteTimeout,
// 		ContextTimeoutEnabled: o.ContextTimeoutEnabled,

// 		PoolFIFO:        o.PoolFIFO,
// 		PoolSize:        o.PoolSize,
// 		PoolTimeout:     o.PoolTimeout,
// 		MinIdleConns:    o.MinIdleConns,
// 		MaxIdleConns:    o.MaxIdleConns,
// 		MaxActiveConns:  o.MaxActiveConns,
// 		ConnMaxIdleTime: o.ConnMaxIdleTime,
// 		ConnMaxLifetime: o.ConnMaxLifetime,

// 		TLSConfig: o.TLSConfig,

// 		DisableIndentity: o.DisableIndentity,
// 		IdentitySuffix:   o.IdentitySuffix,
// 		UnstableResp3:    o.UnstableResp3,
// 	}
// }
