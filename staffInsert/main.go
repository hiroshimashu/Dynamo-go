package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Image struct {
	ID        string `dynamodbav:"id"`
	URL       string `dynamodbav:"url"`
	CreatedAt int64  `dynamodbav:"created_at"`
	State     string `dynamodbav:"state"`
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
		ID:        uuid,
		URL:       url,
		CreatedAt: time.Now().Unix(),
		State:     "active",
	}
	// Post image to dynamoDB
	svc := dynamodb.New(cfg)

	av, err := dynamodbattribute.MarshalMap(image)
	if err != nil {
		fmt.Println(err.Error())
		return "failed", err
	}

	req := svc.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item:      av,
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