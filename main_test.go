package awseventgenerator_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awseventgenerator"
	"github.com/webdestroya/awseventgenerator/internal/testutil/genhelpers"
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

func TestGenerate(t *testing.T) {

	fileData := genhelpers.MustRet(os.ReadFile(testSchemaFile))
	schema := genhelpers.MustRet(awseventgenerator.Parse(string(fileData), nil))

	expectedOut := genhelpers.MustRet(awseventgenerator.GenerateFromSchema(schema, &genEntrypointsTestConfig))

	tables := []struct {
		label    string
		from     any
		errMatch any
	}{
		{
			label: "bytearray",
			from:  fileData,
		},
		{
			label: "filepath",
			from:  testSchemaFile,
		},
		{
			label: "stringdata",
			from:  string(fileData),
		},
		{
			label: "schema",
			from:  schema,
		},
		{
			label: "filehandle",
			from:  genhelpers.MustRet(os.Open(testSchemaFile)),
		},
		{
			label: "StringReader",
			from:  strings.NewReader(string(fileData)),
		},

		// Failures
		{
			label:    "MissingFile",
			from:     "/tmp/fakefile",
			errMatch: "no such file or directory",
		},

		{
			label:    "InvalidFrom",
			from:     'x',
			errMatch: "invalid source type",
		},
	}

	for _, table := range tables {
		t.Run(table.label, func(t *testing.T) {
			out, err := awseventgenerator.Generate(table.from, &genEntrypointsTestConfig)
			if table.errMatch != nil {
				require.Error(t, err)

				switch e := table.errMatch.(type) {
				case string:
					require.ErrorContains(t, err, e)
				default:
					require.FailNowf(t, "unknown error matching type???", "got: %T", table.errMatch)
				}

				return
			}

			require.NoError(t, err)
			require.Len(t, out, len(expectedOut))
			require.Equal(t, expectedOut, out)
		})
	}

}

func TestGenerateAndExport(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		dirname := t.TempDir()
		destFile := path.Join(dirname, "subdir", "gen_export.txt")
		err := awseventgenerator.GenerateAndExport(testSchemaFile, destFile, &genEntrypointsTestConfig)
		require.NoError(t, err)

		require.FileExists(t, destFile)

	})
}
