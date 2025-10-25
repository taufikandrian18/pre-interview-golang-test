package main

import (
	"sync"
	"testing"
	"time"
)

func TestSimpleCache_Set(t *testing.T) {
	cache := NewSimpleCache()

	err := cache.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value, exists, err := cache.Get("key1"); err != nil || !exists || value != "value1" {
		t.Errorf("Expected value1, got %v, exists=%v, err=%v", value, exists, err)
	}

	err = cache.Set("key2", 42)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value, exists, err := cache.Get("key2"); err != nil || !exists || value != 42 {
		t.Errorf("Expected 42, got %v, exists=%v, err=%v", value, exists, err)
	}

	type TestStruct struct {
		Name  string
		Value int
	}
	testStruct := TestStruct{Name: "test", Value: 100}
	err = cache.Set("key3", testStruct)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value, exists, err := cache.Get("key3"); err != nil || !exists {
		t.Error("Expected struct to exist")
	} else if retrieved, ok := value.(TestStruct); !ok || retrieved != testStruct {
		t.Errorf("Expected %+v, got %+v", testStruct, retrieved)
	}
}

func TestSimpleCache_Get(t *testing.T) {
	cache := NewSimpleCache()

	if _, exists, err := cache.Get("nonexistent"); err != nil || exists {
		t.Errorf("Expected false for non-existent key, got exists=%v, err=%v", exists, err)
	}

	err := cache.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value, exists, err := cache.Get("key1"); err != nil || !exists || value != "value1" {
		t.Errorf("Expected value1, got %v, exists=%v, err=%v", value, exists, err)
	}
}

func TestSimpleCache_Delete(t *testing.T) {
	cache := NewSimpleCache()

	err := cache.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if _, exists, err := cache.Get("key1"); err != nil || !exists {
		t.Errorf("Expected key1 to exist after setting, got exists=%v, err=%v", exists, err)
	}

	err = cache.Delete("key1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if _, exists, err := cache.Get("key1"); err != nil || exists {
		t.Errorf("Expected key1 to be deleted, got exists=%v, err=%v", exists, err)
	}

	err = cache.Delete("nonexistent")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestTTLCache_SetAndGet(t *testing.T) {
	cache, err := NewTTLCache(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to create TTLCache: %v", err)
	}
	defer cache.Close()

	err = cache.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value, exists, err := cache.Get("key1"); err != nil || !exists || value != "value1" {
		t.Errorf("Expected value1, got %v, exists=%v, err=%v", value, exists, err)
	}
}

// expired entries test cases
func TestTTLCache_Expiration(t *testing.T) {
	cache, err := NewTTLCache(50 * time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to create TTLCache: %v", err)
	}
	defer cache.Close()

	err = cache.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if _, exists, err := cache.Get("key1"); err != nil || !exists {
		t.Errorf("Expected key1 to exist immediately after setting, got exists=%v, err=%v", exists, err)
	}

	time.Sleep(60 * time.Millisecond)

	if _, exists, err := cache.Get("key1"); err != nil || exists {
		t.Errorf("Expected key1 to be expired, got exists=%v, err=%v", exists, err)
	}
}

func TestTTLCache_Delete(t *testing.T) {
	cache, err := NewTTLCache(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to create TTLCache: %v", err)
	}
	defer cache.Close()

	err = cache.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	err = cache.Delete("key1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if _, exists, err := cache.Get("key1"); err != nil || exists {
		t.Errorf("Expected key1 to be deleted, got exists=%v, err=%v", exists, err)
	}
}

func TestTTLCache_UpdateTTL(t *testing.T) {
	cache, err := NewTTLCache(50 * time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to create TTLCache: %v", err)
	}
	defer cache.Close()

	err = cache.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	time.Sleep(30 * time.Millisecond)

	err = cache.Set("key1", "value2")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	time.Sleep(30 * time.Millisecond)

	if value, exists, err := cache.Get("key1"); err != nil || !exists || value != "value2" {
		t.Errorf("Expected value2, got %v, exists=%v, err=%v", value, exists, err)
	}
}

func TestTTLCache_Cleanup(t *testing.T) {
	cache, err := NewTTLCache(50 * time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to create TTLCache: %v", err)
	}
	defer cache.Close()

	err = cache.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = cache.Set("key2", "value2")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = cache.Set("key3", "value3")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	time.Sleep(60 * time.Millisecond)

	time.Sleep(30 * time.Millisecond)

	if _, exists, err := cache.Get("key1"); err != nil || exists {
		t.Errorf("Expected key1 to be cleaned up, got exists=%v, err=%v", exists, err)
	}
}

// concurrent access test cases
func TestSimpleCache_ConcurrentAccess(t *testing.T) {
	cache := NewSimpleCache()
	var wg sync.WaitGroup

	numGoroutines := 10
	opsPerGoroutine := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := "key" + string(rune(id*opsPerGoroutine+j))
				err := cache.Set(key, id*opsPerGoroutine+j)
				if err != nil {
					t.Errorf("Unexpected error in goroutine %d: %v", id, err)
				}
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < opsPerGoroutine; j++ {
			key := "key" + string(rune(i*opsPerGoroutine+j))
			expectedValue := i*opsPerGoroutine + j
			if value, exists, err := cache.Get(key); err != nil || !exists || value != expectedValue {
				t.Errorf("Concurrent access failed: expected %d, got %v, exists=%v, err=%v", expectedValue, value, exists, err)
			}
		}
	}
}

