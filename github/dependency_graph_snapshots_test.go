// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestDependencyGraphService_CreateSnapshot(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/dependency-graph/snapshots", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"version":0,"sha":"ce587453ced02b1526dfb4cb910479d431683101","ref":"refs/heads/main","job":{"correlator":"yourworkflowname_youractionname","id":"yourrunid","html_url":"https://example.com"},"detector":{"name":"octo-detector","version":"0.0.1","url":"https://github.com/octo-org/octo-repo"},"scanned":"2022-06-14T20:25:00Z","manifests":{"package-lock.json":{"name":"package-lock.json","file":{"source_location":"src/package-lock.json"},"resolved":{"@actions/core":{"package_url":"pkg:/npm/%40actions/core@1.1.9","relationship":"direct","scope":"runtime","dependencies":["@actions/http-client"]},"@actions/http-client":{"package_url":"pkg:/npm/%40actions/http-client@1.0.7","relationship":"indirect","scope":"runtime","dependencies":["tunnel"]},"tunnel":{"package_url":"pkg:/npm/tunnel@0.0.6","relationship":"indirect","scope":"runtime"}}}}}`+"\n")
		fmt.Fprint(w, `{"id":12345,"created_at":"2022-06-14T20:25:01Z","message":"Dependency results for the repo have been successfully updated.","result":"SUCCESS"}`)
	})

	ctx := context.Background()
	snapshot := &DependencyGraphSnapshot{
		Version: 0,
		Sha:     String("ce587453ced02b1526dfb4cb910479d431683101"),
		Ref:     String("refs/heads/main"),
		Job: &DependencyGraphSnapshotJob{
			Correlator: String("yourworkflowname_youractionname"),
			ID:         String("yourrunid"),
			HTMLURL:    String("https://example.com"),
		},
		Detector: &DependencyGraphSnapshotDetector{
			Name:    String("octo-detector"),
			Version: String("0.0.1"),
			URL:     String("https://github.com/octo-org/octo-repo"),
		},
		Scanned: &Timestamp{time.Date(2022, time.June, 14, 20, 25, 00, 0, time.UTC)},
		Manifests: map[string]*DependencyGraphSnapshotManifest{
			"package-lock.json": {
				Name: String("package-lock.json"),
				File: &DependencyGraphSnapshotManifestFile{SourceLocation: String("src/package-lock.json")},
				Resolved: map[string]*DependencyGraphSnapshotResolvedDependency{
					"@actions/core": {
						PackageURL:   String("pkg:/npm/%40actions/core@1.1.9"),
						Relationship: String("direct"),
						Scope:        String("runtime"),
						Dependencies: []string{"@actions/http-client"},
					},
					"@actions/http-client": {
						PackageURL:   String("pkg:/npm/%40actions/http-client@1.0.7"),
						Relationship: String("indirect"),
						Scope:        String("runtime"),
						Dependencies: []string{"tunnel"},
					},
					"tunnel": {
						PackageURL:   String("pkg:/npm/tunnel@0.0.6"),
						Relationship: String("indirect"),
						Scope:        String("runtime"),
					},
				},
			},
		},
	}

	snapshotCreationData, _, err := client.DependencyGraph.CreateSnapshot(ctx, "o", "r", snapshot)
	if err != nil {
		t.Errorf("DependencyGraph.CreateSnapshot returned error: %v", err)
	}

	want := &DependencyGraphSnapshotCreationData{
		ID:        12345,
		CreatedAt: &Timestamp{time.Date(2022, time.June, 14, 20, 25, 01, 0, time.UTC)},
		Message:   String("Dependency results for the repo have been successfully updated."),
		Result:    String("SUCCESS"),
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
