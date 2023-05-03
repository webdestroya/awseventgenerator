package testwriter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"io"
	"sort"
	"strings"

	"github.com/webdestroya/awseventgenerator/internal/testutil/astinspector"
)

const (
	formatCodeOutput = true

	jsonTagPrefix = "`json:\""
	jsonTagSuffix = "\"`"

	combineJsonAndFieldTest = true
)

type TestWriter struct {
	inspectors map[string]*astinspector.Inspector
	folders    map[string]string

	imports map[string]struct{}
}

func NewTestWriter() *TestWriter {
	return &TestWriter{
		inspectors: make(map[string]*astinspector.Inspector),
		folders:    make(map[string]string),
		imports:    make(map[string]struct{}),
	}
}

func (tw *TestWriter) Add(data []byte, pkgName, folderName string) error {
	ins, err := astinspector.NewInspector(data)
	if err != nil {
		return err
	}

	tw.inspectors[pkgName] = ins
	tw.folders[pkgName] = folderName

	return nil
}

func (tw *TestWriter) Generate() ([]byte, error) {
	var buf bytes.Buffer

	tw.imports["testing"] = struct{}{}
	tw.imports["time"] = struct{}{}
	// tw.imports["github.com/jmespath/go-jmespath"] = struct{}{}
	tw.imports["github.com/stretchr/testify/require"] = struct{}{}

	for pkgName := range tw.inspectors {
		fmt.Fprintf(&buf, `func TestGenerated_%s(t *testing.T) {`+"\n", pkgName)
		fmt.Fprintln(&buf)
		err := tw.generateHelperValues(&buf)
		if err != nil {
			return nil, err
		}
		fmt.Fprintln(&buf)
		// fmt.Fprintf(&buf, `  t.Run("%s", func(t *testing.T) {`+"\n", pkgName)
		if err := tw.generateForPackage(&buf, pkgName); err != nil {
			return nil, err
		}
		// fmt.Fprintln(&buf, `  })`)
		fmt.Fprintln(&buf)
		fmt.Fprintln(&buf, `}`)
	}

	var finalBuf bytes.Buffer

	fmt.Fprintln(&finalBuf, `// Code generated by awseventgenerator/internal/generators/testcode. DO NOT EDIT.`)
	fmt.Fprintln(&finalBuf)
	fmt.Fprintln(&finalBuf, `package testsuitegenerated`)
	fmt.Fprintln(&finalBuf)
	fmt.Fprintln(&finalBuf, `import (`)
	for k := range tw.imports {
		fmt.Fprintf(&finalBuf, `  "%s"`+"\n", k)
	}
	fmt.Fprintln(&finalBuf)
	for k, v := range tw.folders {
		fmt.Fprintf(&finalBuf, "  %s \"github.com/webdestroya/awseventgenerator/internal/testcode/%s\"\n", k, v)
	}
	fmt.Fprintln(&finalBuf, `)`)
	fmt.Fprintln(&finalBuf)

	if _, err := buf.WriteTo(&finalBuf); err != nil {
		return nil, err
	}

	// return buf.Bytes(), nil
	if formatCodeOutput {

		formattedBytes, err := format.Source(finalBuf.Bytes())
		if err != nil {
			return nil, err
		}

		return formattedBytes, nil
	}
	return finalBuf.Bytes(), nil
}

func (tw *TestWriter) generateHelperValues(buf io.Writer) error {
	fmt.Fprintln(buf, `strVal := "someString"`)
	fmt.Fprintln(buf, `floatVal := float64(1232.1424)`)
	fmt.Fprintln(buf, `intVal := int64(1232)`)
	fmt.Fprintln(buf, `timeVal := time.Now().UTC()`)
	fmt.Fprintln(buf, `trueVal := true`)
	fmt.Fprintln(buf, "anyVal := struct{Thing string `json:\"thinger\"`}{Thing: \"anywayanyday\"}")
	fmt.Fprintln(buf)
	fmt.Fprintln(buf, `require.IsType(t, *new(string), strVal)`)
	fmt.Fprintln(buf, `require.IsType(t, *new(float64), floatVal)`)
	fmt.Fprintln(buf, `require.IsType(t, *new(int64), intVal)`)
	fmt.Fprintln(buf, `require.IsType(t, *new(time.Time), timeVal)`)
	fmt.Fprintln(buf, `require.IsType(t, *new(bool), trueVal)`)
	fmt.Fprintln(buf, `_ = anyVal`)
	fmt.Fprintln(buf)

	return nil
}

