package trredis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/RangelReale/trcache"
	"github.com/RangelReale/trcache/codec"
	"github.com/RangelReale/trcache/mocks"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	ctx := context.Background()

	redis, mockRedis := redismock.NewClientMock()

	mockRedis.ExpectSet("a", "12", time.Minute).SetVal("12")
	mockRedis.ExpectGet("a").SetVal("12")
	mockRedis.ExpectGet("a").RedisNil() // simulate expiration
	mockRedis.ExpectGet("z").RedisNil()

	c, err := New[string, string](redis,
		WithValueCodec[string, string](codec.NewForwardCodec[string]()),
		WithDefaultDuration[string, string](time.Minute),
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

func TestCacheValidator(t *testing.T) {
	ctx := context.Background()

	redis, mockRedis := redismock.NewClientMock()
	mockValidator := mocks.NewValidator[string](t)

	mockRedis.ExpectSet("a", "12", time.Minute).SetVal("12")
	mockRedis.ExpectGet("a").SetVal("12")

	mockValidator.EXPECT().
		ValidateGet(mock.Anything, "12").
		Return(trcache.ErrNotFound).
		Once()

	c, err := New[string, string](redis,
		WithValueCodec[string, string](codec.NewForwardCodec[string]()),
		WithValidator[string, string](mockValidator),
		WithDefaultDuration[string, string](time.Minute),
	)
	require.NoError(t, err)

	err = c.Set(ctx, "a", "12")
	require.NoError(t, err)

	_, err = c.Get(ctx, "a")
	require.ErrorIs(t, err, trcache.ErrNotFound)

	if err := mockRedis.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacheCodecError(t *testing.T) {
	ctx := context.Background()

	redis, mockRedis := redismock.NewClientMock()
	mockCodec := mocks.NewCodec[string](t)

	mockRedis.ExpectGet("a").RedisNil()

	mockCodec.EXPECT().
		Marshal(mock.Anything, "12").
		Return(nil, errors.New("my error"))

	c, err := New[string, string](redis,
		WithValueCodec[string, string](mockCodec),
		WithDefaultDuration[string, string](time.Minute),
	)
	require.NoError(t, err)

	err = c.Set(ctx, "a", "12")
	require.ErrorAs(t, err, &trcache.CodecError{})

	_, err = c.Get(ctx, "a")
	require.ErrorIs(t, err, trcache.ErrNotFound)

	if err := mockRedis.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacheJSONCodec(t *testing.T) {
	ctx := context.Background()

	redis, mockRedis := redismock.NewClientMock()

	mockRedis.ExpectSet("a", `"12"`, time.Minute).SetVal(`"12"`)
	mockRedis.ExpectGet("a").SetVal(`"12"`)

	c, err := New[string, string](redis,
		WithValueCodec[string, string](codec.NewJSONCodec[string]()),
		WithDefaultDuration[string, string](time.Minute),
	)
	require.NoError(t, err)

	err = c.Set(ctx, "a", "12")
	require.NoError(t, err)

	v, err := c.Get(ctx, "a")
	require.NoError(t, err)
	require.Equal(t, "12", v)

	if err := mockRedis.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacheJSONCodecInt(t *testing.T) {
	ctx := context.Background()

	redis, mockRedis := redismock.NewClientMock()

	mockRedis.ExpectSet("a", "12", time.Minute).SetVal("12")
	mockRedis.ExpectGet("a").SetVal("12")

	c, err := New[string, int](redis,
		WithValueCodec[string, int](codec.NewJSONCodec[int]()),
		WithDefaultDuration[string, int](time.Minute),
	)
	require.NoError(t, err)

	err = c.Set(ctx, "a", 12)
	require.NoError(t, err)

	v, err := c.Get(ctx, "a")
	require.NoError(t, err)
	require.Equal(t, 12, v)

	if err := mockRedis.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacheFuncCodecInt(t *testing.T) {
	ctx := context.Background()

	redis, mockRedis := redismock.NewClientMock()

	mockRedis.ExpectSet("a", "12", time.Minute).SetVal("12")
	mockRedis.ExpectGet("a").SetVal("12")

	c, err := New[string, int](redis,
		WithValueCodec[string, int](codec.NewFuncCodec[int](
			func(ctx context.Context, data int) (any, error) {
				return fmt.Sprint(data), nil
			}, func(ctx context.Context, data any) (int, error) {
				return strconv.Atoi(fmt.Sprint(data))
			})),
		WithDefaultDuration[string, int](time.Minute),
	)
	require.NoError(t, err)

	err = c.Set(ctx, "a", 12)
	require.NoError(t, err)

	v, err := c.Get(ctx, "a")
	require.NoError(t, err)
	require.Equal(t, 12, v)

	if err := mockRedis.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacheCodecInvalidInt(t *testing.T) {
	ctx := context.Background()

	redis, mockRedis := redismock.NewClientMock()

	mockRedis.ExpectSet("a", 12, time.Minute).SetVal("12")
	mockRedis.ExpectGet("a").SetVal("12")

	c, err := New[string, int](redis,
		WithValueCodec[string, int](codec.NewForwardCodec[int]()),
		WithDefaultDuration[string, int](time.Minute),
	)
	require.NoError(t, err)

	err = c.Set(ctx, "a", 12)
	require.NoError(t, err)

	_, err = c.Get(ctx, "a")
	require.ErrorAs(t, err, new(*trcache.InvalidValueTypeError))

	if err := mockRedis.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacheRefresh(t *testing.T) {
	ctx := context.Background()

	redis, mockRedis := redismock.NewClientMock()

	mockRedis.ExpectGet("a").RedisNil()
	mockRedis.ExpectSet("a", "abc123", time.Minute).SetVal("abc123")

	c, err := NewRefresh[string, string, int](redis,
		WithValueCodec[string, string](codec.NewForwardCodec[string]()),
		WithDefaultDuration[string, string](time.Minute),
		trcache.WithDefaultRefreshFunc[string, string, int](func(ctx context.Context, key string, options trcache.RefreshFuncOptions[int]) (string, error) {
			return fmt.Sprintf("abc%d", options.Data), nil
		}),
	)
	require.NoError(t, err)

	value, err := c.GetOrRefresh(ctx, "a", trcache.WithRefreshData[string, string, int](123))
	require.NoError(t, err)
	require.Equal(t, "abc123", value)

	if err := mockRedis.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
