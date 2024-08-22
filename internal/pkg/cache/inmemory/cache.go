package inmemory

import (
	"HomeWork_1/internal/model/errs"
	"sync"
	"time"
)

type CacheType string

const (
	LruCache = "lru"
	LfuCache = "lfu"
)

type CacheImpl[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V)
}

type InMemoryCache[K comparable, V any] struct {
	cache CacheImpl[K, V]
	mu    sync.RWMutex
}

func NewInMemoryCache[K comparable, V any](capacity int, ttlCache time.Duration, cacheType CacheType) (*InMemoryCache[K, V], error) {
	var cacheImpl CacheImpl[K, V]
	switch cacheType {
	case LruCache:
		cacheImpl = NewLRU[K, V](capacity, ttlCache)
	case LfuCache:
		cacheImpl = NewLFU[K, V](capacity, ttlCache)
	default:
		return nil, errs.ErrUnknownCacheType
	}

	return &InMemoryCache[K, V]{
		cache: cacheImpl,
		mu:    sync.RWMutex{},
	}, nil
}

func (c *InMemoryCache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Put(key, value)
}

func (c *InMemoryCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cache.Get(key)
}
