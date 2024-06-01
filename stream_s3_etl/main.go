package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	//main1()
	main2()
}

func main1() {
	args := getArgs()
	args.validate()

	r := getS3ObjectBody(args)
	defer r.Close()

	if err := configurePipe().PipeBatchJson(os.Stdout, r); err != nil {
		// TODO: rollback writer results
		log.Fatal(fmt.Errorf("failed to process the S3 object: %w", err))
	}

	log.Println("Successfully processed data from input to output")
}

func main2() {
	args := getArgs()
	args.validate()

	r := getS3ObjectBody(args)
	defer r.Close()

	pipe, err := NewJsonBatchPipe(r, handleTransformLogic)

	if err != nil {
		log.Fatal(fmt.Errorf("failed to create new json batch pipe: %w", err))
	}

	fd, err := os.Create("output.json")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create output file: %w", err))
	}

	n, err := io.Copy(fd, pipe)

	if err != nil {
		log.Fatal(fmt.Errorf("failed to copy pipe to stdout: %w", err))
	}

	log.Printf("Successfully processed %d bytes from input to output", n)
}
