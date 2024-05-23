package main

import (
	"context"
	"fmt"
	"lru_cache/cache"
	"sync"
)

func main() {
	c := createCache()
	myCachedFunction := cache.ReadThroughCache(mySlowFunction, c)

	ctx := context.Background()

	wg := &sync.WaitGroup{}

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, _ = myCachedFunction(ctx, inputParam{idx})
		}(i)
	}

	wg.Wait()

	for i := 0; i < 200; i++ {
		res, _ := myCachedFunction(ctx, inputParam{i})
		fmt.Printf("output: %+v\n", res)
	}
}

func createCache() cache.Cache[inputParam, output] {
	//return cache.NewInfInMemoryCache[inputParam, output]()
	return cache.NewLRUCache[inputParam, output](198)
}
