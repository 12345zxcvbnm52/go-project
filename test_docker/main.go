package main

import (
	"fmt"
	"time"

	"github.com/allegro/bigcache"
)

func main() {
	fmt.Println("hello")
	d := bigcache.DefaultConfig(10 * time.Minute)
	cache, _ := bigcache.NewBigCache(d)
	cache.Set("key", []byte("woani"))
	fmt.Println(cache.Get("key"))
}
