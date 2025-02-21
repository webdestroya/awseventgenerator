package testcode

import (
	"testing"

	"github.com/stretchr/testify/require"
	abandoned "github.com/webdestroya/awseventgenerator/internal/testcode/normal_gen/abandoned_gen"
)

func TestAbandoned(t *testing.T) {
	// this just tests the name generation works correctly
	r := abandoned.Root{
		Name:      "jonson",
		Abandoned: "test",
	}
	// the test is the presence of the Abandoned field
	require.NotNil(t, r.Abandoned, "thats the test")
}
