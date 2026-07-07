package server

import "sync"

// Cache is an in-memory key-value cache
type Cache struct {
	mu    sync.RWMutex
	items map[string]string
}

// NewCache creates a new Cache
func NewCache() *Cache {
	return &Cache{items: make(map[string]string)}
}

// Set stores a value in the cache
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.items[key]
	return v, ok
}

// Delete removes a key from the cache - dead code
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache - dead code
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]string)
}

// Len returns the number of items in the cache - dead code
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
