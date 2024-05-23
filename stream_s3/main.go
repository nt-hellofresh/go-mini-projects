package main

import (
	"context"
	"flag"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := configureHandler()

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}

func configureHandler() *fiber.App {
	awsProfileName := flag.String("profile", "default", "The name of the profile in the shared credentials file")
	flag.Parse()

	client := newS3Client(awsProfileName)

	handler := NewStreamingHandler(client)

	app := fiber.New()
	app.Get("/:bucket/*", handler.Serve)
	return app
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
