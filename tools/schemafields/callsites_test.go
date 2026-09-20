// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// cmpTypeChange compares a type change, whose fields are unexported.
var cmpTypeChange = cmp.AllowUnexported(typeChange{})

// TestParseStructLitErr checks the compiler lines that -fix acts on, and that the lines it
// must ignore, such as the package header of a failed build, are ignored.
func TestParseStructLitErr(t *testing.T) {
	t.Parallel()
	const src = "package github\n\nvar X = T{Name: new(\"a\")}\n"
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "callers.go"), src)

	tests := []struct {
		name     string
		line     string
		wantNil  bool
		wantFile string
		wantLine int
		wantDest string
		wantExpr string
	}{
		{
			name:     "the error the compiler reports",
			line:     `callers.go:3:17: cannot use new("a") (value of type *string) as string value in struct literal`,
			wantFile: "callers.go", wantLine: 3, wantDest: "string", wantExpr: `new("a")`,
		},
		{
			name:    "the package header of a failed build is ignored",
			line:    "# github.com/google/go-github/v92/github",
			wantNil: true,
		},
		{
			name:    "the FAIL line is ignored",
			line:    "FAIL\tgithub.com/google/go-github/v92/github [build failed]",
			wantNil: true,
		},
		{
			name:    "an error about something other than a struct literal is ignored",
			line:    "callers.go:3:9: undefined: Q",
			wantNil: true,
		},
		{
			name:    "a file that does not exist is ignored",
			line:    `absent.go:3:17: cannot use new("a") (value of type *string) as string value in struct literal`,
			wantNil: true,
		},
		{
			name:    "a position past the end of the file is ignored",
			line:    `callers.go:99:17: cannot use new("a") (value of type *string) as string value in struct literal`,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parseStructLitErr(dir, tt.line)
			if tt.wantNil {
				if got != nil {
					t.Fatalf("parseStructLitErr(%q) = %+v, want nil", tt.line, got)
				}
				return
			}
			if got == nil {
				t.Fatalf("parseStructLitErr(%q) = nil", tt.line)
			}
			assertEqual(t, tt.wantFile, got.file)
			assertEqual(t, tt.wantLine, got.line)
			assertEqual(t, tt.wantDest, got.destType)
			// The offset has to land on the expression the compiler named, because that is
			// where the repair is written.
			assertEqual(t, tt.wantExpr, src[got.offset:got.offset+len(tt.wantExpr)])
		})
	}
}

func TestOffsetOf(t *testing.T) {
	t.Parallel()
	src := []byte("one\ntwo\nthree\n")
	tests := []struct {
		name      string
		line, col int
		want      int
	}{
		{"first line", 1, 1, 0},
		{"first line, last column", 1, 3, 2},
		{"the line after a newline", 2, 1, 4},
		{"second line, second column", 2, 2, 5},
		{"third line", 3, 1, 8},
		{"line 0", 0, 1, -1},
		{"column 0", 1, 0, -1},
		{"a line past the end", 4, 1, -1},
		{"a column past the end of the last line", 3, 7, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertEqual(t, tt.want, offsetOf(src, tt.line, tt.col))
		})
	}
}

// TestLocate checks that the repair finds the expression the compiler named and the literal
// that gives it its type, for the shapes a call site takes.
func TestLocate(t *testing.T) {
	t.Parallel()
	src := `package github

func Callers() []any {
	return []any{
		CreateRunnerGroupRequest{Name: new("a"), Visibility: new("all")},
		UpdateRunnerGroupRequest{Name: name},
		[]*UpdateRunnerGroupRequest{{Name: new("c")}},
		pkg.T{Name: new("d")},
	}
}
`
	pf := parseSource(t, "callers.go", src)
	tests := []struct {
		name        string
		expr        string
		wantField   string
		wantLitType string
	}{
		{"a field named by the literal", `new("a")`, "Name", "CreateRunnerGroupRequest"},
		{"the other field of the same literal", `new("all")`, "Visibility", "CreateRunnerGroupRequest"},
		{"a value rather than a call", "name", "Name", "UpdateRunnerGroupRequest"},
		{"a literal that takes its type from the slice", `new("c")`, "Name", "UpdateRunnerGroupRequest"},
		{"a type from another package", `new("d")`, "Name", "T"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			offset := strings.Index(src, tt.expr)
			if offset < 0 {
				t.Fatalf("%q is not in the source", tt.expr)
			}
			lit, err := pf.locate(offset)
			assertNilError(t, err)
			assertEqual(t, tt.wantField, lit.field)
			assertEqual(t, tt.wantLitType, lit.litType)
			assertEqual(t, tt.expr, pf.text(lit.expr))
		})
	}

	// An offset that holds no expression is an error rather than a wrong repair.
	if _, err := pf.locate(strings.Index(src, "Callers") + 1); err == nil {
		t.Error("locate accepted an offset that holds no expression")
	}
}

