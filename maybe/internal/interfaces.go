package internal

import (
	"maybe/pkg/option"
)

type Repository[T any] interface {
	GetById(id string) option.Maybe[T]
	Save(value T) error
}
