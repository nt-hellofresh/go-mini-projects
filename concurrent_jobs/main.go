package main

import (
	"concurrent_jobs/external"
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	wg := &sync.WaitGroup{}
	wg.Add(2) // Wait for 2 tasks to call wg.Done()

	// task 1
	respCh1 := make(chan int)
	defer close(respCh1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		respCh1 <- external.GetValueLongRunningTask()
	}(wg)

	// task 2
	respCh2 := make(chan int)
	defer close(respCh2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		respCh2 <- external.WowSuperLongRunningFunction()
	}(wg)

	v1 := <-respCh1
	v2 := <-respCh2

	wg.Wait()

	fmt.Println("GetValueLongRunningTask response", v1)
	fmt.Println("WowSuperLongRunningFunction response", v2)

	elapsed := time.Since(start)
	fmt.Println("Total duration:", elapsed)
}
