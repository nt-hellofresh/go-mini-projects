package main

import (
	"context"
	"fmt"
	"log"
	"promise/external"
	"promise/future"
	"sync"
	"time"
)

func errHandler(err error) {
	panic(err)
}

func run() error {
	start := time.Now()

	ctx := context.TODO()
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// store promise as variable before invoking
	p1 := future.Deferred(func() (int, error) {
		return external.GetValueLongRunningTask(), nil
	}).Then(func(r future.Results[int]) {
		defer wg.Done()
		fmt.Println("GetValueLongRunningTask response", r.Value)
	}).Error(errHandler)

	p1.ExecuteAsync(ctx)

	// inline style of invoking promise
	future.Deferred(func() (int, error) {
		return external.WowSuperLongRunningFunction(), nil
	}).
		Then(func(r future.Results[int]) {
			defer wg.Done()
			fmt.Println("WowSuperLongRunningFunction response", r.Value)
		}).
		Error(errHandler).
		ExecuteAsync(ctx)

	// Blocks until all promises are completed
	wg.Wait()

	// TODO: examples promise chaining

	elapsed := time.Since(start)
	fmt.Println("Total duration:", elapsed)

	return nil
}

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
	}
}
