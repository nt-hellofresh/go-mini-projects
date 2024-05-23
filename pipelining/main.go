package main

import (
	"fmt"
	"log"
	"pipelining/pipeline"
)

func main() {
	// Declare pipeline
	pl := makePipeline()

	toResults := Convert(toOrder)

	results := make([]order, 0)
	input := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	for _, value := range input {
		res := pl.Process(value)
		results = append(results, toResults(res))
	}

	fmt.Printf("results=%v\n", results)
}

func makePipeline() *pipeline.Pipeline[int] {
	//pl := pipeline.StartWith(Map(square))
	//pl.Then(Map(timesnine))
	//pl.Then(Map(plusone))
	//
	//return pl

	return pipeline.Define(
		pipeline.WithStep(Map(square)),
		pipeline.WithStep(Map(timesnine)),
		pipeline.WithStep(Map(plusone)),
	)

}

func Convert[A any, B any](fn func(value A) (B, error)) func(value A) B {
	return func(value A) B {
		res, err := fn(value)

		if err != nil {
			fail(err)
		}

		return res
	}
}

func Map(fn func(value int) (int, error)) func(value int) int {
	return func(value int) int {
		res, err := fn(value)

		if err != nil {
			fail(err)
		}

		return res
	}
}

func fail(err error) {
	log.Fatal(err)
}
