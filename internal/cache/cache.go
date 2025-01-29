package cache

type Cache[K comparable, D interface{}] interface {
	Clear()
	Get(key K) (D, bool)
	Empty() bool
	Len() int
	Remove(key K)
	Range(handler func(key K, data D))
	Set(key K, value D) bool
}
