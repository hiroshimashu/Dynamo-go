package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
)

// Image Payload format
type Image struct {
	FileName      string `json:"filename"`
	String64      string `json:"base64Image"`
	FileExtension string `json:"extension"`
}

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// var image Image
	fmt.Println(request.Body)
	/*
		err := json.Unmarshal([]byte(request.Body), &image)
		if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
			fmt.Printf("%#v", typeErr)
		}
	*/
	body, decodeError := decodeBase64(request.Body, "png")
	fmt.Println(body)
	if decodeError != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid payload",
		}, decodeError
	}

	result, uploadError := uploadS3(body)
	fmt.Println(result)
	if uploadError != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Upload failed",
		}, uploadError
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: result.Location,
	}, nil
}

func uploadS3(imageBody []byte) (*s3manager.UploadOutput, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}
	// Create S3 service client with a specific Region.
	svc := s3.New(cfg)
	uploader := s3manager.NewUploaderWithClient(svc)
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("miraikan"),
		Key:         aws.String(os.Getenv("URL") + "maru" + ".png"),
		Body:        bytes.NewReader(imageBody),
		ContentType: aws.String("image/png"),
	})

	return res, nil
}

func decodeBase64(imageBase64 string, fileExtension string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}

func main() {
	lambda.Start(handleRequest)
}
