package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
}

func NewRedisStorage(rdb *redis.Client) *Storage {
	return &Storage{client: rdb}
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *Storage) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return s.client.Set(ctx, key, value, ttl).Err()
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

func (s *Storage) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, cur, err := s.client.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			pipe := s.client.Pipeline()
			for _, k := range keys {
				pipe.Del(ctx, k)
			}
			if _, err := pipe.Exec(ctx); err != nil {
				return err
			}
		}
		cursor = cur
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (s *Storage) DeleteByPrefixes(ctx context.Context, prefixes ...string) error {
	for _, p := range prefixes {
		if err := s.DeleteByPattern(ctx, p+"*"); err != nil {
			return err
		}
	}
	return nil
}
