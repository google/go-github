// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestIssueArtifactWithVersions_String(t *testing.T) {
	t.Parallel()
	v := IssueArtifactWithVersions{Artifact: &IssueArtifact{}}
	want := `github.IssueArtifactWithVersions{Artifact:github.IssueArtifact{ID:0, Number:0, Type:"", Name:"", Body:"", Description:""}}`
	if got := v.String(); got != want {
		t.Errorf("IssueArtifactWithVersions.String = %v, want %v", got, want)
	}
}

func TestIssueArtifact_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &IssueArtifact{}, `{"id":0,"number":0,"type":"","name":"","body":"","body_html":null,"description":"","creator":null,"updated_by_actor":null,"created_at":null,"updated_at":null}`)
}

func TestIssueArtifactVersion_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &IssueArtifactVersion{}, `{"version":0,"name":"","body":"","body_html":null,"created_at":null,"actor":null}`)
}

func TestIssuesService_ListArtifacts(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/issues/42/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2", "per_page": "50", "type": "link"})
		fmt.Fprint(w, `[{"id":9001,"number":3,"type":"link","name":"Runbook","body":"https://example.com"}]`)
	})

	opts := &IssueArtifactListOptions{Type: "link", ListOptions: ListOptions{Page: 2, PerPage: 50}}
	got, _, err := client.Issues.ListArtifacts(t.Context(), "o", "r", 42, opts)
	if err != nil {
		t.Fatalf("Issues.ListArtifacts returned error: %v", err)
	}
	want := []*IssueArtifact{{ID: 9001, Number: 3, Type: "link", Name: "Runbook", Body: "https://example.com"}}
	if !cmp.Equal(got, want) {
		t.Errorf("Issues.ListArtifacts = %+v, want %+v", got, want)
	}

	const methodName = "ListArtifacts"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.Issues.ListArtifacts(t.Context(), "\n", "\n", 42, opts)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.ListArtifacts(t.Context(), "o", "r", 42, opts)
		if got != nil {
			t.Errorf("%v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_ListArtifactsIter_Options(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	requests := 0
	mux.HandleFunc("/repos/o/r/issues/42/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		requests++
		testFormValues(t, r, values{"page": strconv.Itoa(requests), "per_page": "1", "type": "link"})
		if requests == 1 {
			w.Header().Set("Link", `<https://api.github.com/repos/o/r/issues/42/artifacts?page=2>; rel="next"`)
		}
		fmt.Fprint(w, `[{"type":"link"}]`)
	})

	opts := &IssueArtifactListOptions{Type: "link", ListOptions: ListOptions{Page: 1, PerPage: 1}}
	var count int
	for _, err := range client.Issues.ListArtifactsIter(t.Context(), "o", "r", 42, opts) {
		if err != nil {
			t.Fatalf("ListArtifactsIter returned error: %v", err)
		}
		count++
	}
	if count != 2 || requests != 2 {
		t.Errorf("ListArtifactsIter returned %v artifacts in %v requests; want 2 in 2 requests", count, requests)
	}
	want := &IssueArtifactListOptions{Type: "link", ListOptions: ListOptions{Page: 1, PerPage: 1}}
	if !cmp.Equal(opts, want) {
		t.Errorf("ListArtifactsIter mutated options: %+v, want %+v", opts, want)
	}
}

