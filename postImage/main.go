package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
)

// payload format
type Image struct {
	fileName      string `json:"filename"`
	string64      string `json:"base64Image"`
	fileExtension string `json:"extension"`
}

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var image Image
	err := json.Unmarshal([]byte(request.Body), &image)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid payload",
		}, err
	}
	result, uploadError := uploadS3(image.string64, image.fileExtension, image.fileExtension)
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

func uploadS3(imageBase64 string, fileExtension string, fileName string) (*s3manager.UploadOutput, error) {
	fmt.Println(imageBase64, fileExtension, fileName)
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}
	// Create S3 service client with a specific Region.
	svc := s3.New(cfg)
	uploader := s3manager.NewUploaderWithClient(svc)
	data, decodeError := decodeBase64(imageBase64, fileExtension)
	if decodeError != nil {
		return nil, decodeError
	}
	wb := new(bytes.Buffer)
	wb.Write(data)

	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("miraikan"),
		Key:         aws.String(os.Getenv("URL") + fileName + "." + fileExtension),
		Body:        wb,
		ContentType: aws.String("image/" + fileExtension),
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
