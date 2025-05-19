package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSlidingLimiter struct {
	Client *redis.Client
}

func NewRedisSlidingLimiter(client *redis.Client) *RedisSlidingLimiter {
	return &RedisSlidingLimiter{Client: client}
}

func (r *RedisSlidingLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, time.Duration, error) {
	now := time.Now().UnixNano()
	windowStart := now - window.Nanoseconds()

	keyName := fmt.Sprintf("ratelimit:%s", key)

	pipe := r.Client.TxPipeline()

	pipe.ZRemRangeByScore(ctx, keyName, "0", fmt.Sprintf("%d", windowStart))
	pipe.ZAdd(ctx, keyName, redis.Z{Score: float64(now), Member: now})
	count := pipe.ZCard(ctx, keyName)
	pipe.Expire(ctx, keyName, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, err
	}

	n, err := count.Result()
	if err != nil {
		return false, 0, err
	}

	if n > int64(limit) {
		retry := window - time.Duration(now-windowStart)
		return false, retry, nil
	}

	return true, 0, nil
}
