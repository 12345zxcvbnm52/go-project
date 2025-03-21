package cache

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/allegro/bigcache"
// 	"github.com/redis/go-redis/v9"
// )

// var c *MultiCache

// func init() {
// 	addr := []string{
// 		"192.168.199.128:6379",
// 		"192.168.199.128:6380",
// 		"192.168.199.128:6381",
// 		"192.168.199.128:6382",
// 		"192.168.199.128:6383",
// 		"192.168.199.128:6384",
// 	}
// 	conf := bigcache.DefaultConfig(8 * time.Minute)
// 	c = MustNewMultiCache(context.Background(),
// 		MustNewDistributedCache(addr, &redis.ClusterOptions{Addrs: addr, Password: "123"}),
// 		MustNewLocalCache(&conf),
// 	)
// }

// type test struct {
// 	input string
// 	sep   string
// 	want  []string
// }

// func TestCache(t *testing.T) {
// 	c.SetWithMutex()
// }
