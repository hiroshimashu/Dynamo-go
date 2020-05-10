package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestJsonParse(t *testing.T) {
	t.Run("correctly parsing json", func(t *testing.T) {
		want := []UserPost{
			UserPost{
				1,
				"http://dummyimage.com/224x118.jpg/5fa2dd/ffffff",
				1257894000,
				"activez",
			},
		}

		got, err := testReadPosts("./test.json")
		fmt.Println(got)
		if err != nil {
			fmt.Println("parsing error")
		}

		reflect.DeepEqual(got, want)
	})
}

func testReadPosts(fileName string) ([]UserPost, error) {
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
