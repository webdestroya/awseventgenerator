package awseventgenerator_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awseventgenerator"
)

const (
	testSchemaFile = "internal/testdata/awsexample.json"
)

var genEntrypointsTestConfig = awseventgenerator.Config{
	AlwaysPointerize: true,
	GenerateEnums:    true,
}

func TestGenerateFromSchemaFile(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		out, err := awseventgenerator.GenerateFromSchemaFile(testSchemaFile, &genEntrypointsTestConfig)
		require.NotNil(t, out)
		require.NoError(t, err)
	})
}

func TestGenerateFromSchemaString(t *testing.T) {
	t.Run("happy", func(t *testing.T) {

		filedata, err := os.ReadFile(testSchemaFile)
		require.NoError(t, err)

		out, err := awseventgenerator.GenerateFromSchemaString(string(filedata), &genEntrypointsTestConfig)
		require.NotNil(t, out)
		require.NoError(t, err)

		expectedOut, err := awseventgenerator.GenerateFromSchemaFile(testSchemaFile, &genEntrypointsTestConfig)
		require.NoError(t, err)
		require.Equal(t, string(expectedOut), string(out))
	})
}

func TestGenerateFromSchema(t *testing.T) {
	t.Run("happy", func(t *testing.T) {

		filedata, err := os.ReadFile(testSchemaFile)
		require.NoError(t, err)

		schema, err := awseventgenerator.Parse(string(filedata), nil)
		require.NoError(t, err)
		require.Equal(t, "file://stringdata.json", schema.ID())

		out, err := awseventgenerator.GenerateFromSchema(schema, &genEntrypointsTestConfig)
		require.NotNil(t, out)
		require.NoError(t, err)

		expectedOut, err := awseventgenerator.GenerateFromSchemaFile(testSchemaFile, &genEntrypointsTestConfig)
		require.NoError(t, err)
		require.Equal(t, string(expectedOut), string(out))
	})
}
