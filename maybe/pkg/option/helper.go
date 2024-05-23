package option

type Factory[T any] struct{}

func (f Factory[T]) Some(value T) Maybe[T] {
	return Maybe[T]{
		value: value,
		Error: nil,
	}
}

func (f Factory[T]) None(err error) Maybe[T] {
	return Maybe[T]{
		Error: err,
	}
}
