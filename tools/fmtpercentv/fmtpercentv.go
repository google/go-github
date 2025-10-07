// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fmtpercentv is a custom linter to be used by
// golangci-lint to find instances of `%d` or `%s` in
// format strings when `%v` would be more consistent.
package fmtpercentv

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("fmtpercentv", New)
}

// FmtPercentVPlugin is a custom linter plugin for golangci-lint.
type FmtPercentVPlugin struct{}

// New returns an analysis.Analyzer to use with golangci-lint.
func New(_ any) (register.LinterPlugin, error) {
	return &FmtPercentVPlugin{}, nil
}

// BuildAnalyzers builds the analyzers for the FmtPercentVPlugin.
func (f *FmtPercentVPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "fmtpercentv",
			Doc:  "Reports usage of %d or %s in format strings.",
			Run:  run,
		},
	}, nil
}

// GetLoadMode returns the load mode for the FmtPercentVPlugin.
func (f *FmtPercentVPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return false
			}

			switch t := n.(type) {
			case *ast.CallExpr:
				checkCallExpr(t, t.Pos(), pass)
			}

			return true
		})
	}
	return nil, nil
}

func checkCallExpr(expr *ast.CallExpr, tokenPos token.Pos, pass *analysis.Pass) {
	fun, ok := expr.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	funX, ok := fun.X.(*ast.Ident)
	if !ok {
		return
	}
	if funX.Name != "fmt" && funX.Name != "t" {
		return
	}
	if fun.Sel.Name != "Sprintf" && fun.Sel.Name != "Printf" && fun.Sel.Name != "Fprintf" && fun.Sel.Name != "Errorf" {
		return
	}
	fmtStrBasicLit, ok := expr.Args[0].(*ast.BasicLit)
	if !ok {
		return
	}
	fmtStr := fmtStrBasicLit.Value
	hasD := strings.Contains(fmtStr, "%d")
	hasS := strings.Contains(fmtStr, "%s")
	switch {
	case hasD && hasS:
		pass.Reportf(tokenPos, "use %%v instead of %%s and %%d")
	case hasD:
		pass.Reportf(tokenPos, "use %%v instead of %%d")
	case hasS:
		pass.Reportf(tokenPos, "use %%v instead of %%s")
	}
}
