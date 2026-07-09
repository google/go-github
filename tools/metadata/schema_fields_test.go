// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
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
	got := goNameCandidates("projects-v2")
	want := []string{"ProjectsV2", "ProjectV2"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("goNameCandidates mismatch (-want +got):\n%v", diff)
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
	var exemptStruct, exemptField string
	for key := range schemaFieldExceptions {
		if s, f, ok := strings.Cut(key, "."); ok {
			exemptStruct, exemptField = s, f
			break
		}
	}
	if exemptStruct == "" {
		t.Skip("no schema field exceptions configured")
	}

	got := filterAllowedSchemaFieldDiagnostics([]*schemaFieldDiagnostic{
		{GoStruct: exemptStruct, Field: exemptField},
		{GoStruct: "NotExemptStruct", Field: "NotExemptField"},
	})
	if len(got) != 1 || got[0].GoStruct != "NotExemptStruct" {
		t.Errorf("filterAllowedSchemaFieldDiagnostics = %+v, want only NotExemptStruct.NotExemptField", got)
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
