package util

import (
	"errors"

	"github.com/go-redsync/redsync/v4"
	redpool "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	redis "github.com/redis/go-redis/v9"
)

type RedLock struct {
	addr     []string
	password string
	sync     *redsync.Redsync
	//用于检测addr,password是否改变过
	flag bool
}

// // 不推荐吗?
// func (r *RedLock) newClusterClient() *redis.ClusterClient {
// 	client := redis.NewClusterClient(&redis.ClusterOptions{
// 		Addrs:    r.addr,
// 		Password: r.password,
// 	})
// 	return client
// }

func (r *RedLock) newPool() []redpool.Pool {
	var pool []redpool.Pool = make([]redpool.Pool, len(r.addr))
	for i, v := range r.addr {
		cs := redis.NewClient(&redis.Options{
			Addr:     v,
			Password: r.password,
		})
		pool[i] = goredis.NewPool(cs)
	}
	return pool
}

func (r *RedLock) AppendAddr(addr ...string) {
	if r.addr == nil {
		r.addr = make([]string, len(addr))
	}
	r.addr = append(r.addr, addr...)
	r.flag = true
}

func (r *RedLock) SetAddr(addr ...string) {
	if r.addr == nil {
		r.addr = make([]string, len(addr))
	}
	copy(r.addr, addr)
	r.flag = true
}

func (r *RedLock) SetPassword(pass string) {
	r.password = pass
}

func (r *RedLock) NewRedMutexDebug(name string) (*redsync.Mutex, error) {
	if r.sync == nil || r.flag {
		if len(r.addr) == 0 {
			return nil, errors.New("没有与redlock关联的redis节点")
		}
		p := r.newPool()
		r.sync = redsync.New(p...)
		r.flag = false
	}
	return r.sync.NewMutex(name), nil
}

func (r *RedLock) NewMutex(name string) *redsync.Mutex {
	if r.sync == nil || r.flag {
		p := r.newPool()
		r.sync = redsync.New(p...)
		r.flag = false
	}
	return r.sync.NewMutex(name)
}
