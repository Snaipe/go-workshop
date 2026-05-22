package main

import (
	"sync"
)

type Cache[K comparable, V any] interface {
	Set(key K, value V)
	Get(key K) (value V, found bool)
	Delete(key K)
}

type entry[K comparable, V any] struct {
	next *entry[K, V]
	prev *entry[K, V]
	key  K
	val  V
}

type InMemoryCache[K comparable, V any] struct {
	size  int
	items map[K]*entry[K, V]
	head  entry[K, V]
	tail  entry[K, V]
	mu    sync.Mutex
}

func NewInMemoryCache[K comparable, V any](size int) *InMemoryCache[K, V] {
	cache := &InMemoryCache[K, V]{
		size:  size,
		items: make(map[K]*entry[K, V], size),
	}
	cache.head.next = &cache.tail
	cache.tail.prev = &cache.head
	return cache
}

// Set assigns the specified value to the specified key in the cache.
// If the cache is full, it deletes the least recently used entry to make space.
func (cache *InMemoryCache[K, V]) Set(key K, value V) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	e, found := cache.items[key]
	if found {
		e.val = value
		return
	}

	if len(cache.items) >= cache.size {
		cache.evict(cache.tail.prev)
	}
	e = &entry[K, V]{
		next: cache.head.next,
		prev: &cache.head,
		key:  key,
		val:  value,
	}
	cache.items[key] = e
	cache.head.next.prev = e
	cache.head.next = e
}

// Get retrieves the value in the cache for the specified key if it exists,
// as well as whether the value was found.
func (cache *InMemoryCache[K, V]) Get(key K) (value V, found bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	e, found := cache.items[key]
	if !found {
		return value, found
	}

	e.prev.next = e.next
	e.next.prev = e.prev
	e.prev = &cache.head
	e.next = cache.head.next
	cache.head.next.prev = e
	cache.head.next = e
	return e.val, true
}

// Delete removes the value associated with the specified key, if any.
func (cache *InMemoryCache[K, V]) Delete(key K) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	e, found := cache.items[key]
	if !found {
		return
	}

	cache.evict(e)
}

func (cache *InMemoryCache[K, V]) evict(e *entry[K, V]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	delete(cache.items, e.key)
}
