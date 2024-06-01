package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
)

func getS3ObjectBody(args CommandLineArgs) io.ReadCloser {
	client := newS3Client(args.awsProfileName)

	ctx := context.Background()

	output, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: args.bucket,
		Key:    args.prefix,
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
