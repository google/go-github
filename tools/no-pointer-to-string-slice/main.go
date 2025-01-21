// no-pointer-to-string-slice is a custom linter to be used by
// golangci-lint to find instances of `[]*string` and report
// that they should be changed to `[]string`.
package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var Analyzer = &analysis.Analyzer{
	Name: "noPointerToStringSlice",
	Doc:  "reports usage of []*string and suggests using []string instead",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Check for []*string in array types
			if arrType, ok := n.(*ast.ArrayType); ok {
				if starExpr, ok := arrType.Elt.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == "string" {
						msg := "use []string instead of []*string"
						pass.Reportf(arrType.Pos(), msg)
						println("Reported:", msg) // Debugging output
					}
				}
			}

			// Check for []*string in struct fields
			if structType, ok := n.(*ast.StructType); ok {
				for _, field := range structType.Fields.List {
					if arrType, ok := field.Type.(*ast.ArrayType); ok {
						if starExpr, ok := arrType.Elt.(*ast.StarExpr); ok {
							if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == "string" {
								msg := "use []string instead of []*string in struct fields"
								pass.Reportf(field.Pos(), msg)
								println("Reported:", msg) // Debugging output
							}
						}
					}
				}
			}

			return true
		})
	}
	return nil, nil
}

func main() {
	singlechecker.Main(Analyzer)
}
