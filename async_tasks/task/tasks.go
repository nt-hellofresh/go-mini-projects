package task

import (
	"log"
	"sync"
)

type Response struct {
	Value any
	Error error
}

func (r Response) PanicOnError() {
	if r.Error != nil {
		log.Panicln(r.Error)
	}
}

type ForkedTask struct {
	inputFunc    InputFunc
	ResponseChan chan Response
	wg           *sync.WaitGroup
}

type InputFunc func() (any, error)

func NewTask(f InputFunc, wg *sync.WaitGroup) *ForkedTask {
	return &ForkedTask{
		inputFunc:    f,
		ResponseChan: make(chan Response),
		wg:           wg,
	}
}

func (t *ForkedTask) ExecuteAsync() {
	go func() {
		defer t.wg.Done()
		defer close(t.ResponseChan)
		val, err := t.inputFunc()
		t.ResponseChan <- Response{Value: val, Error: err}
	}()
}
