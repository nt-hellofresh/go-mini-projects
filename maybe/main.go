package main

import (
	"log"
	"maybe/internal"
	"maybe/pkg/option"
)

func main() {
	var repo internal.Repository[internal.Thing] = &internal.ThingRepository{}

	printValue(repo.GetById("1"))
	printValue(repo.GetById("2"))
	printValue(repo.GetById("3"))
}

func printValue(result option.Maybe[internal.Thing]) {
	if result.HasError() {
		log.Printf("Error with result: %v\n", result.Error)
	} else {
		log.Printf("Result has value: %+v\n", result.Value())
	}
}
