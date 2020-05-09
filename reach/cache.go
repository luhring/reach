package reach

// Cache is the interface that describes the ability to Put items into a cache and Get items that have been stored in a cache.
type Cache interface {
	Put(key string, value interface{})
	Get(key string) interface{}
}
