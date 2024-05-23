package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type RetryClientOpts func(*RetryClient)

type retryClientConfig struct {
	maxRetries       int
	retryStatusCodes []int
	retryMethods     []string
}

type RetryClient struct {
	*http.Client
	retryClientConfig
}

func defaultClient() *RetryClient {
	return &RetryClient{
		Client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
			},
		},
		retryClientConfig: retryClientConfig{
			maxRetries:       1,
			retryStatusCodes: []int{},
			retryMethods:     []string{http.MethodGet},
		},
	}
}

func WithMaxRetries(maxRetries int) RetryClientOpts {
	return func(c *RetryClient) {
		c.maxRetries = maxRetries
	}
}

func WithRetryStatusCodes(statusCodes ...int) RetryClientOpts {
	return func(c *RetryClient) {
		c.retryStatusCodes = statusCodes
	}
}

func WithAllowedRetryMethods(httpMethods ...string) RetryClientOpts {
	return func(c *RetryClient) {
		c.retryMethods = httpMethods
	}
}

func NewRetryClient(opts ...RetryClientOpts) *RetryClient {
	cl := defaultClient()

	for _, opt := range opts {
		opt(cl)
	}

	return cl
}

func (c *RetryClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for retry := 0; retry <= c.maxRetries; retry++ {
		resp, err = c.Client.Do(req)

		if err == nil && resp != nil && !c.shouldRetry(req.Method, resp.StatusCode) {
			return resp, nil
		}

		if retry > 0 {
			fmt.Printf("Request failed (retry %d): %v\n", retry, err)

			// Calculate exponential backoff with jitter
			backoff := (1 << uint(retry)) * time.Second
			backoff += time.Duration(rand.Intn(1000)) * time.Millisecond

			fmt.Printf("Retrying in %v...\n", backoff)
			time.Sleep(backoff)
		}
	}

	return resp, errors.New("max retries attempted")
}

func (c *RetryClient) shouldRetry(httpMethod string, statusCode int) bool {
	retryMethods := make(map[string]bool)
	for _, method := range c.retryMethods {
		retryMethods[method] = true
	}

	retryStatusCodes := make(map[int]bool)
	for _, code := range c.retryStatusCodes {
		retryStatusCodes[code] = true
	}

	return retryMethods[httpMethod] && retryStatusCodes[statusCode]
}
