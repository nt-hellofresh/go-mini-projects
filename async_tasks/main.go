package main

import (
	"async_tasks/external"
	"async_tasks/task"
	"fmt"
	"sync"
	"time"
)

// Helper function to convert method signature of
// inner functions to desirable signature required
// for task.NewTask
func Wrap(f func() int) task.InputFunc {
	return func() (any, error) {
		result := f()
		return result, nil
	}
}

func main() {
	start := time.Now()

	wg := &sync.WaitGroup{}

	wg.Add(2) // Wait for 2 tasks to call wg.Done()

	t1 := task.NewTask(Wrap(external.GetValueLongRunningTask), wg)
	t1.ExecuteAsync()

	t2 := task.NewTask(Wrap(external.WowSuperLongRunningFunction), wg)
	t2.ExecuteAsync()

	r1 := <-t1.ResponseChan
	r1.PanicOnError()
	fmt.Println("GetValueLongRunningTask response", r1)

	r2 := <-t1.ResponseChan
	r2.PanicOnError()
	fmt.Println("WowSuperLongRunningFunction response", r2)

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("Total duration:", elapsed)
}
