// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/parser"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-cmp/cmp"
)

func TestCheckSchemaFieldsMatchesBySchemaName(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

import "encoding/json"

type Demo struct {
	ID       *int64             `+"`json:\"id,omitempty\"`"+`
	Name     string             `+"`json:\"name\"`"+`
	Note     string             `+"`json:\"note\"`"+`
	Items    []string           `+"`json:\"items\"`"+`
	Metadata map[string]string  `+"`json:\"metadata,omitempty\"`"+`
	Raw      json.RawMessage    `+"`json:\"raw,omitempty\"`"+`
	Extra    *string            `+"`json:\"extra,omitempty\"`"+`
	Internal *string            `+"`json:\"-\"`"+`
}
`)

	result, err := checkSchemaFields(schemaFieldCheckOptions{
		descriptions: []*openapiFile{testOpenAPIFile("descriptions/api.github.com/api.github.com.json", openapi3.Schemas{
			"demo": openapi3.NewSchemaRef("", &openapi3.Schema{
				Required: []string{"id", "name", "items", "metadata"},
				Properties: openapi3.Schemas{
					"id":       openapi3.NewSchemaRef("", openapi3.NewIntegerSchema()),
					"name":     openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
					"note":     openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
					"items":    openapi3.NewSchemaRef("", openapi3.NewArraySchema()),
					"metadata": openapi3.NewSchemaRef("", openapi3.NewObjectSchema()),
					"raw":      openapi3.NewSchemaRef("", openapi3.NewObjectSchema()),
					"missing":  openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
				},
			}),
		})},
		githubDir: githubDir,
		schemaNames: sliceSet([]string{
			"demo",
		}),
	})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff([]*schemaFieldChecked{{
		OpenAPISchema: "demo",
		GoStruct:      "Demo",
		OpenAPIFile:   "descriptions/api.github.com/api.github.com.json",
		MatchReason:   "schema name",
	}}, result.Checked); diff != "" {
		t.Errorf("checked mismatch (-want +got):\n%v", diff)
	}

	var got []string
	for _, diag := range result.Diagnostics {
		got = append(got, diag.JSONName+": "+diag.Message)
	}
	want := []string{
		"extra: field is not present in the OpenAPI schema properties",
		"id: field is required and non-nullable in the OpenAPI schema but is a pointer",
		"metadata: field is required by the OpenAPI schema but has an omit option",
		"missing: OpenAPI schema property is missing from the Go struct",
		"note: field is optional in the OpenAPI schema but is not a pointer, slice, map, interface, or selector type",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("diagnostics mismatch (-want +got):\n%v", diff)
	}
}

func TestCheckSchemaFieldsMatchesByExactFieldSet(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

type ExactFieldsRequest struct {
	ID   *int64  `+"`json:\"id,omitempty\"`"+`
	Name string  `+"`json:\"name\"`"+`
	Note *string `+"`json:\"note,omitempty\"`"+`
}

func (s *svc) Create(ctx context.Context, body *ExactFieldsRequest) {
	s.client.NewRequest(ctx, "POST", "u", body)
}
`)

	result, err := checkSchemaFields(schemaFieldCheckOptions{
		descriptions: []*openapiFile{testOpenAPIFile("descriptions/api.github.com/api.github.com.json", openapi3.Schemas{
			"unrelated-schema-name": openapi3.NewSchemaRef("", &openapi3.Schema{
				Required: []string{"id", "name"},
				Properties: openapi3.Schemas{
					"id":   openapi3.NewSchemaRef("", openapi3.NewIntegerSchema()),
					"name": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
					"note": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
				},
			}),
		})},
		githubDir: githubDir,
	})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff([]*schemaFieldChecked{{
		OpenAPISchema: "unrelated-schema-name",
		GoStruct:      "ExactFieldsRequest",
		OpenAPIFile:   "descriptions/api.github.com/api.github.com.json",
		MatchReason:   "exact JSON field set",
	}}, result.Checked); diff != "" {
		t.Errorf("checked mismatch (-want +got):\n%v\nskipped: %#v", diff, result.Skipped)
	}

	var got []string
	for _, diag := range result.Diagnostics {
		got = append(got, diag.JSONName+": "+diag.Message)
	}
	want := []string{
		"id: field is required and non-nullable in the OpenAPI schema but is a pointer",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("diagnostics mismatch (-want +got):\n%v", diff)
	}
}

