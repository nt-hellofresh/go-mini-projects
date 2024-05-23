package cache

import (
	"context"
)

type WrappedFunc[K, V any] func(ctx context.Context, key K) (V, error)
type CachedFunc[K, V any] func(ctx context.Context, key K) (V, error)

func ReadThroughCache[K, V any](fn WrappedFunc[K, V], c Cache[K, V]) CachedFunc[Comparable[K], V] {
	return func(ctx context.Context, key Comparable[K]) (V, error) {
		option := c.Get(key)

		if option.Error == NotFoundError {
			val, err := fn(ctx, key.Instance())
			c.Put(key, val)
			return val, err
		}

		return option.Value, option.Error
	}
}
