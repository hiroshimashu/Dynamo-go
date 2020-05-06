package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserPost struct {
	ID        string
	URL       string
	CreatedAt time.Time
	State     string
}

func config() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}

	svc := dynamodb.New(cfg)
	params := &dynamodb.ScanInput{
		TableName: aws.String("UserPost"),
	}
	req := svc.ScanRequest(params)
	res, err := req.Send(req.Context())
	if err != nil {
		fmt.Println(err)
	}
	posts := make([]UserPost, 0)
	for _, item := range res.Items {
		posts = append(posts, UserPost{
			ID:    *item["ID"].N,
			URL:   *item["URL"].S,
			State: *item["State"].S,
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
