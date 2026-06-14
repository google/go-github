// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package requestbody is a custom linter for client.NewRequest body parameters.
//
// For client.NewRequest calls using the PATCH, POST, or PUT methods it:
//   - suggests renaming a body parameter to "body"
//   - reports body parameters passed by pointer, which should be passed by value
//   - reports body parameter types with an "Options" suffix, which should use a "Request" suffix instead.
package requestbody

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("requestbody", New)
}

// RequestBodyPlugin is a custom linter plugin for golangci-lint.
type RequestBodyPlugin struct {
	// allowedPointerTypes are body type names exempt from the by-value rule.
	allowedPointerTypes map[string]bool
	// allowedWrongNames are body type names exempt from the "Options" suffix rule.
	allowedWrongNames map[string]bool
}

// New returns an analysis.Analyzer to use with golangci-lint.
func New(cfg any) (register.LinterPlugin, error) {
	allowedPointerTypes := map[string]bool{}
	allowedWrongNames := map[string]bool{}

	if cfg != nil {
		if settingsMap, ok := cfg.(map[string]any); ok {
			if exceptionsRaw, ok := settingsMap["allowed-pointer-types"]; ok {
				if exceptionsList, ok := exceptionsRaw.([]any); ok {
					for _, item := range exceptionsList {
						if exception, ok := item.(string); ok {
							allowedPointerTypes[exception] = true
						}
					}
				}
			}

			if exceptionsRaw, ok := settingsMap["allowed-wrong-names"]; ok {
				if exceptionsList, ok := exceptionsRaw.([]any); ok {
					for _, item := range exceptionsList {
						if exception, ok := item.(string); ok {
							allowedWrongNames[exception] = true
						}
					}
				}
			}
		}
	}
	return &RequestBodyPlugin{
		allowedPointerTypes: allowedPointerTypes,
		allowedWrongNames:   allowedWrongNames,
	}, nil
}

// BuildAnalyzers builds the analyzers for the RequestBodyPlugin.
func (p *RequestBodyPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "requestbody",
			Doc:  "Reports issues with request body parameters in client.NewRequest PATCH/POST/PUT calls.",
			Run:  p.run,
		},
	}, nil
}

// GetLoadMode returns the load mode for the RequestBodyPlugin.
func (p *RequestBodyPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

// usage records which body types triggered each rule within a package, so that allow-list entries that are no longer needed can be reported.
type usage struct {
	pointerTypes map[string]bool
	optionsTypes map[string]bool
}

func (p *RequestBodyPlugin) run(pass *analysis.Pass) (any, error) {
	used := usage{
		pointerTypes: map[string]bool{},
		optionsTypes: map[string]bool{},
	}
	declared := map[string]*ast.Ident{}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if d.Tok != token.TYPE {
					continue
				}
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						declared[ts.Name.Name] = ts.Name
					}
				}
			case *ast.FuncDecl:
				if d.Body != nil {
					p.analyzeFunc(pass, d, &used)
				}
			}
		}
	}

	p.reportUnusedSettings(pass, declared, &used)
	return nil, nil
}

func (p *RequestBodyPlugin) analyzeFunc(pass *analysis.Pass, fn *ast.FuncDecl, used *usage) {
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || !isClientNewRequest(call) || !isMutatingMethod(call) {
			return true
		}

		bodyIdent, ok := call.Args[3].(*ast.Ident)
		if !ok {
			return true
		}

		field, name := findParam(fn, bodyIdent.Name)
		if field == nil {
			return true
		}

		reportRename(pass, fn, name)
		p.reportByValue(pass, field, name, used)
		p.reportTypeSuffix(pass, field, used)
		return true
	})
}

// reportUnusedSettings reports allow-list entries whose type is declared in this package but never triggered the rule they exempt,
// meaning the exception can be removed.
// Both the type declarations and their NewRequest usages live in the same package,
// so a type declared here that is never recorded as used is genuinely no longer in violation.
func (p *RequestBodyPlugin) reportUnusedSettings(pass *analysis.Pass, declared map[string]*ast.Ident, used *usage) {
	for name := range p.allowedPointerTypes {
		ident, ok := declared[name]
		if !ok || used.pointerTypes[name] {
			continue
		}
		pass.Report(analysis.Diagnostic{
			Pos:     ident.Pos(),
			End:     ident.End(),
			Message: fmt.Sprintf("unused requestbody exception: type %q in allowed-pointer-types is never passed by pointer to client.NewRequest", name),
		})
	}

	for name := range p.allowedWrongNames {
		ident, ok := declared[name]
		if !ok || used.optionsTypes[name] {
			continue
		}
		pass.Report(analysis.Diagnostic{
			Pos:     ident.Pos(),
			End:     ident.End(),
			Message: fmt.Sprintf("unused requestbody exception: type %q in allowed-wrong-names is never passed as a request body to client.NewRequest", name),
		})
	}
}

func reportRename(pass *analysis.Pass, fn *ast.FuncDecl, name *ast.Ident) {
	if name.Name == "body" {
		return
	}

	diag := analysis.Diagnostic{
		Pos:     name.Pos(),
		End:     name.End(),
		Message: fmt.Sprintf("rename request body parameter %q to \"body\"", name.Name),
	}
	if edits := renameEdits(fn, name.Name); edits != nil {
		diag.SuggestedFixes = []analysis.SuggestedFix{
			{
				Message:   `Rename to "body"`,
				TextEdits: edits,
			},
		}
	}
	pass.Report(diag)
}

func (p *RequestBodyPlugin) reportByValue(pass *analysis.Pass, field *ast.Field, name *ast.Ident, used *usage) {
	if _, ok := field.Type.(*ast.StarExpr); !ok {
		return
	}
	if ident := typeNameIdent(field.Type); ident != nil {
		used.pointerTypes[ident.Name] = true
		if p.allowedPointerTypes[ident.Name] {
			return
		}
	}
	pass.Report(analysis.Diagnostic{
		Pos:     name.Pos(),
		End:     name.End(),
		Message: fmt.Sprintf("pass request body %q by value, not by pointer", name.Name),
	})
}

// reportTypeSuffix reports body parameter types whose name ends with "Options", which should use a "Request" suffix instead
// (e.g. UserSuspendOptions -> UserSuspendRequest).
// It is report-only because renaming a type affects its declaration and every use across the codebase.
func (p *RequestBodyPlugin) reportTypeSuffix(pass *analysis.Pass, field *ast.Field, used *usage) {
	ident := typeNameIdent(field.Type)
	if ident == nil || !strings.HasSuffix(ident.Name, "Options") {
		return
	}
	used.optionsTypes[ident.Name] = true
	if p.allowedWrongNames[ident.Name] {
		return
	}
	pass.Report(analysis.Diagnostic{
		Pos:     ident.Pos(),
		End:     ident.End(),
		Message: fmt.Sprintf("request body type %q should use a \"Request\" suffix, not \"Options\"", ident.Name),
	})
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

// renameEdits builds the text edits that rename every reference to the old parameter within fn to "body".
// It returns nil (no auto-fix) when "body" is already used as an identifier in fn or when old is shadowed by a local declaration,
// since either case could make the rename incorrect.
func renameEdits(fn *ast.FuncDecl, old string) []analysis.TextEdit {
	const newName = "body"

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
