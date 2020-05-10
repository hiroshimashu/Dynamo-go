package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
)

type UserPost struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	State     string `json:"state"`
}

func config() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}

	filt := expression.Name("created_at").GreaterThan(expression.Value(1257894000)).And(expression.Name("state").Equal(expression.Value("deactive")))
	proj := expression.NamesList(expression.Name("created_at"), expression.Name("state"), expression.Name("id"), expression.Name("url"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		fmt.Println(err)
	}

	svc := dynamodb.New(cfg)
	params := &dynamodb.ScanInput{
		TableName:                 aws.String("UserPost"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}
	req := svc.ScanRequest(params)
	res, err := req.Send(req.Context())
	if err != nil {
		fmt.Println(err)
	}
	posts := make([]UserPost, 0)
	for _, item := range res.Items {
		posts = append(posts, UserPost{
			ID:        *item["id"].N,
			CreatedAt: *item["created_at"].N,
			URL:       *item["url"].S,
			State:     *item["state"].S,
		})
	}

	response, err := json.Marshal(posts)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(response))

}

func main() {
	config()
}
