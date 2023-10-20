package main

import (
	"container/heap"
	"sync"
	"time"
)

type Cache[K comparable, V any] interface {
	Set(key K, value V, ttl time.Duration)
	Get(key K) (value V, found bool)
	Expire(key K)
}

type InMemoryCache[K comparable, V any] struct {
	cache      map[K]*cacheBucket[K, V]
	expireList cacheExpireList[K, V]
	mux        sync.RWMutex
}

func NewInMemoryCache[K comparable, V any]() *InMemoryCache[K, V] {
	return &InMemoryCache[K, V]{
		cache: make(map[K]*cacheBucket[K, V]),
	}
}

// Set assigns the specified value to the specified key in the cache, with
// an expiration of ttl.
func (cache *InMemoryCache[K, V]) Set(key K, value V, ttl time.Duration) {
	cache.mux.Lock()
	defer cache.mux.Unlock()

	cache.flush()

	bucket := &cacheBucket[K, V]{
		expiry: time.Now().Add(ttl),
		key:    key,
		val:    value,
	}
	cache.cache[key] = bucket
	heap.Push(&cache.expireList, bucket)
}

// Get retrieves the value in the cache for the specified key if it exists,
// as well as whether the value was found.
func (cache *InMemoryCache[K, V]) Get(key K) (value V, found bool) {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	bucket, found := cache.cache[key]
	if found {
		value = bucket.val
	}
	return value, found
}

// Expire expires the value associated with the specified key, if any.
func (cache *InMemoryCache[K, V]) Expire(key K) {
	cache.mux.Lock()
	defer cache.mux.Unlock()

	bucket, found := cache.cache[key]
	if !found {
		return
	}
	delete(cache.cache, key)
	heap.Remove(&cache.expireList, bucket.idx)
}

// flush removes all expired keys from the cache
func (cache *InMemoryCache[K, V]) flush() {
	now := time.Now()
	for {
		bucket, ok := cache.expireList.Peek()
		if !ok || bucket.expiry.After(now) {
			break
		}
		delete(cache.cache, bucket.key)
		heap.Remove(&cache.expireList, bucket.idx)
	}
}

type cacheBucket[K, V any] struct {
	expiry time.Time
	key    K
	val    V
	idx    int // cache buckets know their position in the expire list
}

type cacheExpireList[K, V any] struct {
	elts []*cacheBucket[K, V]
}

func (c *cacheExpireList[K, V]) Peek() (*cacheBucket[K, V], bool) {
	if len(c.elts) > 0 {
		return c.elts[0], true
	}
	return nil, false
}

// cacheExpireList must implement sort.Interface and container/heap.Interface

func (c *cacheExpireList[K, V]) Len() int {
	return len(c.elts)
}

func (c *cacheExpireList[K, V]) Less(i, j int) bool {
	return c.elts[i].expiry.Before(c.elts[j].expiry)
}

func (c *cacheExpireList[K, V]) Swap(i, j int) {
	c.elts[i], c.elts[j] = c.elts[j], c.elts[i]
	// Fix the expire list indices
	c.elts[i].idx, c.elts[j].idx = i, j
}

func (c *cacheExpireList[K, V]) Push(x any) {
	c.elts = append(c.elts, x.(*cacheBucket[K, V]))
}

func (c *cacheExpireList[K, V]) Pop() (val any) {
	c.elts, val = c.elts[:len(c.elts)-1], c.elts[len(c.elts)-1]
	return val
}
