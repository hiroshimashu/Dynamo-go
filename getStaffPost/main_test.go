package main

import (
	"encoding/json"
	"testing"
)

type Image struct {
	ID        string `dynamodbav:"id"`
	URL       string `dynamodbav:"url"`
	CreatedAt int64  `dynamodbav:"created_at"`
	State     string `dynamodbav:"state"`
}

func TestFindAll(t *testing.T) {
	t.Run("Collectly fetch only the status is active records", func(t *testing.T) {
		res, err := findAll()
		if err != nil {
			t.Error(err)
		}
		var got []Image
		err = json.Unmarshal([]byte(res.Body), &got)
		if err != nil {
			t.Error(err)
		}
		want := 1
		if len(got) != want {
			t.Errorf("want %d but got %d", want, len(got))
		}
	})
}
