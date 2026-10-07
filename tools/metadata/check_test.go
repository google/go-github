// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCheckProblems(t *testing.T) {
	t.Parallel()
	generated := []*operation{
		{Name: "GET /a/{a_id}", DocumentationURL: "https://docs.github.com/rest/a/a#get-a"},
	}
	for _, td := range []struct {
		name    string
		opsFile *operationsFile
		want    []string
	}{
		{
			name:    "nothing listed",
			opsFile: &operationsFile{},
		},
		{
			// An override that decorates the URL is what the generated comments
			// need, because normalizeDocURL leaves a URL that has its own version.
			name: "override changes the documentation URL",
			opsFile: &operationsFile{
				OpenapiOps: generated,
				OverrideOps: []*operation{{
					Name:             "GET /a/{a_id}",
					DocumentationURL: "https://docs.github.com/rest/a/a?apiVersion=2026-03-10#get-a",
				}},
			},
		},
		{
			name: "override changes the description files",
			opsFile: &operationsFile{
				OpenapiOps: generated,
				OverrideOps: []*operation{{
					Name:         "GET /a/{a_id}",
					OpenAPIFiles: []string{"descriptions/ghec/ghec.json"},
				}},
			},
		},
		{
			// An override of a hand-written operation is not dead, even though the
			// generated section has no operation of that name.
			name: "override of a hand-written operation",
			opsFile: &operationsFile{
				ManualOps: []*operation{{Name: "GET /manual/{manual_id}"}},
				OverrideOps: []*operation{{
					Name:             "GET /manual/{manual_id}",
					DocumentationURL: "https://docs.github.com/rest/m/manual#get-manual",
				}},
			},
		},
		{
			name: "override of no operation",
			opsFile: &operationsFile{
				OpenapiOps: generated,
				OverrideOps: []*operation{{
					Name:             "DELETE /gone/{id}",
					DocumentationURL: "https://docs.github.com/rest/d/d#delete-d",
				}},
			},
			want: []string{
				`operation_overrides: "DELETE /gone/{id}" overrides no operation in openapi_operations or operations, so it is ignored`,
			},
		},
		{
			name: "override that repeats the generated URL",
			opsFile: &operationsFile{
				OpenapiOps: generated,
				OverrideOps: []*operation{{
					Name:             "GET /a/{a_id}",
					DocumentationURL: "https://docs.github.com/rest/a/a#get-a",
				}},
			},
			want: []string{
				`operation_overrides: "GET /a/{a_id}" sets what its operation already has, so it can be deleted`,
			},
		},
		{
			// The generated comments render both URLs the same way, so the override
			// does nothing even though the two strings differ.
			name: "override that repeats the URL as the comments render it",
			opsFile: &operationsFile{
				OpenapiOps: generated,
				OverrideOps: []*operation{{
					Name:             "GET /a/{a_id}",
					DocumentationURL: "https://docs.github.com/rest/a/a?apiVersion=" + metadataDocsAPIVersion + "#get-a",
				}},
			},
			want: []string{
				`operation_overrides: "GET /a/{a_id}" sets what its operation already has, so it can be deleted`,
			},
		},
		{
			// resolve() applies the documentation URL and the description files and
			// nothing else, so an override of any other field does nothing.
			name: "override that sets only a field that resolve ignores",
			opsFile: &operationsFile{
				OpenapiOps: generated,
				OverrideOps: []*operation{{
					Name:       "GET /a/{a_id}",
					Deprecated: true,
				}},
			},
			want: []string{
				`operation_overrides: "GET /a/{a_id}" sets what its operation already has, so it can be deleted`,
			},
		},
		{
			name: "name in the hand-written and generated sections",
			opsFile: &operationsFile{
				ManualOps:  []*operation{{Name: "GET /a/{a_id}"}},
				OpenapiOps: generated,
			},
			want: []string{
				`operations: "GET /a/{a_id}" is also in openapi_operations, so the hand-written entry replaces the generated one`,
			},
		},
		{
			name: "name listed twice in a hand-written section",
			opsFile: &operationsFile{
				ManualOps: []*operation{{Name: "GET /b/{b_id}"}, {Name: "GET /b/{b_id}"}},
				OpenapiOps: []*operation{{
					Name:             "GET /c/{c_id}",
					DocumentationURL: "https://docs.github.com/rest/c/c#get-c",
				}},
				OverrideOps: []*operation{
					{Name: "GET /c/{c_id}", DocumentationURL: "https://docs.github.com/rest/c/c#get-c-new"},
					{Name: "GET /c/{c_id}", DocumentationURL: "https://docs.github.com/rest/c/c#get-c-newer"},
				},
			},
			want: []string{
				`operation_overrides: "GET /c/{c_id}" is listed more than once, so the entries overlap and should be merged into one`,
				`operations: "GET /b/{b_id}" is listed more than once, so only the last entry is used`,
			},
		},
	} {
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			assertEqualStrings(t, joinProblems(td.want), joinProblems(checkProblems(td.opsFile)))
		})
	}
}

// TestCheckProblemsRepoFile checks the file this repository ships. Every other case is
// synthetic, and this is the one that fails when an entry that does nothing reaches the file,
// which is what keeps the file minimal rather than merely sorted.
func TestCheckProblemsRepoFile(t *testing.T) {
	t.Parallel()
	opsFile, err := loadOperationsFile(filepath.Join("..", "..", "openapi_operations.yaml"))
	assertNilError(t, err)
	assertEqualStrings(t, "", joinProblems(checkProblems(opsFile)))
}

// TestResolveAppliesOnlyDocURLAndFiles checks the assumption that checkProblems() relies on:
// an override patches the documentation URL and the description files, and no other field.
func TestResolveAppliesOnlyDocURLAndFiles(t *testing.T) {
	t.Parallel()
	generated := &operation{
		Name:             "GET /a/{a_id}",
		DocumentationURL: "https://docs.github.com/rest/a/a#get-a",
		OpenAPIFiles:     []string{"descriptions/api.github.com/api.github.com.json"},
	}
	opsFile := &operationsFile{
		OpenapiOps: []*operation{generated.clone()},
		OverrideOps: []*operation{{
			Name:             generated.Name,
			DocumentationURL: "https://docs.github.com/rest/a/a?apiVersion=2026-03-10#get-a",
			OpenAPIFiles:     []string{"descriptions/ghec/ghec.json"},
			Deprecated:       true,
		}},
	}
	opsFile.resolve()
	want := &operation{
		Name:             generated.Name,
		DocumentationURL: "https://docs.github.com/rest/a/a?apiVersion=2026-03-10#get-a",
		OpenAPIFiles:     []string{"descriptions/ghec/ghec.json"},
		Deprecated:       generated.Deprecated,
	}
	if diff := cmp.Diff(want, opsFile.resolvedOps[generated.Name]); diff != "" {
		t.Error(diff)
	}
	// The same override is not a problem, because it changes both fields that a
	// generated documentation link and the schema check depend on.
	assertEqualStrings(t, "", strings.Join(checkProblems(opsFile), "\n"))
}

// joinProblems renders problems as one comparable string, one per line.
func joinProblems(problems []string) string {
	return strings.Join(problems, "\n")
}
