package main

import (
	"testing"
)

func TestInMemoryCache(t *testing.T) {
	const maxSize = 16

	cache := NewInMemoryCache[int, int](maxSize)

	cache.Set(0, 0)
	cache.Set(0, 1)
	cache.Set(0, 2)

	val, found := cache.Get(0)
	if !found {
		t.Fatal("cache should have item 0: got not found instead")
	}
	if val != 2 {
		t.Fatalf("unexpected value for item 0: expected 2, got %d", val)
	}

	n := 0
	for e := cache.head.next; e != &cache.tail; e = e.next {
		n++
	}
	if n != 1 {
		t.Fatalf("unexpected lru size: expected 1, got %d", n)
	}

	for i := range maxSize + 1 {
		cache.Set(i, i)
	}

	if len(cache.items) > maxSize {
		t.Fatalf("cache grew beyond its maximum size: expected %d, got %d", maxSize, cache.size)
	}

	val, found = cache.Get(0)
	if found {
		t.Fatalf("cache kept least-recently-used item: expected not found, got item %d", val)
	}

	val, found = cache.Get(1)
	if !found {
		t.Fatal("cache should have kept item 1: got not found instead")
	}

	if cache.items[1].prev != &cache.head {
		t.Fatal("cache.Get(1) should have moved element to front")
	}
}
