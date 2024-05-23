package option

type MaybeFactory[T any] interface {
	Some(value T) Maybe[T]
	None(err error) Maybe[T]
}
