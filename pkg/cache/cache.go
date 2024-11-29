package cache

import (
	"sync"
)

// The Cache interface defines the methods required for caching.
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
}

// memoryCache is an in-memory implementation of the Cache interface.
type memoryCache struct {
	mu    sync.RWMutex
	store map[string]item
}

// item represents the value stored in the cache
type item struct {
	value interface{}
}

// NewMemoryCache returns a new instance of memoryCache.
func NewMemoryCache() Cache {
	return &memoryCache{
		store: make(map[string]item),
	}
}

// Get retrieves a cached value by its key. If the key is not found, it returns nil and false.
func (m *memoryCache) Get(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	itm, exists := m.store[key]
	if !exists {
		return nil, false
	}

	return itm.value, true
}

// Set adds a value in the cache with a key.
func (m *memoryCache) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] = item{
		value: value,
	}
}

// Delete removes a key-value pair from the cache.
func (m *memoryCache) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, key)
}
