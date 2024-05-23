package task

import "context"

type Task[T any] struct {
	inputFunc    InputFunc[T]
	ResponseChan chan Response[T]
}

func NewTask[T any](f InputFunc[T]) *Task[T] {
	return &Task[T]{
		inputFunc:    f,
		ResponseChan: make(chan Response[T]),
	}
}

func (t *Task[T]) ExecuteAsync(ctx context.Context) {
	go func() {
		defer close(t.ResponseChan)
		val, err := t.inputFunc(ctx)
		t.ResponseChan <- Response[T]{Value: val, Error: err}
	}()
}
