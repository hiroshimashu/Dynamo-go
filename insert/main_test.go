package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-lambda-go/events"
)

func TestInsert_ValidPayload(t *testing.T) {
	input := events.APIGatewayProxyRequest{
		Body:            "{\"id\":\"41\", \"url\":\"https://miraikan.s3-ap-northeast-1.amazonaws.com/images/bastu.png\"}",
		IsBase64Encoded: true,
	}
	expected := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	response, _ := insert(input)
	assert.Equal(t, expected, response)
}
