package cache

type Interface interface {
	Set(key string, value interface{}, ttl ...uint) bool
	Get(key string) interface{}
	Delete(key string) bool
	Exist(key string) bool
	Ttl(key string) uint
	Flush() bool
	GC()
}
