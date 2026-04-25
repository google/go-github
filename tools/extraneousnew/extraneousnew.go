// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package extraneousnew is a custom linter to be used by
// golangci-lint to find instances where the Go code could
// replace extraneous and problematic usage of `new` or `&SomeStruct{}`
// with an initialized pointer in very specific instances.
// It promotes the idiomatic Go concept that the zero-value is useful.
package extraneousnew

import (
	"go/ast"
	"go/token"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("extraneousnew", New)
}

// ExtraneousNewPlugin is a custom linter plugin for golangci-lint.
type ExtraneousNewPlugin struct {
	ignoredMethods map[string]bool
}

// New returns an analysis.Analyzer to use with golangci-lint.
func New(cfg any) (register.LinterPlugin, error) {
	ignoredMethods := map[string]bool{}
	if cfg != nil {
		if data, ok := cfg.(map[string]any); ok {
			if ignored, ok := data["ignored-methods"].([]any); ok {
				for _, m := range ignored {
					if s, ok := m.(string); ok {
						ignoredMethods[s] = true
					}
				}
			}
		}
	}

	return &ExtraneousNewPlugin{
		ignoredMethods: ignoredMethods,
	}, nil
}

// BuildAnalyzers builds the analyzers for the ExtraneousNewPlugin.
func (f *ExtraneousNewPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "extraneousnew",
			Doc: `Reports problematic usage of 'new' or '&SomeStruct{}' when a
more idiomatic 'var' pointer would be more appropriate.
It encourages use of the idiomatic Go concept that the zero-value is useful.`,
			Run: func(pass *analysis.Pass) (any, error) {
				return run(pass, f.ignoredMethods)
			},
		},
	}, nil
}

// GetLoadMode returns the load mode for the ExtraneousNewPlugin.
func (f *ExtraneousNewPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func run(pass *analysis.Pass, ignoredMethods map[string]bool) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			if !fn.Name.IsExported() {
				return false
			}

			// Check if method should be ignored.
			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				var recvName string
				recvType := fn.Recv.List[0].Type
				if star, ok := recvType.(*ast.StarExpr); ok {
					recvType = star.X
				}
				if ident, ok := recvType.(*ast.Ident); ok {
					recvName = ident.Name
				}

				if recvName != "" {
					fullName := recvName + "." + fn.Name.Name
					if ignoredMethods[fullName] {
						return false
					}
				}
			}

			if fn.Body != nil {
				inspectAllBlocks(pass, fn.Body)
			}
			return true
		})
	}
	return nil, nil
}

func inspectAllBlocks(pass *analysis.Pass, root ast.Node) {
	ast.Inspect(root, func(n ast.Node) bool {
		block, ok := n.(*ast.BlockStmt)
		if !ok {
			return true
		}
		inspectBlock(pass, block)
		return true
	})
}

