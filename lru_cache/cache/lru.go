package cache

import (
	"container/list"
	"errors"
	"log"
	"sync"
)

var NotFoundError = errors.New("unable to find item from lru cache")

// LRUCache is a generic LRU cache implementation.
type LRUCache[K, V any] struct {
	capacity int
	cache    map[Comparable[K]]*list.Element
	order    *list.List
	mu       *sync.Mutex
}

// entry represents an entry in the cache.
type entry[K, V any] struct {
	key   Comparable[K]
	value V
}

// NewLRUCache creates a new LRUCache with the given capacity.
func NewLRUCache[K, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		cache:    make(map[Comparable[K]]*list.Element),
		order:    list.New(),
		mu:       &sync.Mutex{},
	}
}

// Get retrieves the value associated with the given key from the cache.
func (lru *LRUCache[K, V]) Get(key Comparable[K]) Option[V] {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if elem, ok := lru.cache[key]; ok {
		lru.order.MoveToFront(elem)
		return Option[V]{
			Value: elem.Value.(*entry[K, V]).value,
			Error: nil,
		}
	}
	return Option[V]{
		Error: NotFoundError,
	}
}

// Put adds or updates the key-value pair in the cache.
func (lru *LRUCache[K, V]) Put(key Comparable[K], value V) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if elem, ok := lru.cache[key]; ok {
		elem.Value = value
		lru.order.MoveToFront(elem)
	} else {
		if len(lru.cache) >= lru.capacity {
			lastElem := lru.order.Back()
			delete(lru.cache, lastElem.Value.(*entry[K, V]).key)
			log.Printf("evicted %+v\n", lastElem.Value)
			lru.order.Remove(lastElem)
		}

		newElem := lru.order.PushFront(&entry[K, V]{key, value})
		lru.cache[key] = newElem
	}
}