func (tw *TestWriter) generateForPackage(buf io.Writer, pkgName string) error {

	ins := tw.inspectors[pkgName]

	constants := make(map[string]*ast.ValueSpec)
	enumTypes := make(map[string]map[string]*ast.ValueSpec)
	aliasTypes := make(map[string]*ast.TypeSpec)
	structTypes := make(map[string]*ast.StructType)

	for tname, scopeObj := range ins.File.Scope.Objects {
		switch v := scopeObj.Decl.(type) {
		case *ast.ValueSpec:
			if v.Type != nil {

				// fmt.Println("DEBUG ENUM", tname)
				// ast.Print(ins.FSet, v.Type)

				enumName := v.Type.(*ast.Ident).Name

				if _, ok := enumTypes[enumName]; !ok {
					enumTypes[enumName] = make(map[string]*ast.ValueSpec)
				}

				enumTypes[enumName][tname] = v
			} else {
				constants[tname] = v
			}
		case *ast.TypeSpec:
			if v.Assign.IsValid() {
				aliasTypes[tname] = v
				continue
			}

			switch vt := v.Type.(type) {
			case *ast.StructType:
				structTypes[tname] = vt
			case *ast.Ident:
				// ignore, this is probably an Enum type?
			default:
				fmt.Printf("UNHANDLED TYPESPEC: %s = %T\n", tname, vt)
				ast.Print(ins.FSet, v)
			}
		}
	}

	if len(constants) > 0 {
		fmt.Fprintln(buf)
		fmt.Fprintln(buf, `t.Run("constants", func(t *testing.T){`)
		for _, cname := range getSortedKeys(constants) {
			con := constants[cname]
			conVal := con.Values[0].(*ast.BasicLit)

			if conVal.Kind != token.STRING {
				fmt.Fprintf(buf, `  // Unknown type?? %s for %s.%s`+"\n", conVal.Kind.String(), pkgName, cname)
				continue
			}
			fmt.Fprintf(buf, `  require.Equal(t, %s, %s.%s)`+"\n", conVal.Value, pkgName, cname)

			// fmt.Println("DEBUG CONSTANT", cname)
			// ast.Print(ins.FSet, con)
		}
		fmt.Fprintln(buf, `})`)
	}

	if len(enumTypes) > 0 {
		fmt.Fprintln(buf)
		fmt.Fprintln(buf, `t.Run("enums", func(t *testing.T){`)
		for _, ename := range getSortedKeys(enumTypes) {
			enumType := enumTypes[ename]

			fmt.Fprintf(buf, `t.Run("%s", func(t *testing.T){`+"\n", ename)
			for _, cname := range getSortedKeys(enumTypes[ename]) {
				con := enumType[cname]
				conVal := con.Values[0].(*ast.BasicLit)

				fmt.Fprintf(buf, `  require.Equal(t, %s, string(%s.%s))`+"\n", conVal.Value, pkgName, cname)
				fmt.Fprintf(buf, `  require.Contains(t, %s.%s.Values(), %s.%s)`+"\n", pkgName, cname, pkgName, cname)
				fmt.Fprintln(buf)

				// fmt.Println("DEBUG CONSTANT", cname)
				// ast.Print(ins.FSet, con)

			}
			fmt.Fprintln(buf, `})`)
		}
		fmt.Fprintln(buf, `})`)
	}

	if len(aliasTypes) > 0 {
		fmt.Fprintln(buf)
		fmt.Fprintln(buf, `t.Run("aliases", func(t *testing.T){`)
		for _, k := range getSortedKeys(aliasTypes) {
			fmt.Fprintf(buf, `  require.IsType(t, *new(%s), *new(%s.%s))`+"\n", tw.genTypeName(ins, pkgName, aliasTypes[k].Type), pkgName, k)
		}
		fmt.Fprintln(buf, `})`)
	}

	if len(structTypes) > 0 {
		fmt.Fprintln(buf)
		fmt.Fprintln(buf, `t.Run("structs", func(t *testing.T){`)
		for _, k := range getSortedKeys(structTypes) {
			fmt.Fprintf(buf, `t.Run("%s", func(t *testing.T){`+"\n", k)
			v := structTypes[k]

			fmt.Fprint(buf, `genStruct := &`)
			tw.genStructFakeValue(buf, ins, pkgName, k, v, true)
			fmt.Fprintln(buf)

			if combineJsonAndFieldTest {
				fmt.Fprintln(buf)
				fmt.Fprintln(buf, `// JSON tests`)
				fmt.Fprintln(buf, `{`)
			} else {
				fmt.Fprintln(buf, `t.Run("json", func(t *testing.T){`)
			}
			tw.genStructMarshalTests(buf, ins, pkgName, k, v)
			if combineJsonAndFieldTest {
				fmt.Fprintln(buf, `}`)
			} else {
				fmt.Fprintln(buf, `})`)
			}

			if combineJsonAndFieldTest {
				fmt.Fprintln(buf)
				fmt.Fprintln(buf, `// Struct Field Tests`)
				fmt.Fprintln(buf, `{`)
			} else {
				fmt.Fprintln(buf, `t.Run("fields", func(t *testing.T){`)
			}
			for _, f := range v.Fields.List {
				tw.genStructFieldTest(buf, ins, pkgName, k, v, f)
			}
			if combineJsonAndFieldTest {
				fmt.Fprintln(buf, `}`)
			} else {
				fmt.Fprintln(buf, `})`)
			}

			fmt.Fprintln(buf, `})`)
			fmt.Fprintln(buf)
		}
		fmt.Fprintln(buf, `})`)
	}

	return nil
}

