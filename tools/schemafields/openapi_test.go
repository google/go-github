// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
)

// flattenDesc is a description that exercises every composition this tool flattens.
const flattenDesc = `{
  "components": {
    "schemas": {
      "base": {
        "type": "object",
        "required": ["name"],
        "properties": {"name": {"type": "string"}}
      },
      "all-of": {
        "allOf": [
          {"$ref": "#/components/schemas/base"},
          {"type": "object", "required": ["extra"], "properties": {"extra": {"type": "string"}}}
        ]
      },
      "one-of": {
        "oneOf": [
          {"required": ["kind", "kind_tag", "size"], "properties": {"kind": {"type": "string"}, "kind_tag": {"type": "string"}, "size": {"type": "integer"}}},
          {"required": ["kind", "kind_tag"], "properties": {"kind": {"type": "string"}, "kind_tag": {"type": "string"}}}
        ]
      },
      "any-of": {
        "anyOf": [
          {"required": ["a"], "properties": {"a": {"type": "string"}}},
          {"required": ["b"], "properties": {"b": {"type": "string"}}}
        ]
      },
      "metas": {
        "type": "object",
        "properties": {
          "ro": {"type": "string", "readOnly": true},
          "wo": {"type": "string", "writeOnly": true},
          "nul": {"type": "string", "nullable": true},
          "plain": {"type": "string"}
        }
      },
      "circular": {"$ref": "#/components/schemas/circular"},
      "dangling": {"$ref": "#/components/schemas/nowhere"}
    }
  }
}`

func descFromJSON(t *testing.T, src string) *description {
	t.Helper()
	d := newDescription()
	assertNilError(t, json.Unmarshal([]byte(src), d))
	return d
}

func flattenNamed(t *testing.T, d *description, name string) *flat {
	t.Helper()
	f, err := d.flatten(&schema{Ref: "#/components/schemas/" + name}, map[string]bool{})
	assertNilError(t, err)
	return f
}

func TestFlattenAllOf(t *testing.T) {
	t.Parallel()
	f := flattenNamed(t, descFromJSON(t, flattenDesc), "all-of")
	// allOf is a union: the referenced schema's required properties apply too.
	assertEqual(t, map[string]bool{"name": true, "extra": true}, f.required)
	assertEqual(t, []string{"extra", "name"}, sortedProps(f))
}

func TestFlattenOneOf(t *testing.T) {
	t.Parallel()
	f := flattenNamed(t, descFromJSON(t, flattenDesc), "one-of")
	// Only a property that EVERY variant requires is unconditionally required; one that
	// only some variants require is conditionally required, and is left to judgement.
	assertEqual(t, map[string]bool{"kind": true, "kind_tag": true}, f.required)
	assertEqual(t, []string{"kind", "kind_tag", "size"}, sortedProps(f))
}

func TestFlattenAnyOf(t *testing.T) {
	t.Parallel()
	f := flattenNamed(t, descFromJSON(t, flattenDesc), "any-of")
	// Neither variant's required property is required by the other, so neither is
	// unconditionally required.
	assertEqual(t, map[string]bool{}, f.required)
	assertEqual(t, []string{"a", "b"}, sortedProps(f))
}

func TestFlattenPropMeta(t *testing.T) {
	t.Parallel()
	f := flattenNamed(t, descFromJSON(t, flattenDesc), "metas")
	assertEqual(t, map[string]propMeta{
		"ro":    {readOnly: true},
		"wo":    {writeOnly: true},
		"nul":   {nullable: true},
		"plain": {},
	}, f.props, cmpopts.EquateComparable(propMeta{}))
}

