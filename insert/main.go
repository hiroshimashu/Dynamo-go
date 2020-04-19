package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

type Image struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// HandleRequest :handle sqs's queue
func HandleRequest(ctx context.Context, evt events.SQSEvent) (string, error) {
	s3 := events.S3Event{}
	fmt.Println(evt)
	for _, item := range evt.Records {
		fmt.Printf("The message %s for event source %s = %s \n", item.MessageId, item.EventSource, item.Body)
		if err := json.Unmarshal([]byte(item.Body), &s3); err != nil {
			fmt.Printf("***error*** %#v\n", err)
			return "error", nil
		}
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return "Configuration error", err
	}

	uuid, err := createuuid()
	if err != nil {
		return "creating ID has failed", err
	}

	url := os.Getenv("PREFIX") + s3.Records[0].S3.Object.Key
	fmt.Printf("url variable: %s", url)

	image := Image{
		ID:  uuid,
		URL: url,
	}

	// Post image to dynamoDB
	svc := dynamodb.New(cfg)
	req := svc.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item: map[string]dynamodb.AttributeValue{
			"ID": dynamodb.AttributeValue{S: aws.String(image.ID)},
			"URL": dynamodb.AttributeValue{
				S: aws.String(image.URL),
			},
		},
	})

	_, err = req.Send(req.Context())

	if err != nil {
		return "Insertion Error", err
	}

	return "success", nil
}

func createuuid() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	uu := u.String()
	return uu, nil
}

func main() {
	lambda.Start(HandleRequest)
}
