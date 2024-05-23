package cache

import "sync"

type InfInMemoryCache[K, V any] struct {
	cache map[Comparable[K]]Option[V]
	mu    *sync.RWMutex
}

// NewInfInMemoryCache returns a new instance of InfInMemoryCache
func NewInfInMemoryCache[K, V any]() *InfInMemoryCache[K, V] {
	return &InfInMemoryCache[K, V]{
		cache: map[Comparable[K]]Option[V]{},
		mu:    &sync.RWMutex{},
	}
}

// Get retrieves the value associated with the given key from the cache.
func (c *InfInMemoryCache[K, V]) Get(key Comparable[K]) Option[V] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	res, exists := c.cache[key]

	if !exists {
		return Option[V]{
			Error: NotFoundError,
		}
	}

	return res
}

// Put adds or updates the key-value pair in the cache.
func (c *InfInMemoryCache[K, V]) Put(key Comparable[K], value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = Option[V]{
		Value: value,
	}
}
