package cache

type Option[V any] struct {
	Value V
	Error error
}

type Comparable[T any] interface {
	Equals(other T) bool
	LessThan(other T) bool
	Instance() T
}

type Cache[K, V any] interface {
	Get(key Comparable[K]) Option[V]
	Put(key Comparable[K], value V)
}
