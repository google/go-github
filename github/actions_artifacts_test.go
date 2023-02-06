// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestActionsService_ListArtifacts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/artifacts", func(w http.ResponseWriter, r *http.Request) {
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
	ctx := context.Background()
	artifacts, _, err := client.Actions.ListArtifacts(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListArtifacts returned error: %v", err)
	}

	want := &ArtifactList{TotalCount: Int64(1), Artifacts: []*Artifact{{ID: Int64(1)}}}
	if !cmp.Equal(artifacts, want) {
		t.Errorf("Actions.ListArtifacts returned %+v, want %+v", artifacts, want)
	}

	const methodName = "ListArtifacts"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListArtifacts(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListArtifacts(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_ListArtifacts_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Actions.ListArtifacts(ctx, "%", "r", nil)
	testURLParseError(t, err)
}

func TestActionsService_ListArtifacts_invalidRepo(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Actions.ListArtifacts(ctx, "o", "%", nil)
	testURLParseError(t, err)
}

func TestActionsService_ListArtifacts_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	artifacts, resp, err := client.Actions.ListArtifacts(ctx, "o", "r", nil)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Actions.ListArtifacts return status %d, want %d", got, want)
	}
	if artifacts != nil {
		t.Errorf("Actions.ListArtifacts return %+v, want nil", artifacts)
	}
}

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
	ctx := context.Background()
	artifacts, _, err := client.Actions.ListWorkflowRunArtifacts(ctx, "o", "r", 1, opts)
	if err != nil {
		t.Errorf("Actions.ListWorkflowRunArtifacts returned error: %v", err)
	}

	want := &ArtifactList{TotalCount: Int64(1), Artifacts: []*Artifact{{ID: Int64(1)}}}
	if !cmp.Equal(artifacts, want) {
		t.Errorf("Actions.ListWorkflowRunArtifacts returned %+v, want %+v", artifacts, want)
	}

	const methodName = "ListWorkflowRunArtifacts"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListWorkflowRunArtifacts(ctx, "\n", "\n", -1, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListWorkflowRunArtifacts(ctx, "o", "r", 1, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_ListWorkflowRunArtifacts_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Actions.ListWorkflowRunArtifacts(ctx, "%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestActionsService_ListWorkflowRunArtifacts_invalidRepo(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Actions.ListWorkflowRunArtifacts(ctx, "o", "%", 1, nil)
	testURLParseError(t, err)
}

func TestActionsService_ListWorkflowRunArtifacts_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runs/1/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	artifacts, resp, err := client.Actions.ListWorkflowRunArtifacts(ctx, "o", "r", 1, nil)
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

	ctx := context.Background()
	artifact, _, err := client.Actions.GetArtifact(ctx, "o", "r", 1)
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
	if !cmp.Equal(artifact, want) {
		t.Errorf("Actions.GetArtifact returned %+v, want %+v", artifact, want)
	}

	const methodName = "GetArtifact"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetArtifact(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetArtifact(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetArtifact_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Actions.GetArtifact(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

func TestActionsService_GetArtifact_invalidRepo(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Actions.GetArtifact(ctx, "o", "%", 1)
	testURLParseError(t, err)
}

func TestActionsService_GetArtifact_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/artifacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	artifact, resp, err := client.Actions.GetArtifact(ctx, "o", "r", 1)
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

	ctx := context.Background()
	url, resp, err := client.Actions.DownloadArtifact(ctx, "o", "r", 1, true)
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

	const methodName = "DownloadArtifact"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.DownloadArtifact(ctx, "\n", "\n", -1, true)
		return err
	})

	// Add custom round tripper
	client.client.Transport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("failed to download artifact")
	})
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.DownloadArtifact(ctx, "o", "r", 1, true)
		return err
	})
}

func TestActionsService_DownloadArtifact_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Actions.DownloadArtifact(ctx, "%", "r", 1, true)
	testURLParseError(t, err)
}

func TestActionsService_DownloadArtifact_invalidRepo(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Actions.DownloadArtifact(ctx, "o", "%", 1, true)
	testURLParseError(t, err)
}

