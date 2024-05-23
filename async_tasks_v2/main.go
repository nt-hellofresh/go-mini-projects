package main

import (
	"async_tasks_v2/external"
	"async_tasks_v2/task"
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 9900*time.Millisecond)
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(2) // Wait for 2 tasks to call wg.Done()

	t1 := task.NewTask(func(ctx context.Context) (int, error) {
		defer wg.Done()
		result := external.GetValueLongRunningTask()
		return result, nil
	})
	t1.ExecuteAsync(ctx)

	t2 := task.NewTask(func(ctx context.Context) (int, error) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return -1, ctx.Err()
		default:
			return external.WowSuperLongRunningFunction(), nil
		}
	})
	t2.ExecuteAsync(ctx)

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
