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

// KeyData represents an API key with metadata
type KeyData struct {
	ID        string
	Value     string
	Scopes    []Scope
	Owner     string
	Namespace string
	Enabled   bool
	ExpiresAt *time.Time
}

// Cache is a thread-safe in-memory cache for KeyData
type Cache struct {
	mu       sync.RWMutex
	byID     map[string]*KeyData
	byValue  map[string]*KeyData
}

// NewCache creates a new cache instance.
func NewCache() *Cache {
	return &Cache{
		byID:    make(map[string]*KeyData),
		byValue: make(map[string]*KeyData),
	}
}

// GetByID retrieves a KeyData by ID, returning a copy
func (c *Cache) GetByID(id string) *KeyData {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if data, exists := c.byID[id]; exists {
		return &KeyData{
			ID:        data.ID,
			Value:     data.Value,
			Scopes:    append([]Scope{}, data.Scopes...),
			Owner:     data.Owner,
			Namespace: data.Namespace,
			Enabled:   data.Enabled,
			ExpiresAt: data.ExpiresAt,
		}
	}
	return nil
}

// GetByValue retrieves a KeyData by value, returning a copy
func (c *Cache) GetByValue(value string) *KeyData {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if data, exists := c.byValue[value]; exists {
		return &KeyData{
			ID:        data.ID,
			Value:     data.Value,
			Scopes:    append([]Scope{}, data.Scopes...),
			Owner:     data.Owner,
			Namespace: data.Namespace,
			Enabled:   data.Enabled,
			ExpiresAt: data.ExpiresAt,
		}
	}
	return nil
}

// Upsert inserts or updates a KeyData in the cache
func (c *Cache) Upsert(data *KeyData) {
	if data == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	// Handle value collision: remove old entry if same value exists with different ID
	if existing, exists := c.byValue[data.Value]; exists && existing.ID != data.ID {
		delete(c.byID, existing.ID)
	}
	
	// Remove old value if updating
	if existing, exists := c.byID[data.ID]; exists {
		delete(c.byValue, existing.Value)
	}
	
	c.byID[data.ID] = data
	c.byValue[data.Value] = data
}

// Delete removes a KeyData from the cache
func (c *Cache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if data, exists := c.byID[id]; exists {
		delete(c.byValue, data.Value)
		delete(c.byID, id)
	}
}

// Len returns the number of items in the cache
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.byID)
}
