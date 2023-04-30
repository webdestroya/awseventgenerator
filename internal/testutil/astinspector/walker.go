package astinspector

import (
	"fmt"
	"go/ast"
)

type walker struct {
}

func (w *walker) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	fmt.Printf("Walking: %T\n", node)
	return w
}
