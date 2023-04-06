package trredis

import (
	"context"
	"errors"
	"time"

	"github.com/RangelReale/trcache"
	"github.com/redis/go-redis/v9"
)

type RedisGetFunc[K comparable, V any] interface {
	Get(ctx context.Context, c *Cache[K, V], keyValue string, customParams any) (string, error)
}

type RedisSetFunc[K comparable, V any] interface {
	Set(ctx context.Context, c *Cache[K, V], keyValue string, value any, expiration time.Duration, customParams any) error
}

type RedisDelFunc[K comparable, V any] interface {
	Delete(ctx context.Context, c *Cache[K, V], keyValue string, customParams any) error
}

// Interface funcs

type RedisGetFuncFunc[K comparable, V any] func(ctx context.Context, c *Cache[K, V], keyValue string, customParams any) (string, error)

func (o RedisGetFuncFunc[K, V]) Get(ctx context.Context, c *Cache[K, V], keyValue string, customParams any) (string, error) {
	return o(ctx, c, keyValue, customParams)
}

type RedisSetFuncFunc[K comparable, V any] func(ctx context.Context, c *Cache[K, V], keyValue string, value any, expiration time.Duration, customParams any) error

func (o RedisSetFuncFunc[K, V]) Set(ctx context.Context, c *Cache[K, V], keyValue string, value any, expiration time.Duration, customParams any) error {
	return o(ctx, c, keyValue, value, expiration, customParams)
}

type RedisDelFuncFunc[K comparable, V any] func(ctx context.Context, c *Cache[K, V], keyValue string, customParams any) error

func (o RedisDelFuncFunc[K, V]) Delete(ctx context.Context, c *Cache[K, V], keyValue string, customParams any) error {
	return o(ctx, c, keyValue, customParams)
}

// Default

type DefaultRedisGetFunc[K comparable, V any] struct {
}

func (f DefaultRedisGetFunc[K, V]) Get(ctx context.Context, c *Cache[K, V], keyValue string, _ any) (string, error) {
	value, err := c.Handle().Get(ctx, keyValue).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", trcache.ErrNotFound
		}
		return "", err
	}
	return value, nil
}

type DefaultRedisSetFunc[K comparable, V any] struct {
}

func (f DefaultRedisSetFunc[K, V]) Set(ctx context.Context, c *Cache[K, V], keyValue string, value any,
	expiration time.Duration, _ any) error {
	return c.Handle().Set(ctx, keyValue, value, expiration).Err()
}

type DefaultRedisDelFunc[K comparable, V any] struct {
}

func (f DefaultRedisDelFunc[K, V]) Delete(ctx context.Context, c *Cache[K, V], keyValue string, _ any) error {
	return c.Handle().Del(ctx, keyValue).Err()
}