func TestCheckSchemaFieldsSkipsAmbiguousExactFieldSet(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

type FirstMatchRequest struct {
	ID   int64  `+"`json:\"id\"`"+`
	Name string `+"`json:\"name\"`"+`
}

type SecondMatchRequest struct {
	ID   int64  `+"`json:\"id\"`"+`
	Name string `+"`json:\"name\"`"+`
}

func (s *svc) First(ctx context.Context, body *FirstMatchRequest) {
	s.client.NewRequest(ctx, "POST", "u", body)
}

func (s *svc) Second(ctx context.Context, body *SecondMatchRequest) {
	s.client.NewRequest(ctx, "POST", "u", body)
}
`)

	result, err := checkSchemaFields(schemaFieldCheckOptions{
		descriptions: []*openapiFile{testOpenAPIFile("descriptions/api.github.com/api.github.com.json", openapi3.Schemas{
			"unrelated-schema-name": openapi3.NewSchemaRef("", &openapi3.Schema{
				Required: []string{"id", "name"},
				Properties: openapi3.Schemas{
					"id":   openapi3.NewSchemaRef("", openapi3.NewIntegerSchema()),
					"name": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
				},
			}),
		})},
		githubDir: githubDir,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Checked) != 0 {
		t.Fatalf("checked = %v, want none", result.Checked)
	}
	if len(result.Diagnostics) != 0 {
		t.Fatalf("diagnostics = %v, want none", result.Diagnostics)
	}
	if len(result.Skipped) != 1 {
		t.Fatalf("skipped = %v, want one skip", result.Skipped)
	}
	if got := result.Skipped[0].Reason; !strings.Contains(got, "ambiguous Go struct field-set match") {
		t.Errorf("skip reason = %q, want ambiguous field-set match", got)
	}
}

func TestCheckSchemaFieldsAllowsRequiredNullablePointer(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "nullable.go"), `package github

type NullableDemo struct {
	ID *int64 `+"`json:\"id\"`"+`
}
`)

	nullableInteger := openapi3.NewIntegerSchema()
	nullableInteger.Nullable = true
	result, err := checkSchemaFields(schemaFieldCheckOptions{
		descriptions: []*openapiFile{testOpenAPIFile("descriptions/api.github.com/api.github.com.json", openapi3.Schemas{
			"nullable-demo": openapi3.NewSchemaRef("", &openapi3.Schema{
				Required: []string{"id"},
				Properties: openapi3.Schemas{
					"id": openapi3.NewSchemaRef("", nullableInteger),
				},
			}),
		})},
		githubDir:   githubDir,
		schemaNames: sliceSet([]string{"nullable-demo"}),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Diagnostics) != 0 {
		t.Errorf("diagnostics = %v, want none", result.Diagnostics)
	}
}

func TestGoNameCandidates(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{name: "plural and version tokens", in: "projects-v2", want: []string{"ProjectsV2", "ProjectV2"}},
		{name: "singular unchanged", in: "repository", want: []string{"Repository"}},
		{name: "trailing plural", in: "teams", want: []string{"Teams", "Team"}},
		{name: "multiple words no plural", in: "code-scanning-alert", want: []string{"CodeScanningAlert"}},
		{name: "initialism token", in: "api", want: []string{"API"}},
		{name: "empty", in: "", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if diff := cmp.Diff(tt.want, goNameCandidates(tt.in)); diff != "" {
				t.Errorf("goNameCandidates(%q) mismatch (-want +got):\n%v", tt.in, diff)
			}
		})
	}
}

