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
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestActionsService_ListWorkflowRunsByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/workflows/29679449/runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"workflow_runs":[{"id":399444496,"run_number":296,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"id":399444497,"run_number":296,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListWorkflowRunsOptions{ListOptions: ListOptions{Page: 2, PerPage: 2}}
	ctx := context.Background()
	runs, _, err := client.Actions.ListWorkflowRunsByID(ctx, "o", "r", 29679449, opts)
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
	if !cmp.Equal(runs, want) {
		t.Errorf("Actions.ListWorkflowRunsByID returned %+v, want %+v", runs, want)
	}

	const methodName = "ListWorkflowRunsByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListWorkflowRunsByID(ctx, "\n", "\n", 29679449, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListWorkflowRunsByID(ctx, "o", "r", 29679449, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_ListWorkflowRunsFileName(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/workflows/29679449/runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"workflow_runs":[{"id":399444496,"run_number":296,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"id":399444497,"run_number":296,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListWorkflowRunsOptions{ListOptions: ListOptions{Page: 2, PerPage: 2}}
	ctx := context.Background()
	runs, _, err := client.Actions.ListWorkflowRunsByFileName(ctx, "o", "r", "29679449", opts)
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
	if !cmp.Equal(runs, want) {
		t.Errorf("Actions.ListWorkflowRunsByFileName returned %+v, want %+v", runs, want)
	}

	const methodName = "ListWorkflowRunsByFileName"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListWorkflowRunsByFileName(ctx, "\n", "\n", "29679449", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListWorkflowRunsByFileName(ctx, "o", "r", "29679449", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetWorkflowRunByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/29679449", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":399444496,"run_number":296,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}}`)
	})

	ctx := context.Background()
	runs, _, err := client.Actions.GetWorkflowRunByID(ctx, "o", "r", 29679449)
	if err != nil {
		t.Errorf("Actions.GetWorkflowRunByID returned error: %v", err)
	}

	want := &WorkflowRun{
		ID:        Int64(399444496),
		RunNumber: Int(296),
		CreatedAt: &Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
		UpdatedAt: &Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
	}

	if !cmp.Equal(runs, want) {
		t.Errorf("Actions.GetWorkflowRunByID returned %+v, want %+v", runs, want)
	}

	const methodName = "GetWorkflowRunByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetWorkflowRunByID(ctx, "\n", "\n", 29679449)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetWorkflowRunByID(ctx, "o", "r", 29679449)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetWorkflowRunAttempt(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/29679449/attempts/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"exclude_pull_requests": "true"})
		fmt.Fprint(w, `{"id":399444496,"run_number":296,"run_attempt":3,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}}`)
	})

	opts := &WorkflowRunAttemptOptions{ExcludePullRequests: Bool(true)}
	ctx := context.Background()
	runs, _, err := client.Actions.GetWorkflowRunAttempt(ctx, "o", "r", 29679449, 3, opts)
	if err != nil {
		t.Errorf("Actions.GetWorkflowRunAttempt returned error: %v", err)
	}

	want := &WorkflowRun{
		ID:         Int64(399444496),
		RunNumber:  Int(296),
		RunAttempt: Int(3),
		CreatedAt:  &Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
		UpdatedAt:  &Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
	}

	if !cmp.Equal(runs, want) {
		t.Errorf("Actions.GetWorkflowRunAttempt returned %+v, want %+v", runs, want)
	}

	const methodName = "GetWorkflowRunAttempt"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetWorkflowRunAttempt(ctx, "\n", "\n", 29679449, 3, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetWorkflowRunAttempt(ctx, "o", "r", 29679449, 3, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetWorkflowRunAttemptLogs(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/399444496/attempts/2/logs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/a", http.StatusFound)
	})

	ctx := context.Background()
	url, resp, err := client.Actions.GetWorkflowRunAttemptLogs(ctx, "o", "r", 399444496, 2, 1)
	if err != nil {
		t.Errorf("Actions.GetWorkflowRunAttemptLogs returned error: %v", err)
	}
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Actions.GetWorkflowRunAttemptLogs returned status: %d, want %d", resp.StatusCode, http.StatusFound)
	}
	want := "http://github.com/a"
	if url.String() != want {
		t.Errorf("Actions.GetWorkflowRunAttemptLogs returned %+v, want %+v", url.String(), want)
	}

	const methodName = "GetWorkflowRunAttemptLogs"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetWorkflowRunAttemptLogs(ctx, "\n", "\n", 399444496, 2, 1)
		return err
	})
}

func TestActionsService_GetWorkflowRunAttemptLogs_StatusMovedPermanently_dontFollowRedirects(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/399444496/attempts/2/logs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/a", http.StatusMovedPermanently)
	})

	ctx := context.Background()
	_, resp, _ := client.Actions.GetWorkflowRunAttemptLogs(ctx, "o", "r", 399444496, 2, 0)
	if resp.StatusCode != http.StatusMovedPermanently {
		t.Errorf("Actions.GetWorkflowRunAttemptLogs returned status: %d, want %d", resp.StatusCode, http.StatusMovedPermanently)
	}
}

func TestActionsService_GetWorkflowRunAttemptLogs_StatusMovedPermanently_followRedirects(t *testing.T) {
	t.Parallel()
	client, mux, serverURL := setup(t)

	// Mock a redirect link, which leads to an archive link
	mux.HandleFunc("/repos/o/r/actions/runs/399444496/attempts/2/logs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		redirectURL, _ := url.Parse(serverURL + baseURLPath + "/redirect")
		http.Redirect(w, r, redirectURL.String(), http.StatusMovedPermanently)
	})

	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/a", http.StatusFound)
	})

	ctx := context.Background()
	url, resp, err := client.Actions.GetWorkflowRunAttemptLogs(ctx, "o", "r", 399444496, 2, 1)
	if err != nil {
		t.Errorf("Actions.GetWorkflowRunAttemptLogs returned error: %v", err)
	}

	if resp.StatusCode != http.StatusFound {
		t.Errorf("Actions.GetWorkflowRunAttemptLogs returned status: %d, want %d", resp.StatusCode, http.StatusFound)
	}

	want := "http://github.com/a"
	if url.String() != want {
		t.Errorf("Actions.GetWorkflowRunAttemptLogs returned %+v, want %+v", url.String(), want)
	}

	const methodName = "GetWorkflowRunAttemptLogs"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetWorkflowRunAttemptLogs(ctx, "\n", "\n", 399444496, 2, 1)
		return err
	})
}

func TestActionsService_RerunWorkflowRunByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/3434/rerun", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
	})

	ctx := context.Background()
	resp, err := client.Actions.RerunWorkflowByID(ctx, "o", "r", 3434)
	if err != nil {
		t.Errorf("Actions.RerunWorkflowByID returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Actions.RerunWorkflowRunByID returned status: %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	const methodName = "RerunWorkflowByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RerunWorkflowByID(ctx, "\n", "\n", 3434)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.RerunWorkflowByID(ctx, "o", "r", 3434)
	})
}

func TestActionsService_RerunFailedJobsByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/3434/rerun-failed-jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
	})

	ctx := context.Background()
	resp, err := client.Actions.RerunFailedJobsByID(ctx, "o", "r", 3434)
	if err != nil {
		t.Errorf("Actions.RerunFailedJobsByID returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Actions.RerunFailedJobsByID returned status: %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	const methodName = "RerunFailedJobsByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RerunFailedJobsByID(ctx, "\n", "\n", 3434)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.RerunFailedJobsByID(ctx, "o", "r", 3434)
	})
}

func TestActionsService_RerunJobByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/jobs/3434/rerun", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
	})

	ctx := context.Background()
	resp, err := client.Actions.RerunJobByID(ctx, "o", "r", 3434)
	if err != nil {
		t.Errorf("Actions.RerunJobByID returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Actions.RerunJobByID returned status: %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	const methodName = "RerunJobByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RerunJobByID(ctx, "\n", "\n", 3434)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.RerunJobByID(ctx, "o", "r", 3434)
	})
}

func TestActionsService_CancelWorkflowRunByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/3434/cancel", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusAccepted)
	})

	ctx := context.Background()
	resp, err := client.Actions.CancelWorkflowRunByID(ctx, "o", "r", 3434)
	if _, ok := err.(*AcceptedError); !ok {
		t.Errorf("Actions.CancelWorkflowRunByID returned error: %v (want AcceptedError)", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Actions.CancelWorkflowRunByID returned status: %d, want %d", resp.StatusCode, http.StatusAccepted)
	}

	const methodName = "CancelWorkflowRunByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.CancelWorkflowRunByID(ctx, "\n", "\n", 3434)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.CancelWorkflowRunByID(ctx, "o", "r", 3434)
	})
}

func TestActionsService_GetWorkflowRunLogs(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/399444496/logs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/a", http.StatusFound)
	})

	ctx := context.Background()
	url, resp, err := client.Actions.GetWorkflowRunLogs(ctx, "o", "r", 399444496, 1)
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

	const methodName = "GetWorkflowRunLogs"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetWorkflowRunLogs(ctx, "\n", "\n", 399444496, 1)
		return err
	})
}

func TestActionsService_GetWorkflowRunLogs_StatusMovedPermanently_dontFollowRedirects(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/399444496/logs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/a", http.StatusMovedPermanently)
	})

	ctx := context.Background()
	_, resp, _ := client.Actions.GetWorkflowRunLogs(ctx, "o", "r", 399444496, 0)
	if resp.StatusCode != http.StatusMovedPermanently {
		t.Errorf("Actions.GetWorkflowJobLogs returned status: %d, want %d", resp.StatusCode, http.StatusMovedPermanently)
	}
}

func TestActionsService_GetWorkflowRunLogs_StatusMovedPermanently_followRedirects(t *testing.T) {
	t.Parallel()
	client, mux, serverURL := setup(t)

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

	ctx := context.Background()
	url, resp, err := client.Actions.GetWorkflowRunLogs(ctx, "o", "r", 399444496, 1)
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

	const methodName = "GetWorkflowRunLogs"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetWorkflowRunLogs(ctx, "\n", "\n", 399444496, 1)
		return err
	})
}

func TestActionService_ListRepositoryWorkflowRuns(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":2,
		"workflow_runs":[
			{"id":298499444,"run_number":301,"created_at":"2020-04-11T11:14:54Z","updated_at":"2020-04-11T11:14:54Z"},
			{"id":298499445,"run_number":302,"created_at":"2020-04-11T11:14:54Z","updated_at":"2020-04-11T11:14:54Z"}]}`)
	})

	opts := &ListWorkflowRunsOptions{ListOptions: ListOptions{Page: 2, PerPage: 2}}
	ctx := context.Background()
	runs, _, err := client.Actions.ListRepositoryWorkflowRuns(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListRepositoryWorkflowRuns returned error: %v", err)
	}

	expected := &WorkflowRuns{
		TotalCount: Int(2),
		WorkflowRuns: []*WorkflowRun{
			{ID: Int64(298499444), RunNumber: Int(301), CreatedAt: &Timestamp{time.Date(2020, time.April, 11, 11, 14, 54, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.April, 11, 11, 14, 54, 0, time.UTC)}},
			{ID: Int64(298499445), RunNumber: Int(302), CreatedAt: &Timestamp{time.Date(2020, time.April, 11, 11, 14, 54, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.April, 11, 11, 14, 54, 0, time.UTC)}},
		},
	}

	if !cmp.Equal(runs, expected) {
		t.Errorf("Actions.ListRepositoryWorkflowRuns returned %+v, want %+v", runs, expected)
	}

	const methodName = "ListRepositoryWorkflowRuns"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListRepositoryWorkflowRuns(ctx, "\n", "\n", opts)

		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListRepositoryWorkflowRuns(ctx, "o", "r", opts)

		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionService_DeleteWorkflowRun(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/399444496", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Actions.DeleteWorkflowRun(ctx, "o", "r", 399444496); err != nil {
		t.Errorf("DeleteWorkflowRun returned error: %v", err)
	}

	const methodName = "DeleteWorkflowRun"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteWorkflowRun(ctx, "\n", "\n", 399444496)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteWorkflowRun(ctx, "o", "r", 399444496)
	})
}