func TestFlattenErrors(t *testing.T) {
	t.Parallel()
	d := descFromJSON(t, flattenDesc)
	if _, err := d.flatten(&schema{Ref: "#/components/schemas/circular"}, map[string]bool{}); err == nil {
		t.Error("flatten accepted a circular $ref")
	}
	if _, err := d.flatten(&schema{Ref: "#/components/schemas/dangling"}, map[string]bool{}); err == nil {
		t.Error("flatten accepted an unresolved $ref")
	}
	// A schema with no content at all is not an error: it just says nothing.
	f, err := d.flatten(nil, map[string]bool{})
	assertNilError(t, err)
	assertEqual(t, 0, len(f.props))
}

// requestDesc is a description whose single path has one operation of every shape.
const requestDesc = `{
  "paths": {
    "/a": {
      "post": {"requestBody": {"content": {"application/json": {"schema": {"$ref": "#/components/schemas/a"}}}}},
      "get": {},
      "put": {"requestBody": {"content": {"text/plain": {"schema": {"type": "string"}}}}},
      "patch": {"requestBody": {"content": {"application/json": {"schema": {"type": "object", "properties": {}}}}}},
      "delete": {"requestBody": {"content": {"application/json": {"$ref": "#/components/schemas/a"}}}}
    }
  },
  "components": {
    "schemas": {"a": {"type": "object", "required": ["x"], "properties": {"x": {"type": "string"}}}}
  }
}`

func TestRequestSchema(t *testing.T) {
	t.Parallel()
	d := descFromJSON(t, requestDesc)
	tests := []struct {
		name       string
		op         opRef
		wantReason bodyReason
	}{
		{"json body", opRef{method: "POST", path: "/a"}, bodyFound},
		{"lowercase method", opRef{method: "post", path: "/a"}, bodyFound},
		{"no request body", opRef{method: "GET", path: "/a"}, bodyNoJSONBody},
		{"no json content", opRef{method: "PUT", path: "/a"}, bodyNoJSONBody},
		{"no properties", opRef{method: "PATCH", path: "/a"}, bodyNoJSONBody},
		{"no schema", opRef{method: "DELETE", path: "/a"}, bodyNoJSONBody},
		{"unknown path", opRef{method: "POST", path: "/b"}, bodyNotInPlan},
		{"unknown method", opRef{method: "TRACE", path: "/a"}, bodyNotInPlan},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f, reason, err := d.requestSchema(&tt.op)
			assertNilError(t, err)
			assertEqual(t, tt.wantReason, reason)
			if reason == bodyFound {
				assertEqual(t, map[string]bool{"x": true}, f.required)
			}
		})
	}

	// The result is cached, so a repeated lookup does not re-flatten, and a repeated miss
	// keeps the reason it was recorded with.
	first, _, err := d.requestSchema(&opRef{method: "POST", path: "/a"})
	assertNilError(t, err)
	second, _, err := d.requestSchema(&opRef{method: "POST", path: "/a"})
	assertNilError(t, err)
	assertEqual(t, true, first == second)
	_, firstReason, err := d.requestSchema(&opRef{method: "GET", path: "/a"})
	assertNilError(t, err)
	_, secondReason, err := d.requestSchema(&opRef{method: "GET", path: "/a"})
	assertNilError(t, err)
	assertEqual(t, firstReason, secondReason)
}

func TestDescriptionsPlanOrder(t *testing.T) {
	t.Parallel()
	// The first plan that documents an operation wins, so GHEC-only operations resolve
	// even though api.github.com is consulted first.
	api := descFromJSON(t, requestDesc)
	c := descFromJSON(t, `{"paths": {"/c": {"post": {"requestBody": {"content":
		{"application/json": {"schema": {"$ref": "#/components/schemas/c"}}}}}}},
		"components": {"schemas": {"c": {"type": "object", "required": ["y"], "properties": {"y": {"type": "string"}}}}}}`)
	ds := &descriptions{files: []*description{api, c}}

	f, reason, err := ds.requestSchema(&opRef{method: "POST", path: "/a"})
	assertNilError(t, err)
	assertEqual(t, bodyFound, reason)
	assertEqual(t, map[string]bool{"x": true}, f.required)

	f, reason, err = ds.requestSchema(&opRef{method: "POST", path: "/c"})
	assertNilError(t, err)
	assertEqual(t, bodyFound, reason)
	assertEqual(t, map[string]bool{"y": true}, f.required)

	// A plan that documents the operation without a JSON body outvotes the plans that do not
	// document it at all: only an operation no plan documents is missing from the revision.
	_, reason, err = ds.requestSchema(&opRef{method: "GET", path: "/a"})
	assertNilError(t, err)
	assertEqual(t, bodyNoJSONBody, reason)

	_, reason, err = ds.requestSchema(&opRef{method: "POST", path: "/nowhere"})
	assertNilError(t, err)
	assertEqual(t, bodyNotInPlan, reason)
}

