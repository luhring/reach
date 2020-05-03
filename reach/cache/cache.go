package cache

import (
	"sync"
)

type Cache struct {
	l     sync.RWMutex
	store map[string]interface{}
}

func New() Cache {
	store := make(map[string]interface{})

	return Cache{
		l:     sync.RWMutex{},
		store: store,
	}
}

func (c *Cache) Get(key string) interface{} {
	c.l.RLock()
	defer c.l.RUnlock()

	return c.store[key]
}

func (c *Cache) Put(key string, value interface{}) {
	c.l.Lock()
	defer c.l.Unlock()

	c.store[key] = value
}