func TestBaseTypeName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  string
		want string
	}{
		{"T", "T"},
		{"*T", "T"},
		{"[]T", "T"},
		{"[]*T", "T"},
		{"map[string]T", "T"},
		{"(*T)", "T"},
		{"pkg.T", "T"},
		{"[]map[string]*pkg.T", "T"},
		{"[3]T", "T"},
		{"func()", ""},
	}
	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			t.Parallel()
			pf := parseSource(t, "t.go", "package github\n\nvar x []"+tt.src+"\n")
			// Reach the composite literal type through the declaration, so that the test
			// exercises the ast.Expr the scanner produces.
			assertEqual(t, tt.want, baseTypeName(arrayElt(t, pf)))
		})
	}
}

// TestRepairedExpr checks the two mechanical repairs and that everything else is refused,
// because a wrong guess here rewrites a caller's code.
func TestRepairedExpr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		expr      string
		toPointer bool
		want      string
		wantErr   bool
	}{
		{"wrap a string", `"a"`, true, `new("a")`, false},
		{"wrap a number", "3", true, "new(3)", false},
		{"unwrap new", `new("a")`, false, `"a"`, false},
		{"unwrap Ptr", `Ptr("a")`, false, `"a"`, false},
		{"unwrap the address of a variable", "&name", false, "name", false},
		{"a variable is left alone", "name", false, "", true},
		{"a call is left alone", "fn()", false, "", true},
		{"a literal is left alone", "T{}", false, "", true},
		{"a value in a slice is left alone", "[]string{}", false, "", true},
		{"new with two arguments is left alone", "new(a, b)", false, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fset := token.NewFileSet()
			expr, err := parser.ParseExprFrom(fset, "expr.go", tt.expr, 0)
			assertNilError(t, err)
			pf := &parsedFile{src: []byte(tt.expr), fset: fset}
			got, err := repairedExpr(pf, expr, tt.toPointer)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("repairedExpr(%q) = %q, want an error", tt.expr, got)
				}
				assertContains(t, err.Error(), "not mechanical")
				return
			}
			assertNilError(t, err)
			assertEqual(t, tt.want, got)
		})
	}
}

