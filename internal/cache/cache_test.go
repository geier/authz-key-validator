package cache

import (
	"testing"
)

func TestCacheUpsertAndGet(t *testing.T) {
	cache := NewCache()

	data := &KeyData{
		ID:    "test-id",
		Value: "test-value",
		Owner: "test-owner",
	}
	cache.Upsert(data)

	// Test GetByID
	result := cache.GetByID("test-id")
	if result == nil {
		t.Fatal("expected to find item by ID")
	}
	if result.ID != "test-id" {
		t.Errorf("expected ID test-id, got %s", result.ID)
	}

	// Test GetByValue
	result2 := cache.GetByValue("test-value")
	if result2 == nil {
		t.Fatal("expected to find item by value")
	}
	if result2.ID != "test-id" {
		t.Errorf("expected ID test-id, got %s", result2.ID)
	}
}

func TestCacheDelete(t *testing.T) {
	cache := NewCache()

	cache.Upsert(&KeyData{ID: "test-id", Value: "test-value"})
	if cache.Len() != 1 {
		t.Fatalf("expected length 1, got %d", cache.Len())
	}

	cache.Delete("test-id")
	if cache.Len() != 0 {
		t.Errorf("expected empty cache, got %d items", cache.Len())
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	cache := NewCache()
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			cache.Upsert(&KeyData{ID: "test-id", Value: "test-value"})
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if cache.Len() != 1 {
		t.Errorf("expected cache length 1, got %d", cache.Len())
	}
}

func BenchmarkCacheLookup(b *testing.B) {
	cache := NewCache()
	cache.Upsert(&KeyData{ID: "bench-id", Value: "bench-value"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetByID("bench-id")
	}
}
