package dataloader

import (
	"sync"
	"time"
)

// Loader is a generic data loader for batching and caching requests.
type Loader[T any, key comparable] struct {
	maxBatch int
	wait     time.Duration
	fetch    func([]key) ([]*T, []error)
	ttl      time.Duration

	cache        map[key]cacheEntry[T]
	currentBatch *loaderBatch[T, key]
	mu           sync.Mutex
}

// cacheEntry represents a cached item with a value and a timestamp.
type cacheEntry[T any] struct {
	value   *T
	addedAt time.Time
}

// loaderBatch holds the batch of keys being loaded.
type loaderBatch[T any, key comparable] struct {
	keys    []key
	data    []*T
	errors  []error
	closing bool
	done    chan struct{}
}

// NewLoader creates a new Loader with the given configuration.
func NewLoader[T any, key comparable](
	maxBatch int,
	wait time.Duration,
	fetch func([]key) ([]*T, []error),
	ttl time.Duration,
	cacheCleanerInterval time.Duration,
) *Loader[T, key] {
	loader := &Loader[T, key]{
		maxBatch: maxBatch,
		wait:     wait,
		fetch:    fetch,
		ttl:      ttl,
		cache:    make(map[key]cacheEntry[T]),
	}
	if ttl > 0 && cacheCleanerInterval > 0 {
		go loader.startCacheCleaner(cacheCleanerInterval)
	}
	return loader
}

// Load adds a key to the current batch and returns the cached value if available.
func (l *Loader[T, key]) Load(k key) (*T, error) {
	l.mu.Lock()

	// Check cache
	if entry, found := l.cache[k]; found && time.Since(entry.addedAt) <= l.ttl {
		l.mu.Unlock()
		return entry.value, nil
	}

	// Initialize a new batch if necessary
	if l.currentBatch == nil {
		l.currentBatch = &loaderBatch[T, key]{done: make(chan struct{})}
		go l.startTimer()
	}

	// Add key to batch
	batch := l.currentBatch
	batch.keys = append(batch.keys, k)
	if len(batch.keys) >= l.maxBatch {
		l.dispatch()
	}

	l.mu.Unlock()
	<-batch.done

	// Retrieve the result
	for i, key := range batch.keys {
		if key == k {
			return batch.data[i], batch.errors[i]
		}
	}
	return nil, nil
}

// startTimer triggers a batch dispatch after the wait duration.
func (l *Loader[T, key]) startTimer() {
	time.Sleep(l.wait)
	l.mu.Lock()
	if l.currentBatch != nil && !l.currentBatch.closing {
		l.dispatch()
	}
	l.mu.Unlock()
}

// dispatch sends the current batch to the fetch function.
func (l *Loader[T, key]) dispatch() {
	batch := l.currentBatch
	l.currentBatch = nil

	batch.closing = true
	l.mu.Unlock()

	batch.data, batch.errors = l.fetch(batch.keys)
	l.mu.Lock()

	// Cache the results
	for i, key := range batch.keys {
		if batch.errors[i] == nil && i < len(batch.data) {
			l.cache[key] = cacheEntry[T]{
				value:   batch.data[i],
				addedAt: time.Now(),
			}
		}
	}

	close(batch.done)
}

// Close forces the dispatch of any remaining batch.
func (l *Loader[T, key]) Close() {
	l.mu.Lock()
	if l.currentBatch != nil && !l.currentBatch.closing {
		l.dispatch()
	}
	l.mu.Unlock()
}

// startCacheCleaner periodically removes stale items from the cache.
func (l *Loader[T, key]) startCacheCleaner(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		l.cleanUpCache()
	}
}

// cleanUpCache removes items from the cache that have exceeded the TTL.
func (l *Loader[T, key]) cleanUpCache() {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	for k, entry := range l.cache {
		if now.Sub(entry.addedAt) > l.ttl {
			delete(l.cache, k)
		}
	}
}
