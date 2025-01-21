// sliceofpointers is a custom linter to be used by
// golangci-lint to find instances of `[]*string` and
// slices of structs without pointers and report them.
// See: https://github.com/google/go-github/issues/180
package main

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var analyzer = &analysis.Analyzer{
	Name: "sliceofpointers",
	Doc:  "reports usage of []*string and slices of structs without pointers",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return false
			}

			switch t := n.(type) {
			case *ast.ArrayType:
				checkArrayType(t, t.Pos(), pass)
			}

			return true
		})
	}
	return nil, nil
}

func checkArrayType(arrType *ast.ArrayType, tokenPos token.Pos, pass *analysis.Pass) {
	if starExpr, ok := arrType.Elt.(*ast.StarExpr); ok {
		if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == "string" {
			const msg = "use []string instead of []*string"
			pass.Reportf(tokenPos, msg)
		}
	} else if ident, ok := arrType.Elt.(*ast.Ident); ok && ident.Obj != nil {
		if _, ok := ident.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType); ok {
			pass.Reportf(tokenPos, "use []*%v instead of []%[1]v", ident.Name)
		}
	}
}

func main() {
	singlechecker.Main(analyzer)
}
