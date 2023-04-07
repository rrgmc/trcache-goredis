package trredis

import (
	"context"

	"github.com/RangelReale/trcache"
	"github.com/RangelReale/trcache/refresh"
	"github.com/redis/go-redis/v9"
)

type RefreshCache[K comparable, V any] struct {
	*Cache[K, V]
	helper *refresh.Helper[K, V]
}

var _ trcache.RefreshCache[string, string] = &RefreshCache[string, string]{}

func NewRefresh[K comparable, V any](redis *redis.Client,
	options ...trcache.RootOption) (*RefreshCache[K, V], error) {
	checker := trcache.NewOptionChecker(options)

	c, err := New[K, V](redis, trcache.ForwardOptionsChecker(checker)...)
	if err != nil {
		return nil, err
	}

	helper, err := refresh.NewHelper[K, V](trcache.ForwardOptionsChecker(checker)...)
	if err != nil {
		return nil, err
	}

	if err = checker.CheckCacheError(); err != nil {
		return nil, err
	}

	ret := &RefreshCache[K, V]{
		Cache:  c,
		helper: helper,
	}
	return ret, nil
}

func (c *RefreshCache[K, V]) GetOrRefresh(ctx context.Context, key K, options ...trcache.RefreshOption) (V, error) {
	return c.helper.GetOrRefresh(ctx, c, key, options...)
}
