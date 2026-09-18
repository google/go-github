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
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
)

// fixSource is a struct whose fields exercise every tag shape the fixer must handle: no tag,
// a tag without a json key, a json tag with and without an omit option, a pointer, and a
// declaration that gives several fields one type and one tag.
const fixSource = "package github\n" +
	"\n" +
	"type T struct {\n" +
	"\tA string\n" +
	"\tB string `yaml:\"b\"`\n" +
	"\tC string `json:\"c\"`\n" +
	"\tD string `json:\"-\"`\n" +
	"\tE, F string\n" +
	"\tG *string `json:\"g,omitempty\"`\n" +
	"\tI int `json:\"i\"`\n" +
	"}\n"

// goField returns the field with the given Go name, if the struct has one.
func (s *structInfo) goField(goName string) *fieldInfo {
	for _, f := range s.fields {
		if f.goName == goName {
			return f
		}
	}
	return nil
}

func TestFixTagShapes(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "github", "t.go"), fixSource)
	info, err := scanRepo(dir)
	assertNilError(t, err)
	si, ok := info.structs["T"]
	if !ok {
		t.Fatal("T was not scanned")
	}

	diag := func(goName string, action *fixAction) *diagnostic {
		return &diagnostic{
			file: "github/t.go", owner: "T", field: goName,
			info: si.goField(goName), action: action,
		}
	}
	fixed, written, notes, err := applyFixes(dir, []*diagnostic{
		diag("A", &fixAction{addOmit: "omitzero"}),                     // No struct tag at all.
		diag("B", &fixAction{addOmit: "omitempty"}),                    // A tag without a json key.
		diag("C", &fixAction{addOmit: "omitzero"}),                     // A json tag.
		diag("E", &fixAction{addOmit: "omitzero"}),                     // Shared with F.
		diag("G", &fixAction{unomit: true, unwrap: true}),              // Required, not nullable.
		diag("I", &fixAction{makePointer: true, addOmit: "omitempty"}), // Optional value type.
	})
	assertNilError(t, err)
	assertEqual(t, 5, fixed)
	assertEqual(t, 1, written)
	assertEqual(t, []string{"T.E: the declaration lists several fields, which share one type and tag"}, notes)

	got, readErr := os.ReadFile(filepath.Join(dir, "github", "t.go"))
	assertNilError(t, readErr)
	want := "package github\n" +
		"\n" +
		"type T struct {\n" +
		"\tA    string `json:\"A,omitzero\"`\n" +
		"\tB    string `json:\"B,omitempty\" yaml:\"b\"`\n" +
		"\tC    string `json:\"c,omitzero\"`\n" +
		"\tD    string `json:\"-\"`\n" +
		"\tE, F string\n" +
		"\tG    string `json:\"g\"`\n" +
		"\tI    *int   `json:\"i,omitempty\"`\n" +
		"}\n"
	assertEqual(t, want, string(got))
}

func TestApplyFixesUnrepairable(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "github", "t.go"), fixSource)
	// A finding without an action, and one whose file is gone, change nothing.
	fixed, written, notes, err := applyFixes(dir, []*diagnostic{
		{file: "github/t.go", owner: "T", field: "A"},
		{file: "github/gone.go", owner: "T", field: "B", info: &fieldInfo{}},
	})
	assertNilError(t, err)
	assertEqual(t, 0, fixed)
	assertEqual(t, 0, written)
	assertEqual(t, []string(nil), notes)
}

func TestApplyEditsRejectsOverlap(t *testing.T) {
	t.Parallel()
	// Edits are applied from the end of the file backwards, so that the offsets of the
	// ones that have not been applied yet stay valid.
	got, err := applyEdits([]byte("abcdef"), []*edit{{start: 1, end: 3, text: "X"}, insert(4, "Y")})
	assertNilError(t, err)
	assertEqual(t, "aXdYef", string(got))

	if _, err := applyEdits([]byte("abc"), []*edit{{start: 2, end: 9, text: "X"}}); err == nil {
		t.Error("applyEdits accepted an edit past the end of the file")
	}
}

