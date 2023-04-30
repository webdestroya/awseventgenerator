package astinspector

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
)

type StructFieldInfo struct {
	Struct string
	Field  string
	Type   string

	InnerType string
	Pointer   bool
	Map       bool
	Array     bool
}

type Inspector struct {
	FSet *token.FileSet
	File *ast.File

	structFields map[string]StructFieldInfo
}

func NewInspector(src string) (*Inspector, error) {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, err
	}

	ins := &Inspector{
		FSet:         fset,
		File:         f,
		structFields: make(map[string]StructFieldInfo),
	}

	// ins.Print()

	ins.walk()

	return ins, nil
}

func (i *Inspector) HasExport(name string) bool {
	return i.File.Scope.Lookup(name) != nil
}

func (i *Inspector) GetStructField(key string) *StructFieldInfo {
	if val, ok := i.structFields[key]; ok {
		return &val
	}
	return nil
}

func (i *Inspector) PackageName() string {
	return i.File.Name.Name
}

func (i *Inspector) Print() {
	ast.Print(i.FSet, i.File)
}

func (i *Inspector) DumpFile(w io.Writer) {
	printer.Fprint(w, i.FSet, i.File)
}

func (i *Inspector) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	// fmt.Printf("Walking: %T\n", node)

	switch obj := node.(type) {
	case *ast.TypeSpec:
		switch inner := obj.Type.(type) {
		case *ast.StructType:
			return i.visitStructType(obj, inner)
		default:
			fmt.Printf("  innertype: %T\n", inner)
		}
	case *ast.StarExpr:
		printer.Fprint(os.Stdout, i.FSet, obj)
	default:
		_ = obj
	}

	return i
}

func (i *Inspector) walk() {
	ast.Walk(i, i.File)
}

func (i *Inspector) printNode(node ast.Node) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, i.FSet, node)
	return buf.String()
}

func (i *Inspector) visitStructType(node *ast.TypeSpec, s *ast.StructType) ast.Visitor {
	structName := node.Name.Name
	_ = structName

	if len(s.Fields.List) == 0 {
		return nil
	}

	var buf bytes.Buffer

	for _, field := range s.Fields.List {
		buf.Reset()

		fieldName := field.Names[0].Name

		if len(field.Names) > 1 {
			fmt.Printf("FOUND MULTIPLE NAMES!!: %s %s\n", structName, fieldName)
		}

		// printer.Fprint(&buf, i.FSet, field.Type)

		finfo := StructFieldInfo{
			Struct: structName,
			Field:  fieldName,
			Type:   buf.String(),
		}

		switch ftype := field.Type.(type) {
		case *ast.ArrayType:
			finfo.Array = true
			finfo.InnerType = i.printNode(ftype.Elt)

		case *ast.MapType:

		case *ast.StarExpr:
			finfo.Pointer = true
			finfo.InnerType = i.printNode(ftype.X)
		}

		i.structFields[structName+"."+fieldName] = finfo

	}

	return nil
}
