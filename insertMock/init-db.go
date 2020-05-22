package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

type Timestamp int64
type UserPost struct {
	ID        string    `json:"id,omitempty"`
	URL       string    `json:"url,omitempyt"`
	CreatedAt Timestamp `json:"created_at,int,omitempty"`
	State     string    `json:"state,omitempty"`
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	ts := t
	stamp := fmt.Sprint(ts)
	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*t = Timestamp(ts)
	return nil
}

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}

	up, err := readPosts("./generated.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range up {
		fmt.Println("Inserting:", v)
		err = insertPost(cfg, v)
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
	if err != nil {
		return posts, err
	}

	return posts, nil
}

func insertPost(cfg aws.Config, up UserPost) error {
	item, err := dynamodbattribute.MarshalMap(up)
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
