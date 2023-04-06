package trredis

import (
	"time"

	"github.com/RangelReale/trcache"
)

// Option

//troptgen:root
type options[K comparable, V any] interface {
	trcache.Options[K, V]
	trcache.NameOptions[K, V]
	trcache.CallDefaultOptions[K, V]
	OptKeyCodec(keyCodec trcache.KeyCodec[K])
	OptValueCodec(valueCodec trcache.Codec[V])
	OptValidator(validator trcache.Validator[V])
	OptDefaultDuration(duration time.Duration)
	OptRedisGetFunc(redisGetFunc RedisGetFunc[K, V])
	OptRedisSetFunc(redisSetFunc RedisSetFunc[K, V])
	OptRedisDelFunc(redisDelFunc RedisDelFunc[K, V])
}

// Cache get options

//troptgen:get
type getOptions[K comparable, V any] interface {
	trcache.GetOptions[K, V]
	OptCustomParams(customParams any)
	OptRedisGetFunc(redisGetFunc RedisGetFunc[K, V])
}

// helpers

func WithGetRedisGetFuncFunc[K comparable, V any](redisGetFunc RedisGetFuncFunc[K, V]) trcache.GetOption {
	return WithGetRedisGetFunc[K, V](redisGetFunc)
}

// Cache set options

//troptgen:set
type setOptions[K comparable, V any] interface {
	trcache.SetOptions[K, V]
	OptCustomParams(customParams any)
	OptRedisSetFunc(redisSetFunc RedisSetFunc[K, V])
}

// helpers

func WithSetRedisSetFuncFunc[K comparable, V any](redisSetFuncFunc RedisSetFuncFunc[K, V]) trcache.SetOption {
	return WithSetRedisSetFunc[K, V](redisSetFuncFunc)
}

// Cache delete options

//troptgen:delete
type deleteOptions[K comparable, V any] interface {
	trcache.DeleteOptions[K, V]
	OptCustomParams(customParams any)
	OptRedisDelFunc(redisDelFunc RedisDelFunc[K, V])
}

// helpers

func WithDeleteRedisDelFuncFunc[K comparable, V any](redisDelFunc RedisDelFuncFunc[K, V]) trcache.DeleteOption {
	return WithDeleteRedisDelFunc[K, V](redisDelFunc)
}

//go:generate troptgen
