package internal

import (
	"maybe/pkg/option"
)

var things option.MaybeFactory[Thing] = option.Factory[Thing]{}

type Thing struct {
	Id string
}
