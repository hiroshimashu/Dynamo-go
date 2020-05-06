package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

type UserPost struct {
	ID        int
	URL       string
	CreatedAt time.Time
	State     string
}

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}

	up, err := readPosts("./MOCK_DATA.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range up {
		fmt.Println("Inserting:", p)
		err = insertPost(cfg, p)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func readPosts(fileName string) ([]UserPost, error) {
	posts := make([]UserPost, 0)

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return posts, err
	}

	err = json.Unmarshal(data, &posts)
	for _, v := range posts {
		v.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(v.CreatedAt))
	}
	if err != nil {
		return posts, err
	}

	return posts, nil
}

func insertPost(cfg aws.Config, up UserPost) error {
	item, err := dynamodbattribute.MarshalMap(up)
	fmt.Println(item)
	if err != nil {
		return err
	}

	svc := dynamodb.New(cfg)
	req := svc.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String("UserPost"),
		Item:      item,
	})
	_, err = req.Send(req.Context())
	if err != nil {
		return err
	}
	return nil
}
