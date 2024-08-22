package inmemory

import "time"

type CacheValue[V any] struct {
	Item V
	TTL  time.Time
}

type LRU[K comparable, V any] struct {
	values   map[K]CacheValue[V]
	keys     []K
	capacity int
	ttlCache time.Duration
}

func NewLRU[K comparable, V any](capacity int, ttlCache time.Duration) *LRU[K, V] {
	return &LRU[K, V]{
		values:   make(map[K]CacheValue[V], capacity),
		keys:     make([]K, 0, capacity),
		capacity: capacity,
		ttlCache: ttlCache,
	}
}

func (l *LRU[K, V]) Get(key K) (V, bool) {
	val, ok := l.values[key]
	if !ok {
		return *new(V), false
	}

	if val.TTL.Before(time.Now()) {
		for i, k := range l.keys {
			if k != key {
				continue
			}

			l.keys = append(l.keys[:i], l.keys[:i]...)

			break
		}
		delete(l.values, key)
		return *new(V), false
	}

	for i, k := range l.keys {
		if k != key {
			continue
		}

		l.keys = append(l.keys[:i], append(l.keys[i+1:], l.keys[i])...)

		break
	}

	return val.Item, true
}

func (l *LRU[K, V]) Put(key K, value V) {
	_, ok := l.values[key]

	if !ok {
		if len(l.keys) == l.capacity {
			keyToDelete := l.keys[0]
			l.keys = l.keys[1:]
			delete(l.values, keyToDelete)
		}

		l.keys = append(l.keys, key)
		l.values[key] = CacheValue[V]{
			Item: value,
			TTL:  time.Now().Add(l.ttlCache),
		}
	} else {
		_, _ = l.Get(key)
		l.values[key] = CacheValue[V]{
			Item: value,
			TTL:  time.Now().Add(l.ttlCache),
		}
	}
}
