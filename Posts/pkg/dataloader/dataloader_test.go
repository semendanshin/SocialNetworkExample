package dataloader

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

// Mock fetch function
func mockFetch(keys []int) ([]*string, []error) {
	results := make([]*string, len(keys))
	errs := make([]error, len(keys))

	for i, key := range keys {
		val := "value" + strconv.Itoa(key)
		results[i] = &val
		errs[i] = nil
	}

	return results, errs
}

func TestLoader_Load(t *testing.T) {
	loader := NewLoader[string, int](3, 10*time.Millisecond, mockFetch, 1*time.Minute, 1*time.Minute)

	val, err := loader.Load(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val == nil || *val != "value1" {
		t.Fatalf("unexpected value: %v", val)
	}
}

func TestLoader_BatchLoad(t *testing.T) {
	loader := NewLoader[string, int](3, 10*time.Millisecond, mockFetch, 1*time.Minute, 1*time.Minute)

	var wg sync.WaitGroup
	results := make([]*string, 3)
	errs := make([]error, 3)

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			val, err := loader.Load(i)
			results[i-1] = val
			errs[i-1] = err
		}(i)
	}

	wg.Wait()

	for i := 1; i <= 3; i++ {
		if errs[i-1] != nil {
			t.Fatalf("unexpected error for key %d: %v", i, errs[i-1])
		}
		if results[i-1] == nil || *results[i-1] != "value"+strconv.Itoa(i) {
			t.Fatalf("unexpected value for key %d: %v", i, results[i-1])
		}
	}
}

func TestLoader_Cache(t *testing.T) {
	loader := NewLoader[string, int](3, 10*time.Millisecond, mockFetch, 1*time.Minute, 1*time.Minute)

	val1, err1 := loader.Load(1)
	if err1 != nil {
		t.Fatalf("unexpected error: %v", err1)
	}
	if val1 == nil || *val1 != "value1" {
		t.Fatalf("unexpected value: %v", val1)
	}

	val2, err2 := loader.Load(1)
	if err2 != nil {
		t.Fatalf("unexpected error: %v", err2)
	}
	if val2 == nil || *val2 != *val1 {
		t.Fatalf("cache miss: %v != %v", val2, val1)
	}
}

func TestLoader_CacheExpiry(t *testing.T) {
	loader := NewLoader[string, int](3, 10*time.Millisecond, mockFetch, 50*time.Millisecond, 50*time.Millisecond)

	val1, err1 := loader.Load(1)
	if err1 != nil {
		t.Fatalf("unexpected error: %v", err1)
	}
	if val1 == nil || *val1 != "value1" {
		t.Fatalf("unexpected value: %v", val1)
	}

	time.Sleep(100 * time.Millisecond) // wait for cache to expire

	val2, err2 := loader.Load(1)
	if err2 != nil {
		t.Fatalf("unexpected error: %v", err2)
	}
	if val2 == nil || *val2 != "value1" {
		t.Fatalf("unexpected value after cache expiry: %v", val2)
	}
	if val2 == val1 {
		t.Fatalf("cache did not expire: %v == %v", val2, val1)
	}
}

func TestLoader_ErrorHandling(t *testing.T) {
	// Mock fetch function with error
	fetchWithError := func(keys []int) ([]*string, []error) {
		results := make([]*string, len(keys))
		errs := make([]error, len(keys))

		for i := range keys {
			results[i] = nil
			errs[i] = fmt.Errorf("fetch error")
		}

		return results, errs
	}

	loader := NewLoader[string, int](3, 10*time.Millisecond, fetchWithError, 1*time.Minute, 1*time.Minute)

	_, err := loader.Load(1)
	if err == nil || err.Error() != "fetch error" {
		t.Fatalf("expected fetch error, got: %v", err)
	}
}
