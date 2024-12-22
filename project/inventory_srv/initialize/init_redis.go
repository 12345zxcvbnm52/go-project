package initialize

import (
	"fmt"
	gb "inventory_srv/global"
	"net"
	"runtime"
	"time"

	"github.com/go-redis/redis"
)

var DefaultRedisOpts *redis.Options

func InitRedis() {
	DefaultRedisOpts = &redis.Options{
		Network:  gb.ServerConfig.RedisConfig.Host,
		Addr:     fmt.Sprintf("%s:%d", gb.ServerConfig.RedisConfig.Host, gb.ServerConfig.RedisConfig.Port),
		Password: gb.ServerConfig.RedisConfig.Password,
		DB:       0,
		//PoolSize侧面表示了连通的socket数
		PoolSize: 4 * runtime.NumCPU(),
		//指定最少长期维持idle状态的连接数不少于该数量,这里选择一半的cpu数
		MinIdleConns: 2 * runtime.NumCPU(),
		DialTimeout:  4 * time.Second,
		//读超时的最大时间
		ReadTimeout: 3 * time.Second,
		//写超时
		WriteTimeout: 3 * time.Second,
		//所有连接都处在繁忙状态时,客户端等待可用连接的最大等待时长
		PoolTimeout: 4 * time.Second,
		//对所有连接的健康检查周期,若为-1则默认只在客户端请求时检查
		IdleCheckFrequency: 1 * time.Minute,
		//指定一个连接最长的空闲时间,超过则关闭连接,直至没有连接数或达到指定的MinidelConns
		//为-1时不进行空闲关闭,
		IdleTimeout: 5 * time.Minute,
		//一个连接最大允许存活的时间,无论是否空闲,为0时则永不关闭
		MaxConnAge: 0 * time.Second,
		//命令失败时最大重试次数
		MaxRetries: 1,
		//命令失败到下次尝试重试之间的最小时间间隔(即等等再重试)
		MinRetryBackoff: 8 * time.Microsecond,
		//命令失败到下次尝试重试之间的最大时间间隔
		MaxRetryBackoff: 512 * time.Microsecond,
		//自定义的连接函数
		Dialer: func() (conn net.Conn, err error) {
			dialer := &net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Minute,
			}
			return dialer.Dial("tcp", fmt.Sprintf("%s:%d", gb.ServerConfig.RedisConfig.Host, gb.ServerConfig.RedisConfig.Port))
		},
		//hook,当客户端执行命令需要从连接池获取连接且连接池需要新建连接时则会调用此钩子函数
		//OnConnect: ,
	}
	gb.RedisConn = redis.NewClient(DefaultRedisOpts)
}