func TestSingularize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in, want string
	}{
		{"projects", "project"},
		{"policies", "policy"},
		{"statuses", "status"},
		{"boxes", "box"},
		{"branches", "branch"},
		{"buses", "bus"},
		{"keys", "key"},
		{"class", "class"}, // "ss" is not treated as a plural "s"
		{"v2", "v2"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := singularize(tt.in); got != tt.want {
			t.Errorf("singularize(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestGoName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{name: "initialism", in: []string{"api"}, want: "API"},
		{name: "url initialism", in: []string{"url"}, want: "URL"},
		{name: "special case oauth", in: []string{"oauth"}, want: "OAuth"},
		{name: "version token", in: []string{"projects", "v2"}, want: "ProjectsV2"},
		{name: "plain word", in: []string{"repository"}, want: "Repository"},
		{name: "skips empty tokens", in: []string{"code", "", "scanning"}, want: "CodeScanning"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := goName(tt.in); got != tt.want {
				t.Errorf("goName(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestIsVersionToken(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   string
		want bool
	}{
		{"v2", true},
		{"v10", true},
		{"v", false},
		{"vx", false},
		{"2", false},
		{"version", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := isVersionToken(tt.in); got != tt.want {
			t.Errorf("isVersionToken(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestSplitOpenAPIName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   string
		want []string
	}{
		{"projects-v2", []string{"projects", "v2"}},
		{"api.github.com", []string{"api", "github", "com"}},
		{"foo_bar-baz", []string{"foo", "bar", "baz"}},
		{"single", []string{"single"}},
	}
	for _, tt := range tests {
		if diff := cmp.Diff(tt.want, splitOpenAPIName(tt.in)); diff != "" {
			t.Errorf("splitOpenAPIName(%q) mismatch (-want +got):\n%v", tt.in, diff)
		}
	}
	if got := splitOpenAPIName(""); len(got) != 0 {
		t.Errorf("splitOpenAPIName(\"\") = %v, want empty", got)
	}
}

func TestParseJSONTag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		tag         string
		wantName    string
		wantOmit    bool
		wantIgnored bool
	}{
		{name: "name only", tag: "name", wantName: "name"},
		{name: "omitempty", tag: "name,omitempty", wantName: "name", wantOmit: true},
		{name: "omitzero", tag: "id,omitzero", wantName: "id", wantOmit: true},
		{name: "ignored", tag: "-", wantIgnored: true},
		{name: "empty", tag: ""},
		{name: "empty name with omit", tag: ",omitempty", wantOmit: true},
		{name: "extra options", tag: "name,omitempty,string", wantName: "name", wantOmit: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			name, omit, ignored := parseJSONTag(tt.tag)
			if name != tt.wantName || omit != tt.wantOmit || ignored != tt.wantIgnored {
				t.Errorf("parseJSONTag(%q) = (%q, %v, %v), want (%q, %v, %v)",
					tt.tag, name, omit, ignored, tt.wantName, tt.wantOmit, tt.wantIgnored)
			}
		})
	}
}

func TestIsPointerTypeAndCanBeOmitted(t *testing.T) {
	t.Parallel()
	tests := []struct {
		expr          string
		wantPointer   bool
		wantOmittable bool
	}{
		{expr: "*int", wantPointer: true, wantOmittable: true},
		{expr: "[]string", wantPointer: false, wantOmittable: true},
		{expr: "map[string]int", wantPointer: false, wantOmittable: true},
		{expr: "interface{}", wantPointer: false, wantOmittable: true},
		{expr: "any", wantPointer: false, wantOmittable: true},
		{expr: "pkg.Type", wantPointer: false, wantOmittable: true},
		{expr: "string", wantPointer: false, wantOmittable: false},
		{expr: "int", wantPointer: false, wantOmittable: false},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			t.Parallel()
			e, err := parser.ParseExpr(tt.expr)
			if err != nil {
				t.Fatalf("ParseExpr(%q): %v", tt.expr, err)
			}
			if got := isPointerType(e); got != tt.wantPointer {
				t.Errorf("isPointerType(%q) = %v, want %v", tt.expr, got, tt.wantPointer)
			}
			if got := canBeOmitted(e); got != tt.wantOmittable {
				t.Errorf("canBeOmitted(%q) = %v, want %v", tt.expr, got, tt.wantOmittable)
			}
		})
	}
}

func TestSliceSet(t *testing.T) {
	t.Parallel()
	got := sliceSet([]string{"a", "b", "a"})
	want := map[string]bool{"a": true, "b": true}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("sliceSet mismatch (-want +got):\n%v", diff)
	}
	if got := sliceSet(nil); len(got) != 0 {
		t.Errorf("sliceSet(nil) = %v, want empty", got)
	}
}

func TestDiagLocation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		filename string
		line     int
		want     string
	}{
		{filename: "", line: 0, want: ""},
		{filename: "f.go", line: 0, want: "f.go"},
		{filename: "f.go", line: 12, want: "f.go:12"},
	}
	for _, tt := range tests {
		if got := diagLocation(tt.filename, tt.line); got != tt.want {
			t.Errorf("diagLocation(%q, %d) = %q, want %q", tt.filename, tt.line, got, tt.want)
		}
	}
}

