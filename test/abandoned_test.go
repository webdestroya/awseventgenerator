package test

import (
	"testing"

	abandoned "github.com/webdestroya/awseventgenerator/test/abandoned_gen"
)

func TestAbandoned(t *testing.T) {
	// this just tests the name generation works correctly
	str := "jonson"
	r := abandoned.Root{
		Name:      &str,
		Abandoned: &abandoned.Abandoned{},
	}
	// the test is the presence of the Abandoned field
	if r.Abandoned == nil {
		t.Fatal("thats the test")
	}
}
