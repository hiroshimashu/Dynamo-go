package main

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-lambda-go/events"
)

func encodeImageToBase64() string {

	file, _ := os.Open("maru.png")
	defer file.Close()

	fi, _ := file.Stat() //FileInfo interface
	size := fi.Size()    //ファイルサイズ

	data := make([]byte, size)
	file.Read(data)

	return base64.StdEncoding.EncodeToString(data)
}

func TestInsert_ValidPayload(t *testing.T) {
	input := events.APIGatewayProxyRequest{
		Body:            "{\"filename\":\"sample\",\"base64Image\":\"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\",\"extension\":\"png\"}",
		IsBase64Encoded: true,
	}
	expected := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	response, _ := handleRequest(input)
	assert.Equal(t, expected, response)
}
