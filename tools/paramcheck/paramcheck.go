// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package paramcheck is a custom linter for parameter naming and type conventions.
//
// For body parameters it:
//   - suggests renaming a body parameter to "body"
//   - reports body parameters passed by pointer, which should be passed by value
//   - reports body parameter types with an "Options" suffix, which should use a "Request" suffix instead.
//
// For query parameters it:
//   - suggests renaming the options parameter to "opts"
//   - reports options parameters passed by value, which should be passed by pointer
//   - reports options parameter types without an "Options" suffix.
package paramcheck

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("paramcheck", New)
}

// ParamCheckPlugin is a custom linter plugin for golangci-lint.
type ParamCheckPlugin struct {
	// bodyAllowedPointerTypes are body type names exempt from the by-value rule.
	bodyAllowedPointerTypes map[string]bool
	// bodyAllowedWrongNames are body type names exempt from the "Options" suffix rule.
	bodyAllowedWrongNames map[string]bool
}

// New returns an analysis.Analyzer to use with golangci-lint.
func New(cfg any) (register.LinterPlugin, error) {
	bodyAllowedPointerTypes := map[string]bool{}
	bodyAllowedWrongNames := map[string]bool{}

	if cfg != nil {
		if settingsMap, ok := cfg.(map[string]any); ok {
			if exceptionsRaw, ok := settingsMap["body-allowed-pointer-types"]; ok {
				if exceptionsList, ok := exceptionsRaw.([]any); ok {
					for _, item := range exceptionsList {
						if exception, ok := item.(string); ok {
							bodyAllowedPointerTypes[exception] = true
						}
					}
				}
			}

			if exceptionsRaw, ok := settingsMap["body-allowed-wrong-names"]; ok {
				if exceptionsList, ok := exceptionsRaw.([]any); ok {
					for _, item := range exceptionsList {
						if exception, ok := item.(string); ok {
							bodyAllowedWrongNames[exception] = true
						}
					}
				}
			}
		}
	}
	return &ParamCheckPlugin{
		bodyAllowedPointerTypes: bodyAllowedPointerTypes,
		bodyAllowedWrongNames:   bodyAllowedWrongNames,
	}, nil
}

// BuildAnalyzers builds the analyzers for the ParamCheckPlugin.
func (p *ParamCheckPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "paramcheck",
			Doc:  "Reports parameter naming and type convention issues in endpoint method calls.",
			Run:  p.run,
		},
	}, nil
}

// GetLoadMode returns the load mode for the ParamCheckPlugin.
func (p *ParamCheckPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func (p *ParamCheckPlugin) run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			p.analyzeFunc(pass, fn)
		}
	}
	return nil, nil
}

func (p *ParamCheckPlugin) analyzeFunc(pass *analysis.Pass, fn *ast.FuncDecl) {
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isClientNewRequest(call) && isMutatingMethod(call) {
			bodyIdent, ok := call.Args[3].(*ast.Ident)
			if !ok {
				return true
			}
			field, name := findParam(fn, bodyIdent.Name)
			if field == nil {
				return true
			}
			bodyReportName(pass, fn, name)
			p.bodyReportPass(pass, field, name)
			p.bodyReportSuffix(pass, field)
		}

		if isAddOptions(call) {
			optsIdent, ok := call.Args[1].(*ast.Ident)
			if !ok {
				return true
			}
			field, name := findParam(fn, optsIdent.Name)
			if field == nil {
				return true
			}
			queryReportName(pass, fn, name)
			p.queryReportPass(pass, field, name)
			p.queryReportSuffix(pass, field)
		}

		return true
	})
}

func bodyReportName(pass *analysis.Pass, fn *ast.FuncDecl, name *ast.Ident) {
	if name.Name == "body" {
		return
	}

	diag := analysis.Diagnostic{
		Pos:     name.Pos(),
		End:     name.End(),
		Message: fmt.Sprintf("rename request body parameter %q to \"body\"", name.Name),
	}
	if edits := renameEdits(fn, name.Name, "body"); edits != nil {
		diag.SuggestedFixes = []analysis.SuggestedFix{
			{
				Message:   `Rename to "body"`,
				TextEdits: edits,
			},
		}
	}
	pass.Report(diag)
}

func queryReportName(pass *analysis.Pass, fn *ast.FuncDecl, name *ast.Ident) {
	if name.Name == "opts" {
		return
	}

	diag := analysis.Diagnostic{
		Pos:     name.Pos(),
		End:     name.End(),
		Message: fmt.Sprintf("rename addOptions parameter %q to \"opts\"", name.Name),
	}
	if edits := renameEdits(fn, name.Name, "opts"); edits != nil {
		diag.SuggestedFixes = []analysis.SuggestedFix{
			{
				Message:   `Rename to "opts"`,
				TextEdits: edits,
			},
		}
	}
	pass.Report(diag)
}

