package cache

import (
	"sync"
	"time"
)

// Scope represents a scope with type and resource
type Scope struct {
	Type     string
	Resource string
}

// KeyData represents an API key with its metadata
type KeyData struct {
	ID        string
	Value     string
	Scopes    []Scope
	Owner     string
	Namespace string
	Enabled   bool
	ExpiresAt *time.Time
}

// Cache provides a thread-safe in-memory cache for KeyData
type Cache struct {
	mu       sync.RWMutex
	byID    map[string]*KeyData
	byValue  map[string]*KeyData
}

// NewCache creates a new cache instance
func NewCache() *Cache {
	return &Cache{
		byID:    make(map[string]*KeyData),
		byValue: make(map[string]*KeyData),
	}
}

// GetByID retrieves a KeyData by ID
func (c *Cache) GetByID(id string) *KeyData {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.byID[id]
}

// GetByValue retrieves a KeyData by its key value
func (c *Cache) GetByValue(value string) *KeyData {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.byValue[value]
}

// Upsert inserts or updates the KeyData in the cache
func (c *Cache) Upsert(data *KeyData) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.byID[data.ID] = data
	c.byValue[data.Value] = data
}

// Delete removes a KeyData from the cache
func (c *Cache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if data, exists := c.byID[id]; exists {
		delete(c.byValue, data.Value)
	}
	delete(c.byID, id)
}

// Len returns the number of items in the cache
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.byID)
}
