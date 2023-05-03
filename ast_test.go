package awseventgenerator

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/webdestroya/awseventgenerator/internal/testutil/astinspector"

	"github.com/stretchr/testify/require"
)

const testDataRoot = "./internal/testdata"

func TestAstSimple(t *testing.T) {
	fset := token.NewFileSet()

	printAst := func(node ast.Node) string {
		var buf bytes.Buffer
		err := printer.Fprint(&buf, fset, node)
		require.NoError(t, err)
		return buf.String()
	}

	require.Equal(t, "[]SomeStruct", printAst(&ast.ArrayType{
		Elt: &ast.Ident{Name: "SomeStruct"},
	}))

	require.Equal(t, "SomeStruct", printAst(&ast.Ident{Name: "SomeStruct"}))
	require.Equal(t, "SomeStruct", printAst(&ast.Ident{Name: "SomeStruct"}))
	require.Equal(t, "*SomeStruct", printAst(&ast.StarExpr{X: &ast.Ident{Name: "SomeStruct"}}))
	require.Equal(t, "type Thinger = interface{}", printAst(
		&ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name:   &ast.Ident{Name: "Thinger"},
					Assign: 1,
					Type: &ast.InterfaceType{Methods: &ast.FieldList{
						Opening: 1,
						Closing: 2,
					}},
				},
			},
		}))

	require.Equal(t, "type Thinger = map[string]interface{}", printAst(
		&ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name:   &ast.Ident{Name: "Thinger"},
					Assign: 1,
					Type: &ast.MapType{
						Key:   &ast.Ident{Name: "string"},
						Value: &ast.Ident{Name: "interface{}"},
					},
				},
			},
		}))
}

func TestASTAll(t *testing.T) {

	files, err := os.ReadDir(testDataRoot)
	require.NoError(t, err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if path.Ext(file.Name()) != ".json" {
			continue
		}

		packageName := strings.TrimSuffix(strings.ToLower(file.Name()), ".json")

		t.Run(packageName, func(t *testing.T) {

			config := &Config{
				AddTestHelpers:          false,
				AlwaysPointerize:        true,
				GenerateEnumValueMethod: true,
			}

			data, err := GenerateFromSchemaFile(path.Join(testDataRoot, file.Name()), config)
			require.NoError(t, err)

			ins, err := astinspector.NewInspector(string(data))
			require.NoError(t, err)

			require.True(t, ins.HasExport("AwsEvent") || ins.HasExport("Root"))
		})

	}
}

func TestASTEnum(t *testing.T) {
	config := &Config{
		AddTestHelpers:          false,
		GenerateEnumValueMethod: true,
	}
	// data, err := GenerateFromSchemaFile("./internal/testdata/enum.json", config)
	data, err := GenerateFromSchemaFile("./internal/testdata/simpleenum.json", config)
	require.NoError(t, err)

	ins, err := astinspector.NewInspector(string(data))
	require.NoError(t, err)
	require.NotNil(t, ins)
	// ins.DumpFile(os.Stdout)
}

func TestASTEnumPointer(t *testing.T) {
	config := &Config{
		AddTestHelpers:          false,
		AlwaysPointerize:        true,
		GenerateEnumValueMethod: true,
	}
	// data, err := GenerateFromSchemaFile("./internal/testdata/enum.json", config)
	data, err := GenerateFromSchemaFile("./internal/testdata/simpleenum.json", config)
	require.NoError(t, err)

	ins, err := astinspector.NewInspector(string(data))
	require.NoError(t, err)
	require.NotNil(t, ins)
}

func TestASTObjAdditional(t *testing.T) {
	config := &Config{
		AddTestHelpers:          false,
		GenerateEnumValueMethod: true,
	}
	data, err := GenerateFromSchemaFile("./internal/testdata/obj_vs_additional.json", config)
	require.NoError(t, err)

	ins, err := astinspector.NewInspector(string(data))
	require.NoError(t, err)
	require.NotNil(t, ins)
	// ins.DumpFile(os.Stdout)
	// ins.Print()
}

func TestASTAdditional2(t *testing.T) {
	config := &Config{
		AddTestHelpers:          false,
		GenerateEnumValueMethod: true,
	}
	data, err := GenerateFromSchemaFile("./internal/testdata/additionalProperties2.json", config)
	require.NoError(t, err)

	ins, err := astinspector.NewInspector(string(data))
	require.NoError(t, err)
	require.NotNil(t, ins)
	// ins.DumpFile(os.Stdout)
	// ins.Print()
}