func (tw *TestWriter) genStructMarshalTests(buf io.Writer, ins *astinspector.Inspector, pkgName, structName string, strct *ast.StructType) {
	tw.imports["encoding/json"] = struct{}{}

	fmt.Fprintln(buf, `jsonOut, err := json.Marshal(genStruct)`)
	fmt.Fprintln(buf, `require.NoError(t, err)`)
	fmt.Fprintln(buf)

	fmt.Fprintf(buf, `unmarObj := &%s.%s{}`+"\n", pkgName, structName)
	fmt.Fprintln(buf, `require.NoError(t, json.Unmarshal(jsonOut, unmarObj))`)
	fmt.Fprintln(buf)
	fmt.Fprintln(buf, `jsonOut2, err := json.Marshal(unmarObj)`)
	fmt.Fprintln(buf, `require.NoError(t, err)`)
	fmt.Fprintln(buf, `require.JSONEq(t, string(jsonOut), string(jsonOut2))`)
	fmt.Fprintln(buf)
	fmt.Fprintln(buf, `var jsearch interface{}`)
	fmt.Fprintln(buf, `require.NoError(t, json.Unmarshal(jsonOut, &jsearch))`)

	for _, f := range strct.Fields.List {
		fname := f.Names[0].Name
		if f.Tag == nil {
			continue
		}

		if !strings.HasPrefix(f.Tag.Value, jsonTagPrefix) {
			fmt.Println("ABORT", f.Tag.Value)
			continue
		}

		tagKey, _, _ := strings.Cut(strings.TrimPrefix(strings.TrimSuffix(f.Tag.Value, jsonTagSuffix), jsonTagPrefix), `,`)

		if tagKey == "-" {
			continue
		}

		ftypeStr := printNode(ins, f.Type)
		switch ftypeStr {
		case "*string", "string":
			fmt.Fprintf(buf, "requireJmesMatch(t, jsearch, `\"%s\"`, strVal, \"%s\")\n", tagKey, fname)
		case "*float64", "float64":
			fmt.Fprintf(buf, "requireJmesMatch(t, jsearch, `\"%s\"`, floatVal, \"%s\")\n", tagKey, fname)
		case "*int64", "int64":
			fmt.Fprintf(buf, "requireJmesMatch(t, jsearch, `\"%s\"`, intVal, \"%s\")\n", tagKey, fname)
		case "*time.Time", "time.Time":
			fmt.Fprintf(buf, "requireJmesMatch(t, jsearch, `\"%s\"`, string(mustRet(timeVal.MarshalText())), \"%s\")\n", tagKey, fname)
		case "*bool", "bool":
			fmt.Fprintf(buf, "requireJmesMatch(t, jsearch, `\"%s\"`, trueVal, \"%s\")\n", tagKey, fname)
		default:
			fmt.Fprintf(buf, "require.GreaterOrEqual(t, jmesMatch(t, jsearch, `length(\"%s\")`).(float64), 1.0)\n", tagKey)

			// fmt.Fprintf(buf, `/* %s */`+"\n", fname)
		}
	}
	fmt.Fprintln(buf)

}

