package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"io"
)

type StreamingHandler struct {
	s3Client *s3.Client
}

func NewStreamingHandler(client *s3.Client) *StreamingHandler {
	return &StreamingHandler{
		s3Client: client,
	}
}

func (h *StreamingHandler) Serve(c *fiber.Ctx) error {
	ctx := c.Context()
	bucket := c.Params("bucket")
	prefix := c.Params("*")

	output, err := h.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(prefix),
	})

	if err != nil {
		return c.SendString(err.Error())
	}
	defer output.Body.Close()

	_, err = io.Copy(c, output.Body)
	return err
}