func TestParseStructTag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		lit         string
		wantName    string
		wantOmit    string
		wantIgnored bool
	}{
		{"json with omitempty", "`json:\"a,omitempty\"`", "a", "omitempty", false},
		{"json with omitzero", "`json:\"a,omitzero\"`", "a", "omitzero", false},
		{"json without a name", "`json:\",omitempty\"`", "", "omitempty", false},
		{"json without options", "`json:\"a\"`", "a", "", false},
		{"ignored", "`json:\"-\"`", "", "", true},
		{"ignored with options", "`json:\"-,omitempty\"`", "", "", true},
		{"no json key", "`yaml:\"a\"`", "", "", false},
		{"json first", "`json:\"a,omitempty\" yaml:\"b\"`", "a", "omitempty", false},
		{"json second", "`yaml:\"b\" json:\"a\"`", "a", "", false},
		{"empty tag", "``", "", "", false},
		{"unterminated", "`json:\"a", "", "", false},
		{"escaped quote", "`json:\"a\\\"b,omitempty\"`", "a\\\"b", "omitempty", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			name, omit, ignored, _, _ := parseStructTag(tt.lit)
			assertEqual(t, tt.wantName, name)
			assertEqual(t, tt.wantOmit, omit)
			assertEqual(t, tt.wantIgnored, ignored)
		})
	}
}

// TestParseStructTagOffsets pins the offsets down, because the fixer inserts at valEnd.
func TestParseStructTagOffsets(t *testing.T) {
	t.Parallel()
	const lit = "`json:\"a,omitempty\" yaml:\"b\"`"
	_, _, _, valStart, valEnd := parseStructTag(lit)
	assertEqual(t, "a,omitempty", lit[valStart:valEnd])
}

func TestOmits(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		f    fieldInfo
		want bool
	}{
		{"no omit option", fieldInfo{}, false},
		{"omitempty on a scalar", fieldInfo{hasOmit: true}, true},
		{"omitempty on a struct", fieldInfo{hasOmit: true, isStruct: true}, false},
		{"omitzero on a struct", fieldInfo{hasOmit: true, omitZero: true, isStruct: true}, true},
		{"omitzero on a scalar", fieldInfo{hasOmit: true, omitZero: true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertEqual(t, tt.want, tt.f.omits())
		})
	}
}

func TestOmitOption(t *testing.T) {
	t.Parallel()
	assertEqual(t, "", (&fieldInfo{}).omitOption())
	assertEqual(t, "omitempty", (&fieldInfo{hasOmit: true}).omitOption())
	assertEqual(t, "omitzero", (&fieldInfo{hasOmit: true, omitZero: true}).omitOption())
}

func TestParseMetaOps(t *testing.T) {
	t.Parallel()
	fn := parseFunc(t, "package github\n\n// A comment.\n"+
		"//meta:operation POST /a\n"+
		"// meta:operation patch /b\n"+
		"//meta:operation bogus\n"+
		"func (s *S) M(body T) {}\n")
	got := parseMetaOps(fn.Doc)
	assertEqual(t, []*opRef{{method: "POST", path: "/a"}, {method: "PATCH", path: "/b"}}, got, cmpopts.EquateComparable(opRef{}))

	assertEqual(t, []*opRef(nil), parseMetaOps(nil), cmpopts.EquateComparable(opRef{}))
	assertEqual(t, "POST /a", got[0].String())
}

func TestFindBodyParam(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		params      string
		wantName    string
		wantPointer bool
	}{
		{"by value", "ctx context.Context, owner string, body T", "T", false},
		{"by pointer", "ctx context.Context, body *T", "T", true},
		{"renamed", "ctx context.Context, payload T", "", false},
		{"no body", "ctx context.Context", "", false},
		{"reader", "ctx context.Context, body io.Reader", "", false},
		{"slice", "ctx context.Context, body []byte", "", false},
		{"qualified type", "ctx context.Context, body pkg.T", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fn := parseFunc(t, "package github\n\n\nfunc (s *S) M("+tt.params+") {}\n")
			name, pointer := findBodyParam(fn)
			assertEqual(t, tt.wantName, name)
			assertEqual(t, tt.wantPointer, pointer)
		})
	}

	// A function without parameters has no body to check.
	fn := parseFunc(t, "package github\n\nfunc F() {}\n")
	if name, _ := findBodyParam(fn); name != "" {
		t.Errorf("findBodyParam returned %q for a function with no parameters", name)
	}
}