// TestTypeChanges checks the type changes that -fix reads from the sources before it repairs
// them, which is what the call site repair matches the compiler errors against.
func TestTypeChanges(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "github", "runner_groups.go"), `package github

type CreateRunnerGroupRequest struct {
	Name *string `+"`json:\"name,omitempty\"`"+`
}

type ValueTypeRequest struct {
	Count int `+"`json:\"count\"`"+`
}
`)
	info := func(goName string, typOff, typEnd, starOff int) *fieldInfo {
		return &fieldInfo{
			file: "github/runner_groups.go", goName: goName,
			typOff: typOff, typEnd: typEnd, starOff: starOff,
		}
	}
	src, err := os.ReadFile(filepath.Join(dir, "github", "runner_groups.go"))
	assertNilError(t, err)
	pointerOff := strings.Index(string(src), "*string")
	countOff := strings.Index(string(src), "int `json:\"count\"`")

	tests := []struct {
		name string
		diag *diagnostic
		want []*typeChange
	}{
		{
			name: "a type change to a value",
			diag: &diagnostic{
				owner: "CreateRunnerGroupRequest", field: "Name", action: &fixAction{unomit: true, unwrap: true},
				info: info("Name", pointerOff, pointerOff+len("*string"), pointerOff),
			},
			want: []*typeChange{{
				file: "github/runner_groups.go", structName: "CreateRunnerGroupRequest",
				fieldName: "Name", newType: "string",
			}},
		},
		{
			name: "a type change to a pointer",
			diag: &diagnostic{
				owner: "ValueTypeRequest", field: "Count", action: &fixAction{makePointer: true, addOmit: "omitempty"},
				info: info("Count", countOff, countOff+3, -1),
			},
			want: []*typeChange{{
				file: "github/runner_groups.go", structName: "ValueTypeRequest",
				fieldName: "Count", newType: "*int", toPointer: true,
			}},
		},
		{
			name: "a tag-only repair changes no type",
			diag: &diagnostic{
				owner: "ValueTypeRequest", field: "Count", action: &fixAction{addOmit: "omitzero"},
				info: info("Count", countOff, countOff+3, -1),
			},
		},
		{
			name: "a finding with no repair",
			diag: &diagnostic{owner: "ValueTypeRequest", field: "Count", info: info("Count", countOff, countOff+3, -1)},
		},
		{
			name: "a finding with no field",
			diag: &diagnostic{owner: "ValueTypeRequest", field: "Count", action: &fixAction{unwrap: true}},
		},
		{
			name: "a field whose offsets do not fit the file",
			diag: &diagnostic{
				owner: "ValueTypeRequest", field: "Count", action: &fixAction{unwrap: true},
				info: info("Count", len(src)+1, len(src)+2, len(src)+1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := typeChanges(dir, []*diagnostic{tt.diag})
			assertEqual(t, tt.want, got, cmpTypeChange)
		})
	}

	// The same change from two findings is one change, because the call sites are repaired
	// once either way.
	unwrap := &fixAction{unwrap: true}
	one := &diagnostic{
		owner: "CreateRunnerGroupRequest", field: "Name", action: unwrap,
		info: info("Name", pointerOff, pointerOff+len("*string"), pointerOff),
	}
	two := &diagnostic{
		owner: "CreateRunnerGroupRequest", field: "Name", action: &fixAction{unomit: true, unwrap: true},
		info: info("Name", pointerOff, pointerOff+len("*string"), pointerOff),
	}
	assertEqual(t, 1, len(typeChanges(dir, []*diagnostic{one, two})))
}

