// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestIssueImportService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	createdAt := time.Date(2020, time.August, 11, 15, 30, 0, 0, time.UTC)
	input := &IssueImportRequest{
		IssueImport: IssueImport{
			Assignee:  String("developer"),
			Body:      "Dummy description",
			CreatedAt: &createdAt,
			Labels:    []string{"l1", "l2"},
			Milestone: Int(1),
			Title:     "Dummy Issue",
		},
		Comments: []*Comment{{
			CreatedAt: &createdAt,
			Body:      "Comment body",
		}},
	}

	mux.HandleFunc("/repos/o/r/import/issues", func(w http.ResponseWriter, r *http.Request) {
		v := new(IssueImportRequest)
		json.NewDecoder(r.Body).Decode(v)
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeIssueImportAPI)
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write(issueImportResponseJSON)
	})

	ctx := context.Background()
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

func TestIssueImportService_Create_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.IssueImport.Create(ctx, "%", "r", nil)
	testURLParseError(t, err)
}

func TestIssueImportService_CheckStatus(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/import/issues/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeIssueImportAPI)
		w.WriteHeader(http.StatusOK)
		w.Write(issueImportResponseJSON)
	})

	ctx := context.Background()
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
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.IssueImport.CheckStatus(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

func TestIssueImportService_CheckStatusSince(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/import/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeIssueImportAPI)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("[%s]", issueImportResponseJSON)))
	})

	ctx := context.Background()
	got, _, err := client.IssueImport.CheckStatusSince(ctx, "o", "r", time.Now())
	if err != nil {
		t.Errorf("CheckStatusSince returned error: %v", err)
	}

	want := []*IssueImportResponse{wantIssueImportResponse}
	if !cmp.Equal(want, got) {
		t.Errorf("CheckStatusSince = %v, want = %v", got, want)
	}

	const methodName = "CheckStatusSince"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.IssueImport.CheckStatusSince(ctx, "\n", "\n", time.Now())
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.IssueImport.CheckStatusSince(ctx, "o", "r", time.Now())
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssueImportService_CheckStatusSince_badResponse(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/import/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeIssueImportAPI)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{badly-formed JSON"))
	})

	ctx := context.Background()
	if _, _, err := client.IssueImport.CheckStatusSince(ctx, "o", "r", time.Now()); err == nil {
		t.Errorf("CheckStatusSince returned no error, want JSON err")
	}
}

func TestIssueImportService_CheckStatusSince_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.IssueImport.CheckStatusSince(ctx, "%", "r", time.Now())
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
	ID:              Int(3),
	Status:          String("pending"),
	URL:             String("https://api.github.com/repos/o/r/import/issues/3"),
	ImportIssuesURL: String("https://api.github.com/repos/o/r/import/issues"),
	RepositoryURL:   String("https://api.github.com/repos/o/r"),
}

func TestIssueImportError_Marshal(t *testing.T) {
	testJSONMarshal(t, &IssueImportError{}, "{}")

	u := &IssueImportError{
		Location: String("loc"),
		Resource: String("res"),
		Field:    String("field"),
		Value:    String("value"),
		Code:     String("code"),
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
	testJSONMarshal(t, &IssueImportResponse{}, "{}")

	u := &IssueImportResponse{
		ID:               Int(1),
		Status:           String("status"),
		URL:              String("url"),
		ImportIssuesURL:  String("iiu"),
		RepositoryURL:    String("ru"),
		CreatedAt:        &referenceTime,
		UpdatedAt:        &referenceTime,
		Message:          String("msg"),
		DocumentationURL: String("durl"),
		Errors: []*IssueImportError{
			{
				Location: String("loc"),
				Resource: String("res"),
				Field:    String("field"),
				Value:    String("value"),
				Code:     String("code"),
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
	testJSONMarshal(t, &Comment{}, "{}")

	u := &Comment{
		CreatedAt: &referenceTime,
		Body:      "body",
	}

	want := `{
		"created_at": ` + referenceTimeStr + `,
		"body": "body"
	}`

	testJSONMarshal(t, u, want)
}

func TestIssueImport_Marshal(t *testing.T) {
	testJSONMarshal(t, &IssueImport{}, "{}")

	u := &IssueImport{
		Title:     "title",
		Body:      "body",
		CreatedAt: &referenceTime,
		ClosedAt:  &referenceTime,
		UpdatedAt: &referenceTime,
		Assignee:  String("a"),
		Milestone: Int(1),
		Closed:    Bool(false),
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
	testJSONMarshal(t, &IssueImportRequest{}, "{}")

	u := &IssueImportRequest{
		IssueImport: IssueImport{
			Title:     "title",
			Body:      "body",
			CreatedAt: &referenceTime,
			ClosedAt:  &referenceTime,
			UpdatedAt: &referenceTime,
			Assignee:  String("a"),
			Milestone: Int(1),
			Closed:    Bool(false),
			Labels:    []string{"l"},
		},
		Comments: []*Comment{
			{
				CreatedAt: &referenceTime,
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
