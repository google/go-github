// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package gogithub is a custom linter to be used by golangci-lint.

It finds:

 1. instances of `[]*string` and slices of structs without pointers and report them. (See https://github.com/google/go-github/issues/180)
 2. struct fields with `json:"omitempty"` tags that are not pointers, slices, maps, or interfaces.
*/
package gogithub

import (
	"go/ast"
	"go/token"
	"reflect"
	"slices"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("gogithub", New)
}

// GoGithubPlugin is a custom linter plugin for golangci-lint.
type GoGithubPlugin struct{}

// New returns an analysis.Analyzer to use with golangci-lint.
func New(_ any) (register.LinterPlugin, error) {
	return &GoGithubPlugin{}, nil
}

// BuildAnalyzers builds the analyzers for the GoGithubPlugin.
func (f *GoGithubPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "gogithub",
			Doc:  "Reports usage of []*string and slices of structs without pointers.",
			Run:  run,
		},
	}, nil
}

// GetLoadMode returns the load mode for the GoGithubPlugin.
func (f *GoGithubPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return false
			}

			switch t := n.(type) {
			case *ast.ArrayType:
				checkArrayType(t, t.Pos(), pass)
			case *ast.StructType:
				checkStructType(t, pass)
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

// allowedQualifiedTypes is a allowlist of qualified types (package.Type) that are allowed
// to use omitempty without pointers because they are reference types or have special semantics.
var allowedQualifiedTypes = [][2]string{
	{"json", "RawMessage"}, // json.RawMessage is []byte, can be nil
}

func checkStructType(structType *ast.StructType, pass *analysis.Pass) {
	if structType.Fields == nil {
		return
	}

	for _, field := range structType.Fields.List {
		if field.Tag == nil {
			continue
		}

		// Parse struct tag properly using reflect.StructTag
		tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
		if jsonTag := tag.Get("json"); !strings.Contains(jsonTag, "omitempty") {
			continue
		}

		switch t := field.Type.(type) {
		case *ast.Ident:
			if t.Name == "any" {
				break
			}

			pass.Reportf(field.Pos(), "using json:\"omitempty\" tag will cause zero values to be unexpectedly omitted")
		case *ast.SelectorExpr:
			if x, ok := t.X.(*ast.Ident); ok && slices.Contains(allowedQualifiedTypes, [2]string{x.Name, t.Sel.Name}) {
				break
			}

			pass.Reportf(field.Pos(), "using json:\"omitempty\" tag will cause zero values to be unexpectedly omitted")
		default:
			// *ast.StarExpr: pointers (can be nil)
			// *ast.ArrayType: slices/arrays (slices can be nil)
			// *ast.MapType: maps (can be nil)
			// *ast.InterfaceType: interfaces (can be nil)
			// All other types: safe by default
			break
		}
	}
}
