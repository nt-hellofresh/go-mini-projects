package main

import (
	"context"
	"flag"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
)

func getS3ObjectReader() io.ReadCloser {
	awsProfileName := flag.String("profile", "default", "The name of the profile in the shared credentials file")
	bucket := flag.String("bucket", "", "The name of the bucket")
	prefix := flag.String("prefix", "", "The prefix of the object")

	flag.Parse()

	if *awsProfileName == "" || *bucket == "" || *prefix == "" {
		log.Fatal("aws profile, bucket and prefix command line arguments are required")
	}

	client := newS3Client(awsProfileName)

	ctx := context.Background()

	output, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(*bucket),
		Key:    aws.String(*prefix),
	})

	//client.PutObject(ctx, &s3.PutObjectInput{
	//	Bucket: aws.String(*bucket),
	//	Key:    aws.String("output"),
	//	Body:   output.Body,
	//})

	if err != nil {
		log.Fatal(err)
	}
	return output.Body
}

func newS3Client(awsProfileName *string) *s3.Client {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(*awsProfileName))

	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	return client
}
