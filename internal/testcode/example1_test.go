package testcode

import (
	"encoding/json"
	"testing"

	example1 "github.com/webdestroya/awseventgenerator/internal/testcode/example1_gen"
)

func TestExample1(t *testing.T) {
	params := []struct {
		Name           string
		Data           string
		ExpectedResult bool
	}{
		{
			Name: "Blue Sky",
			Data: `{
				"id": 1,
				"name": "Unbridled Optimism 2.0",
				"price": 99.99,
				"tags": [ "happy" ] }`,
			ExpectedResult: true,
		},
		{
			Name: "Missing Price",
			Data: `{
				"id": 1,
				"name": "Unbridled Optimism 2.0",
				"tags": [ "happy" ] }`,
			ExpectedResult: false,
		},
	}

	for _, param := range params {

		prod := &example1.Root{}
		if err := json.Unmarshal([]byte(param.Data), &prod); err != nil {
			if param.ExpectedResult {
				t.Fatal(err)
			}
		} else {
			if !param.ExpectedResult {
				t.Fatal("Expected failure, got success: " + param.Name)
			}
		}
	}
}
