package trredis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RangelReale/trcache"
	"github.com/RangelReale/trcache/codec"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestFuncGet(t *testing.T) {
	ctx := context.Background()

	redisClient, mockRedis := redismock.NewClientMock()

	mockRedis.ExpectHSet("a", "f1", "12", time.Minute).SetVal(1)
	mockRedis.ExpectHGet("a", "f1").SetVal("12")
	mockRedis.ExpectHGet("a", "f1").RedisNil() // simulate expiration
	mockRedis.ExpectHGet("z", "f1").RedisNil()

	c, err := New[string, string](redisClient,
		WithValueCodec[string, string](codec.NewForwardCodec[string]()),
		WithDefaultDuration[string, string](time.Minute),
		trcache.WithCallDefaultGetOptions[string, string](
			WithGetRedisGetFuncFunc[string, string](func(ctx context.Context, c *Cache[string, string], keyValue string, customParams any) (string, error) {
				value, err := c.Handle().HGet(ctx, keyValue, "f1").Result()
				if err != nil {
					if errors.Is(err, redis.Nil) {
						return "", trcache.ErrNotFound
					}
					return "", err
				}
				return value, nil
			}),
		),
		trcache.WithCallDefaultSetOptions[string, string](
			WithSetRedisSetFuncFunc[string, string](func(ctx context.Context, c *Cache[string, string], keyValue string, value any, expiration time.Duration, customParams any) error {
				return c.Handle().HSet(ctx, keyValue, "f1", value, expiration).Err()
			}),
		),
		trcache.WithCallDefaultDeleteOptions[string, string](
			WithDeleteRedisDelFuncFunc[string, string](func(ctx context.Context, c *Cache[string, string], keyValue string, customParams any) error {
				return c.Handle().HDel(ctx, keyValue, "f1").Err()
			}),
		),
	)
	require.NoError(t, err)

	err = c.Set(ctx, "a", "12")
	require.NoError(t, err)

	v, err := c.Get(ctx, "a")
	require.NoError(t, err)
	require.Equal(t, "12", v)

	v, err = c.Get(ctx, "a")
	require.ErrorIs(t, err, trcache.ErrNotFound)

	v, err = c.Get(ctx, "z")
	require.ErrorIs(t, err, trcache.ErrNotFound)

	if err := mockRedis.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
