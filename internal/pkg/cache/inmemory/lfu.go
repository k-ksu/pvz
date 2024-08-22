package inmemory

import "time"

type Item[K comparable] struct {
	key          K
	lastUsedTime int64
	freq         int
}

type LFU[K comparable, V any] struct {
	data     []*Item[K]
	values   map[K]CacheValue[V]
	m        map[K]int
	capacity int
	ttlCache time.Duration
}

func NewLFU[K comparable, V any](capacity int, ttlCache time.Duration) *LFU[K, V] {
	return &LFU[K, V]{
		data:     make([]*Item[K], 0, capacity),
		values:   make(map[K]CacheValue[V], capacity),
		m:        make(map[K]int, capacity),
		capacity: capacity,
		ttlCache: ttlCache,
	}
}

func (l *LFU[K, V]) Get(key K) (V, bool) {
	index, ok := l.m[key]
	if !ok {
		return *new(V), false
	}

	cacheVal := l.values[key]
	if cacheVal.TTL.Before(time.Now()) {
		l.remove(key, index)
		return *new(V), false
	}

	item := l.data[index]
	item.freq++
	item.lastUsedTime = time.Now().UnixNano()

	return cacheVal.Item, true
}

func (l *LFU[K, V]) Put(key K, value V) {
	if index, ok := l.m[key]; ok {
		l.updateValue(key, value, index)
	} else {
		if len(l.data) == l.capacity {
			l.evict()
		}
		l.addValue(key, value)
	}
}

func (l *LFU[K, V]) updateValue(key K, value V, index int) {
	l.values[key] = CacheValue[V]{
		Item: value,
		TTL:  time.Now().Add(l.ttlCache),
	}
	item := l.data[index]
	item.freq++
	item.lastUsedTime = time.Now().UnixNano()
}

func (l *LFU[K, V]) addValue(key K, value V) {
	cacheVal := CacheValue[V]{
		Item: value,
		TTL:  time.Now().Add(l.ttlCache),
	}
	l.values[key] = cacheVal
	item := &Item[K]{key: key, freq: 1, lastUsedTime: time.Now().UnixNano()}
	l.data = append(l.data, item)
	l.m[key] = len(l.data) - 1
}

func (l *LFU[K, V]) evict() {
	var minFreq = int(^uint(0) >> 1) // Max int
	var minTime int64
	var evictIndex int

	for i, item := range l.data {
		if item.freq < minFreq || (item.freq == minFreq && item.lastUsedTime < minTime) {
			minFreq = item.freq
			minTime = item.lastUsedTime
			evictIndex = i
		}
	}

	evictItem := l.data[evictIndex]
	delete(l.values, evictItem.key)
	l.data = append(l.data[:evictIndex], l.data[evictIndex+1:]...)
	delete(l.m, evictItem.key)
	for i := evictIndex; i < len(l.data); i++ {
		l.m[l.data[i].key] = i
	}
}

func (l *LFU[K, V]) remove(key K, index int) {
	delete(l.values, key)
	l.data = append(l.data[:index], l.data[index+1:]...)
	delete(l.m, key)
	for i := index; i < len(l.data); i++ {
		l.m[l.data[i].key] = i
	}
}
