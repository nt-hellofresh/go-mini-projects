package future

type InputFunc[T any] func() (T, error)
