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

func TestCacheUpdateValue(t *testing.T) {
	cache := NewCache()

	// Insert initial value
	cache.Upsert(&KeyData{ID: "test-id", Value: "value1"})

	// Update with new value
	cache.Upsert(&KeyData{ID: "test-id", Value: "value2"})

	// Old value should not be found
	if item := cache.GetByValue("value1"); item != nil {
		t.Error("expected old value to be removed")
	}

	// New value should be found
	item := cache.GetByValue("value2")
	if item == nil {
		t.Fatal("expected to find item by new value")
	}
	if item.ID != "test-id" {
		t.Errorf("expected ID test-id, got %s", item.ID)
	}

	// Old value lookup should fail
	if cache.GetByValue("value1") != nil {
		t.Error("expected old value to be removed from index")
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

func TestCacheNilInput(t *testing.T) {
	cache := NewCache()
	cache.Upsert(nil) // should not panic
	if cache.Len() != 0 {
		t.Errorf("expected empty cache, got %d items", cache.Len())
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
