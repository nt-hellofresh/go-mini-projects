package main

import (
	"context"
	"time"
)

type inputParam struct {
	idx int
}

func (ip inputParam) Equals(other inputParam) bool {
	return ip.idx == other.idx
}

func (ip inputParam) LessThan(other inputParam) bool {
	return ip.idx < other.idx
}

func (ip inputParam) Instance() inputParam {
	return ip
}

type output struct {
	input inputParam
}

func mySlowFunction(ctx context.Context, param inputParam) (output, error) {
	time.Sleep(2 * time.Second)

	return output{
		input: param,
	}, nil
}