func TestActionService_DeleteWorkflowRunLogs(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/399444496/logs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Actions.DeleteWorkflowRunLogs(ctx, "o", "r", 399444496); err != nil {
		t.Errorf("DeleteWorkflowRunLogs returned error: %v", err)
	}

	const methodName = "DeleteWorkflowRunLogs"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteWorkflowRunLogs(ctx, "\n", "\n", 399444496)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteWorkflowRunLogs(ctx, "o", "r", 399444496)
	})
}

func TestPendingDeployment_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PendingDeployment{}, "{}")

	u := &PendingDeployment{
		Environment: &PendingDeploymentEnvironment{
			ID:      Int64(1),
			NodeID:  String("nid"),
			Name:    String("n"),
			URL:     String("u"),
			HTMLURL: String("hu"),
		},
		WaitTimer:             Int64(100),
		WaitTimerStartedAt:    &Timestamp{referenceTime},
		CurrentUserCanApprove: Bool(false),
		Reviewers: []*RequiredReviewer{
			{
				Type: String("User"),
				Reviewer: &User{
					Login: String("l"),
				},
			},
			{
				Type: String("Team"),
				Reviewer: &Team{
					Name: String("n"),
				},
			},
		},
	}
	want := `{
		"environment": {
			"id": 1,
			"node_id": "nid",
			"name": "n",
			"url": "u",
			"html_url": "hu"
		},
		"wait_timer": 100,
		"wait_timer_started_at": ` + referenceTimeStr + `,
		"current_user_can_approve": false,
		"reviewers": [
			{
				"type": "User",
				"reviewer": {
					"login": "l"
				}
			},
			{
				"type": "Team",
				"reviewer": {
					"name": "n"
				}
			}
		]
	}`
	testJSONMarshal(t, u, want)
}

