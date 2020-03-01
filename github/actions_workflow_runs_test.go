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
	"time"
)

func TestActionsService_ListWorkflowRunsByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/workflows/29679449/runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"workflow_runs":[{"id":399444496,"run_number":296,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"id":399444497,"run_number":296,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListWorkflowRunsOptions{ListOptions: ListOptions{Page: 2, PerPage: 2}}
	runs, _, err := client.Actions.ListWorkflowRunsByID(context.Background(), "o", "r", 29679449, opts)
	if err != nil {
		t.Errorf("Actions.ListWorkFlowRunsByID returned error: %v", err)
	}

	want := &WorkflowRuns{
		TotalCount: Int(4),
		WorkflowRuns: []*WorkflowRun{
			{ID: Int64(399444496), RunNumber: Int(296), CreatedAt: &Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
			{ID: Int64(399444497), RunNumber: Int(296), CreatedAt: &Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
		},
	}
	if !reflect.DeepEqual(runs, want) {
		t.Errorf("Actions.ListWorkflowRunsByID returned %+v, want %+v", runs, want)
	}
}

func TestActionsService_ListWorkflowRunsFileName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/workflows/29679449/runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"workflow_runs":[{"id":399444496,"run_number":296,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"id":399444497,"run_number":296,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListWorkflowRunsOptions{ListOptions: ListOptions{Page: 2, PerPage: 2}}
	runs, _, err := client.Actions.ListWorkflowRunsByFileName(context.Background(), "o", "r", "29679449", opts)
	if err != nil {
		t.Errorf("Actions.ListWorkFlowRunsByFileName returned error: %v", err)
	}

	want := &WorkflowRuns{
		TotalCount: Int(4),
		WorkflowRuns: []*WorkflowRun{
			{ID: Int64(399444496), RunNumber: Int(296), CreatedAt: &Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
			{ID: Int64(399444497), RunNumber: Int(296), CreatedAt: &Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
		},
	}
	if !reflect.DeepEqual(runs, want) {
		t.Errorf("Actions.ListWorkflowRunsByFileName returned %+v, want %+v", runs, want)
	}
}

func TestActionsService_GetWorkflowRunByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runs/29679449", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":399444496,"run_number":296,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}}`)
	})

	runs, _, err := client.Actions.GetWorkflowRunByID(context.Background(), "o", "r", 29679449)
	if err != nil {
		t.Errorf("Actions.ListWorkFlowRunsByFileName returned error: %v", err)
	}

	want := &WorkflowRun{
		ID:        Int64(399444496),
		RunNumber: Int(296),
		CreatedAt: &Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
		UpdatedAt: &Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
	}

	if !reflect.DeepEqual(runs, want) {
		t.Errorf("Actions.GetWorkflowRunByID returned %+v, want %+v", runs, want)
	}
}

func TestActionsService_RerunWorkflowRunByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runs/3434/rerun", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := client.Actions.RerunWorkflowByID(context.Background(), "o", "r", 3434)
	if err != nil {
		t.Errorf("Actions.RerunWorkflowByID returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Actions.RerunWorkflowRunByID returned status: %d, want %d", resp.StatusCode, http.StatusCreated)
	}
}

func TestActionsService_CancelWorkflowRunByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runs/3434/cancel", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusAccepted)
	})

	resp, err := client.Actions.CancelWorkflowRunByID(context.Background(), "o", "r", 3434)
	if _, ok := err.(*AcceptedError); !ok {
		t.Errorf("Actions.CancelWorkflowRunByID returned error: %v (want AcceptedError)", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Actions.CancelWorkflowRunByID returned status: %d, want %d", resp.StatusCode, http.StatusAccepted)
	}
}

func TestActionsService_GetWorkflowRunLogs(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runs/399444496/logs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/a", http.StatusFound)
	})

	url, resp, err := client.Actions.GetWorkflowRunLogs(context.Background(), "o", "r", 399444496, true)
	if err != nil {
		t.Errorf("Actions.GetWorkflowRunLogs returned error: %v", err)
	}
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Actions.GetWorkflowRunLogs returned status: %d, want %d", resp.StatusCode, http.StatusFound)
	}
	want := "http://github.com/a"
	if url.String() != want {
		t.Errorf("Actions.GetWorkflowRunLogs returned %+v, want %+v", url.String(), want)
	}
}

func TestActionsService_GetWorkflowRunLogs_StatusMovedPermanently_dontFollowRedirects(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runs/399444496/logs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/a", http.StatusMovedPermanently)
	})

	_, resp, _ := client.Actions.GetWorkflowRunLogs(context.Background(), "o", "r", 399444496, false)
	if resp.StatusCode != http.StatusMovedPermanently {
		t.Errorf("Actions.GetWorkflowJobLogs returned status: %d, want %d", resp.StatusCode, http.StatusMovedPermanently)
	}
}

func TestActionsService_GetWorkflowRunLogs_StatusMovedPermanently_followRedirects(t *testing.T) {
	client, mux, serverURL, teardown := setup()
	defer teardown()

	// Mock a redirect link, which leads to an archive link
	mux.HandleFunc("/repos/o/r/actions/runs/399444496/logs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		redirectURL, _ := url.Parse(serverURL + baseURLPath + "/redirect")
		http.Redirect(w, r, redirectURL.String(), http.StatusMovedPermanently)
	})

	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/a", http.StatusFound)
	})

	url, resp, err := client.Actions.GetWorkflowRunLogs(context.Background(), "o", "r", 399444496, true)
	if err != nil {
		t.Errorf("Actions.GetWorkflowJobLogs returned error: %v", err)
	}

	if resp.StatusCode != http.StatusFound {
		t.Errorf("Actions.GetWorkflowJobLogs returned status: %d, want %d", resp.StatusCode, http.StatusFound)
	}

	want := "http://github.com/a"
	if url.String() != want {
		t.Errorf("Actions.GetWorkflowJobLogs returned %+v, want %+v", url.String(), want)
	}
}
