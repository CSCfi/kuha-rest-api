package cache

import (
	"context"
	"encoding/json"
	"time"
)

func SetCacheJSON(ctx context.Context, cache *Storage, key string, value any, ttl time.Duration) {
	if cache == nil {
		return
	}
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	_ = cache.Set(ctx, key, string(data), ttl)
}
