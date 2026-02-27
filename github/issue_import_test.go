// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestIssueImportService_Create(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	createdAt := time.Date(2020, time.August, 11, 15, 30, 0, 0, time.UTC)
	input := &IssueImportRequest{
		IssueImport: IssueImport{
			Assignee:  Ptr("developer"),
			Body:      "Dummy description",
			CreatedAt: &Timestamp{createdAt},
			Labels:    []string{"l1", "l2"},
			Milestone: Ptr(1),
			Title:     "Dummy Issue",
		},
		Comments: []*Comment{{
			CreatedAt: &Timestamp{createdAt},
			Body:      "Comment body",
		}},
	}

	mux.HandleFunc("/repos/o/r/import/issues", func(w http.ResponseWriter, r *http.Request) {
		v := new(IssueImportRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeIssueImportAPI)
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		assertWrite(t, w, issueImportResponseJSON)
	})

	ctx := t.Context()
	got, _, err := client.IssueImport.Create(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Create returned error: %v", err)
	}

	want := wantIssueImportResponse
	if !cmp.Equal(got, want) {
		t.Errorf("Create = %+v, want %+v", got, want)
	}

	const methodName = "Create"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.IssueImport.Create(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.IssueImport.Create(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssueImportService_Create_deferred(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	createdAt := time.Date(2020, time.August, 11, 15, 30, 0, 0, time.UTC)
	input := &IssueImportRequest{
		IssueImport: IssueImport{
			Assignee:  Ptr("developer"),
			Body:      "Dummy description",
			CreatedAt: &Timestamp{createdAt},
			Labels:    []string{"l1", "l2"},
			Milestone: Ptr(1),
			Title:     "Dummy Issue",
		},
		Comments: []*Comment{{
			CreatedAt: &Timestamp{createdAt},
			Body:      "Comment body",
		}},
	}

	mux.HandleFunc("/repos/o/r/import/issues", func(w http.ResponseWriter, r *http.Request) {
		v := new(IssueImportRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeIssueImportAPI)
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusAccepted)
		assertWrite(t, w, issueImportResponseJSON)
	})

	ctx := t.Context()
	got, _, err := client.IssueImport.Create(ctx, "o", "r", input)

	if !errors.As(err, new(*AcceptedError)) {
		t.Errorf("Create returned error: %v (want AcceptedError)", err)
	}

	want := wantIssueImportResponse
	if !cmp.Equal(got, want) {
		t.Errorf("Create = %+v, want %+v", got, want)
	}
}

func TestIssueImportService_Create_badResponse(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	createdAt := time.Date(2020, time.August, 11, 15, 30, 0, 0, time.UTC)
	input := &IssueImportRequest{
		IssueImport: IssueImport{
			Assignee:  Ptr("developer"),
			Body:      "Dummy description",
			CreatedAt: &Timestamp{createdAt},
			Labels:    []string{"l1", "l2"},
			Milestone: Ptr(1),
			Title:     "Dummy Issue",
		},
		Comments: []*Comment{{
			CreatedAt: &Timestamp{createdAt},
			Body:      "Comment body",
		}},
	}

	mux.HandleFunc("/repos/o/r/import/issues", func(w http.ResponseWriter, r *http.Request) {
		v := new(IssueImportRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeIssueImportAPI)
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusAccepted)
		assertWrite(t, w, []byte("{[}"))
	})

	ctx := t.Context()
	_, _, err := client.IssueImport.Create(ctx, "o", "r", input)

	if err == nil || err.Error() != "invalid character '[' looking for beginning of object key string" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestIssueImportService_Create_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.IssueImport.Create(ctx, "%", "r", nil)
	testURLParseError(t, err)
}

func TestIssueImportService_CheckStatus(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/import/issues/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeIssueImportAPI)
		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, issueImportResponseJSON)
	})

	ctx := t.Context()
	got, _, err := client.IssueImport.CheckStatus(ctx, "o", "r", 3)
	if err != nil {
		t.Errorf("CheckStatus returned error: %v", err)
	}

	want := wantIssueImportResponse
	if !cmp.Equal(got, want) {
		t.Errorf("CheckStatus = %+v, want %+v", got, want)
	}

	const methodName = "CheckStatus"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.IssueImport.CheckStatus(ctx, "\n", "\n", -3)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.IssueImport.CheckStatus(ctx, "o", "r", 3)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssueImportService_CheckStatus_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.IssueImport.CheckStatus(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

func TestIssueImportService_CheckStatusSince(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/import/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeIssueImportAPI)
		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, fmt.Appendf(nil, "[%s]", issueImportResponseJSON))
	})

	ctx := t.Context()
	got, _, err := client.IssueImport.CheckStatusSince(ctx, "o", "r", Timestamp{time.Now()})
	if err != nil {
		t.Errorf("CheckStatusSince returned error: %v", err)
	}

	want := []*IssueImportResponse{wantIssueImportResponse}
	if !cmp.Equal(want, got) {
		t.Errorf("CheckStatusSince = %v, want = %v", got, want)
	}

	const methodName = "CheckStatusSince"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.IssueImport.CheckStatusSince(ctx, "\n", "\n", Timestamp{time.Now()})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.IssueImport.CheckStatusSince(ctx, "o", "r", Timestamp{time.Now()})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssueImportService_CheckStatusSince_badResponse(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/import/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeIssueImportAPI)
		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte("{badly-formed JSON"))
	})

	ctx := t.Context()
	if _, _, err := client.IssueImport.CheckStatusSince(ctx, "o", "r", Timestamp{time.Now()}); err == nil {
		t.Error("CheckStatusSince returned no error, want JSON err")
	}
}