func TestIssuesService_GetArtifact(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/repos/o/r/issues/42/artifacts/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"artifact":{
				"id":9001,"number":3,"type":"plan","name":"Plan","body":"# Plan",
				"body_html":"<h1>Plan</h1>","description":"Summary",
				"html_url":"https://github.com/o/r/issues/42/artifacts/3",
				"creator":{"id":12,"login":"author"},
				"updated_by_actor":{"id":13,"login":"editor"},
				"created_at":"2026-09-01T10:00:00Z","updated_at":"2026-09-02T10:00:00Z"
			},
			"versions":[
				{"version":4,"name":"Plan","body":"# Plan","body_html":"<h1>Plan</h1>","created_at":"2026-09-02T10:00:00Z","actor":{"id":13,"login":"editor"}},
				{"version":1,"name":"Draft","body":"","body_html":null,"created_at":null,"actor":null}
			]
		}`)
	})
	created := &Timestamp{time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)}
	updated := &Timestamp{time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)}
	editor := &User{ID: new(int64(13)), Login: new("editor")}
	want := &IssueArtifactWithVersions{
		Artifact: &IssueArtifact{
			ID: 9001, Number: 3, Type: "plan", Name: "Plan", Body: "# Plan",
			BodyHTML: new("<h1>Plan</h1>"), Description: "Summary",
			HTMLURL: new("https://github.com/o/r/issues/42/artifacts/3"),
			Creator: &User{ID: new(int64(12)), Login: new("author")}, UpdatedByActor: editor,
			CreatedAt: created, UpdatedAt: updated,
		},
		Versions: []*IssueArtifactVersion{
			{Version: 4, Name: "Plan", Body: "# Plan", BodyHTML: new("<h1>Plan</h1>"), CreatedAt: updated, Actor: editor},
			{Version: 1, Name: "Draft", Body: ""},
		},
	}
	got, _, err := client.Issues.GetArtifact(t.Context(), "o", "r", 42, 3)
	if err != nil {
		t.Fatalf("Issues.GetArtifact returned error: %v", err)
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Issues.GetArtifact mismatch (-got +want):\n%v", cmp.Diff(got, want))
	}
	const methodName = "GetArtifact"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.GetArtifact(t.Context(), "o", "r", 42, 3)
		if got != nil {
			t.Errorf("%v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCreateIssueArtifactRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, CreateIssueArtifactRequest{Name: "Draft"}, `{"name":"Draft","body":""}`)
}

func TestIssuesService_CreateArtifact(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	input := CreateIssueArtifactRequest{Type: new("plan"), Name: "Draft", Body: ""}
	mux.HandleFunc("/repos/o/r/issues/42/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONBody(t, r, input)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":9001,"number":3,"type":"plan","name":"Draft","body":""}`)
	})
	got, _, err := client.Issues.CreateArtifact(t.Context(), "o", "r", 42, input)
	if err != nil {
		t.Fatalf("Issues.CreateArtifact returned error: %v", err)
	}
	want := &IssueArtifact{ID: 9001, Number: 3, Type: "plan", Name: "Draft", Body: ""}
	if !cmp.Equal(got, want) {
		t.Errorf("Issues.CreateArtifact = %+v, want %+v", got, want)
	}
	const methodName = "CreateArtifact"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.CreateArtifact(t.Context(), "o", "r", 42, input)
		if got != nil {
			t.Errorf("%v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUpdateIssueArtifactRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, UpdateIssueArtifactRequest{}, `{}`)
	testJSONMarshal(t, UpdateIssueArtifactRequest{Name: new("Renamed")}, `{"name":"Renamed"}`)
	testJSONMarshal(t, UpdateIssueArtifactRequest{Body: new("Revised")}, `{"body":"Revised"}`)
	testJSONMarshal(t, UpdateIssueArtifactRequest{Name: new("")}, `{"name":""}`)
	testJSONMarshal(t, UpdateIssueArtifactRequest{Body: new("")}, `{"body":""}`)
	testJSONMarshal(t, UpdateIssueArtifactRequest{Name: new(""), Body: new("")}, `{"name":"","body":""}`)
}

func TestIssuesService_UpdateArtifact(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	input := UpdateIssueArtifactRequest{Name: new("Renamed"), Body: new("Revised")}
	mux.HandleFunc("/repos/o/r/issues/42/artifacts/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{"id":9001,"number":3,"type":"plan","name":"Renamed","body":"Revised"}`)
	})
	got, _, err := client.Issues.UpdateArtifact(t.Context(), "o", "r", 42, 3, input)
	if err != nil {
		t.Fatalf("Issues.UpdateArtifact returned error: %v", err)
	}
	want := &IssueArtifact{ID: 9001, Number: 3, Type: "plan", Name: "Renamed", Body: "Revised"}
	if !cmp.Equal(got, want) {
		t.Errorf("Issues.UpdateArtifact = %+v, want %+v", got, want)
	}
	const methodName = "UpdateArtifact"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.UpdateArtifact(t.Context(), "o", "r", 42, 3, input)
		if got != nil {
			t.Errorf("%v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_DeleteArtifact(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/repos/o/r/issues/42/artifacts/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	_, err := client.Issues.DeleteArtifact(t.Context(), "o", "r", 42, 3)
	if err != nil {
		t.Fatalf("Issues.DeleteArtifact returned error: %v", err)
	}
	testNewRequestAndDoFailure(t, "DeleteArtifact", client, func() (*Response, error) {
		return client.Issues.DeleteArtifact(t.Context(), "o", "r", 42, 3)
	})
}
