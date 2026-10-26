// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestIssueArtifactWithVersions_String(t *testing.T) {
	t.Parallel()
	v := IssueArtifactWithVersions{Artifact: &IssueArtifact{}}
	want := `github.IssueArtifactWithVersions{Artifact:github.IssueArtifact{ID:0, Number:0, Type:"", Name:"", Body:""}}`
	if got := v.String(); got != want {
		t.Errorf("IssueArtifactWithVersions.String = %v, want %v", got, want)
	}
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
	got, resp, err := client.Issues.ListArtifacts(t.Context(), "o", "r", 42, opts)
	if err != nil {
		t.Fatalf("Issues.ListArtifacts returned error: %v", err)
	}
	want := []*IssueArtifact{{ID: 9001, Number: 3, Type: "link", Name: "Runbook", Body: "https://example.com"}}
	if !cmp.Equal(got, want) {
		t.Errorf("Issues.ListArtifacts = %+v, want %+v", got, want)
	}
	if resp.StatusCode != http.StatusOK || resp.NextPage != 0 {
		t.Errorf("response status = %v, NextPage = %v; want 200, 0", resp.StatusCode, resp.NextPage)
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

func TestIssuesService_ListArtifacts_Pages(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name string
		body string
		opts *IssueArtifactListOptions
		want int
	}{
		{name: "empty", body: `[]`, want: 0},
		{name: "full", body: `[{"id":9001,"number":3}]`, opts: &IssueArtifactListOptions{ListOptions: ListOptions{PerPage: 1}}, want: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)
			requests := 0
			mux.HandleFunc("/repos/o/r/issues/42/artifacts", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				query := values{}
				if tt.opts != nil {
					query["per_page"] = "1"
				}
				testFormValues(t, r, query)
				requests++
				fmt.Fprint(w, tt.body)
			})
			got, resp, err := client.Issues.ListArtifacts(t.Context(), "o", "r", 42, tt.opts)
			if err != nil {
				t.Fatalf("Issues.ListArtifacts returned error: %v", err)
			}
			if got == nil || len(got) != tt.want || requests != 1 || resp.NextPage != 0 {
				t.Errorf("got %v, requests = %v, NextPage = %v; want non-nil slice of length %v, 1, 0", got, requests, resp.NextPage, tt.want)
			}
		})
	}
}

func TestIssuesService_GetArtifact(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/repos/o/r/issues/42/artifacts/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
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
			BodyHTML: new("<h1>Plan</h1>"), Description: new("Summary"),
			HTMLURL: new("https://github.com/o/r/issues/42/artifacts/3"),
			Creator: &User{ID: new(int64(12)), Login: new("author")}, UpdatedByActor: editor,
			CreatedAt: created, UpdatedAt: updated,
		},
		Versions: []*IssueArtifactVersion{
			{Version: 4, Name: "Plan", Body: "# Plan", BodyHTML: new("<h1>Plan</h1>"), CreatedAt: updated, Actor: editor},
			{Version: 1, Name: "Draft", Body: ""},
		},
	}
	got, resp, err := client.Issues.GetArtifact(t.Context(), "o", "r", 42, 3)
	if err != nil {
		t.Fatalf("Issues.GetArtifact returned error: %v", err)
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Issues.GetArtifact mismatch (-got +want):\n%v", cmp.Diff(got, want))
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("response status = %v, want 200", resp.StatusCode)
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

func TestIssuesService_GetArtifact_NullableFields(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/repos/o/r/issues/42/artifacts/3", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `{"artifact":{"id":9001,"number":3,"type":"generic","name":"Draft","body":"","body_html":null,"creator":null,"updated_by_actor":null,"created_at":null,"updated_at":null},"versions":[]}`)
	})
	got, _, err := client.Issues.GetArtifact(t.Context(), "o", "r", 42, 3)
	if err != nil {
		t.Fatalf("Issues.GetArtifact returned error: %v", err)
	}
	want := &IssueArtifactWithVersions{
		Artifact: &IssueArtifact{ID: 9001, Number: 3, Type: "generic", Name: "Draft", Body: ""},
		Versions: []*IssueArtifactVersion{},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Issues.GetArtifact mismatch (-got +want):\n%v", cmp.Diff(got, want))
	}
}

