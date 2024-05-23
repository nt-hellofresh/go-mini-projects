package task

import (
	"context"
	"log"
)

type Response[T any] struct {
	Value T
	Error error
}

func (r Response[T]) PanicOnError() {
	if r.Error != nil {
		log.Panicln(r.Error)
	}
}

type InputFunc[T any] func(ctx context.Context) (T, error)
