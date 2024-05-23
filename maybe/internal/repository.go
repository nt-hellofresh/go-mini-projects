package internal

import (
	"errors"
	"fmt"
	"maybe/pkg/option"
)

type ThingRepository struct {
}

func (r *ThingRepository) GetById(id string) option.Maybe[Thing] {
	if id == "1" {
		return things.Some(Thing{
			Id: id,
		})
	} else {
		return things.None(notFoundError(id))
	}
}

func (r *ThingRepository) Save(value Thing) error {
	return nil
}

func notFoundError(id string) error {
	msg := fmt.Sprintf("couldn't find thing by id %v", id)
	return errors.New(msg)
}
