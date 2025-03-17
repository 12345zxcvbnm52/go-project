package cache

import (
	"context"
	"encoding/json"

	"github.com/allegro/bigcache"
)

func (c *MultiCache) getSpan(ctx context.Context) {}

func (c *MultiCache) getVersion(_ context.Context, data []byte) int64 {
	if err := json.Unmarshal(data, v); err != nil {
		return 0
	}
	return v.Version
}

func (c *MultiCache) Stats() bigcache.Stats {
	return c.localCache.Stats()
}
