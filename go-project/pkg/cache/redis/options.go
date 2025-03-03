package redis

// import (
// 	"context"
// 	"crypto/tls"
// 	"net"
// 	"time"

// 	"github.com/redis/go-redis/v9"
// )

// // OptionFunc 定义一个函数类型，用于修改 RedisOptions
// type OptionFunc func(*RedisOptions)

// // NewUniversalOptions 创建 RedisOptions 实例，并应用所有传入的选项
// func MustNewRedisOptions(opts ...OptionFunc) *RedisOptions {
// 	options := &RedisOptions{
// 		Addrs:                 []string{"127.0.0.1:6379"},
// 		DB:                    0,
// 		MaxRetries:            3,
// 		MinRetryBackoff:       8 * time.Millisecond,
// 		MaxRetryBackoff:       512 * time.Millisecond,
// 		DialTimeout:           5 * time.Second,
// 		ReadTimeout:           3 * time.Second,
// 		WriteTimeout:          3 * time.Second,
// 		PoolSize:              10,
// 		PoolTimeout:           4 * time.Second,
// 		MinIdleConns:          0,
// 		MaxIdleConns:          0,
// 		MaxActiveConns:        0,
// 		ConnMaxIdleTime:       30 * time.Minute,
// 		ConnMaxLifetime:       0,
// 		RouteByLatency:        false,
// 		RouteRandomly:         false,
// 		ContextTimeoutEnabled: true,
// 	}

// 	// 应用所有传入的选项
// 	for _, opt := range opts {
// 		opt(options)
// 	}
// 	return options
// }

// // 选项函数列表
// func WithAddrs(addrs ...string) OptionFunc {
// 	return func(o *RedisOptions) { o.Addrs = addrs }
// }

// func WithClientName(name string) OptionFunc {
// 	return func(o *RedisOptions) { o.ClientName = name }
// }

// func WithDB(db int) OptionFunc {
// 	return func(o *RedisOptions) { o.DB = db }
// }

// func WithDialer(dialer func(ctx context.Context, network, addr string) (net.Conn, error)) OptionFunc {
// 	return func(o *RedisOptions) { o.Dialer = dialer }
// }

// func WithOnConnect(hook func(ctx context.Context, cn *redis.Conn) error) OptionFunc {
// 	return func(o *RedisOptions) { o.OnConnect = hook }
// }

// func WithProtocol(protocol int) OptionFunc {
// 	return func(o *RedisOptions) { o.Protocol = protocol }
// }

// func WithAuth(username, password string) OptionFunc {
// 	return func(o *RedisOptions) {
// 		o.Username = username
// 		o.Password = password
// 	}
// }

// func WithSentinelAuth(username, password string) OptionFunc {
// 	return func(o *RedisOptions) {
// 		o.SentinelUsername = username
// 		o.SentinelPassword = password
// 	}
// }

// func WithRetries(maxRetries int, minBackoff, maxBackoff time.Duration) OptionFunc {
// 	return func(o *RedisOptions) {
// 		o.MaxRetries = maxRetries
// 		o.MinRetryBackoff = minBackoff
// 		o.MaxRetryBackoff = maxBackoff
// 	}
// }

// func WithTimeouts(dial, read, write time.Duration) OptionFunc {
// 	return func(o *RedisOptions) {
// 		o.DialTimeout = dial
// 		o.ReadTimeout = read
// 		o.WriteTimeout = write
// 	}
// }

// func WithContextTimeoutEnabled(enabled bool) OptionFunc {
// 	return func(o *RedisOptions) { o.ContextTimeoutEnabled = enabled }
// }

// func WithPoolSettings(poolSize int, poolTimeout time.Duration, minIdle, maxIdle, maxActive int) OptionFunc {
// 	return func(o *RedisOptions) {
// 		o.PoolSize = poolSize
// 		o.PoolTimeout = poolTimeout
// 		o.MinIdleConns = minIdle
// 		o.MaxIdleConns = maxIdle
// 		o.MaxActiveConns = maxActive
// 	}
// }

// func WithPoolFIFO(enabled bool) OptionFunc {
// 	return func(o *RedisOptions) { o.PoolFIFO = enabled }
// }

// func WithConnLifetime(maxIdle, maxLifetime time.Duration) OptionFunc {
// 	return func(o *RedisOptions) {
// 		o.ConnMaxIdleTime = maxIdle
// 		o.ConnMaxLifetime = maxLifetime
// 	}
// }

// func WithTLSConfig(config *tls.Config) OptionFunc {
// 	return func(o *RedisOptions) { o.TLSConfig = config }
// }

// func WithClusterSettings(maxRedirects int, readOnly, routeByLatency, routeRandomly bool) OptionFunc {
// 	return func(o *RedisOptions) {
// 		o.MaxRedirects = maxRedirects
// 		o.ReadOnly = readOnly
// 		o.RouteByLatency = routeByLatency
// 		o.RouteRandomly = routeRandomly
// 	}
// }

// func WithMasterName(name string) OptionFunc {
// 	return func(o *RedisOptions) { o.MasterName = name }
// }

// func WithIdentitySettings(disableIdentity bool, suffix string) OptionFunc {
// 	return func(o *RedisOptions) {
// 		o.DisableIndentity = disableIdentity
// 		o.IdentitySuffix = suffix
// 	}
// }

// func WithUnstableResp3(enabled bool) OptionFunc {
// 	return func(o *RedisOptions) { o.UnstableResp3 = enabled }
// }
