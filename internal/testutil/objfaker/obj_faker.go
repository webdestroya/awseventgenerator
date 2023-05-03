package testwriter

import (
	"encoding/json"

	"github.com/webdestroya/awseventgenerator"
)

type jsonFaker struct {
	root *awseventgenerator.Schema
}

func GenerateFakedJson(schema *awseventgenerator.Schema) (json.RawMessage, error) {
	_ = jsonFaker{}
	return nil, nil
}
