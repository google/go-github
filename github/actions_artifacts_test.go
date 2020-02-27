// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestActionsService_ListWorkflowRunArtifacts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runs/1/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w,
			`{
				"total_count":1, 
				"artifacts":[{"id":1}]
			}`,
		)
	})

	opts := &ListOptions{Page: 2}
	artifacts, _, err := client.Actions.ListWorkflowRunArtifacts(context.Background(), "o", "r", 1, opts)
	if err != nil {
		t.Errorf("Actions.ListWorkflowRunArtifacts returned error: %v", err)
	}

	want := &ArtifactList{TotalCount: Int64(1), Artifacts: []*Artifact{{ID: Int64(1)}}}
	if !reflect.DeepEqual(artifacts, want) {
		t.Errorf("Actions.ListWorkflowRunArtifacts returned %+v, want %+v", artifacts, want)
	}
}

func TestActionsService_ListWorkflowRunArtifacts_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Actions.ListWorkflowRunArtifacts(context.Background(), "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestActionsService_ListWorkflowRunArtifacts_invalidRepo(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Actions.ListWorkflowRunArtifacts(context.Background(), "o", "%", 1, nil)
	testURLParseError(t, err)
}

func TestActionsService_ListWorkflowRunArtifacts_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runs/1/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	artifacts, resp, err := client.Actions.ListWorkflowRunArtifacts(context.Background(), "o", "r", 1, nil)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Actions.ListWorkflowRunArtifacts return status %d, want %d", got, want)
	}
	if artifacts != nil {
		t.Errorf("Actions.ListWorkflowRunArtifacts return %+v, want nil", artifacts)
	}
}

func TestActionsService_GetArtifact(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/artifacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id":1,
			"node_id":"xyz",
			"name":"a",
			"size_in_bytes":5,
			"archive_download_url":"u"
		}`)
	})

	artifact, _, err := client.Actions.GetArtifact(context.Background(), "o", "r", 1)
	if err != nil {
		t.Errorf("Actions.GetArtifact returned error: %v", err)
	}

	want := &Artifact{
		ID:                 Int64(1),
		NodeID:             String("xyz"),
		Name:               String("a"),
		SizeInBytes:        Int64(5),
		ArchiveDownloadURL: String("u"),
	}
	if !reflect.DeepEqual(artifact, want) {
		t.Errorf("Actions.GetArtifact returned %+v, want %+v", artifact, want)
	}
}

func TestActionsService_GetArtifact_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Actions.GetArtifact(context.Background(), "%", "r", 1)
	testURLParseError(t, err)
}

func TestActionsService_GetArtifact_invalidRepo(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Actions.GetArtifact(context.Background(), "o", "%", 1)
	testURLParseError(t, err)
}

func TestActionsService_GetArtifact_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/artifacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	artifact, resp, err := client.Actions.GetArtifact(context.Background(), "o", "r", 1)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Actions.GetArtifact return status %d, want %d", got, want)
	}
	if artifact != nil {
		t.Errorf("Actions.GetArtifact return %+v, want nil", artifact)
	}
}

func TestActionsSerivice_DownloadArtifact(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/artifacts/1/zip", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "https://github.com/artifact", http.StatusFound)
	})

	url, resp, err := client.Actions.DownloadArtifact(context.Background(), "o", "r", 1, true)
	if err != nil {
		t.Errorf("Actions.DownloadArtifact returned error: %v", err)
	}
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Actions.DownloadArtifact returned status: %d, want %d", resp.StatusCode, http.StatusFound)
	}

	want := "https://github.com/artifact"
	if url.String() != want {
		t.Errorf("Actions.DownloadArtifact returned %+v, want %+v", url.String(), want)
	}
}

func TestActionsService_DownloadArtifact_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Actions.DownloadArtifact(context.Background(), "%", "r", 1, true)
	testURLParseError(t, err)
}

func TestActionsService_DownloadArtifact_invalidRepo(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Actions.DownloadArtifact(context.Background(), "o", "%", 1, true)
	testURLParseError(t, err)
}

func TestActionsService_DownloadArtifact_StatusMovedPermanently_dontFollowRedirects(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/artifacts/1/zip", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "https://github.com/artifact", http.StatusMovedPermanently)
	})

	_, resp, _ := client.Actions.DownloadArtifact(context.Background(), "o", "r", 1, false)
	if resp.StatusCode != http.StatusMovedPermanently {
		t.Errorf("Actions.DownloadArtifact return status %d, want %d", resp.StatusCode, http.StatusMovedPermanently)
	}
}

func TestActionsService_DownloadArtifact_StatusMovedPermanently_followRedirects(t *testing.T) {
	client, mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/artifacts/1/zip", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		redirectURL, _ := url.Parse(serverURL + baseURLPath + "/redirect")
		http.Redirect(w, r, redirectURL.String(), http.StatusMovedPermanently)
	})
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/artifact", http.StatusFound)
	})

	url, resp, err := client.Actions.DownloadArtifact(context.Background(), "o", "r", 1, true)
	if err != nil {
		t.Errorf("Actions.DownloadArtifact return error: %v", err)
	}
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Actions.DownloadArtifact return status %d, want %d", resp.StatusCode, http.StatusFound)
	}
	want := "http://github.com/artifact"
	if url.String() != want {
		t.Errorf("Actions.DownloadArtifact returned %+v, want %+v", url.String(), want)
	}
}

func TestActionsService_DeleteArtifact(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/artifacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Actions.DeleteArtifact(context.Background(), "o", "r", 1)
	if err != nil {
		t.Errorf("Actions.DeleteArtifact return error: %v", err)
	}
}

func TestActionsService_DeleteArtifact_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Actions.DeleteArtifact(context.Background(), "%", "r", 1)
	testURLParseError(t, err)
}

func TestActionsService_DeleteArtifact_invalidRepo(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Actions.DeleteArtifact(context.Background(), "o", "%", 1)
	testURLParseError(t, err)
}

func TestActionsService_DeleteArtifact_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/artifacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNotFound)
	})

	resp, err := client.Actions.DeleteArtifact(context.Background(), "o", "r", 1)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Actions.DeleteArtifact return status %d, want %d", got, want)
	}
}