func TestIssuesService_CreateArtifact(t *testing.T) {
	t.Parallel()
	for _, artifactType := range []string{"generic", "plan", "link", "embedded_html"} {
		t.Run(artifactType, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)
			input := CreateIssueArtifactRequest{Type: new(artifactType), Name: "Draft", Body: ""}
			mux.HandleFunc("/repos/o/r/issues/42/artifacts", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				testFormValues(t, r, values{})
				testJSONBodyRaw(t, r, fmt.Sprintf(`{"type":%q,"name":"Draft","body":""}`, artifactType))
				w.WriteHeader(http.StatusCreated)
				fmt.Fprintf(w, `{"id":9001,"number":3,"type":%q,"name":"Draft","body":""}`, artifactType)
			})
			got, resp, err := client.Issues.CreateArtifact(t.Context(), "o", "r", 42, input)
			if err != nil {
				t.Fatalf("Issues.CreateArtifact returned error: %v", err)
			}
			want := &IssueArtifact{ID: 9001, Number: 3, Type: artifactType, Name: "Draft", Body: ""}
			if !cmp.Equal(got, want) {
				t.Errorf("Issues.CreateArtifact = %+v, want %+v", got, want)
			}
			if resp.StatusCode != http.StatusCreated {
				t.Errorf("response status = %v, want 201", resp.StatusCode)
			}
			const methodName = "CreateArtifact"
			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				got, resp, err := client.Issues.CreateArtifact(t.Context(), "o", "r", 42, input)
				if got != nil {
					t.Errorf("%v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestIssuesService_CreateArtifact_DefaultType(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/repos/o/r/issues/42/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testJSONBodyRaw(t, r, `{"name":"Draft","body":""}`)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":9001,"number":3,"type":"generic","name":"Draft","body":""}`)
	})
	got, _, err := client.Issues.CreateArtifact(t.Context(), "o", "r", 42, CreateIssueArtifactRequest{Name: "Draft"})
	if err != nil {
		t.Fatalf("Issues.CreateArtifact returned error: %v", err)
	}
	if got.Type != "generic" {
		t.Errorf("artifact type = %q, want generic", got.Type)
	}
}

func TestIssuesService_UpdateArtifact(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name  string
		input UpdateIssueArtifactRequest
		body  string
	}{
		{name: "name only", input: UpdateIssueArtifactRequest{Name: new("Renamed")}, body: `{"name":"Renamed"}`},
		{name: "body only", input: UpdateIssueArtifactRequest{Body: new("Revised")}, body: `{"body":"Revised"}`},
		{name: "both", input: UpdateIssueArtifactRequest{Name: new("Renamed"), Body: new("Revised")}, body: `{"name":"Renamed","body":"Revised"}`},
		{name: "empty name", input: UpdateIssueArtifactRequest{Name: new("")}, body: `{"name":""}`},
		{name: "empty body", input: UpdateIssueArtifactRequest{Body: new("")}, body: `{"body":""}`},
		{name: "both empty", input: UpdateIssueArtifactRequest{Name: new(""), Body: new("")}, body: `{"name":"","body":""}`},
		{name: "omitted", body: `{}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)
			mux.HandleFunc("/repos/o/r/issues/42/artifacts/3", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "PATCH")
				testFormValues(t, r, values{})
				testJSONBodyRaw(t, r, tt.body)
				fmt.Fprint(w, `{"id":9001,"number":3,"type":"plan","name":"Updated","body":""}`)
			})
			got, resp, err := client.Issues.UpdateArtifact(t.Context(), "o", "r", 42, 3, tt.input)
			if err != nil {
				t.Fatalf("Issues.UpdateArtifact returned error: %v", err)
			}
			want := &IssueArtifact{ID: 9001, Number: 3, Type: "plan", Name: "Updated", Body: ""}
			if !cmp.Equal(got, want) {
				t.Errorf("Issues.UpdateArtifact = %+v, want %+v", got, want)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("response status = %v, want 200", resp.StatusCode)
			}
			const methodName = "UpdateArtifact"
			testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
				got, resp, err := client.Issues.UpdateArtifact(t.Context(), "o", "r", 42, 3, tt.input)
				if got != nil {
					t.Errorf("%v = %#v, want nil", methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestIssuesService_DeleteArtifact(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	mux.HandleFunc("/repos/o/r/issues/42/artifacts/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, values{})
		testPlainBody(t, r, "")
		w.WriteHeader(http.StatusNoContent)
	})
	resp, err := client.Issues.DeleteArtifact(t.Context(), "o", "r", 42, 3)
	if err != nil {
		t.Fatalf("Issues.DeleteArtifact returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("response status = %v, want 204", resp.StatusCode)
	}
	testNewRequestAndDoFailure(t, "DeleteArtifact", client, func() (*Response, error) {
		return client.Issues.DeleteArtifact(t.Context(), "o", "r", 42, 3)
	})
}

func TestIssuesService_Artifacts_APIErrors(t *testing.T) {
	t.Parallel()
	for _, method := range []string{"GET", "POST", "PATCH", "DELETE", "LIST"} {
		t.Run(method, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)
			mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprint(w, `{"message":"Forbidden"}`)
			})
			var resp *Response
			var err error
			switch method {
			case "GET":
				var got *IssueArtifactWithVersions
				got, resp, err = client.Issues.GetArtifact(t.Context(), "o", "r", 42, 3)
				if got != nil {
					t.Errorf("GetArtifact = %+v, want nil", got)
				}
			case "POST":
				var got *IssueArtifact
				got, resp, err = client.Issues.CreateArtifact(t.Context(), "o", "r", 42, CreateIssueArtifactRequest{})
				if got != nil {
					t.Errorf("CreateArtifact = %+v, want nil", got)
				}
			case "PATCH":
				var got *IssueArtifact
				got, resp, err = client.Issues.UpdateArtifact(t.Context(), "o", "r", 42, 3, UpdateIssueArtifactRequest{})
				if got != nil {
					t.Errorf("UpdateArtifact = %+v, want nil", got)
				}
			case "DELETE":
				resp, err = client.Issues.DeleteArtifact(t.Context(), "o", "r", 42, 3)
			case "LIST":
				var got []*IssueArtifact
				got, resp, err = client.Issues.ListArtifacts(t.Context(), "o", "r", 42, nil)
				if got != nil {
					t.Errorf("ListArtifacts = %+v, want nil", got)
				}
			}
			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) || apiErr.Message != "Forbidden" {
				t.Errorf("error = %v, want ErrorResponse with message Forbidden", err)
			}
			if resp == nil || resp.StatusCode != http.StatusForbidden {
				t.Errorf("response = %+v, want status 403", resp)
			}
		})
	}
}
