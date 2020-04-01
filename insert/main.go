package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Movie struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func insert(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var movies Movie
	err := json.Unmarshal([]byte(request.Body), &movies)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid payload",
		}, nil
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while retrieving AWS credentials",
		}, nil
	}
	svc := dynamodb.New(cfg)
	req := svc.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item: map[string]dynamodb.AttributeValue{
			"ID": dynamodb.AttributeValue{
				S: aws.String(movies.ID),
			},
			"Name": dynamodb.AttributeValue{
				S: aws.String(movies.URL),
			},
		},
	})
	_, err = req.Send(req.Context())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while inserting movie to DynamoDB",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(insert)
}