func TestTTLCache_ConcurrentAccess(t *testing.T) {
	cache, err := NewTTLCache(1 * time.Second)
	if err != nil {
		t.Fatalf("Failed to create TTLCache: %v", err)
	}
	defer cache.Close()
	var wg sync.WaitGroup

	numGoroutines := 10
	opsPerGoroutine := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := "key" + string(rune(id*opsPerGoroutine+j))
				err := cache.Set(key, id*opsPerGoroutine+j)
				if err != nil {
					t.Errorf("Unexpected error in goroutine %d: %v", id, err)
				}
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < opsPerGoroutine; j++ {
			key := "key" + string(rune(i*opsPerGoroutine+j))
			expectedValue := i*opsPerGoroutine + j
			if value, exists, err := cache.Get(key); err != nil || !exists || value != expectedValue {
				t.Errorf("Concurrent access failed: expected %d, got %v, exists=%v, err=%v", expectedValue, value, exists, err)
			}
		}
	}
}

func TestTTLCache_ConcurrentExpiration(t *testing.T) {
	cache, err := NewTTLCache(10 * time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to create TTLCache: %v", err)
	}
	defer cache.Close()
	var wg sync.WaitGroup

	err = cache.Set("shared_key", "shared_value")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	numGoroutines := 50
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_, _, _ = cache.Get("shared_key")
		}()
	}

	time.Sleep(5 * time.Millisecond)

	err = cache.Set("shared_key", "new_value")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	wg.Wait()

	if value, exists, err := cache.Get("shared_key"); err != nil || !exists || value != "new_value" {
		t.Errorf("Expected new_value after concurrent operations, got %v, exists=%v, err=%v", value, exists, err)
	}
}

func TestInterfaceCompliance(t *testing.T) {
	var cache Cache

	cache = NewSimpleCache()
	err := cache.Set("test", "value")
	if err != nil {
		t.Fatalf("SimpleCache Set error: %v", err)
	}
	if val, exists, err := cache.Get("test"); err != nil || !exists || val != "value" {
		t.Error("SimpleCache does not implement Cache interface correctly")
	}
	err = cache.Delete("test")
	if err != nil {
		t.Fatalf("SimpleCache Delete error: %v", err)
	}

	// Test TTLCache
	ttlCache, err := NewTTLCache(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to create TTLCache: %v", err)
	}
	defer ttlCache.Close()
	cache = ttlCache
	err = cache.Set("test", "value")
	if err != nil {
		t.Fatalf("TTLCache Set error: %v", err)
	}
	if val, exists, err := cache.Get("test"); err != nil || !exists || val != "value" {
		t.Error("TTLCache does not implement Cache interface correctly")
	}
	err = cache.Delete("test")
	if err != nil {
		t.Fatalf("TTLCache Delete error: %v", err)
	}
}

func BenchmarkSimpleCache_Set(b *testing.B) {
	cache := NewSimpleCache()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(rune(i)), i)
	}
}

func BenchmarkTTLCache_Set(b *testing.B) {
	cache, _ := NewTTLCache(1 * time.Hour)
	defer cache.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key"+string(rune(i)), i)
	}
}

func BenchmarkSimpleCache_Get(b *testing.B) {
	cache := NewSimpleCache()
	cache.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}

func BenchmarkTTLCache_Get(b *testing.B) {
	cache, _ := NewTTLCache(1 * time.Hour)
	defer cache.Close()
	cache.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}