func (p *ParamCheckPlugin) bodyReportPass(pass *analysis.Pass, field *ast.Field, name *ast.Ident) {
	if _, ok := field.Type.(*ast.StarExpr); !ok {
		return
	}
	if ident := typeNameIdent(field.Type); ident != nil && p.bodyAllowedPointerTypes[ident.Name] {
		return
	}
	pass.Report(analysis.Diagnostic{
		Pos:     name.Pos(),
		End:     name.End(),
		Message: fmt.Sprintf("pass request body %q by value, not by pointer", name.Name),
	})
}

func (p *ParamCheckPlugin) queryReportPass(pass *analysis.Pass, field *ast.Field, name *ast.Ident) {
	if _, ok := field.Type.(*ast.StarExpr); ok {
		return
	}
	pass.Report(analysis.Diagnostic{
		Pos:     name.Pos(),
		End:     name.End(),
		Message: fmt.Sprintf("pass query parameter %q by pointer, not by value", name.Name),
	})
}

func (p *ParamCheckPlugin) queryReportSuffix(pass *analysis.Pass, field *ast.Field) {
	ident := typeNameIdent(field.Type)
	if ident == nil || strings.HasSuffix(ident.Name, "Options") {
		return
	}
	pass.Report(analysis.Diagnostic{
		Pos:     ident.Pos(),
		End:     ident.End(),
		Message: fmt.Sprintf("query parameter type %q should use an \"Options\" suffix", ident.Name),
	})
}

// bodyReportSuffix reports body parameter types whose name ends with "Options", which should use a "Request" suffix instead
// (e.g. UserSuspendOptions -> UserSuspendRequest).
// It is report-only because renaming a type affects its declaration and every use across the codebase.
func (p *ParamCheckPlugin) bodyReportSuffix(pass *analysis.Pass, field *ast.Field) {
	ident := typeNameIdent(field.Type)
	if ident == nil || !strings.HasSuffix(ident.Name, "Options") {
		return
	}
	if p.bodyAllowedWrongNames[ident.Name] {
		return
	}
	pass.Report(analysis.Diagnostic{
		Pos:     ident.Pos(),
		End:     ident.End(),
		Message: fmt.Sprintf("request body type %q should use a \"Request\" suffix, not \"Options\"", ident.Name),
	})
}

// isAddOptions reports whether call is of the form addOptions(...) with at least two arguments.
// addOptions is used in endpoint methods to add query parameters, so the second argument is expected to be the options struct.
func isAddOptions(call *ast.CallExpr) bool {
	if len(call.Args) < 2 {
		return false
	}
	ident, ok := call.Fun.(*ast.Ident)
	return ok && ident.Name == "addOptions"
}

// isClientNewRequest reports whether call is of the form x.client.NewRequest(...) or client.NewRequest(...).
func isClientNewRequest(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "NewRequest" {
		return false
	}
	switch x := sel.X.(type) {
	case *ast.SelectorExpr:
		return x.Sel.Name == "client"
	case *ast.Ident:
		return x.Name == "client"
	default:
		return false
	}
}

// isMutatingMethod reports whether the call's method argument is "PATCH", "POST", or "PUT" and a body argument is present.
func isMutatingMethod(call *ast.CallExpr) bool {
	if len(call.Args) < 4 {
		return false
	}
	lit, ok := call.Args[1].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	switch lit.Value {
	case `"PATCH"`, `"POST"`, `"PUT"`:
		return true
	default:
		return false
	}
}

func findParam(fn *ast.FuncDecl, name string) (*ast.Field, *ast.Ident) {
	if fn.Type.Params == nil {
		return nil, nil
	}
	for _, field := range fn.Type.Params.List {
		for _, ident := range field.Names {
			if ident.Name == name {
				return field, ident
			}
		}
	}
	return nil, nil
}

// typeNameIdent returns the base type name identifier of expr, unwrapping a pointer and resolving a qualified (pkg.Type) selector.
func typeNameIdent(expr ast.Expr) *ast.Ident {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return typeNameIdent(t.X)
	case *ast.Ident:
		return t
	case *ast.SelectorExpr:
		return t.Sel
	default:
		return nil
	}
}

// renameEdits builds the text edits that rename every reference to the old parameter within fn to newName.
// It returns nil (no auto-fix) when newName is already used as an identifier in fn or when old is shadowed by a local declaration,
// since either case could make the rename incorrect.
func renameEdits(fn *ast.FuncDecl, old, newName string) []analysis.TextEdit {
	// Idents that are the selector field (x.old) must not be renamed.
	skip := map[*ast.Ident]bool{}
	ast.Inspect(fn, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			skip[sel.Sel] = true
		}
		return true
	})

	var edits []analysis.TextEdit
	conflict := false
	ast.Inspect(fn, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if !ok || skip[ident] {
			return true
		}
		switch ident.Name {
		case newName:
			conflict = true
		case old:
			edits = append(edits, analysis.TextEdit{
				Pos:     ident.Pos(),
				End:     ident.End(),
				NewText: []byte(newName),
			})
		}
		return true
	})
	if conflict {
		return nil
	}
	return edits
}
