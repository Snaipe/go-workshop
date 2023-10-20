package main

import (
	"testing"
	"time"
)

func TestInMemoryCache(t *testing.T) {
	stringCache := NewInMemoryCache[string, string]()
	stringCache.Set("foo", "1", 1*time.Hour)
	stringCache.Set("bar", "2", 1*time.Nanosecond)
	stringCache.Set("baz", "3", 1*time.Hour)

	foo, ok := stringCache.Get("foo")
	if !ok {
		t.Fatalf("expected key foo to be in cache, but it was not (cache: %v)", stringCache)
	}
	if foo != "1" {
		t.Fatalf("expected key foo to have value 1, but got %v", foo)
	}

	_, ok = stringCache.Get("bar")
	if ok {
		t.Fatal("expected key bar to have expired, but it was still present")
	}

	stringCache.Expire("foo")
	_, ok = stringCache.Get("foo")
	if ok {
		t.Fatal("expected key foo to have expired, but it was still present")
	}
}
