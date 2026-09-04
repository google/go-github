// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDependencyGraphService_CreateSnapshot(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	snapshot := &DependencyGraphSnapshot{
		Version: 0,
		SHA:     new("ce587453ced02b1526dfb4cb910479d431683101"),
		Ref:     new("refs/heads/main"),
		Job: &DependencyGraphSnapshotJob{
			Correlator: new("yourworkflowname_youractionname"),
			ID:         new("yourrunid"),
			HTMLURL:    new("https://example.com"),
		},
		Detector: &DependencyGraphSnapshotDetector{
			Name:    new("octo-detector"),
			Version: new("0.0.1"),
			URL:     new("https://github.com/octo-org/octo-repo"),
		},
		Scanned: &referenceTimestamp,
		Metadata: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
		Manifests: map[string]*DependencyGraphSnapshotManifest{
			"package-lock.json": {
				Name: new("package-lock.json"),
				File: &DependencyGraphSnapshotManifestFile{SourceLocation: new("src/package-lock.json")},
				Metadata: map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
				Resolved: map[string]*DependencyGraphSnapshotResolvedDependency{
					"@actions/core": {
						PackageURL:   new("pkg:/npm/%40actions/core@1.1.9"),
						Relationship: new("direct"),
						Scope:        new("runtime"),
						Metadata: map[string]any{
							"licenses": "MIT",
						},
						Dependencies: []string{"@actions/http-client"},
					},
					"@actions/http-client": {
						PackageURL:   new("pkg:/npm/%40actions/http-client@1.0.7"),
						Relationship: new("indirect"),
						Scope:        new("runtime"),
						Dependencies: []string{"tunnel"},
					},
					"tunnel": {
						PackageURL:   new("pkg:/npm/tunnel@0.0.6"),
						Relationship: new("indirect"),
						Scope:        new("runtime"),
					},
				},
			},
		},
	}

	mux.HandleFunc("/repos/o/r/dependency-graph/snapshots", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONBody(t, r, snapshot)
		fmt.Fprint(w, `{"id":12345,"created_at":`+referenceTimeStr+`,"message":"Dependency results for the repo have been successfully updated.","result":"SUCCESS"}`)
	})

	ctx := t.Context()
	snapshotCreationData, _, err := client.DependencyGraph.CreateSnapshot(ctx, "o", "r", snapshot)
	if err != nil {
		t.Errorf("DependencyGraph.CreateSnapshot returned error: %v", err)
	}

	want := &DependencyGraphSnapshotCreationData{
		ID:        12345,
		CreatedAt: &referenceTimestamp,
		Message:   new("Dependency results for the repo have been successfully updated."),
		Result:    new("SUCCESS"),
	}
	if !cmp.Equal(snapshotCreationData, want) {
		t.Errorf("DependencyGraph.CreateSnapshot returned %+v, want %+v", snapshotCreationData, want)
	}

	const methodName = "CreateSnapshot"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.DependencyGraph.CreateSnapshot(ctx, "o", "r", snapshot)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
