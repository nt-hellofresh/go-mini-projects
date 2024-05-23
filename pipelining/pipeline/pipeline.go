package pipeline

type PipelineOpts[T any] func(p *Pipeline[T])

type Pipeline[T any] struct {
	handlers []func(in T) T
}

func Define[T any](opts ...PipelineOpts[T]) *Pipeline[T] {
	pl := &Pipeline[T]{}

	for _, o := range opts {
		o(pl)
	}

	return pl
}

func WithStep[T any](fn func(value T) T) PipelineOpts[T] {
	return func(p *Pipeline[T]) {
		p.handlers = append(p.handlers, fn)
	}
}

func StartWith[T any](fn func(value T) T) *Pipeline[T] {
	return &Pipeline[T]{
		handlers: []func(T) T{
			fn,
		},
	}
}

func (pl *Pipeline[T]) Then(fn func(value T) T) {
	pl.handlers = append(pl.handlers, fn)
}

func (pl *Pipeline[T]) Process(value T) T {
	result := value
	for _, h := range pl.handlers {
		result = h(result)
	}

	return result
}
