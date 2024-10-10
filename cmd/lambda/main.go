package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/WillChangeThisLater/urlify/pkg/urlify"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	region := "us-east-2"
	bucket := "urlify"
	prefix := "urlify"

	log.Printf("Starting handler for endpoint. region=%s, bucket=%s, prefix=%s\n", region, bucket, prefix)

	log.Printf("event.Headers=%+v\n", event.Headers)
	log.Printf("Checking for content-type header")
	contentType := event.Headers["Content-Type"]
	if contentType == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Content-Type header is missing",
		}, nil
	}
	log.Printf("content-type=%s\n", contentType)

	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Printf("Invalid content-type %s: %v", contentType, err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Invalid content-type: %v", err),
		}, nil
	}

	boundary, ok := params["boundary"]
	if !ok {
		log.Println("Boundary parameter missing")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "boundary parameter missing",
		}, nil
	}
	log.Printf("Parsed boundary: %s\n", boundary)

	// Determine if body needs to be decoded
	var bodyBytes []byte
	if event.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(event.Body)
		if err != nil {
			log.Println("Error decoding Base64 body:", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Could not decode base64 body.",
			}, nil
		}
		bodyBytes = decoded
	} else {
		bodyBytes = []byte(event.Body)
	}
	reader := multipart.NewReader(bytes.NewReader(bodyBytes), boundary)

	for {
		log.Println("Reading part")
		part, err := reader.NextPart()
		if err == io.EOF {
			log.Println("Hit EOF")
			break
		}
		if err != nil {
			log.Println("Error reading part:", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Could not read part from multipart request: %v", err),
			}, nil
		}

		if part.FileName() == "" {
			log.Printf("part %v has no file name: skipping\n", part)
			continue
		}

		fileBytes, err := io.ReadAll(part)
		if err != nil {
			log.Println("Error reading file bytes:", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Could not read file: %v", err),
			}, nil
		}

		url, err := urlify.Urlify(bucket, prefix, region, fileBytes)
		if err != nil {
			log.Println("Error urlifying file: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Error urlifying file: %v", err),
			}, nil
		}

		log.Println("Got presigned url:", url)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       url,
		}, nil
	}

	log.Println("Not allowed: no files uploaded")
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       "No file was detected. Please upload a file",
	}, nil

}

func main() {
	lambda.Start(handler)
}
