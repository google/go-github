// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package redundantptr is a custom linter to find redundant github.Ptr calls.
package redundantptr

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
)

func init() {
	register.Plugin("redundantptr", New)
}

// RedundantPtrPlugin is a custom linter plugin for golangci-lint.
type RedundantPtrPlugin struct{}

// New returns an analysis.Analyzer to use with golangci-lint.
func New(_ any) (register.LinterPlugin, error) {
	return &RedundantPtrPlugin{}, nil
}

// BuildAnalyzers builds the analyzers for the RedundantPtrPlugin.
func (p *RedundantPtrPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "redundantptr",
			Doc:  "Reports github.Ptr(x) calls that can be replaced with &x.",
			Run:  run,
		},
	}, nil
}

// GetLoadMode returns the load mode for the RedundantPtrPlugin.
func (p *RedundantPtrPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch fn := n.(type) {
			case *ast.FuncDecl:
				if fn.Body != nil {
					analyzeFunction(pass, fn, fn.Body)
				}
			case *ast.FuncLit:
				analyzeFunction(pass, fn, fn.Body)
			}
			return true
		})
	}

	return nil, nil
}

func analyzeFunction(pass *analysis.Pass, fn ast.Node, body *ast.BlockStmt) {
	if shouldIgnoreDeprecatedPtrWrapper(fn) {
		return
	}

	locals := collectLocals(fn, body)

	astutil.Apply(body, func(cursor *astutil.Cursor) bool {
		call, ok := cursor.Node().(*ast.CallExpr)
		if !ok {
			return true
		}

		rootName, argText := redundantPtrCall(pass, call)
		if rootName == "" || argText == "" {
			return true
		}

		if !locals[rootName] {
			return true
		}

		replacement := "&" + argText
		pass.Report(analysis.Diagnostic{
			Pos:     call.Pos(),
			End:     call.End(),
			Message: "replace github.Ptr(" + argText + ") with " + replacement,
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "Use address-of operator",
					TextEdits: []analysis.TextEdit{
						{Pos: call.Pos(), End: call.End(), NewText: []byte(replacement)},
					},
				},
			},
		})
		return true
	}, nil)
}

func redundantPtrCall(pass *analysis.Pass, call *ast.CallExpr) (rootName, argText string) {
	if len(call.Args) != 1 {
		return "", ""
	}

	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		if fun.Sel == nil || fun.Sel.Name != "Ptr" {
			return "", ""
		}
		qual, ok := fun.X.(*ast.Ident)
		if !ok || qual.Name != "github" {
			return "", ""
		}
	case *ast.Ident:
		if fun.Name != "Ptr" {
			return "", ""
		}
	default:
		return "", ""
	}

	arg := call.Args[0]
	root := rootIdentName(arg)
	if root == "" {
		return "", ""
	}
	return root, exprString(pass, arg)
}

func rootIdentName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return rootIdentName(e.X)
	default:
		return ""
	}
}

func exprString(pass *analysis.Pass, expr ast.Expr) string {
	var b bytes.Buffer
	if err := format.Node(&b, pass.Fset, expr); err != nil {
		return ""
	}
	return b.String()
}

func collectLocals(fn ast.Node, body *ast.BlockStmt) map[string]bool {
	locals := map[string]bool{}

	switch node := fn.(type) {
	case *ast.FuncDecl:
		addFieldListNames(locals, node.Recv)
		addFieldListNames(locals, node.Type.Params)
		addFieldListNames(locals, node.Type.Results)
	case *ast.FuncLit:
		addFieldListNames(locals, node.Type.Params)
		addFieldListNames(locals, node.Type.Results)
	}

	ast.Inspect(body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.DeclStmt:
			gen, ok := stmt.Decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.VAR {
				return true
			}
			for _, spec := range gen.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, name := range valueSpec.Names {
					locals[name.Name] = true
				}
			}
		case *ast.AssignStmt:
			if stmt.Tok == token.DEFINE {
				for _, lhs := range stmt.Lhs {
					ident, ok := lhs.(*ast.Ident)
					if !ok || ident.Name == "_" {
						continue
					}
					locals[ident.Name] = true
				}
			}
		case *ast.RangeStmt:
			if stmt.Tok == token.DEFINE {
				if ident, ok := stmt.Key.(*ast.Ident); ok && ident.Name != "_" {
					locals[ident.Name] = true
				}
				if ident, ok := stmt.Value.(*ast.Ident); ok && ident.Name != "_" {
					locals[ident.Name] = true
				}
			}
		}
		return true
	})

	return locals
}

func addFieldListNames(locals map[string]bool, fields *ast.FieldList) {
	if fields == nil {
		return
	}
	for _, field := range fields.List {
		for _, name := range field.Names {
			if name.Name != "_" {
				locals[name.Name] = true
			}
		}
	}
}

func shouldIgnoreDeprecatedPtrWrapper(fn ast.Node) bool {
	decl, ok := fn.(*ast.FuncDecl)
	if !ok {
		return false
	}

	if decl.Recv != nil || decl.Name == nil {
		return false
	}

	switch decl.Name.Name {
	case "Bool", "Int", "Int64", "String":
		return true
	default:
		return false
	}
}
