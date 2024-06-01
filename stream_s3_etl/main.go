package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"os"
)

func main() {
	//main1()
	main2()
	//main3()
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

// This example doesn't work
func main3() {
	args := getArgs()
	args.validate()

	r := getS3ObjectBody(args)
	defer r.Close()

	pipe, err := NewJsonBatchPipe(r, handleTransformLogic)

	if err != nil {
		log.Fatal(fmt.Errorf("failed to create new json batch pipe: %w", err))
	}

	client := newS3Client(args.awsProfileName)

	// Can't stream an io.Reader of unknown size to S3, which is a shame.
	// Reference issue: https://github.com/aws/aws-sdk-go/issues/122#issuecomment-76826724
	result, err := client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: args.bucket,
		Key:    aws.String("output/transformed.json"),
		Body:   pipe,
	})

	if err != nil {
		log.Println(fmt.Errorf("failed to put object: %w", err))
	}

	log.Println(result)
}
