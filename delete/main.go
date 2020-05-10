package main

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
)

var activeLimit int = 10
var numOfDeletion int = 10

type UserPost struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	State     string `json:"state"`
}

type UserPosts []UserPost

func (a UserPosts) Len() int           { return len(a) }
func (a UserPosts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UserPosts) Less(i, j int) bool { return a[i].CreatedAt < a[j].CreatedAt }

func config() (int, []UserPost) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}

	filt := expression.Name("state").Equal(expression.Value("active"))
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

	return len(posts), posts

}

func checkNeedUpdate(lenPost int) bool {
	if lenPost > activeLimit {
		return true
	}
	return false
}

// Get sorted slice of UserPost by created_at, Set the older state
func DeactivatePost(ups UserPosts) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}
	svc := dynamodb.New(cfg)
	update := expression.Set(expression.Name("state"), expression.Value("ruuning"))
	condition := expression.Equal(expression.Name("state"), expression.Value("active"))

	expr, err := expression.NewBuilder().WithCondition(condition).WithUpdate(update).Build()
	if err != nil {
		fmt.Println("Condition parsing error")
	}

	// Step1. Get candidate posts and turn that state into RUUNING
	for i := 0; i < numOfDeletion; i++ {
		fmt.Println(ups[i].ID)
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ConditionExpression:       expr.Condition(),
			Key: map[string]dynamodb.AttributeValue{
				"id": {
					N: aws.String(ups[i].ID),
				},
				"created_at": {
					N: aws.String(ups[i].CreatedAt),
				},
			},
			TableName: aws.String("UserPost"),
		}

		req := svc.UpdateItemRequest(input)
		result, err := req.Send(context.Background())
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeConditionalCheckFailedException:
					fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
				case dynamodb.ErrCodeProvisionedThroughputExceededException:
					fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				case dynamodb.ErrCodeResourceNotFoundException:
					fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
					fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			return
		}
		fmt.Println(result)
	}

}
func sortUserPost(ups UserPosts) {
	sort.Sort(ups)
	fmt.Println(ups)
}

func main() {
	lenPosts, posts := config()
	fmt.Println(lenPosts)
	if checkNeedUpdate(lenPosts) {
		sortUserPost(posts)
		DeactivatePost(posts)
	}
}
