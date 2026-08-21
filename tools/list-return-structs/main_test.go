// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// parseExpr parses a Go type expression string, for testing baseStructName.
func parseExpr(t *testing.T, src string) ast.Expr {
	t.Helper()
	expr, err := parser.ParseExpr(src)
	if err != nil {
		t.Fatalf("parser.ParseExpr(%q): %v", src, err)
	}
	return expr
}

func TestBaseStructName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  string
		want string
	}{
		{"Widget", "Widget"},
		{"*Widget", "Widget"},
		{"[]*Widget", "Widget"},
		{"[...]Widget", "Widget"},
		{"map[string]*Widget", "Widget"},
		{"map[string]Widget", "Widget"},
		{"List[Widget]", "List"},
		{"Map[K, V]", "Map"},
		{"(*Widget)", "Widget"},
		// Qualified and anonymous types are not resolvable as package-local
		// struct names.
		{"time.Time", ""},
		{"json.RawMessage", ""},
		{"struct{ X int }", ""},
		{"int", "int"}, // base name only; callers intersect with the structs map
		{"[]int", "int"},
	}
	for _, tc := range tests {
		t.Run(tc.src, func(t *testing.T) {
			t.Parallel()
			got := baseStructName(parseExpr(t, tc.src))
			if got != tc.want {
				t.Errorf("baseStructName(%q) = %q, want %q", tc.src, got, tc.want)
			}
		})
	}
}

func TestBaseStructNameEllipsis(t *testing.T) {
	t.Parallel()
	// Variadic params are the only place a standalone *ast.Ellipsis occurs.
	// Parse a function type and inspect its parameter type.
	expr, err := parser.ParseExpr("func(x ...*Widget)")
	if err != nil {
		t.Fatalf("parser.ParseExpr: %v", err)
	}
	fn, ok := expr.(*ast.FuncType)
	if !ok {
		t.Fatalf("expected *ast.FuncType, got %T", expr)
	}
	field := fn.Params.List[0]
	if got := baseStructName(field.Type); got != "Widget" {
		t.Errorf("baseStructName of variadic param = %q, want %q", got, "Widget")
	}
}

func TestStructTagLookup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		tag  string
		key  string
		want string
		ok   bool
	}{
		{`json:"name,omitempty"`, "json", "name,omitempty", true},
		{`json:"name,omitempty"`, "xml", "", false},
		{`xml:"x" json:"name,omitzero"`, "json", "name,omitzero", true},
		{`json:"-"`, "json", "-", true},
		{`url:"page"`, "json", "", false},
		{``, "json", "", false},
	}
	for _, tc := range tests {
		t.Run(tc.tag, func(t *testing.T) {
			t.Parallel()
			got, ok := structTagLookup(tc.tag, tc.key)
			if got != tc.want || ok != tc.ok {
				t.Errorf("structTagLookup(%q, %q) = (%q, %v), want (%q, %v)",
					tc.tag, tc.key, got, ok, tc.want, tc.ok)
			}
		})
	}
}

func TestHasJSONOmitempty(t *testing.T) {
	t.Parallel()
	// Build a struct from source and inspect it.
	src := `package p
type T struct {
	A string ` + "`json:\"a\"`" + `
	B string ` + "`json:\"b,omitempty\"`" + `
	C int    ` + "`json:\"c,omitzero\"`" + `
	D string ` + "`url:\"d\"`" + `
	E string
}`
	file := mustParse(t, src)
	st := findStruct(t, file, "T")

	if !hasJSONOmitempty(st) {
		t.Error("hasJSONOmitempty(T) = false, want true (B has omitempty)")
	}
}

func TestHasJSONOmitemptyNone(t *testing.T) {
	t.Parallel()
	src := `package p
type T struct {
	A string ` + "`json:\"a\"`" + `
	B int    ` + "`url:\"b\"`" + `
	C string
}`
	file := mustParse(t, src)
	st := findStruct(t, file, "T")

	if hasJSONOmitempty(st) {
		t.Error("hasJSONOmitempty(T) = true, want false")
	}
}

func TestAnalyze(t *testing.T) {
	t.Parallel()
	structs, err := analyze("testdata/github", false)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}

	got := names(structs)
	// Response is a struct returned by every method and never an input, so it
	// is correctly reported; it is excluded by the -omitempty filter below.
	want := []string{"Embedded", "Inner", "Response", "Widget", "WidgetSpec"}
	if !cmp.Equal(got, want) {
		t.Errorf("analyze omitempty=false names =\n  %v\nwant\n  %v", got, want)
	}
}

func TestAnalyzeOmitempty(t *testing.T) {
	t.Parallel()
	structs, err := analyze("testdata/github", true)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}

	got := names(structs)
	// Widget (Name omitempty) and Inner (X omitempty) qualify. WidgetSpec,
	// Embedded have no omitempty/omitzero json-tagged field.
	want := []string{"Inner", "Widget"}
	if !cmp.Equal(got, want) {
		t.Errorf("analyze omitempty=true names =\n  %v\nwant\n  %v", got, want)
	}
}

func TestAnalyzeFormat(t *testing.T) {
	t.Parallel()
	structs, err := analyze("testdata/github", false)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}

	// The Widget struct declaration must be reported in the grep-style format
	// with the package directory prefix preserved.
	var found string
	for _, s := range structs {
		if s.name == "Widget" {
			found = formatLine(s)
		}
	}
	if found == "" {
		t.Fatalf("Widget not found in results: %v", names(structs))
	}
	if !strings.HasPrefix(found, "github/sample.go:") {
		t.Errorf("Widget line %q does not start with package prefix", found)
	}
	if !strings.HasSuffix(found, ":type Widget struct {") {
		t.Errorf("Widget line %q does not have expected suffix", found)
	}
}

// names returns the sorted struct names from a slice of structInfo.
func names(structs []*structInfo) []string {
	out := make([]string, 0, len(structs))
	for _, s := range structs {
		out = append(out, s.name)
	}
	slices.Sort(out)
	return out
}

// formatLine renders a structInfo the way main prints it.
func formatLine(s *structInfo) string {
	return s.relPath + ":" + strconv.Itoa(s.line) + ":type " + s.name + " struct {"
}

// mustParse parses a single source file.
func mustParse(t *testing.T, src string) *ast.File {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		t.Fatalf("parser.ParseFile: %v", err)
	}
	return f
}

// findStruct returns the StructType for the named type in the file.
func findStruct(t *testing.T, file *ast.File, name string) *ast.StructType {
	t.Helper()
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts.Name.Name != name {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				t.Fatalf("%s is not a struct type", name)
			}
			return st
		}
	}
	t.Fatalf("struct %s not found", name)
	return nil
}