func TestSelectPlans(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		files []string
		want  []string
	}{
		{"empty", nil, nil},
		{
			"one of each",
			[]string{
				"descriptions/api.github.com/api.github.com.json",
				"descriptions/ghec/ghec.json",
				"descriptions/ghes-3.9/ghes-3.9.json",
			},
			[]string{
				"descriptions/api.github.com/api.github.com.json",
				"descriptions/ghec/ghec.json",
				"descriptions/ghes-3.9/ghes-3.9.json",
			},
		},
		{
			// The newest GHES release wins, and the older ones are dropped.
			"newest ghes",
			[]string{
				"descriptions/ghes-3.10/ghes-3.10.json",
				"descriptions/ghes-3.9/ghes-3.9.json",
				"descriptions/ghes-3.21/ghes-3.21.json",
				"descriptions/ghes-3.2/ghes-3.2.json",
			},
			[]string{"descriptions/ghes-3.21/ghes-3.21.json"},
		},
		{
			"ignores other plans and malformed paths",
			[]string{
				"descriptions/ghes-2.20/ghes-2.20.json",
				"descriptions/ghes-3/ghes-3.json",
				"descriptions/ghes-x.y/ghes-x.y.json",
				"descriptions/other/other.json",
				"descriptions/api.github.com/README.md",
				"openapi_operations.yaml",
			},
			nil,
		},
		{
			"removes duplicates",
			[]string{
				"descriptions/ghec/ghec.json",
				"descriptions/ghec/ghec.json",
			},
			[]string{"descriptions/ghec/ghec.json"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertEqual(t, tt.want, selectPlans(tt.files))
		})
	}
}

func TestParseVersion(t *testing.T) {
	t.Parallel()
	major, minor, err := parseVersion("3.21")
	assertNilError(t, err)
	assertEqual(t, 3, major)
	assertEqual(t, 21, minor)

	if _, _, err := parseVersion("3"); err == nil {
		t.Error("parseVersion accepted a version without a minor number")
	}
	if _, _, err := parseVersion("x.y"); err == nil {
		t.Error("parseVersion accepted a non-numeric version")
	}
}

func TestPinnedDescriptions(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, operationsFile),
		"openapi_commit: a1b2c3\n"+
			"openapi_operations:\n"+
			"- descriptions/ghes-3.9/ghes-3.9.json\n"+
			"- descriptions/api.github.com/api.github.com.json\n"+
			"- descriptions/ghec/ghec.json\n"+
			"other: ignored\n")
	ref, plans, err := pinnedDescriptions(dir)
	assertNilError(t, err)
	assertEqual(t, "a1b2c3", ref)
	assertEqual(t, []string{
		"descriptions/api.github.com/api.github.com.json",
		"descriptions/ghec/ghec.json",
		"descriptions/ghes-3.9/ghes-3.9.json",
	}, plans)

	// A file without a pinned commit cannot be used to download anything.
	mustWriteFile(t, filepath.Join(dir, operationsFile), "openapi_operations: []\n")
	if _, _, err := pinnedDescriptions(dir); err == nil {
		t.Error("pinnedDescriptions accepted a file without openapi_commit")
	}
	if _, _, err := pinnedDescriptions(t.TempDir()); err == nil {
		t.Error("pinnedDescriptions accepted a checkout without openapi_operations.yaml")
	}
}