func TestSchemaFieldDiagnosticString(t *testing.T) {
	t.Parallel()
	withLoc := schemaFieldDiagnostic{
		OpenAPISchema: "sch", GoStruct: "S", Field: "F", JSONName: "j",
		Message: "msg", Filename: "f.go", Line: 3, OpenAPIFile: "api.json",
	}
	if got, want := withLoc.String(), "f.go:3: S.F (j from sch): msg [api.json]"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	noLoc := schemaFieldDiagnostic{
		OpenAPISchema: "sch", GoStruct: "S", Field: "F", JSONName: "j", Message: "msg",
	}
	if got, want := noLoc.String(), "S.F (j from sch): msg"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestCanCheckOptionality(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		prop openapiSchemaProperty
		want bool
	}{
		{name: "plain", prop: openapiSchemaProperty{}, want: true},
		{name: "readOnly", prop: openapiSchemaProperty{readOnly: true}, want: false},
		{name: "writeOnly", prop: openapiSchemaProperty{writeOnly: true}, want: false},
	}
	for _, tt := range tests {
		if got := tt.prop.canCheckOptionality(); got != tt.want {
			t.Errorf("%s: canCheckOptionality() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

// fieldSet and structWith build the minimal shapes the field-matching helpers need.
func fieldSet[T any](names ...string) map[string]T {
	m := make(map[string]T, len(names))
	for _, name := range names {
		var zero T
		m[name] = zero
	}
	return m
}

func structWith(fields ...string) *goStructInfo {
	return &goStructInfo{fields: fieldSet[goStructField](fields...)}
}

func TestSharedFieldCount(t *testing.T) {
	t.Parallel()
	got := sharedFieldCount(
		&openapiSchemaFields{properties: fieldSet[openapiSchemaProperty]("a", "b", "c")},
		structWith("b", "c", "d"),
	)
	if got != 2 {
		t.Errorf("sharedFieldCount = %d, want 2", got)
	}
}

func TestSameJSONFieldSet(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		schema []string
		strct  []string
		want   bool
	}{
		{name: "equal", schema: []string{"a", "b"}, strct: []string{"a", "b"}, want: true},
		{name: "different length", schema: []string{"a", "b"}, strct: []string{"a"}, want: false},
		{name: "same length different members", schema: []string{"a", "b"}, strct: []string{"a", "c"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := sameJSONFieldSet(
				&openapiSchemaFields{properties: fieldSet[openapiSchemaProperty](tt.schema...)},
				structWith(tt.strct...),
			)
			if got != tt.want {
				t.Errorf("sameJSONFieldSet = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasEnoughSharedFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		schema []string
		strct  []string
		want   bool
	}{
		{name: "three shared", schema: []string{"a", "b", "c"}, strct: []string{"a", "b", "c"}, want: true},
		{name: "tiny type all shared", schema: []string{"a", "b"}, strct: []string{"a", "b"}, want: true},
		{name: "tiny type partial", schema: []string{"a", "b"}, strct: []string{"a", "x"}, want: false},
		{name: "percentage threshold met", schema: []string{"a", "b", "c"}, strct: []string{"a", "b", "x"}, want: true},
		{name: "percentage threshold missed", schema: []string{"a", "b", "c", "d"}, strct: []string{"a", "b", "x", "y"}, want: false},
		{name: "no overlap", schema: []string{"a", "b", "c"}, strct: []string{"x", "y", "z"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := hasEnoughSharedFields(
				&openapiSchemaFields{properties: fieldSet[openapiSchemaProperty](tt.schema...)},
				structWith(tt.strct...),
			)
			if got != tt.want {
				t.Errorf("hasEnoughSharedFields = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasUnsupportedComposition(t *testing.T) {
	t.Parallel()
	str := openapi3.NewSchemaRef("", openapi3.NewStringSchema())
	tests := []struct {
		name   string
		schema *openapi3.Schema
		want   bool
	}{
		{name: "plain object", schema: openapi3.NewObjectSchema(), want: false},
		{name: "oneOf", schema: &openapi3.Schema{OneOf: openapi3.SchemaRefs{str}}, want: true},
		{name: "anyOf", schema: &openapi3.Schema{AnyOf: openapi3.SchemaRefs{str}}, want: true},
		{name: "not", schema: &openapi3.Schema{Not: str}, want: true},
	}
	for _, tt := range tests {
		if got := hasUnsupportedComposition(tt.schema); got != tt.want {
			t.Errorf("%s: hasUnsupportedComposition = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestFlattenObjectSchema(t *testing.T) {
	t.Parallel()

	t.Run("plain object is returned unchanged", func(t *testing.T) {
		t.Parallel()
		obj := openapi3.NewObjectSchema()
		obj.Required = []string{"a"}
		got, reason, err := flattenObjectSchema(obj)
		if err != nil || reason != "" {
			t.Fatalf("flattenObjectSchema = (_, %q, %v)", reason, err)
		}
		if got != obj {
			t.Errorf("flattenObjectSchema returned a different schema for a plain object")
		}
	})

	t.Run("unsupported composition is skipped with a reason", func(t *testing.T) {
		t.Parallel()
		schema := &openapi3.Schema{OneOf: openapi3.SchemaRefs{openapi3.NewSchemaRef("", openapi3.NewStringSchema())}}
		got, reason, err := flattenObjectSchema(schema)
		if err != nil || got != nil || reason == "" {
			t.Fatalf("flattenObjectSchema = (%v, %q, %v), want (nil, non-empty reason, nil)", got, reason, err)
		}
	})

	t.Run("allOf is merged", func(t *testing.T) {
		t.Parallel()
		part := &openapi3.Schema{
			Required:   []string{"b"},
			Properties: openapi3.Schemas{"b": openapi3.NewSchemaRef("", openapi3.NewStringSchema())},
		}
		base := &openapi3.Schema{
			Required:   []string{"a"},
			Properties: openapi3.Schemas{"a": openapi3.NewSchemaRef("", openapi3.NewStringSchema())},
			AllOf:      openapi3.SchemaRefs{openapi3.NewSchemaRef("", part)},
		}
		got, reason, err := flattenObjectSchema(base)
		if err != nil || reason != "" {
			t.Fatalf("flattenObjectSchema = (_, %q, %v)", reason, err)
		}
		if _, ok := got.Properties["a"]; !ok {
			t.Error("merged schema missing property a")
		}
		if _, ok := got.Properties["b"]; !ok {
			t.Error("merged schema missing property b")
		}
		if len(got.Required) != 2 {
			t.Errorf("merged Required = %v, want a and b", got.Required)
		}
	})

	t.Run("nil schema is an error", func(t *testing.T) {
		t.Parallel()
		if _, _, err := flattenObjectSchema(nil); err == nil {
			t.Error("flattenObjectSchema(nil) = nil error, want error")
		}
	})
}

func TestSchemaProperties(t *testing.T) {
	t.Parallel()
	nullable := openapi3.NewStringSchema()
	nullable.Nullable = true
	readOnly := openapi3.NewStringSchema()
	readOnly.ReadOnly = true
	got := schemaProperties(openapi3.Schemas{
		"n":   openapi3.NewSchemaRef("", nullable),
		"r":   openapi3.NewSchemaRef("", readOnly),
		"nil": nil,
	})
	if !got["n"].nullable {
		t.Error("property n should be nullable")
	}
	if !got["r"].readOnly {
		t.Error("property r should be readOnly")
	}
	if _, ok := got["nil"]; !ok {
		t.Error("nil property ref should still yield a zero-value entry")
	}
}

//nolint:paralleltest // cannot use t.Parallel() when helper calls t.Setenv
func TestCheckSchemaFieldsCommand(t *testing.T) {
	testServer := newTestServer(t, "schema-ref", map[string]any{
		"api.github.com/api.github.com.json": openapi3.T{
			Components: &openapi3.Components{
				Schemas: openapi3.Schemas{
					"demo": openapi3.NewSchemaRef("", &openapi3.Schema{
						Required: []string{"id", "name"},
						Properties: openapi3.Schemas{
							"id":   openapi3.NewSchemaRef("", openapi3.NewIntegerSchema()),
							"name": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
							"note": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
						},
					}),
				},
			},
		},
	})

	res := runTest(t, "testdata/check-schema-fields", "check-schema-fields", "--github-url", testServer.URL)
	res.assertOutput("Found 0 schema field issues\nChecked 1 OpenAPI schema/Go struct pairs; skipped 0 OpenAPI schemas", "")
	res.assertNoErr()
	res.checkGolden()
}

func TestCheckSchemaFieldsSkipsStructMatchingMultipleSchemas(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

type ItemRequest struct {
	ID   int64  `+"`json:\"id\"`"+`
	Type string `+"`json:\"type\"`"+`
}

func (s *svc) Add(ctx context.Context, body *ItemRequest) {
	s.client.NewRequest(ctx, "POST", "u", body)
}
`)

	itemSchema := func() *openapi3.SchemaRef {
		return openapi3.NewSchemaRef("", &openapi3.Schema{
			Required: []string{"id", "type"},
			Properties: openapi3.Schemas{
				"id":   openapi3.NewSchemaRef("", openapi3.NewIntegerSchema()),
				"type": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
			},
		})
	}
	result, err := checkSchemaFields(schemaFieldCheckOptions{
		descriptions: []*openapiFile{testOpenAPIFile("descriptions/api.github.com/api.github.com.json", openapi3.Schemas{
			"schema-a": itemSchema(),
			"schema-b": itemSchema(),
		})},
		githubDir: githubDir,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Checked) != 0 {
		t.Fatalf("checked = %v, want none", result.Checked)
	}
	if len(result.Diagnostics) != 0 {
		t.Fatalf("diagnostics = %v, want none", result.Diagnostics)
	}
	var dropped bool
	for _, skip := range result.Skipped {
		if strings.Contains(skip.Reason, "matches multiple schemas by field set") {
			dropped = true
		}
	}
	if !dropped {
		t.Errorf("skipped = %#v, want a \"matches multiple schemas by field set\" reason", result.Skipped)
	}
}

func TestFilterAllowedSchemaFieldDiagnostics(t *testing.T) {
	t.Parallel()
	exceptions := sliceSet([]string{"ExemptStruct.ExemptField"})
	got := filterAllowedSchemaFieldDiagnostics([]*schemaFieldDiagnostic{
		{GoStruct: "ExemptStruct", Field: "ExemptField"},
		{GoStruct: "NotExemptStruct", Field: "NotExemptField"},
	}, exceptions)
	if len(got) != 1 || got[0].GoStruct != "NotExemptStruct" {
		t.Errorf("filterAllowedSchemaFieldDiagnostics = %+v, want only NotExemptStruct.NotExemptField", got)
	}
}

func TestLoadSchemaFieldExceptions(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "schema_field_exceptions.yaml")
	writeFile(t, path, `# comment
exceptions:
  - StructA.FieldA
  - StructB.FieldB # TODO: fix
`)

	got, err := loadSchemaFieldExceptions(path, false)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"StructA.FieldA": true, "StructB.FieldB": true}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("loadSchemaFieldExceptions mismatch (-want +got):\n%v", diff)
	}

	// A missing optional file yields an empty set and no error.
	missing := filepath.Join(dir, "does-not-exist.yaml")
	got, err = loadSchemaFieldExceptions(missing, true)
	if err != nil {
		t.Fatalf("optional missing file: unexpected error %v", err)
	}
	if len(got) != 0 {
		t.Errorf("optional missing file: got %v, want empty", got)
	}

	// A missing required file is an error.
	if _, err := loadSchemaFieldExceptions(missing, false); err == nil {
		t.Error("required missing file: got nil error, want error")
	}
}

// TestSchemaFieldExceptionsFileParses guards the committed exceptions file so a
// malformed edit is caught by unit tests rather than only in CI.
func TestSchemaFieldExceptionsFileParses(t *testing.T) {
	t.Parallel()
	// Tests run with the package directory as the working directory, so the committed
	// file (tools/metadata/schema_field_exceptions.yaml) is reachable by its basename.
	got, err := loadSchemaFieldExceptions("schema_field_exceptions.yaml", false)
	if err != nil {
		t.Fatal(err)
	}
	for key := range got {
		if _, _, ok := strings.Cut(key, "."); !ok {
			t.Errorf("exception %q is not in Struct.Field form", key)
		}
	}
}

func testOpenAPIFile(filename string, schemas openapi3.Schemas) *openapiFile {
	return &openapiFile{
		filename: filename,
		description: &openapi3.T{
			Components: &openapi3.Components{
				Schemas: schemas,
			},
		},
	}
}

func writeFile(t *testing.T, filename, content string) {
	t.Helper()
	if err := os.WriteFile(filename, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
