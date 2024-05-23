package future

import "context"

type Results[T any] struct {
	Value T
}

type Promise[T any] struct {
	inputFunc   InputFunc[T]
	successChan chan Results[T]
	errorChan   chan error
}

func Deferred[T any](f InputFunc[T]) *Promise[T] {
	return &Promise[T]{
		inputFunc:   f,
		successChan: make(chan Results[T]),
		errorChan:   make(chan error),
	}
}

func (p *Promise[T]) ExecuteAsync(ctx context.Context) {
	go func() {
		defer close(p.successChan)
		defer close(p.errorChan)

		// TODO: pass ctx to wrapped func
		val, err := p.inputFunc()

		if err != nil {
			p.errorChan <- err
		} else {
			p.successChan <- Results[T]{Value: val}
		}
	}()
}

func (p *Promise[T]) Then(successHandler func(results Results[T])) *Promise[T] {
	go func() {
		successHandler(<-p.successChan)
	}()
	return p
}

func (p *Promise[T]) Error(errorHandler func(err error)) *Promise[T] {
	go func() {
		err := <-p.errorChan
		if err != nil {
			errorHandler(err)
		}
	}()
	return p
}
