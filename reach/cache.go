package reach

type Cache interface {
	Put(key string, value interface{})
	Get(key string) interface{}
}
