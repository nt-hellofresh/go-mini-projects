package option

type Maybe[T any] struct {
	value T
	Error error
}

func (m Maybe[T]) Value() T {
	if m.HasError() {
		panic(m.Error)
	}
	return m.value
}

func (m Maybe[T]) HasError() bool {
	return m.Error != nil
}