func (tw *TestWriter) genStructFakeValue(buf io.Writer, ins *astinspector.Inspector, pkgName, structName string, strct *ast.StructType, prefixType bool) {
	if prefixType {
		fmt.Fprintf(buf, `%s.%s`, pkgName, structName)
	}
	fmt.Fprintln(buf, `{`)
	for _, f := range strct.Fields.List {
		fname := f.Names[0].Name
		fmt.Fprintf(buf, `%s: `, fname)
		tw.genBaseFakeValue(buf, ins, pkgName, f.Type, false)

		fmt.Fprintln(buf, `,`)
	}
	fmt.Fprint(buf, `}`)
}

func (tw *TestWriter) genTypeName(ins *astinspector.Inspector, pkgName string, ftype ast.Node) string {
	switch t := ftype.(type) {
	case *ast.StarExpr:
		return fmt.Sprintf(`*%s`, tw.genTypeName(ins, pkgName, t.X))
	case *ast.Ident:
		if t.Obj == nil {
			return printNode(ins, t)
		}
		return fmt.Sprintf(`%s.%s`, pkgName, t.Name)
	case *ast.MapType:
		return fmt.Sprintf(`map[%s]%s`, tw.genTypeName(ins, pkgName, t.Key), tw.genTypeName(ins, pkgName, t.Value))
	case *ast.ArrayType:
		return fmt.Sprintf(`[%s]%s`, printNode(ins, t.Len), tw.genTypeName(ins, pkgName, t.Elt))
	case *ast.InterfaceType:
		if t.Methods.NumFields() == 0 {
			return "interface{}"
		}
		return fmt.Sprintf(`/* genTypeName:IFACE_WITH_METHODS? %T */`, t)
	default:
		return fmt.Sprintf(`/* genTypeName:? %T */`, t)
	}
}

// the inner most type
func (tw *TestWriter) genBaseFakeValue(buf io.Writer, ins *astinspector.Inspector, pkgName string, ftype ast.Node, insideArray bool) {

	isPointer := false
	if innerType, ok := ftype.(*ast.StarExpr); ok {
		isPointer = true
		ftype = innerType.X
	}

	if isPointer {
		fmt.Fprint(buf, `&`)
	}

	ftypeStr := printNode(ins, ftype)
	switch tv := ftype.(type) {
	case *ast.SelectorExpr:
		switch ftypeStr {
		case "time.Time":
			fmt.Fprint(buf, `timeVal`)
		}

	case *ast.ArrayType:
		fmt.Fprintf(buf, `%s{`, tw.genTypeName(ins, pkgName, tv))
		tw.genBaseFakeValue(buf, ins, pkgName, tv.Elt, true)
		fmt.Fprint(buf, `}`)
	case *ast.InterfaceType:
		if tv.Methods.NumFields() == 0 {
			// this is the generic interface{}
			fmt.Fprint(buf, `anyVal`)

		} else {
			fmt.Fprintf(buf, `??MULTI_FIELD_INTERFACE::%T /* %s */`, tv, ftypeStr)
		}

	case *ast.MapType:

		if insideArray {
			fmt.Fprintln(buf, `{`)
		} else {
			fmt.Fprintf(buf, `%s{`+"\n", tw.genTypeName(ins, pkgName, tv))
		}
		tw.genBaseFakeValue(buf, ins, pkgName, tv.Key, false)
		fmt.Fprint(buf, `: `)
		tw.genBaseFakeValue(buf, ins, pkgName, tv.Value, true)
		fmt.Fprintln(buf, `,`)
		fmt.Fprint(buf, `}`)

	case *ast.Ident:
		switch tv.Name {
		case "string":
			fmt.Fprint(buf, `strVal`)
		case "float64":
			fmt.Fprint(buf, `floatVal`)
		case "int64", "int":
			fmt.Fprint(buf, `intVal`)
		case "bool":
			fmt.Fprint(buf, `trueVal`)
		default:
			if tv.Obj != nil && tv.Obj.Decl != nil {
				if tspec, ok := tv.Obj.Decl.(*ast.TypeSpec); ok {
					tspecName := tspec.Name.Name
					switch subspec := tspec.Type.(type) {
					case *ast.StructType:
						tw.genStructFakeValue(buf, ins, pkgName, tspec.Name.Name, subspec, !insideArray)
					case *ast.Ident:
						fmt.Fprintf(buf, `%s.%s("FAKE")`, pkgName, tspecName)
					case *ast.MapType:
						tw.genBaseFakeValue(buf, ins, pkgName, subspec, false)
					case *ast.InterfaceType:
						tw.genBaseFakeValue(buf, ins, pkgName, subspec, false)
					default:
						fmt.Fprintf(buf, `??%T /* %s || %T */`, subspec, ftypeStr, tv)
					}
				} else {
					fmt.Fprintf(buf, `"Route1:%T"`, tv.Obj.Decl)
				}
			} else {
				ast.Print(ins.FSet, tv)
				fmt.Fprintf(buf, `"Route2:%T"`, tv)
			}
		}

	default:
		fmt.Fprintf(buf, `"Route5:%T"`, tv)
		_ = tv
	}
}