func TestActionsService_ReviewCustomDeploymentProtectionRule(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/9444496/deployment_protection_rule", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		w.WriteHeader(http.StatusNoContent)
	})

	request := ReviewCustomDeploymentProtectionRuleRequest{
		EnvironmentName: "production",
		State:           "approved",
		Comment:         "Approve deployment",
	}

	ctx := context.Background()
	if _, err := client.Actions.ReviewCustomDeploymentProtectionRule(ctx, "o", "r", 9444496, &request); err != nil {
		t.Errorf("ReviewCustomDeploymentProtectionRule returned error: %v", err)
	}

	const methodName = "ReviewCustomDeploymentProtectionRule"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.ReviewCustomDeploymentProtectionRule(ctx, "\n", "\n", 9444496, &request)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.ReviewCustomDeploymentProtectionRule(ctx, "o", "r", 9444496, &request)
	})
}

func TestReviewCustomDeploymentProtectionRuleRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ReviewCustomDeploymentProtectionRuleRequest{}, "{}")

	r := &ReviewCustomDeploymentProtectionRuleRequest{
		EnvironmentName: "e",
		State:           "rejected",
		Comment:         "c",
	}
	want := `{
		"environment_name": "e",
		"state": "rejected",
		"comment": "c"
	}`
	testJSONMarshal(t, r, want)
}

func TestActionsService_GetWorkflowRunUsageByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/29679449/timing", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"billable":{"UBUNTU":{"total_ms":180000,"jobs":1,"job_runs":[{"job_id":1,"duration_ms":60000}]},"MACOS":{"total_ms":240000,"jobs":2,"job_runs":[{"job_id":2,"duration_ms":30000},{"job_id":3,"duration_ms":10000}]},"WINDOWS":{"total_ms":300000,"jobs":2}},"run_duration_ms":500000}`)
	})

	ctx := context.Background()
	workflowRunUsage, _, err := client.Actions.GetWorkflowRunUsageByID(ctx, "o", "r", 29679449)
	if err != nil {
		t.Errorf("Actions.GetWorkflowRunUsageByID returned error: %v", err)
	}

	want := &WorkflowRunUsage{
		Billable: &WorkflowRunBillMap{
			"UBUNTU": &WorkflowRunBill{
				TotalMS: Int64(180000),
				Jobs:    Int(1),
				JobRuns: []*WorkflowRunJobRun{
					{
						JobID:      Int(1),
						DurationMS: Int64(60000),
					},
				},
			},
			"MACOS": &WorkflowRunBill{
				TotalMS: Int64(240000),
				Jobs:    Int(2),
				JobRuns: []*WorkflowRunJobRun{
					{
						JobID:      Int(2),
						DurationMS: Int64(30000),
					},
					{
						JobID:      Int(3),
						DurationMS: Int64(10000),
					},
				},
			},
			"WINDOWS": &WorkflowRunBill{
				TotalMS: Int64(300000),
				Jobs:    Int(2),
			},
		},
		RunDurationMS: Int64(500000),
	}

	if !cmp.Equal(workflowRunUsage, want) {
		t.Errorf("Actions.GetWorkflowRunUsageByID returned %+v, want %+v", workflowRunUsage, want)
	}

	const methodName = "GetWorkflowRunUsageByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetWorkflowRunUsageByID(ctx, "\n", "\n", 29679449)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetWorkflowRunUsageByID(ctx, "o", "r", 29679449)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestWorkflowRun_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &WorkflowRun{}, "{}")

	u := &WorkflowRun{
		ID:         Int64(1),
		Name:       String("n"),
		NodeID:     String("nid"),
		HeadBranch: String("hb"),
		HeadSHA:    String("hs"),
		Path:       String("p"),
		RunNumber:  Int(1),
		RunAttempt: Int(1),
		Event:      String("e"),
		Status:     String("s"),
		Conclusion: String("c"),
		WorkflowID: Int64(1),
		URL:        String("u"),
		HTMLURL:    String("h"),
		PullRequests: []*PullRequest{
			{
				URL:    String("u"),
				ID:     Int64(1),
				Number: Int(1),
				Head: &PullRequestBranch{
					Ref: String("r"),
					SHA: String("s"),
					Repo: &Repository{
						ID:   Int64(1),
						URL:  String("s"),
						Name: String("n"),
					},
				},
				Base: &PullRequestBranch{
					Ref: String("r"),
					SHA: String("s"),
					Repo: &Repository{
						ID:   Int64(1),
						URL:  String("u"),
						Name: String("n"),
					},
				},
			},
		},
		CreatedAt:          &Timestamp{referenceTime},
		UpdatedAt:          &Timestamp{referenceTime},
		RunStartedAt:       &Timestamp{referenceTime},
		JobsURL:            String("j"),
		LogsURL:            String("l"),
		CheckSuiteURL:      String("c"),
		ArtifactsURL:       String("a"),
		CancelURL:          String("c"),
		RerunURL:           String("r"),
		PreviousAttemptURL: String("p"),
		HeadCommit: &HeadCommit{
			Message: String("m"),
			Author: &CommitAuthor{
				Name:  String("n"),
				Email: String("e"),
				Login: String("l"),
			},
			URL:       String("u"),
			Distinct:  Bool(false),
			SHA:       String("s"),
			ID:        String("i"),
			TreeID:    String("tid"),
			Timestamp: &Timestamp{referenceTime},
			Committer: &CommitAuthor{
				Name:  String("n"),
				Email: String("e"),
				Login: String("l"),
			},
		},
		WorkflowURL: String("w"),
		Repository: &Repository{
			ID:   Int64(1),
			URL:  String("u"),
			Name: String("n"),
		},
		HeadRepository: &Repository{
			ID:   Int64(1),
			URL:  String("u"),
			Name: String("n"),
		},
		Actor: &User{
			Login:           String("l"),
			ID:              Int64(1),
			AvatarURL:       String("a"),
			GravatarID:      String("g"),
			Name:            String("n"),
			Company:         String("c"),
			Blog:            String("b"),
			Location:        String("l"),
			Email:           String("e"),
			Hireable:        Bool(true),
			Bio:             String("b"),
			TwitterUsername: String("t"),
			PublicRepos:     Int(1),
			Followers:       Int(1),
			Following:       Int(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
			URL:             String("u"),
		},
		TriggeringActor: &User{
			Login:           String("l2"),
			ID:              Int64(2),
			AvatarURL:       String("a2"),
			GravatarID:      String("g2"),
			Name:            String("n2"),
			Company:         String("c2"),
			Blog:            String("b2"),
			Location:        String("l2"),
			Email:           String("e2"),
			Hireable:        Bool(false),
			Bio:             String("b2"),
			TwitterUsername: String("t2"),
			PublicRepos:     Int(2),
			Followers:       Int(2),
			Following:       Int(2),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
			URL:             String("u2"),
		},
		ReferencedWorkflows: []*ReferencedWorkflow{
			{
				Path: String("rwfp"),
				SHA:  String("rwfsha"),
				Ref:  String("rwfref"),
			},
		},
	}

	want := `{
		"id": 1,
		"name": "n",
		"node_id": "nid",
		"head_branch": "hb",
		"head_sha": "hs",
		"path": "p",
		"run_number": 1,
		"run_attempt": 1,
		"event": "e",
		"status": "s",
		"conclusion": "c",
		"workflow_id": 1,
		"url": "u",
		"html_url": "h",
		"pull_requests": [
			{
				"id":1,
				"number":1,
				"url":"u",
				"head":{
					"ref":"r",
					"sha":"s",
					"repo": {
						"id":1,
						"name":"n",
						"url":"s"
						}
					},
					"base": {
						"ref":"r",
						"sha":"s",
						"repo": {
							"id":1,
							"name":"n",
							"url":"u"
						}
					}
			}
		],
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"run_started_at": ` + referenceTimeStr + `,
		"jobs_url": "j",
		"logs_url": "l",
		"check_suite_url": "c",
		"artifacts_url": "a",
		"cancel_url": "c",
		"rerun_url": "r",
		"previous_attempt_url": "p",
		"head_commit": {
			"message": "m",
			"author": {
				"name": "n",
				"email": "e",
				"username": "l"
			},
			"url": "u",
			"distinct": false,
			"sha": "s",
			"id": "i",
			"tree_id": "tid",
			"timestamp": ` + referenceTimeStr + `,
			"committer": {
				"name": "n",
				"email": "e",
				"username": "l"
			}
		},
		"workflow_url": "w",
		"repository": {
			"id": 1,
			"url": "u",
			"name": "n"
		},
		"head_repository": {
			"id": 1,
			"url": "u",
			"name": "n"
		},
		"actor": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"triggering_actor": {
			"login": "l2",
			"id": 2,
			"avatar_url": "a2",
			"gravatar_id": "g2",
			"name": "n2",
			"company": "c2",
			"blog": "b2",
			"location": "l2",
			"email": "e2",
			"hireable": false,
			"bio": "b2",
			"twitter_username": "t2",
			"public_repos": 2,
			"followers": 2,
			"following": 2,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u2"
		},
		"referenced_workflows": [
			{
				"path": "rwfp",
				"sha": "rwfsha",
				"ref": "rwfref"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestWorkflowRuns_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &WorkflowRuns{}, "{}")

	u := &WorkflowRuns{
		TotalCount: Int(1),
		WorkflowRuns: []*WorkflowRun{
			{
				ID:         Int64(1),
				Name:       String("n"),
				NodeID:     String("nid"),
				HeadBranch: String("hb"),
				HeadSHA:    String("hs"),
				RunNumber:  Int(1),
				RunAttempt: Int(1),
				Event:      String("e"),
				Status:     String("s"),
				Conclusion: String("c"),
				WorkflowID: Int64(1),
				URL:        String("u"),
				HTMLURL:    String("h"),
				PullRequests: []*PullRequest{
					{
						URL:    String("u"),
						ID:     Int64(1),
						Number: Int(1),
						Head: &PullRequestBranch{
							Ref: String("r"),
							SHA: String("s"),
							Repo: &Repository{
								ID:   Int64(1),
								URL:  String("s"),
								Name: String("n"),
							},
						},
						Base: &PullRequestBranch{
							Ref: String("r"),
							SHA: String("s"),
							Repo: &Repository{
								ID:   Int64(1),
								URL:  String("u"),
								Name: String("n"),
							},
						},
					},
				},
				CreatedAt:          &Timestamp{referenceTime},
				UpdatedAt:          &Timestamp{referenceTime},
				RunStartedAt:       &Timestamp{referenceTime},
				JobsURL:            String("j"),
				LogsURL:            String("l"),
				CheckSuiteURL:      String("c"),
				ArtifactsURL:       String("a"),
				CancelURL:          String("c"),
				RerunURL:           String("r"),
				PreviousAttemptURL: String("p"),
				HeadCommit: &HeadCommit{
					Message: String("m"),
					Author: &CommitAuthor{
						Name:  String("n"),
						Email: String("e"),
						Login: String("l"),
					},
					URL:       String("u"),
					Distinct:  Bool(false),
					SHA:       String("s"),
					ID:        String("i"),
					TreeID:    String("tid"),
					Timestamp: &Timestamp{referenceTime},
					Committer: &CommitAuthor{
						Name:  String("n"),
						Email: String("e"),
						Login: String("l"),
					},
				},
				WorkflowURL: String("w"),
				Repository: &Repository{
					ID:   Int64(1),
					URL:  String("u"),
					Name: String("n"),
				},
				HeadRepository: &Repository{
					ID:   Int64(1),
					URL:  String("u"),
					Name: String("n"),
				},
				Actor: &User{
					Login:           String("l"),
					ID:              Int64(1),
					AvatarURL:       String("a"),
					GravatarID:      String("g"),
					Name:            String("n"),
					Company:         String("c"),
					Blog:            String("b"),
					Location:        String("l"),
					Email:           String("e"),
					Hireable:        Bool(true),
					Bio:             String("b"),
					TwitterUsername: String("t"),
					PublicRepos:     Int(1),
					Followers:       Int(1),
					Following:       Int(1),
					CreatedAt:       &Timestamp{referenceTime},
					SuspendedAt:     &Timestamp{referenceTime},
					URL:             String("u"),
				},
				TriggeringActor: &User{
					Login:           String("l2"),
					ID:              Int64(2),
					AvatarURL:       String("a2"),
					GravatarID:      String("g2"),
					Name:            String("n2"),
					Company:         String("c2"),
					Blog:            String("b2"),
					Location:        String("l2"),
					Email:           String("e2"),
					Hireable:        Bool(false),
					Bio:             String("b2"),
					TwitterUsername: String("t2"),
					PublicRepos:     Int(2),
					Followers:       Int(2),
					Following:       Int(2),
					CreatedAt:       &Timestamp{referenceTime},
					SuspendedAt:     &Timestamp{referenceTime},
					URL:             String("u2"),
				},
			},
		},
	}

	want := `{
		"total_count": 1,
		"workflow_runs": [
			{
				"id": 1,
				"name": "n",
				"node_id": "nid",
				"head_branch": "hb",
				"head_sha": "hs",
				"run_number": 1,
				"run_attempt": 1,
				"event": "e",
				"status": "s",
				"conclusion": "c",
				"workflow_id": 1,
				"url": "u",
				"html_url": "h",
				"pull_requests": [
					{
						"id":1,
						"number":1,
						"url":"u",
						"head":{
							"ref":"r",
							"sha":"s",
							"repo": {
								"id":1,
								"name":"n",
								"url":"s"
								}
							},
							"base": {
								"ref":"r",
								"sha":"s",
								"repo": {
									"id":1,
									"name":"n",
									"url":"u"
								}
							}
					}
				],
				"created_at": ` + referenceTimeStr + `,
				"updated_at": ` + referenceTimeStr + `,
				"run_started_at": ` + referenceTimeStr + `,
				"jobs_url": "j",
				"logs_url": "l",
				"check_suite_url": "c",
				"artifacts_url": "a",
				"cancel_url": "c",
				"rerun_url": "r",
				"previous_attempt_url": "p",
				"head_commit": {
					"message": "m",
					"author": {
						"name": "n",
						"email": "e",
						"username": "l"
					},
					"url": "u",
					"distinct": false,
					"sha": "s",
					"id": "i",
					"tree_id": "tid",
					"timestamp": ` + referenceTimeStr + `,
					"committer": {
						"name": "n",
						"email": "e",
						"username": "l"
					}
				},
				"workflow_url": "w",
				"repository": {
					"id": 1,
					"url": "u",
					"name": "n"
				},
				"head_repository": {
					"id": 1,
					"url": "u",
					"name": "n"
				},
				"actor": {
					"login": "l",
					"id": 1,
					"avatar_url": "a",
					"gravatar_id": "g",
					"name": "n",
					"company": "c",
					"blog": "b",
					"location": "l",
					"email": "e",
					"hireable": true,
					"bio": "b",
					"twitter_username": "t",
					"public_repos": 1,
					"followers": 1,
					"following": 1,
					"created_at": ` + referenceTimeStr + `,
					"suspended_at": ` + referenceTimeStr + `,
					"url": "u"
				},
				"triggering_actor": {
					"login": "l2",
					"id": 2,
					"avatar_url": "a2",
					"gravatar_id": "g2",
					"name": "n2",
					"company": "c2",
					"blog": "b2",
					"location": "l2",
					"email": "e2",
					"hireable": false,
					"bio": "b2",
					"twitter_username": "t2",
					"public_repos": 2,
					"followers": 2,
					"following": 2,
					"created_at": ` + referenceTimeStr + `,
					"suspended_at": ` + referenceTimeStr + `,
					"url": "u2"
				}
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestWorkflowRunBill_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &WorkflowRunBill{}, "{}")

	u := &WorkflowRunBill{
		TotalMS: Int64(1),
		Jobs:    Int(1),
	}

	want := `{
		"total_ms": 1,
		"jobs": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestWorkflowRunBillMap_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &WorkflowRunBillMap{}, "{}")

	u := &WorkflowRunBillMap{
		"UBUNTU": &WorkflowRunBill{
			TotalMS: Int64(1),
			Jobs:    Int(1),
		},
		"MACOS": &WorkflowRunBill{
			TotalMS: Int64(1),
			Jobs:    Int(1),
		},
		"WINDOWS": &WorkflowRunBill{
			TotalMS: Int64(1),
			Jobs:    Int(1),
		},
	}

	want := `{
		"UBUNTU": {
			"total_ms": 1,
			"jobs": 1
		},
		"MACOS": {
			"total_ms": 1,
			"jobs": 1
		},
		"WINDOWS": {
			"total_ms": 1,
			"jobs": 1
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestWorkflowRunUsage_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &WorkflowRunUsage{}, "{}")

	u := &WorkflowRunUsage{
		Billable: &WorkflowRunBillMap{
			"UBUNTU": &WorkflowRunBill{
				TotalMS: Int64(1),
				Jobs:    Int(1),
			},
			"MACOS": &WorkflowRunBill{
				TotalMS: Int64(1),
				Jobs:    Int(1),
			},
			"WINDOWS": &WorkflowRunBill{
				TotalMS: Int64(1),
				Jobs:    Int(1),
			},
		},
		RunDurationMS: Int64(1),
	}

	want := `{
		"billable": {
			"UBUNTU": {
				"total_ms": 1,
				"jobs": 1
			},
			"MACOS": {
				"total_ms": 1,
				"jobs": 1
			},
			"WINDOWS": {
				"total_ms": 1,
				"jobs": 1
			}
		},
		"run_duration_ms": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestActionService_PendingDeployments(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &PendingDeploymentsRequest{EnvironmentIDs: []int64{3, 4}, State: "approved", Comment: ""}

	mux.HandleFunc("/repos/o/r/actions/runs/399444496/pending_deployments", func(w http.ResponseWriter, r *http.Request) {
		v := new(PendingDeploymentsRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	ctx := context.Background()
	deployments, _, err := client.Actions.PendingDeployments(ctx, "o", "r", 399444496, input)
	if err != nil {
		t.Errorf("Actions.PendingDeployments returned error: %v", err)
	}

	want := []*Deployment{{ID: Int64(1)}, {ID: Int64(2)}}
	if !cmp.Equal(deployments, want) {
		t.Errorf("Actions.PendingDeployments returned %+v, want %+v", deployments, want)
	}

	const methodName = "PendingDeployments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.PendingDeployments(ctx, "\n", "\n", 399444496, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.PendingDeployments(ctx, "o", "r", 399444496, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionService_GetPendingDeployments(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/399444496/pending_deployments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{
				"environment": {
					"id": 1,
					"node_id": "nid",
					"name": "n",
					"url": "u",
					"html_url": "hu"
				},
				"wait_timer": 0,
				"wait_timer_started_at": `+referenceTimeStr+`,
				"current_user_can_approve": false,
				"reviewers": []
			},
			{
				"environment": {
					"id": 2,
					"node_id": "nid",
					"name": "n",
					"url": "u",
					"html_url": "hu"
				},
				"wait_timer": 13,
				"wait_timer_started_at": `+referenceTimeStr+`,
				"current_user_can_approve": true,
				"reviewers": [
					{
						"type": "User",
						"reviewer": {
							"login": "l"
						}
					},
					{
						"type": "Team",
						"reviewer": {
							"name": "t",
							"slug": "s"
						}
					}
				]
			}
		]`)
	})

	ctx := context.Background()
	deployments, _, err := client.Actions.GetPendingDeployments(ctx, "o", "r", 399444496)
	if err != nil {
		t.Errorf("Actions.GetPendingDeployments returned error: %v", err)
	}

	want := []*PendingDeployment{
		{
			Environment: &PendingDeploymentEnvironment{
				ID:      Int64(1),
				NodeID:  String("nid"),
				Name:    String("n"),
				URL:     String("u"),
				HTMLURL: String("hu"),
			},
			WaitTimer:             Int64(0),
			WaitTimerStartedAt:    &Timestamp{referenceTime},
			CurrentUserCanApprove: Bool(false),
			Reviewers:             []*RequiredReviewer{},
		},
		{
			Environment: &PendingDeploymentEnvironment{
				ID:      Int64(2),
				NodeID:  String("nid"),
				Name:    String("n"),
				URL:     String("u"),
				HTMLURL: String("hu"),
			},
			WaitTimer:             Int64(13),
			WaitTimerStartedAt:    &Timestamp{referenceTime},
			CurrentUserCanApprove: Bool(true),
			Reviewers: []*RequiredReviewer{
				{
					Type: String("User"),
					Reviewer: &User{
						Login: String("l"),
					},
				},
				{
					Type: String("Team"),
					Reviewer: &Team{
						Name: String("t"),
						Slug: String("s"),
					},
				},
			},
		},
	}

	if !cmp.Equal(deployments, want) {
		t.Errorf("Actions.GetPendingDeployments returned %+v, want %+v", deployments, want)
	}

	const methodName = "GetPendingDeployments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetPendingDeployments(ctx, "\n", "\n", 399444496)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetPendingDeployments(ctx, "o", "r", 399444496)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