func inspectBlock(pass *analysis.Pass, block *ast.BlockStmt) {
	// Track pointers that are currently nil.
	nilPointers := make(map[string]*ast.Ident)

	for i, stmt := range block.List {
		// 1. Check for `var v *T` or `var v *struct{...}`
		if decl, ok := stmt.(*ast.DeclStmt); ok {
			if gen, ok := decl.Decl.(*ast.GenDecl); ok && gen.Tok == token.VAR {
				for _, spec := range gen.Specs {
					if vSpec, ok := spec.(*ast.ValueSpec); ok {
						if _, ok := vSpec.Type.(*ast.StarExpr); ok && len(vSpec.Values) == 0 {
							for _, name := range vSpec.Names {
								nilPointers[name.Name] = name
							}
						}
					}
				}
			}
		}

		// 2. Check for `v = new(T)` or `v := new(T)`
		var assignLHS *ast.Ident
		var isNewT bool
		var typeName string

		if assign, ok := stmt.(*ast.AssignStmt); ok && len(assign.Lhs) == 1 && len(assign.Rhs) == 1 {
			if lhs, ok := assign.Lhs[0].(*ast.Ident); ok {
				assignLHS = lhs
				// Any assignment to v means it's no longer a "nil pointer" for our simple tracking.
				delete(nilPointers, lhs.Name)

				// Check for v := new(T) or v := &T{}
				if call, ok := assign.Rhs[0].(*ast.CallExpr); ok {
					if fun, ok := call.Fun.(*ast.Ident); ok && fun.Name == "new" && len(call.Args) == 1 {
						isNewT = true
						if typeIdent, ok := call.Args[0].(*ast.Ident); ok {
							typeName = typeIdent.Name
						}
					}
				}
				if unary, ok := assign.Rhs[0].(*ast.UnaryExpr); ok && unary.Op == token.AND {
					if composite, ok := unary.X.(*ast.CompositeLit); ok && len(composite.Elts) == 0 {
						isNewT = true
						if typeIdent, ok := composite.Type.(*ast.Ident); ok {
							typeName = typeIdent.Name
						}
					}
				}
			}
		}

		if isNewT && assignLHS != nil {
			lookAhead(pass, block, i, assignLHS, typeName)
			continue
		}

		// If it's a regular assignment (possibly with multiple variables), it might initialize a nil pointer.
		if assign, ok := stmt.(*ast.AssignStmt); ok {
			for _, lhs := range assign.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					delete(nilPointers, ident.Name)
				}
			}
		}

		// 3. Check if a nil pointer is passed to Do/Decode.
		ast.Inspect(stmt, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			fnName := getFunctionName(call.Fun)
			if fnName == "" {
				return true
			}

			var targetArg ast.Expr
			if fnName == "Do" && len(call.Args) == 2 {
				targetArg = call.Args[1]
			} else if fnName == "Decode" && len(call.Args) == 1 {
				targetArg = call.Args[0]
			}

			if targetArg != nil {
				if ident, ok := targetArg.(*ast.Ident); ok {
					if _, isNil := nilPointers[ident.Name]; isNil {
						pass.Reportf(ident.Pos(), "pass '&%v' instead", ident.Name)
					}
				}
			}
			return true
		})
	}
}

func getFunctionName(expr ast.Expr) string {
	switch f := expr.(type) {
	case *ast.SelectorExpr:
		return f.Sel.Name
	case *ast.Ident:
		return f.Name
	}
	return ""
}

func lookAhead(pass *analysis.Pass, block *ast.BlockStmt, startIndex int, lhsIdent *ast.Ident, typeName string) {
	var foundProperUse bool
	var foundOtherUse bool

	for j := startIndex + 1; j < len(block.List); j++ {
		nextStmt := block.List[j]
		ast.Inspect(nextStmt, func(un ast.Node) bool {
			if foundProperUse || foundOtherUse {
				return false
			}

			// Check if lhsIdent is used here.
			ident, ok := un.(*ast.Ident)
			if !ok || ident.Name != lhsIdent.Name {
				return true
			}

			// Found a use of lhsIdent. Is it the target argument in a call to Do or Decode?
			isSafe := false
			ast.Inspect(nextStmt, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				fnName := getFunctionName(call.Fun)
				var targetArg ast.Expr
				if fnName == "Do" && len(call.Args) == 2 {
					targetArg = call.Args[1]
				} else if fnName == "Decode" && len(call.Args) == 1 {
					targetArg = call.Args[0]
				}

				if targetArg != nil {
					if isIdentOrAddressOfIdent(targetArg, lhsIdent.Name) {
						isSafe = true
						return false
					}
				}
				return true
			})

			if isSafe {
				if typeName != "" {
					pass.Reportf(ident.Pos(), "use 'var %v *%v' and pass '&%v' instead", lhsIdent.Name, typeName, lhsIdent.Name)
				} else {
					pass.Reportf(ident.Pos(), "use 'var %v *T' and pass '&%v' instead", lhsIdent.Name, lhsIdent.Name)
				}
				foundProperUse = true
			} else {
				foundOtherUse = true
			}
			return false
		})
		if foundProperUse || foundOtherUse {
			break
		}
	}
}

func isIdentOrAddressOfIdent(expr ast.Expr, name string) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == name
	}
	if unary, ok := expr.(*ast.UnaryExpr); ok && unary.Op == token.AND {
		if ident, ok := unary.X.(*ast.Ident); ok {
			return ident.Name == name
		}
	}
	return false
}
