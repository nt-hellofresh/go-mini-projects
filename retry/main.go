package main

import (
	"errors"
	"log"
	"net/http"
	"sync"
)

type Maybe[T any] struct {
	Error error
	Value T
}

func main() {
	go runServer(90)

	client := NewRetryClient(
		WithMaxRetries(3),
		WithRetryStatusCodes(500, 502, 504),
		WithAllowedRetryMethods(http.MethodGet, http.MethodDelete),
	)

	results := make([]Maybe[string], 1000)
	wg := &sync.WaitGroup{}

	for i := 0; i < len(results); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = doRequest(client)
		}(i)
	}

	wg.Wait()

	var failedCount, total int
	for _, result := range results {
		if result.Error != nil {
			failedCount++
		}
		total++
		//log.Println(result)
	}

	log.Printf("percent success: %f", 100*(1-float64(failedCount)/float64(total)))
}

func doRequest(client *RetryClient) Maybe[string] {

	req, err := http.NewRequest(http.MethodGet, "http://localhost:3000/", nil)

	if err != nil {
		return Maybe[string]{
			Error: err,
			Value: "",
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		return Maybe[string]{
			Error: err,
			Value: "",
		}
	} else if resp.StatusCode == http.StatusInternalServerError {
		return Maybe[string]{
			Error: errors.New("internal server error"),
			Value: "",
		}
	} else if resp.StatusCode == http.StatusOK {
		return Maybe[string]{
			Error: nil,
			Value: "success",
		}
	} else {
		return Maybe[string]{
			Error: errors.New("unknown status code"),
			Value: "",
		}
	}
}
