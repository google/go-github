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

func TestCheckSchemaFieldsAnnotatedStruct(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

import "encoding/json"

// Demo is a demo request body.
//
//meta:schema request POST /demo
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
		descriptions: []*openapiFile{
			testOpenAPIFile("descriptions/api.github.com/api.github.com.json",
				"/demo", &openapi3.PathItem{Post: testRequestBodyOperation(&openapi3.Schema{
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
				})}),
		},
		githubDir: githubDir,
	})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff([]*schemaFieldChecked{{
		Annotation:  "request POST /demo",
		GoStruct:    "Demo",
		OpenAPIFile: "descriptions/api.github.com/api.github.com.json",
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

func TestCheckSchemaFieldsMultipleAnnotations(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

// DemoRequest is used by two operations.
//
//meta:schema request POST /repos/{owner}/{repo}/demos
//meta:schema request PATCH /repos/{owner}/{repo}/demos/{demo_id}
type DemoRequest struct {
	Body string `+"`json:\"body\"`"+`
}
`)

	bodySchema := func() *openapi3.Schema {
		return &openapi3.Schema{
			Required:   []string{"body"},
			Properties: openapi3.Schemas{"body": openapi3.NewSchemaRef("", openapi3.NewStringSchema())},
		}
	}
	descriptions := []*openapiFile{
		testOpenAPIFile("descriptions/api.github.com/api.github.com.json",
			"/repos/{owner}/{repo}/demos", &openapi3.PathItem{Post: testRequestBodyOperation(bodySchema())}),
		testOpenAPIFile("descriptions/ghec/ghec.json",
			"/repos/{owner}/{repo}/demos/{demo_id}", &openapi3.PathItem{Patch: testRequestBodyOperation(bodySchema())}),
	}

	result, err := checkSchemaFields(schemaFieldCheckOptions{descriptions: descriptions, githubDir: githubDir})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff([]*schemaFieldChecked{
		{Annotation: "request PATCH /repos/{owner}/{repo}/demos/{demo_id}", GoStruct: "DemoRequest", OpenAPIFile: "descriptions/ghec/ghec.json"},
		{Annotation: "request POST /repos/{owner}/{repo}/demos", GoStruct: "DemoRequest", OpenAPIFile: "descriptions/api.github.com/api.github.com.json"},
	}, result.Checked); diff != "" {
		t.Errorf("checked mismatch (-want +got):\n%v", diff)
	}
	if len(result.Diagnostics) != 0 {
		t.Errorf("diagnostics = %v, want none", result.Diagnostics)
	}
	if result.Summary.AnnotatedStructs != 1 {
		t.Errorf("AnnotatedStructs = %v, want 1", result.Summary.AnnotatedStructs)
	}
}

func TestCheckSchemaFieldsUnresolvedAnnotation(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

//meta:schema request POST /missing
type Demo struct {
	Body string `+"`json:\"body\"`"+`
}
`)

	result, err := checkSchemaFields(schemaFieldCheckOptions{
		descriptions: []*openapiFile{testOpenAPIFile("descriptions/api.github.com/api.github.com.json",
			"/demo", &openapi3.PathItem{Post: testRequestBodyOperation(openapi3.NewObjectSchema())})},
		githubDir: githubDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Diagnostics) != 1 {
		t.Fatalf("diagnostics = %v, want one", result.Diagnostics)
	}
	if got, want := result.Diagnostics[0].Message, "could not find operation POST /missing in any OpenAPI description"; got != want {
		t.Errorf("message = %q, want %q", got, want)
	}
}

func TestCheckSchemaFieldsInvalidAnnotation(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

//meta:schema POST /demo
type Demo struct{}

//meta:schema body POST /demo
type Demo2 struct{}
`)

	result, err := checkSchemaFields(schemaFieldCheckOptions{
		descriptions: []*openapiFile{testOpenAPIFile("descriptions/api.github.com/api.github.com.json",
			"/demo", &openapi3.PathItem{Post: testRequestBodyOperation(openapi3.NewObjectSchema())})},
		githubDir: githubDir,
	})
	if err != nil {
		t.Fatal(err)
	}

	var got []string
	for _, diag := range result.Diagnostics {
		got = append(got, diag.GoStruct+": "+diag.Message)
	}
	want := []string{
		`Demo: invalid annotation "meta:schema POST /demo": want //meta:schema <request|response> <METHOD> <path>`,
		`Demo2: invalid annotation "meta:schema body POST /demo": unknown role "body"; want request or response`,
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("diagnostics mismatch (-want +got):\n%v", diff)
	}
}

func TestCheckSchemaFieldsResponseRole(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

//meta:schema response GET /demo
type Demo struct {
	ID int64 `+"`json:\"id\"`"+`
}
`)

	op := &openapi3.Operation{Responses: openapi3.NewResponses(openapi3.WithStatus(200, &openapi3.ResponseRef{
		Value: &openapi3.Response{Content: openapi3.NewContentWithJSONSchema(&openapi3.Schema{
			Required:   []string{"id"},
			Properties: openapi3.Schemas{"id": openapi3.NewSchemaRef("", openapi3.NewIntegerSchema())},
		})},
	}))}
	result, err := checkSchemaFields(schemaFieldCheckOptions{
		descriptions: []*openapiFile{testOpenAPIFile("descriptions/api.github.com/api.github.com.json",
			"/demo", &openapi3.PathItem{Get: op})},
		githubDir: githubDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Checked) != 1 || len(result.Diagnostics) != 0 {
		t.Errorf("checked = %v, diagnostics = %v; want one check and no diagnostics", result.Checked, result.Diagnostics)
	}
}

func TestCheckSchemaFieldsUnsupportedComposition(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

//meta:schema request POST /demo
type Demo struct{}
`)

	oneOf := &openapi3.Schema{OneOf: openapi3.SchemaRefs{openapi3.NewSchemaRef("", openapi3.NewStringSchema())}}
	result, err := checkSchemaFields(schemaFieldCheckOptions{
		descriptions: []*openapiFile{testOpenAPIFile("descriptions/api.github.com/api.github.com.json",
			"/demo", &openapi3.PathItem{Post: testRequestBodyOperation(oneOf)})},
		githubDir: githubDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Diagnostics) != 1 {
		t.Fatalf("diagnostics = %v, want one", result.Diagnostics)
	}
	if !strings.Contains(result.Diagnostics[0].Message, "annotated schema cannot be checked") {
		t.Errorf("message = %q, want an unsupported-composition diagnostic", result.Diagnostics[0].Message)
	}
}

func TestCheckSchemaFieldsIgnoresUnannotatedStructs(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

// Unannotated has fields that would fail the check if it were annotated.
type Unannotated struct {
	Name *string `+"`json:\"name\"`"+`
}
`)

	result, err := checkSchemaFields(schemaFieldCheckOptions{
		descriptions: []*openapiFile{testOpenAPIFile("descriptions/api.github.com/api.github.com.json",
			"/demo", &openapi3.PathItem{Post: testRequestBodyOperation(openapi3.NewObjectSchema())})},
		githubDir: githubDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Checked) != 0 || len(result.Diagnostics) != 0 || result.Summary.AnnotatedStructs != 0 {
		t.Errorf("result = %+v, want nothing checked for an unannotated struct", result)
	}
}

func TestCheckSchemaFieldsAllowsRequiredNullablePointer(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

//meta:schema request POST /demo
type Demo struct {
	Name *string `+"`json:\"name\"`"+`
}
`)

	nullable := openapi3.NewStringSchema()
	nullable.Nullable = true
	result, err := checkSchemaFields(schemaFieldCheckOptions{
		descriptions: []*openapiFile{testOpenAPIFile("descriptions/api.github.com/api.github.com.json",
			"/demo", &openapi3.PathItem{Post: testRequestBodyOperation(&openapi3.Schema{
				Required:   []string{"name"},
				Properties: openapi3.Schemas{"name": openapi3.NewSchemaRef("", nullable)},
			})})},
		githubDir: githubDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Diagnostics) != 0 {
		t.Errorf("diagnostics = %v, want none for a required nullable pointer field", result.Diagnostics)
	}
}

func TestResolveSchemaAnnotation(t *testing.T) {
	t.Parallel()
	schema := &openapi3.Schema{Properties: openapi3.Schemas{"body": openapi3.NewSchemaRef("", openapi3.NewStringSchema())}}
	descriptions := []*openapiFile{
		testOpenAPIFile("descriptions/api.github.com/api.github.com.json",
			"/repos/{owner}/{repo}/demos", &openapi3.PathItem{Post: testRequestBodyOperation(schema)}),
		testOpenAPIFile("descriptions/ghes-3.21/ghes-3.21.json",
			"/admin/demos", &openapi3.PathItem{Post: testRequestBodyOperation(schema)}),
	}

	t.Run("path parameter names are normalized", func(t *testing.T) {
		t.Parallel()
		got, file, problem := resolveSchemaAnnotation(descriptions, &schemaAnnotation{role: "request", method: "POST", path: "/repos/{o}/{r}/demos"})
		if got == nil || problem != "" || file != "descriptions/api.github.com/api.github.com.json" {
			t.Errorf("resolveSchemaAnnotation = (%v, %q, %q), want the api.github.com schema", got, file, problem)
		}
	})

	t.Run("operation only in a later description is found", func(t *testing.T) {
		t.Parallel()
		got, file, problem := resolveSchemaAnnotation(descriptions, &schemaAnnotation{role: "request", method: "POST", path: "/admin/demos"})
		if got == nil || problem != "" || file != "descriptions/ghes-3.21/ghes-3.21.json" {
			t.Errorf("resolveSchemaAnnotation = (%v, %q, %q), want the ghes schema", got, file, problem)
		}
	})

	t.Run("wrong method is not found", func(t *testing.T) {
		t.Parallel()
		got, _, problem := resolveSchemaAnnotation(descriptions, &schemaAnnotation{role: "request", method: "PATCH", path: "/repos/{owner}/{repo}/demos"})
		if got != nil || !strings.Contains(problem, "could not find operation") {
			t.Errorf("resolveSchemaAnnotation = (%v, %q), want a not-found problem", got, problem)
		}
	})

	t.Run("missing request body is a problem", func(t *testing.T) {
		t.Parallel()
		noBody := []*openapiFile{testOpenAPIFile("descriptions/api.github.com/api.github.com.json",
			"/demo", &openapi3.PathItem{Post: &openapi3.Operation{}})}
		got, _, problem := resolveSchemaAnnotation(noBody, &schemaAnnotation{role: "request", method: "POST", path: "/demo"})
		if got != nil || problem != "operation has no request body" {
			t.Errorf("resolveSchemaAnnotation = (%v, %q), want a no-request-body problem", got, problem)
		}
	})
}

func TestParseSchemaAnnotationsViaCollectGoStructs(t *testing.T) {
	t.Parallel()
	githubDir := t.TempDir()
	writeFile(t, filepath.Join(githubDir, "demo.go"), `package github

// Demo has annotations in mixed case with surrounding doc text.
//
//meta:schema Request post /demo
//meta:schema response GET /demo
type Demo struct{}

type Group struct{}
`)

	structs, err := collectGoStructs(githubDir)
	if err != nil {
		t.Fatal(err)
	}

	demo := structs["Demo"]
	if demo == nil {
		t.Fatal("Demo struct not collected")
	}
	var got []string
	for _, ann := range demo.annotations {
		got = append(got, ann.String())
	}
	want := []string{"request POST /demo", "response GET /demo"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("annotations mismatch (-want +got):\n%v", diff)
	}
	if len(demo.annotationProblems) != 0 {
		t.Errorf("annotationProblems = %v, want none", demo.annotationProblems)
	}
	if group := structs["Group"]; group == nil || len(group.annotations) != 0 {
		t.Errorf("Group = %+v, want collected with no annotations", group)
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
			t.Errorf("diagLocation(%q, %v) = %q, want %q", tt.filename, tt.line, got, tt.want)
		}
	}
}

func TestSchemaFieldDiagnosticString(t *testing.T) {
	t.Parallel()
	withLoc := schemaFieldDiagnostic{
		Annotation: "request POST /demo", GoStruct: "S", Field: "F", JSONName: "j",
		Message: "msg", Filename: "f.go", Line: 3, OpenAPIFile: "api.json",
	}
	if got, want := withLoc.String(), "f.go:3: S.F (j from request POST /demo): msg [api.json]"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	noLoc := schemaFieldDiagnostic{
		Annotation: "request POST /demo", GoStruct: "S", Field: "F", JSONName: "j", Message: "msg",
	}
	if got, want := noLoc.String(), "S.F (j from request POST /demo): msg"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	annotationLevel := schemaFieldDiagnostic{
		GoStruct: "S", Message: "msg", Filename: "f.go", Line: 3,
	}
	if got, want := annotationLevel.String(), "f.go:3: S: msg"; got != want {
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
			t.Errorf("%v: canCheckOptionality() = %v, want %v", tt.name, got, tt.want)
		}
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
			t.Errorf("%v: hasUnsupportedComposition = %v, want %v", tt.name, got, tt.want)
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
			t.Error("flattenObjectSchema returned a different schema for a plain object")
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
			Paths: openapi3.NewPaths(openapi3.WithPath("/demo", &openapi3.PathItem{
				Post: testRequestBodyOperation(&openapi3.Schema{
					Required: []string{"id", "name"},
					Properties: openapi3.Schemas{
						"id":   openapi3.NewSchemaRef("", openapi3.NewIntegerSchema()),
						"name": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
						"note": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
					},
				}),
			})),
		},
	})

	res := runTest(t, "testdata/check-schema-fields", "check-schema-fields", "--github-url", testServer.URL)
	res.assertOutput("Found 0 schema field issues\nChecked 1 annotations on 1 annotated structs", "")
	res.assertNoErr()
	res.checkGolden()
}

func TestFilterAllowedSchemaFieldDiagnostics(t *testing.T) {
	t.Parallel()
	exceptions := []string{"ExemptStruct.ExemptField"}
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
	metadataDir := filepath.Join(dir, "tools", "metadata")
	if err := os.MkdirAll(metadataDir, 0o700); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(metadataDir, "schema_field_exceptions.yaml"), `# comment
exceptions:
  - StructA.FieldA
  - StructB.FieldB # TODO: fix
`)

	got, err := loadSchemaFieldExceptions(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"StructA.FieldA", "StructB.FieldB"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("loadSchemaFieldExceptions mismatch (-want +got):\n%v", diff)
	}

	// A missing file yields no exceptions and no error.
	got, err = loadSchemaFieldExceptions(t.TempDir())
	if err != nil {
		t.Fatalf("missing file: unexpected error %v", err)
	}
	if len(got) != 0 {
		t.Errorf("missing file: got %v, want empty", got)
	}
}

// TestSchemaFieldExceptionsFileParses guards the committed exceptions file so a
// malformed edit is caught by unit tests rather than only in CI.
func TestSchemaFieldExceptionsFileParses(t *testing.T) {
	t.Parallel()
	// Tests run with the package directory as the working directory, so the repository root that
	// loadSchemaFieldExceptions joins with the fixed relative path is "../..".
	got, err := loadSchemaFieldExceptions(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range got {
		if _, _, ok := strings.Cut(key, "."); !ok {
			t.Errorf("exception %q is not in Struct.Field form", key)
		}
	}
}

func testOpenAPIFile(filename, path string, pathItem *openapi3.PathItem) *openapiFile {
	return &openapiFile{
		filename: filename,
		description: &openapi3.T{
			Paths: openapi3.NewPaths(openapi3.WithPath(path, pathItem)),
		},
	}
}

func testRequestBodyOperation(schema *openapi3.Schema) *openapi3.Operation {
	return &openapi3.Operation{
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Content: openapi3.NewContentWithJSONSchema(schema),
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