func TestActionsService_DownloadArtifact_StatusMovedPermanently_dontFollowRedirects(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/artifacts/1/zip", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "https://github.com/artifact", http.StatusMovedPermanently)
	})

	ctx := context.Background()
	_, resp, _ := client.Actions.DownloadArtifact(ctx, "o", "r", 1, false)
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

	ctx := context.Background()
	url, resp, err := client.Actions.DownloadArtifact(ctx, "o", "r", 1, true)
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

	ctx := context.Background()
	_, err := client.Actions.DeleteArtifact(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Actions.DeleteArtifact return error: %v", err)
	}

	const methodName = "DeleteArtifact"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteArtifact(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteArtifact(ctx, "o", "r", 1)
	})
}

func TestActionsService_DeleteArtifact_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Actions.DeleteArtifact(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

func TestActionsService_DeleteArtifact_invalidRepo(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Actions.DeleteArtifact(ctx, "o", "%", 1)
	testURLParseError(t, err)
}

func TestActionsService_DeleteArtifact_notFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/artifacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	resp, err := client.Actions.DeleteArtifact(ctx, "o", "r", 1)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Actions.DeleteArtifact return status %d, want %d", got, want)
	}
}

func TestArtifact_Marshal(t *testing.T) {
	testJSONMarshal(t, &Artifact{}, "{}")

	u := &Artifact{
		ID:                 Int64(1),
		NodeID:             String("nid"),
		Name:               String("n"),
		SizeInBytes:        Int64(1),
		URL:                String("u"),
		ArchiveDownloadURL: String("a"),
		Expired:            Bool(false),
		CreatedAt:          &Timestamp{referenceTime},
		UpdatedAt:          &Timestamp{referenceTime},
		ExpiresAt:          &Timestamp{referenceTime},
		WorkflowRun: &ArtifactWorkflowRun{
			ID:               Int64(1),
			RepositoryID:     Int64(1),
			HeadRepositoryID: Int64(1),
			HeadBranch:       String("b"),
			HeadSHA:          String("s"),
		},
	}

	want := `{
		"id": 1,
		"node_id": "nid",
		"name": "n",
		"size_in_bytes": 1,
		"url": "u",
		"archive_download_url": "a",
		"expired": false,
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"expires_at": ` + referenceTimeStr + `,
		"workflow_run": {
			"id": 1,
			"repository_id": 1,
			"head_repository_id": 1,
			"head_branch": "b",
			"head_sha": "s"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestArtifactList_Marshal(t *testing.T) {
	testJSONMarshal(t, &ArtifactList{}, "{}")

	u := &ArtifactList{
		TotalCount: Int64(1),
		Artifacts: []*Artifact{
			{
				ID:                 Int64(1),
				NodeID:             String("nid"),
				Name:               String("n"),
				SizeInBytes:        Int64(1),
				URL:                String("u"),
				ArchiveDownloadURL: String("a"),
				Expired:            Bool(false),
				CreatedAt:          &Timestamp{referenceTime},
				UpdatedAt:          &Timestamp{referenceTime},
				ExpiresAt:          &Timestamp{referenceTime},
				WorkflowRun: &ArtifactWorkflowRun{
					ID:               Int64(1),
					RepositoryID:     Int64(1),
					HeadRepositoryID: Int64(1),
					HeadBranch:       String("b"),
					HeadSHA:          String("s"),
				},
			},
		},
	}

	want := `{
		"total_count": 1,
		"artifacts": [{
			"id": 1,
			"node_id": "nid",
			"name": "n",
			"size_in_bytes": 1,
			"url": "u",
			"archive_download_url": "a",
			"expired": false,
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"expires_at": ` + referenceTimeStr + `,
			"workflow_run": {
				"id": 1,
				"repository_id": 1,
				"head_repository_id": 1,
				"head_branch": "b",
				"head_sha": "s"
			}
		}]
	}`

	testJSONMarshal(t, u, want)
}