// TestScanRepoFixture pins the shape of the scanned fixture: the counts the summary prints,
// and the field metadata that the checker and the fixer both depend on.
func TestScanRepoFixture(t *testing.T) {
	t.Parallel()
	info, err := scanRepo(fixtureRepo)
	assertNilError(t, err)
	assertEqual(t, 3, info.files)
	assertEqual(t, 16, len(info.methods))

	names := make([]string, 0, len(info.structs))
	for name := range info.structs {
		names = append(names, name)
	}
	slices.Sort(names)
	assertEqual(t, []string{
		"AllOfRequest", "CreateCommentRequest", "CreateEnterpriseRunnerGroupRequest",
		"CreateRunnerGroupRequest", "DualRequest", "EnterpriseRunnerGroupsService",
		"InnerConfig", "IssuesService", "MissingRequest", "NullableRequest",
		"OneOfRequest", "RunnerGroupsService", "StructTypeRequest", "UntrackedRequest",
		"UpdateCommentRequest", "UpdateRunnerGroupRequest", "ValueTypeRequest",
	}, names)

	si := info.structs["CreateRunnerGroupRequest"]
	if si == nil {
		t.Fatal("CreateRunnerGroupRequest was not scanned")
	}
	assertEqual(t, "github/runner_groups.go", si.file)
	name := si.field("name")
	assertEqual(t, "Name", name.goName)
	assertEqual(t, true, name.isPointer)
	assertEqual(t, true, name.hasOmit)
	assertEqual(t, false, name.omitZero)

	// A nested struct is only recognized as one after every file has been read.
	assertEqual(t, true, info.structs["StructTypeRequest"].field("inner").isStruct)
	assertEqual(t, false, info.structs["ValueTypeRequest"].field("count").isStruct)

	// The field is looked up by JSON name, which the tag overrides.
	if f := si.field("response_only"); f == nil {
		t.Fatal("ResponseOnly was not scanned")
	} else {
		assertEqual(t, "ResponseOnly", f.goName)
	}
	if f := si.field("does_not_exist"); f != nil {
		t.Errorf("a property the struct does not have was found: %v", f.goName)
	}

	// UntrackedRequest has no annotation, but Converted takes it by pointer.
	methods := map[string]*methodInfo{}
	for _, m := range info.methods {
		methods[m.funcName] = m
	}
	assertEqual(t, 0, len(methods["Untracked"].ops))
	assertEqual(t, "UntrackedRequest", methods["Untracked"].bodyType)
	assertEqual(t, true, methods["Converted"].bodyPtr)
	assertEqual(t, 1, len(methods["CreateRunnerGroup"].ops))
	assertEqual(t, "POST /orgs/{org}/actions/runner-groups", methods["CreateRunnerGroup"].ops[0].String())
}

// TestScanRepoIgnoresJSONDash checks that a field with json:"-" is not treated as a field of
// the request body, and that a declaration sharing a type is marked as shared.
func TestScanRepoIgnoresJSONDash(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "github", "t.go"), fixSource)
	info, err := scanRepo(dir)
	assertNilError(t, err)
	si := info.structs["T"]
	if f := si.field("D"); f != nil {
		t.Error("a field with json:\"-\" must not be checked as a field of the request body")
	}
	assertEqual(t, true, si.field("E").shared)
	assertEqual(t, false, si.field("A").shared)
}

func TestScanRepoMissingGitHubDir(t *testing.T) {
	t.Parallel()
	if _, err := scanRepo(filepath.Join(t.TempDir(), "absent")); err == nil {
		t.Error("scanRepo accepted a checkout without a github directory")
	}
}

// parseFunc parses a one-method source file and returns the method declaration.
func parseFunc(t *testing.T, src string) *ast.FuncDecl {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), "src.go", src, parser.ParseComments)
	assertNilError(t, err)
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv != nil {
			return fn
		}
	}
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			return fn
		}
	}
	t.Fatal("no function declaration in source")
	return nil
}
