package genhelpers

import (
	"testing"

	"github.com/jmespath/go-jmespath"
	"github.com/stretchr/testify/require"
)

type AnyValStruct struct {
	Thing string `json:"thinger"`
}

// func Ptrize[T any](v T) *T {
// 	return &v
// }

// func Jsonify(v any) string {
// 	data, err := json.Marshal(v)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return string(data)
// }

func RequireJmesMatch(t *testing.T, data any, expr string, expected any, locationNote string) {
	t.Helper()

	result, err := jmespath.Search(expr, data)
	require.NoErrorf(t, err, locationNote, data)

	require.IsType(t, expected, result, locationNote, data)
	require.EqualValues(t, expected, result, locationNote, data)
}

func JmesMatch(t *testing.T, data any, expr string) any {
	t.Helper()

	result, err := jmespath.Search(expr, data)
	require.NoError(t, err)
	return result
}

func MustRet[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
