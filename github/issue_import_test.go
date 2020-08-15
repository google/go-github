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
	"reflect"
	"testing"
	"time"
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
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write(issueImportResponseJSON)
	})

	got, _, err := client.IssueImport.Create(context.Background(), "o", "r", input)
	if err != nil {
		t.Errorf("Create returned error: %v", err)
	}

	want := wantIssueImportResponse
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Create = %+v, want %+v", got, want)
	}
}

func TestIssueImportService_Create_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.IssueImport.Create(context.Background(), "%", "r", nil)
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

	got, _, err := client.IssueImport.CheckStatus(context.Background(), "o", "r", 3)
	if err != nil {
		t.Errorf("CheckStatus returned error: %v", err)
	}

	want := wantIssueImportResponse
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CheckStatus = %+v, want %+v", got, want)
	}
}

func TestIssueImportService_CheckStatus_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.IssueImport.CheckStatus(context.Background(), "%", "r", 1)
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

	got, _, err := client.IssueImport.CheckStatusSince(context.Background(), "o", "r", time.Now())
	if err != nil {
		t.Errorf("CheckStatusSince returned error: %v", err)
	}

	want := []*IssueImportResponse{wantIssueImportResponse}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("CheckStatusSince = %v, want = %v", got, want)
	}
}

func TestIssueImportService_CheckStatusSince_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.IssueImport.CheckStatusSince(context.Background(), "%", "r", time.Now())
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