// TestMatchChange checks that a compiler error is repaired only when the literal sets a field
// whose type the repair changed to the type the compiler reports.
func TestMatchChange(t *testing.T) {
	t.Parallel()
	name := &typeChange{file: "github/x.go", structName: "Req", fieldName: "Name", newType: "string"}
	other := &typeChange{file: "github/x.go", structName: "Other", fieldName: "Name", newType: "string"}
	count := &typeChange{file: "github/x.go", structName: "Req", fieldName: "Count", newType: "*int", toPointer: true}

	tests := []struct {
		name     string
		lit      *litAt
		destType string
		changes  []*typeChange
		want     *typeChange
	}{
		{"the literal names the struct", &litAt{field: "Name", litType: "Req"}, "string", []*typeChange{name, other}, name},
		{"the literal takes its type from the slice", &litAt{field: "Name"}, "string", []*typeChange{name}, name},
		{"the field name differs", &litAt{field: "Title", litType: "Req"}, "string", []*typeChange{name}, nil},
		{"the type differs", &litAt{field: "Name", litType: "Req"}, "*string", []*typeChange{name}, nil},
		{"the literal is another type with a field of the same name", &litAt{field: "Name", litType: "Elsewhere"}, "string", []*typeChange{name}, name},
		{"no change matches", &litAt{field: "Name", litType: "Req"}, "string", nil, nil},
		{"a change to a pointer", &litAt{field: "Count", litType: "Req"}, "*int", []*typeChange{count}, count},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := matchChange(tt.lit, tt.destType, tt.changes)
			if tt.want == nil {
				if got != nil {
					t.Fatalf("matchChange() = %+v, want nil", got)
				}
				return
			}
			if got != tt.want {
				t.Errorf("matchChange() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

// TestModulesToCheck checks the modules that -fix builds: the one that holds the changed field
// and every module in the checkout that requires it, which is where a caller can live.
func TestModulesToCheck(t *testing.T) {
	t.Parallel()
	repo := t.TempDir()
	mustWriteFile(t, filepath.Join(repo, "go.mod"), "module example.com/m\n\ngo 1.26.0\n")
	mustWriteFile(t, filepath.Join(repo, "github", "x.go"), "package github\n")
	// A nested module, so that the walk up to the nearest go.mod has something to find.
	mustWriteFile(t, filepath.Join(repo, "tools", "schemafields", "go.mod"), "module tools/schemafields\n")
	mustWriteFile(t, filepath.Join(repo, "tools", "schemafields", "y.go"), "package main\n")
	// A module that requires the root one, so that its call sites are checked too.
	mustWriteFile(t, filepath.Join(repo, "example", "go.mod"), "module example.com/example\n\nrequire example.com/m v0.0.0\n")
	mustWriteFile(t, filepath.Join(repo, "example", "main.go"), "package main\n")
	// Modules that no change can reach.
	mustWriteFile(t, filepath.Join(repo, "unrelated", "go.mod"), "module example.com/unrelated\n")
	// A module whose own path starts with the path that changed is not a module that requires
	// it: only a require or replace line is.
	mustWriteFile(t, filepath.Join(repo, "nested", "go.mod"), "module example.com/m/nested\n")
	mustWriteFile(t, filepath.Join(repo, "testdata", "repo", "go.mod"), "module example.com/m\n")
	mustWriteFile(t, filepath.Join(repo, ".hidden", "go.mod"), "module example.com/m\n")
	mustWriteFile(t, filepath.Join(repo, "vendor", "go.mod"), "module example.com/m\n")

	assertEqual(t, repo, moduleOf(repo, filepath.Join(repo, "github", "x.go")))
	// A file in a nested module belongs to that one, not to the root.
	assertEqual(t, filepath.Join(repo, "tools", "schemafields"),
		moduleOf(repo, filepath.Join(repo, "tools", "schemafields", "y.go")))
	// A file outside the checkout is in no module, so there is nothing to build for it.
	assertEqual(t, "", moduleOf(repo, filepath.Join(filepath.Dir(repo), "absent", "x.go")))
	assertEqual(t, "example.com/m", modulePath(repo))
	assertEqual(t, true, requiresModule(filepath.Join(repo, "example"), "example.com/m"))
	assertEqual(t, false, requiresModule(filepath.Join(repo, "unrelated"), "example.com/m"))
	assertEqual(t, false, requiresModule(filepath.Join(repo, "nested"), "example.com/m"))
	assertEqual(t, false, requiresModule(filepath.Join(repo, "example"), ""))

	// The module that holds the changed file, and the one that requires it, which is where a
	// caller can live. The unrelated module, the nested one and the skipped directories are
	// left out.
	got := modulesToCheck(repo, []string{"github/x.go"})
	assertEqual(t, []string{repo, filepath.Join(repo, "example")}, got)

	// A change in a module other than the root leaves that module's callers alone, unless
	// something else requires the module the field is declared in.
	got = modulesToCheck(repo, []string{"tools/schemafields/y.go"})
	assertEqual(t, []string{filepath.Join(repo, "tools", "schemafields")}, got)

	dirs := moduleDirsUnder(repo)
	assertEqual(t, []string{
		repo, filepath.Join(repo, "example"), filepath.Join(repo, "nested"),
		filepath.Join(repo, "tools", "schemafields"), filepath.Join(repo, "unrelated"),
	}, dirs)
}

// parseSource parses a source into the form the repairs read, for the tests that exercise the
// expression lookup.
func parseSource(t *testing.T, name, src string) *parsedFile {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, name, src, parser.SkipObjectResolution)
	assertNilError(t, err)
	return &parsedFile{src: []byte(src), fset: fset, file: file}
}

// arrayElt returns the element type of the slice type in the parsed source, which is the
// ast.Expr that a composite literal of that type is declared with.
func arrayElt(t *testing.T, pf *parsedFile) ast.Expr {
	t.Helper()
	for _, decl := range pf.file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			value, ok := spec.(*ast.ValueSpec)
			if !ok || value.Type == nil {
				continue
			}
			array, ok := value.Type.(*ast.ArrayType)
			if !ok {
				continue
			}
			return array.Elt
		}
	}
	t.Fatal("the source holds no slice type")
	return nil
}