func (tw *TestWriter) genStructFieldTest(buf io.Writer, ins *astinspector.Inspector, pkgName, structName string, strct *ast.StructType, field *ast.Field) {
	fname := field.Names[0].Name

	// ftype := field.Type
	// isPointer := false
	// if innerType, ok := ftype.(*ast.StarExpr); ok {
	// 	ftype = innerType
	// }

	ftype := printNode(ins, field.Type)

	switch ftype {
	case "*string":
		fmt.Fprintf(buf, `  require.Equal(t, strVal, *%s.%s{%s: &strVal}.%s)`+"\n", pkgName, structName, fname, fname)
	case "string":
		fmt.Fprintf(buf, `  require.Equal(t, strVal, %s.%s{%s: strVal}.%s)`+"\n", pkgName, structName, fname, fname)
	case "[]string":
		fmt.Fprintf(buf, `  require.Contains(t, %s.%s{%s: []string{strVal}}.%s, strVal)`+"\n", pkgName, structName, fname, fname)

	case "*float64":
		fmt.Fprintf(buf, `  require.Equal(t, floatVal, *%s.%s{%s: &floatVal}.%s)`+"\n", pkgName, structName, fname, fname)
	case "float64":
		fmt.Fprintf(buf, `  require.Equal(t, floatVal, %s.%s{%s: floatVal}.%s)`+"\n", pkgName, structName, fname, fname)

	case "*int64":
		fmt.Fprintf(buf, `  require.Equal(t, intVal, *%s.%s{%s: &intVal}.%s)`+"\n", pkgName, structName, fname, fname)
	case "int64":
		fmt.Fprintf(buf, `  require.Equal(t, intVal, %s.%s{%s: intVal}.%s)`+"\n", pkgName, structName, fname, fname)

	case "*time.Time":
		fmt.Fprintf(buf, `  require.Equal(t, timeVal, *%s.%s{%s: &timeVal}.%s)`+"\n", pkgName, structName, fname, fname)
	case "time.Time":
		fmt.Fprintf(buf, `  require.Equal(t, timeVal, %s.%s{%s: timeVal}.%s)`+"\n", pkgName, structName, fname, fname)

		// TODO: setup case for Enums
		// simpleenum.Root.Color == ColorType

	default:
		// fmt.Fprintf(buf, `  // Lazily Tested: %s.%s.%s == %s`+"\n", pkgName, structName, fname, ftype)
		fmt.Fprintf(buf, `  require.NotNil(t, genStruct.%s) // Lazily Tested: %s.%s.%s == %s`+"\n", fname, pkgName, structName, fname, ftype)

		// ast.Print(ins.FSet, field)
	}

	// fmt.Fprintf(buf, `  require.IsType(t, *new(%s), *new(%s.%s))`+"\n", printNode(ins, v.Type), pkgName, tname)
}

func printNode(ins *astinspector.Inspector, node ast.Node) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, ins.FSet, node)
	return buf.String()
}

func getSortedKeys[T any](m map[string]T) []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}