// TestDownload checks that a description is fetched from the raw GitHub URL layout into its
// destination, and that a failure leaves nothing behind.
func TestDownload(t *testing.T) {
	t.Parallel()
	const path = "descriptions/api.github.com/api.github.com.json"
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		if r.URL.Path != "/"+descriptionsOwner+"/"+descriptionsRepo+"/REF/"+path {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, `{"cached": true}`)
	}))
	defer srv.Close()

	dir := t.TempDir()
	dest := filepath.Join(dir, "out.json")
	assertNilError(t, download(t.Context(), srv.Client(), srv.URL, "REF", path, dest))
	assertEqual(t, "/github/rest-api-description/REF/"+path, gotPath)
	data, err := os.ReadFile(dest)
	assertNilError(t, err)
	assertEqual(t, `{"cached": true}`, string(data))

	// A 404 is an error, and the partial file is not left in place.
	dest = filepath.Join(dir, "missing.json")
	err = download(t.Context(), srv.Client(), srv.URL, "REF", "missing.json", dest)
	if err == nil {
		t.Fatal("download accepted a 404 response")
	}
	assertContains(t, err.Error(), "404")
	if _, err := os.Stat(dest); err == nil {
		t.Error("download left a file behind after a failed request")
	}
	entries, err := os.ReadDir(dir)
	assertNilError(t, err)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".tmp") {
			t.Errorf("download left the temporary file %v behind", e.Name())
		}
	}
}

// TestFetchDescriptionsCached checks that a description already in the cache is reused,
// which is what keeps a repeated run offline.
func TestFetchDescriptionsCached(t *testing.T) {
	t.Parallel()
	cacheDir := t.TempDir()
	const path = "descriptions/ghec/ghec.json"
	dest := filepath.Join(cacheDir, "REF", "descriptions_ghec_ghec.json")
	mustWriteFile(t, dest, `{"cached": true}`)

	got, err := fetchDescriptions(t.Context(), "REF", []string{path}, cacheDir)
	assertNilError(t, err)
	assertEqual(t, []string{dest}, got)
}

// TestFetchDescriptionsPrunesTheCache checks that the cache does not grow with every revision
// the repository is pinned to: the descriptions of the revision a run uses are kept, and those
// of every other revision are dropped by the end of that run.
func TestFetchDescriptionsPrunesTheCache(t *testing.T) {
	t.Parallel()
	cacheDir := t.TempDir()
	const path = "descriptions/ghec/ghec.json"
	cached := "descriptions_ghec_ghec.json"
	// The revision in use, and two that a run has no use for.
	for _, ref := range []string{"in-use", "older", "oldest"} {
		mustWriteFile(t, filepath.Join(cacheDir, ref, cached), `{"cached": true}`)
	}
	// A directory that this tool did not fill is not the cache's to remove.
	foreign := filepath.Join(cacheDir, "notes", "notes.txt")
	mustWriteFile(t, foreign, "mine\n")

	// The requested file is in the cache, so nothing is downloaded.
	got, err := fetchDescriptions(t.Context(), "in-use", []string{path}, cacheDir)
	assertNilError(t, err)
	assertEqual(t, []string{filepath.Join(cacheDir, "in-use", cached)}, got)

	tests := []struct {
		name string
		want bool
	}{
		{"in-use", true},
		{"older", false},
		{"oldest", false},
		{"notes", true},
	}
	for _, tt := range tests {
		_, err := os.Stat(filepath.Join(cacheDir, tt.name))
		if got := err == nil; got != tt.want {
			t.Errorf("%v is in the cache = %v, want %v", tt.name, got, tt.want)
		}
	}
}

// sortedProps returns the property names of f in sorted order.
func sortedProps(f *flat) []string {
	names := make([]string, 0, len(f.props))
	for name := range f.props {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}
