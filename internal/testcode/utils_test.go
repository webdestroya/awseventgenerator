package testcode

import (
	"encoding/json"
)

func ptr[T any](v T) *T {
	return &v
}

func jsonify(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}
