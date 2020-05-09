package cache

import (
	"sync"
)

// Cache is a simple in-memory implementation of the interface reach.Cache.
type Cache struct {
	l     sync.RWMutex
	store map[string]interface{}
}

// New returns a new instance of Cache, ready to use.
func New() Cache {
	store := make(map[string]interface{})

	return Cache{
		l:     sync.RWMutex{},
		store: store,
	}
}

// Get returns the object from the cache for the matching key, if one exists.
func (c *Cache) Get(key string) interface{} {
	c.l.RLock()
	defer c.l.RUnlock()

	return c.store[key]
}

// Put inserts or updates an object in the cache for the given key.
func (c *Cache) Put(key string, value interface{}) {
	c.l.Lock()
	defer c.l.Unlock()

	c.store[key] = value
}
