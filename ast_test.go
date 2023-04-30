package awseventgenerator

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/webdestroya/awseventgenerator/internal/testutil/astinspector"

	"github.com/stretchr/testify/require"
)

const testDataRoot = "./internal/testdata"

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

			require.True(t, ins.HasExport("AwsEvent"))
		})

	}
}

func TestBasicAST(t *testing.T) {
	config := &Config{
		AddTestHelpers: true,
	}
	// data, err := GenerateFromSchemaFile("./internal/testdata/enum.json", config)
	data, err := GenerateFromSchemaFile("./internal/testdata/ecstaskchange.json", config)
	require.NoError(t, err)

	ins, _ := astinspector.NewInspector(string(data))
	ins.DumpFile(os.Stdout)
	// ins.Print()

	require.True(t, ins.HasExport("AwsEvent"))
	require.True(t, ins.HasExport("ECSTaskStateChange"))
	require.True(t, ins.HasExport("Details"))

	require.Equal(t, "[]NetworkBindingDetails", ins.GetStructField("ContainerDetails.NetworkBindings").Type)
	require.Equal(t, "*time.Time", ins.GetStructField("AwsEvent.Time").Type)

	imp := importer.Default()
	timePkg, _ := imp.Import("time")
	conf := types.Config{
		Importer: imp,
	}
	pkg, err := conf.Check(ins.PackageName(), ins.FSet, []*ast.File{ins.File}, nil)
	require.NoError(t, err)
	pkg.SetImports([]*types.Package{
		timePkg,
	})

	// pkg := types.NewPackage(ins.PackageName(), ins.PackageName())
	exprStr := `*(AwsEvent{Account: ptr("test")}).Account == "test"`
	exprStr = `"test"`
	exprStr = `var string thing = "test"\nthing`
	exprStr = `ternaryStr(*(AwsEvent{Account: ptr("test")}).Account == "test", "yes", "no")`

	// from go/types/eval_test.go:250
	expr, err := parser.ParseExprFrom(ins.FSet, "eval", exprStr, 0)
	require.NoError(t, err, "ParseExprFrom")

	info := &types.Info{
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}
	if err := types.CheckExpr(ins.FSet, pkg, token.NoPos, expr, info); err != nil {
		// require.NoError(t, err, "CheckExpr")
	}
	fmt.Println("TYPES", *info)

	res, err := types.Eval(ins.FSet, pkg, token.NoPos, exprStr)
	// res, err := types.Eval(ins.FSet, pkg, token.NoPos, `*(AwsEvent{Account: ptr("test")}).Account`)
	// res, err := types.Eval(ins.FSet, pkg, token.NoPos, "AwsEventSource")
	require.NoError(t, err, "Eval")
	t.Logf("EVAL: %T %v", res, res)
	t.Logf("EVAL: Kind=%v, ValueType=%T", res.Type.String(), res.Value)
	if res.Value != nil {
		t.Logf("EVAL: Value=%v", res.Value.String())
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

	ins, _ := astinspector.NewInspector(string(data))
	ins.DumpFile(os.Stdout)
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

	ins, _ := astinspector.NewInspector(string(data))
	ins.DumpFile(os.Stdout)
}