func TestIssueImportService_CheckStatusSince_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.IssueImport.CheckStatusSince(ctx, "%", "r", Timestamp{time.Now()})
	testURLParseError(t, err)
}

var issueImportResponseJSON = []byte(`{
	"id": 3,
	"status": "pending",
	"url": "https://api.github.com/repos/o/r/import/issues/3",
	"import_issues_url": "https://api.github.com/repos/o/r/import/issues",
	"repository_url": "https://api.github.com/repos/o/r"
}`)

var wantIssueImportResponse = &IssueImportResponse{
	ID:              Ptr(3),
	Status:          Ptr("pending"),
	URL:             Ptr("https://api.github.com/repos/o/r/import/issues/3"),
	ImportIssuesURL: Ptr("https://api.github.com/repos/o/r/import/issues"),
	RepositoryURL:   Ptr("https://api.github.com/repos/o/r"),
}

func TestIssueImportError_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &IssueImportError{}, "{}")

	u := &IssueImportError{
		Location: Ptr("loc"),
		Resource: Ptr("res"),
		Field:    Ptr("field"),
		Value:    Ptr("value"),
		Code:     Ptr("code"),
	}

	want := `{
		"location": "loc",
		"resource": "res",
		"field": "field",
		"value": "value",
		"code": "code"
	}`

	testJSONMarshal(t, u, want)
}

func TestIssueImportResponse_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &IssueImportResponse{}, "{}")

	u := &IssueImportResponse{
		ID:               Ptr(1),
		Status:           Ptr("status"),
		URL:              Ptr("url"),
		ImportIssuesURL:  Ptr("iiu"),
		RepositoryURL:    Ptr("ru"),
		CreatedAt:        &Timestamp{referenceTime},
		UpdatedAt:        &Timestamp{referenceTime},
		Message:          Ptr("msg"),
		DocumentationURL: Ptr("durl"),
		Errors: []*IssueImportError{
			{
				Location: Ptr("loc"),
				Resource: Ptr("res"),
				Field:    Ptr("field"),
				Value:    Ptr("value"),
				Code:     Ptr("code"),
			},
		},
	}

	want := `{
		"id": 1,
		"status": "status",
		"url": "url",
		"import_issues_url": "iiu",
		"repository_url": "ru",
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"message": "msg",
		"documentation_url": "durl",
		"errors": [
			{
				"location": "loc",
				"resource": "res",
				"field": "field",
				"value": "value",
				"code": "code"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestComment_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Comment{}, `{"body": ""}`)

	u := &Comment{
		CreatedAt: &Timestamp{referenceTime},
		Body:      "body",
	}

	want := `{
		"created_at": ` + referenceTimeStr + `,
		"body": "body"
	}`

	testJSONMarshal(t, u, want)
}

func TestIssueImport_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &IssueImport{}, `{
		"title": "",
		"body": ""
	}`)

	u := &IssueImport{
		Title:     "title",
		Body:      "body",
		CreatedAt: &Timestamp{referenceTime},
		ClosedAt:  &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
		Assignee:  Ptr("a"),
		Milestone: Ptr(1),
		Closed:    Ptr(false),
		Labels:    []string{"l"},
	}

	want := `{
		"title": "title",
		"body": "body",
		"created_at": ` + referenceTimeStr + `,
		"closed_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"assignee": "a",
		"milestone": 1,
		"closed": false,
		"labels": [
			"l"
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestIssueImportRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &IssueImportRequest{}, `{
		"issue": {
			"title": "",
			"body": ""
		}
	}`)

	u := &IssueImportRequest{
		IssueImport: IssueImport{
			Title:     "title",
			Body:      "body",
			CreatedAt: &Timestamp{referenceTime},
			ClosedAt:  &Timestamp{referenceTime},
			UpdatedAt: &Timestamp{referenceTime},
			Assignee:  Ptr("a"),
			Milestone: Ptr(1),
			Closed:    Ptr(false),
			Labels:    []string{"l"},
		},
		Comments: []*Comment{
			{
				CreatedAt: &Timestamp{referenceTime},
				Body:      "body",
			},
		},
	}

	want := `{
		"issue": {
			"title": "title",
			"body": "body",
			"created_at": ` + referenceTimeStr + `,
			"closed_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"assignee": "a",
			"milestone": 1,
			"closed": false,
			"labels": [
				"l"
			]
		},
		"comments": [
			{
				"created_at": ` + referenceTimeStr + `,
				"body": "body"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}
