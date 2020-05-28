package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
)

type GeneratedImage struct {
	ID        string `json:"id,omitempty"`
	URL       string `json:"url,omitempyt"`
	CreatedAt string `json:"created_at,int,omitempty"`
	State     string `json:"state,omitempty"`
}

func findAll() (events.APIGatewayProxyResponse, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while retrieving AWS credentials",
		}, nil
	}

	svc := dynamodb.New(cfg)

	filt := expression.Name("state").Equal(expression.Value("active"))
	proj := expression.NamesList(expression.Name("created_at"), expression.Name("state"), expression.Name("id"), expression.Name("url"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	req := svc.ScanRequest(&dynamodb.ScanInput{
		TableName:                 aws.String(os.Getenv("TABLE_NAME")),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	})
	res, err := req.Send(req.Context())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while scanning DynamoDB",
		}, nil
	}

	generatedImages := make([]GeneratedImage, 0)
	for _, item := range res.Items {
		generatedImages = append(generatedImages, GeneratedImage{
			ID:        *item["id"].S,
			URL:       *item["url"].S,
			CreatedAt: *item["created_at"].N,
			State:     *item["state"].S,
		})
	}

	response, err := json.Marshal(generatedImages)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while decoding to string value",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(response),
	}, nil
}

func main() {
	lambda.Start(findAll)
}
