package trredis

import (
	"context"
	"errors"
	"fmt"

	"github.com/RangelReale/trcache"
	"github.com/RangelReale/trcache/codec"
	"github.com/redis/go-redis/v9"
)

type Cache[K comparable, V any] struct {
	options rootOptionsImpl[K, V]
	redis   *redis.Client
}

func New[K comparable, V any](redis *redis.Client, options ...trcache.RootOption) (*Cache[K, V], error) {
	ret := &Cache[K, V]{
		redis: redis,
		options: rootOptionsImpl[K, V]{
			defaultDuration: 0, // 0 means default for go-redis
			redisGetFunc:    DefaultRedisGetFunc[K, V]{},
			redisSetFunc:    DefaultRedisSetFunc[K, V]{},
			redisDelFunc:    DefaultRedisDelFunc[K, V]{},
		},
	}
	optErr := trcache.ParseOptions[trcache.RootOption](&ret.options, options)
	if optErr.Err() != nil {
		return nil, optErr.Err()
	}
	if ret.options.valueCodec == nil {
		return nil, errors.New("value codec is required")
	}
	if ret.options.keyCodec == nil {
		ret.options.keyCodec = codec.NewStringKeyCodec[K]()
	}
	return ret, nil
}

func (c *Cache[K, V]) Handle() *redis.Client {
	return c.redis
}

func (c *Cache[K, V]) Name() string {
	return c.options.name
}

func (c *Cache[K, V]) Get(ctx context.Context, key K, options ...trcache.GetOption) (V, error) {
	optns := getOptionsImpl[K, V]{
		redisGetFunc: c.options.redisGetFunc,
	}
	optErr := trcache.ParseOptions[trcache.GetOption](&optns, c.options.callDefaultGetOptions, options)
	if optErr.Err() != nil {
		var empty V
		return empty, optErr.Err()
	}

	keyValue, err := c.parseKey(ctx, key)
	if err != nil {
		var empty V
		return empty, err
	}

	value, err := optns.redisGetFunc.Get(ctx, c, keyValue, optns.customParams)
	if err != nil {
		var empty V
		return empty, err
	}

	dec, err := c.options.valueCodec.Unmarshal(ctx, value)
	if err != nil {
		var empty V
		return empty, trcache.CodecError{err}
	}

	if c.options.validator != nil {
		if err = c.options.validator.ValidateGet(ctx, dec); err != nil {
			var empty V
			return empty, err
		}
	}

	return dec, nil
}

func (c *Cache[K, V]) Set(ctx context.Context, key K, value V, options ...trcache.SetOption) error {
	optns := setOptionsImpl[K, V]{
		duration:     c.options.defaultDuration,
		redisSetFunc: c.options.redisSetFunc,
	}
	optErr := trcache.ParseOptions[trcache.SetOption](&optns, c.options.callDefaultSetOptions, options)
	if optErr.Err() != nil {
		return optErr.Err()
	}

	enc, err := c.options.valueCodec.Marshal(ctx, value)
	if err != nil {
		return trcache.CodecError{err}
	}

	keyValue, err := c.parseKey(ctx, key)
	if err != nil {
		return err
	}

	return optns.redisSetFunc.Set(ctx, c, keyValue, enc, c.options.defaultDuration, optns.customParams)
}

func (c *Cache[K, V]) Delete(ctx context.Context, key K, options ...trcache.DeleteOption) error {
	optns := deleteOptionsImpl[K, V]{
		redisDelFunc: c.options.redisDelFunc,
	}
	optErr := trcache.ParseOptions[trcache.DeleteOption](&optns, c.options.callDefaultDeleteOptions, options)
	if optErr.Err() != nil {
		return optErr.Err()
	}

	keyValue, err := c.parseKey(ctx, key)
	if err != nil {
		return err
	}

	return optns.redisDelFunc.Delete(ctx, c, keyValue, optns.customParams)
}

func (c *Cache[K, V]) parseKey(ctx context.Context, key K) (string, error) {
	keyValue, err := c.options.keyCodec.Convert(ctx, key)
	if err != nil {
		return "", trcache.CodecError{err}
	}

	switch kv := keyValue.(type) {
	case string:
		return kv, nil
	case []byte:
		return string(kv), nil
	default:
		return "", trcache.CodecError{
			&trcache.InvalidValueTypeError{fmt.Sprintf("invalid type '%T' for redis key", keyValue)},
		}
	}
}
