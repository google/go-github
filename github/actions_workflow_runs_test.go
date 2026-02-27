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
	"net/url"
	"strings"
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
	ctx := t.Context()
	runs, _, err := client.Actions.ListWorkflowRunsByID(ctx, "o", "r", 29679449, opts)
	if err != nil {
		t.Errorf("Actions.ListWorkFlowRunsByID returned error: %v", err)
	}

	want := &WorkflowRuns{
		TotalCount: Ptr(4),
		WorkflowRuns: []*WorkflowRun{
			{ID: Ptr(int64(399444496)), RunNumber: Ptr(296), CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
			{ID: Ptr(int64(399444497)), RunNumber: Ptr(296), CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
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
	ctx := t.Context()
	runs, _, err := client.Actions.ListWorkflowRunsByFileName(ctx, "o", "r", "29679449", opts)
	if err != nil {
		t.Errorf("Actions.ListWorkFlowRunsByFileName returned error: %v", err)
	}

	want := &WorkflowRuns{
		TotalCount: Ptr(4),
		WorkflowRuns: []*WorkflowRun{
			{ID: Ptr(int64(399444496)), RunNumber: Ptr(296), CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
			{ID: Ptr(int64(399444497)), RunNumber: Ptr(296), CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
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

	ctx := t.Context()
	runs, _, err := client.Actions.GetWorkflowRunByID(ctx, "o", "r", 29679449)
	if err != nil {
		t.Errorf("Actions.GetWorkflowRunByID returned error: %v", err)
	}

	want := &WorkflowRun{
		ID:        Ptr(int64(399444496)),
		RunNumber: Ptr(296),
		CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)},
		UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)},
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

	opts := &WorkflowRunAttemptOptions{ExcludePullRequests: Ptr(true)}
	ctx := t.Context()
	runs, _, err := client.Actions.GetWorkflowRunAttempt(ctx, "o", "r", 29679449, 3, opts)
	if err != nil {
		t.Errorf("Actions.GetWorkflowRunAttempt returned error: %v", err)
	}

	want := &WorkflowRun{
		ID:         Ptr(int64(399444496)),
		RunNumber:  Ptr(296),
		RunAttempt: Ptr(3),
		CreatedAt:  &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)},
		UpdatedAt:  &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)},
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
	tcs := []struct {
		name              string
		respectRateLimits bool
	}{
		{
			name:              "withoutRateLimits",
			respectRateLimits: false,
		},
		{
			name:              "withRateLimits",
			respectRateLimits: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)
			client.RateLimitRedirectionalEndpoints = tc.respectRateLimits

			mux.HandleFunc("/repos/o/r/actions/runs/399444496/attempts/2/logs", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				http.Redirect(w, r, "https://github.com/a", http.StatusFound)
			})

			ctx := t.Context()
			url, resp, err := client.Actions.GetWorkflowRunAttemptLogs(ctx, "o", "r", 399444496, 2, 1)
			if err != nil {
				t.Errorf("Actions.GetWorkflowRunAttemptLogs returned error: %v", err)
			}
			if resp.StatusCode != http.StatusFound {
				t.Errorf("Actions.GetWorkflowRunAttemptLogs returned status: %v, want %v", resp.StatusCode, http.StatusFound)
			}
			want := "https://github.com/a"
			if url.String() != want {
				t.Errorf("Actions.GetWorkflowRunAttemptLogs returned %+v, want %+v", url, want)
			}

			const methodName = "GetWorkflowRunAttemptLogs"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Actions.GetWorkflowRunAttemptLogs(ctx, "\n", "\n", 399444496, 2, 1)
				return err
			})
		})
	}
}

func TestActionsService_GetWorkflowRunAttemptLogs_StatusMovedPermanently_dontFollowRedirects(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name              string
		respectRateLimits bool
	}{
		{
			name:              "withoutRateLimits",
			respectRateLimits: false,
		},
		{
			name:              "withRateLimits",
			respectRateLimits: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)
			client.RateLimitRedirectionalEndpoints = tc.respectRateLimits

			mux.HandleFunc("/repos/o/r/actions/runs/399444496/attempts/2/logs", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				http.Redirect(w, r, "https://github.com/a", http.StatusMovedPermanently)
			})

			ctx := t.Context()
			_, resp, _ := client.Actions.GetWorkflowRunAttemptLogs(ctx, "o", "r", 399444496, 2, 0)
			if resp.StatusCode != http.StatusMovedPermanently {
				t.Errorf("Actions.GetWorkflowRunAttemptLogs returned status: %v, want %v", resp.StatusCode, http.StatusMovedPermanently)
			}
		})
	}
}

func TestActionsService_GetWorkflowRunAttemptLogs_StatusMovedPermanently_followRedirects(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name              string
		respectRateLimits bool
	}{
		{
			name:              "withoutRateLimits",
			respectRateLimits: false,
		},
		{
			name:              "withRateLimits",
			respectRateLimits: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, mux, serverURL := setup(t)
			client.RateLimitRedirectionalEndpoints = tc.respectRateLimits

			// Mock a redirect link, which leads to an archive link
			mux.HandleFunc("/repos/o/r/actions/runs/399444496/attempts/2/logs", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				redirectURL, _ := url.Parse(serverURL + baseURLPath + "/redirect")
				http.Redirect(w, r, redirectURL.String(), http.StatusMovedPermanently)
			})

			mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				http.Redirect(w, r, "https://github.com/a", http.StatusFound)
			})

			ctx := t.Context()
			url, resp, err := client.Actions.GetWorkflowRunAttemptLogs(ctx, "o", "r", 399444496, 2, 1)
			if err != nil {
				t.Errorf("Actions.GetWorkflowRunAttemptLogs returned error: %v", err)
			}

			if resp.StatusCode != http.StatusFound {
				t.Errorf("Actions.GetWorkflowRunAttemptLogs returned status: %v, want %v", resp.StatusCode, http.StatusFound)
			}

			want := "https://github.com/a"
			if url.String() != want {
				t.Errorf("Actions.GetWorkflowRunAttemptLogs returned %+v, want %+v", url, want)
			}

			const methodName = "GetWorkflowRunAttemptLogs"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Actions.GetWorkflowRunAttemptLogs(ctx, "\n", "\n", 399444496, 2, 1)
				return err
			})
		})
	}
}

func TestActionsService_GetWorkflowRunAttemptLogs_unexpectedCode(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name              string
		respectRateLimits bool
	}{
		{
			name:              "withoutRateLimits",
			respectRateLimits: false,
		},
		{
			name:              "withRateLimits",
			respectRateLimits: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, mux, serverURL := setup(t)
			client.RateLimitRedirectionalEndpoints = tc.respectRateLimits

			// Mock a redirect link, which leads to an archive link
			mux.HandleFunc("/repos/o/r/actions/runs/399444496/attempts/2/logs", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				redirectURL, _ := url.Parse(serverURL + baseURLPath + "/redirect")
				http.Redirect(w, r, redirectURL.String(), http.StatusMovedPermanently)
			})

			mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				w.WriteHeader(http.StatusNoContent)
			})

			ctx := t.Context()
			url, resp, err := client.Actions.GetWorkflowRunAttemptLogs(ctx, "o", "r", 399444496, 2, 1)
			if err == nil {
				t.Fatal("Actions.GetWorkflowRunAttemptLogs should return error on unexpected code")
			}
			if !strings.Contains(err.Error(), "unexpected status code") {
				t.Error("Actions.GetWorkflowRunAttemptLogs should return unexpected status code")
			}
			if got, want := resp.Response.StatusCode, http.StatusNoContent; got != want {
				t.Errorf("Actions.GetWorkflowRunAttemptLogs return status %v, want %v", got, want)
			}
			if url != nil {
				t.Errorf("Actions.GetWorkflowRunAttemptLogs return %+v, want nil", url)
			}
		})
	}
}

func TestActionsService_RerunWorkflowRunByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/3434/rerun", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
	})

	ctx := t.Context()
	resp, err := client.Actions.RerunWorkflowByID(ctx, "o", "r", 3434)
	if err != nil {
		t.Errorf("Actions.RerunWorkflowByID returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Actions.RerunWorkflowRunByID returned status: %v, want %v", resp.StatusCode, http.StatusCreated)
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

	ctx := t.Context()
	resp, err := client.Actions.RerunFailedJobsByID(ctx, "o", "r", 3434)
	if err != nil {
		t.Errorf("Actions.RerunFailedJobsByID returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Actions.RerunFailedJobsByID returned status: %v, want %v", resp.StatusCode, http.StatusCreated)
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

	ctx := t.Context()
	resp, err := client.Actions.RerunJobByID(ctx, "o", "r", 3434)
	if err != nil {
		t.Errorf("Actions.RerunJobByID returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Actions.RerunJobByID returned status: %v, want %v", resp.StatusCode, http.StatusCreated)
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

	ctx := t.Context()
	resp, err := client.Actions.CancelWorkflowRunByID(ctx, "o", "r", 3434)
	if !errors.As(err, new(*AcceptedError)) {
		t.Errorf("Actions.CancelWorkflowRunByID returned error: %v (want AcceptedError)", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Actions.CancelWorkflowRunByID returned status: %v, want %v", resp.StatusCode, http.StatusAccepted)
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
	tcs := []struct {
		name              string
		respectRateLimits bool
	}{
		{
			name:              "withoutRateLimits",
			respectRateLimits: false,
		},
		{
			name:              "withRateLimits",
			respectRateLimits: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)
			client.RateLimitRedirectionalEndpoints = tc.respectRateLimits

			mux.HandleFunc("/repos/o/r/actions/runs/399444496/logs", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				http.Redirect(w, r, "https://github.com/a", http.StatusFound)
			})

			ctx := t.Context()
			url, resp, err := client.Actions.GetWorkflowRunLogs(ctx, "o", "r", 399444496, 1)
			if err != nil {
				t.Errorf("Actions.GetWorkflowRunLogs returned error: %v", err)
			}
			if resp.StatusCode != http.StatusFound {
				t.Errorf("Actions.GetWorkflowRunLogs returned status: %v, want %v", resp.StatusCode, http.StatusFound)
			}
			want := "https://github.com/a"
			if url.String() != want {
				t.Errorf("Actions.GetWorkflowRunLogs returned %+v, want %+v", url, want)
			}

			const methodName = "GetWorkflowRunLogs"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Actions.GetWorkflowRunLogs(ctx, "\n", "\n", 399444496, 1)
				return err
			})
		})
	}
}

func TestActionsService_GetWorkflowRunLogs_StatusMovedPermanently_dontFollowRedirects(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name              string
		respectRateLimits bool
	}{
		{
			name:              "withoutRateLimits",
			respectRateLimits: false,
		},
		{
			name:              "withRateLimits",
			respectRateLimits: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)
			client.RateLimitRedirectionalEndpoints = tc.respectRateLimits

			mux.HandleFunc("/repos/o/r/actions/runs/399444496/logs", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				http.Redirect(w, r, "https://github.com/a", http.StatusMovedPermanently)
			})

			ctx := t.Context()
			_, resp, _ := client.Actions.GetWorkflowRunLogs(ctx, "o", "r", 399444496, 0)
			if resp.StatusCode != http.StatusMovedPermanently {
				t.Errorf("Actions.GetWorkflowJobLogs returned status: %v, want %v", resp.StatusCode, http.StatusMovedPermanently)
			}
		})
	}
}

func TestActionsService_GetWorkflowRunLogs_StatusMovedPermanently_followRedirects(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name              string
		respectRateLimits bool
	}{
		{
			name:              "withoutRateLimits",
			respectRateLimits: false,
		},
		{
			name:              "withRateLimits",
			respectRateLimits: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, mux, serverURL := setup(t)
			client.RateLimitRedirectionalEndpoints = tc.respectRateLimits

			// Mock a redirect link, which leads to an archive link
			mux.HandleFunc("/repos/o/r/actions/runs/399444496/logs", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				redirectURL, _ := url.Parse(serverURL + baseURLPath + "/redirect")
				http.Redirect(w, r, redirectURL.String(), http.StatusMovedPermanently)
			})

			mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				http.Redirect(w, r, "https://github.com/a", http.StatusFound)
			})

			ctx := t.Context()
			url, resp, err := client.Actions.GetWorkflowRunLogs(ctx, "o", "r", 399444496, 1)
			if err != nil {
				t.Errorf("Actions.GetWorkflowJobLogs returned error: %v", err)
			}

			if resp.StatusCode != http.StatusFound {
				t.Errorf("Actions.GetWorkflowJobLogs returned status: %v, want %v", resp.StatusCode, http.StatusFound)
			}

			want := "https://github.com/a"
			if url.String() != want {
				t.Errorf("Actions.GetWorkflowJobLogs returned %+v, want %+v", url, want)
			}

			const methodName = "GetWorkflowRunLogs"
			testBadOptions(t, methodName, func() (err error) {
				_, _, err = client.Actions.GetWorkflowRunLogs(ctx, "\n", "\n", 399444496, 1)
				return err
			})
		})
	}
}

func TestActionsService_GetWorkflowRunLogs_unexpectedCode(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name              string
		respectRateLimits bool
	}{
		{
			name:              "withoutRateLimits",
			respectRateLimits: false,
		},
		{
			name:              "withRateLimits",
			respectRateLimits: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, mux, serverURL := setup(t)
			client.RateLimitRedirectionalEndpoints = tc.respectRateLimits

			// Mock a redirect link, which leads to an archive link
			mux.HandleFunc("/repos/o/r/actions/runs/399444496/logs", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				redirectURL, _ := url.Parse(serverURL + baseURLPath + "/redirect")
				http.Redirect(w, r, redirectURL.String(), http.StatusMovedPermanently)
			})

			mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				w.WriteHeader(http.StatusNoContent)
			})

			ctx := t.Context()
			url, resp, err := client.Actions.GetWorkflowRunLogs(ctx, "o", "r", 399444496, 1)
			if err == nil {
				t.Fatal("Actions.GetWorkflowRunLogs should return error on unexpected code")
			}
			if !strings.Contains(err.Error(), "unexpected status code") {
				t.Error("Actions.GetWorkflowRunLogs should return unexpected status code")
			}
			if got, want := resp.Response.StatusCode, http.StatusNoContent; got != want {
				t.Errorf("Actions.GetWorkflowRunLogs return status %v, want %v", got, want)
			}
			if url != nil {
				t.Errorf("Actions.GetWorkflowRunLogs return %+v, want nil", url)
			}
		})
	}
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
	ctx := t.Context()
	runs, _, err := client.Actions.ListRepositoryWorkflowRuns(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListRepositoryWorkflowRuns returned error: %v", err)
	}

	expected := &WorkflowRuns{
		TotalCount: Ptr(2),
		WorkflowRuns: []*WorkflowRun{
			{ID: Ptr(int64(298499444)), RunNumber: Ptr(301), CreatedAt: &Timestamp{time.Date(2020, time.April, 11, 11, 14, 54, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.April, 11, 11, 14, 54, 0, time.UTC)}},
			{ID: Ptr(int64(298499445)), RunNumber: Ptr(302), CreatedAt: &Timestamp{time.Date(2020, time.April, 11, 11, 14, 54, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.April, 11, 11, 14, 54, 0, time.UTC)}},
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

	ctx := t.Context()
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

	ctx := t.Context()
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
			ID:      Ptr(int64(1)),
			NodeID:  Ptr("nid"),
			Name:    Ptr("n"),
			URL:     Ptr("u"),
			HTMLURL: Ptr("hu"),
		},
		WaitTimer:             Ptr(int64(100)),
		WaitTimerStartedAt:    &Timestamp{referenceTime},
		CurrentUserCanApprove: Ptr(false),
		Reviewers: []*RequiredReviewer{
			{
				Type: Ptr("User"),
				Reviewer: &User{
					Login: Ptr("l"),
				},
			},
			{
				Type: Ptr("Team"),
				Reviewer: &Team{
					Name: Ptr("n"),
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

	ctx := t.Context()
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
	testJSONMarshal(t, &ReviewCustomDeploymentProtectionRuleRequest{}, `{
		"comment": "",
		"environment_name": "",
		"state": ""
	}`)

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

	ctx := t.Context()
	workflowRunUsage, _, err := client.Actions.GetWorkflowRunUsageByID(ctx, "o", "r", 29679449)
	if err != nil {
		t.Errorf("Actions.GetWorkflowRunUsageByID returned error: %v", err)
	}

	want := &WorkflowRunUsage{
		Billable: &WorkflowRunBillMap{
			"UBUNTU": &WorkflowRunBill{
				TotalMS: Ptr(int64(180000)),
				Jobs:    Ptr(1),
				JobRuns: []*WorkflowRunJobRun{
					{
						JobID:      Ptr(1),
						DurationMS: Ptr(int64(60000)),
					},
				},
			},
			"MACOS": &WorkflowRunBill{
				TotalMS: Ptr(int64(240000)),
				Jobs:    Ptr(2),
				JobRuns: []*WorkflowRunJobRun{
					{
						JobID:      Ptr(2),
						DurationMS: Ptr(int64(30000)),
					},
					{
						JobID:      Ptr(3),
						DurationMS: Ptr(int64(10000)),
					},
				},
			},
			"WINDOWS": &WorkflowRunBill{
				TotalMS: Ptr(int64(300000)),
				Jobs:    Ptr(2),
			},
		},
		RunDurationMS: Ptr(int64(500000)),
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
		ID:         Ptr(int64(1)),
		Name:       Ptr("n"),
		NodeID:     Ptr("nid"),
		HeadBranch: Ptr("hb"),
		HeadSHA:    Ptr("hs"),
		Path:       Ptr("p"),
		RunNumber:  Ptr(1),
		RunAttempt: Ptr(1),
		Event:      Ptr("e"),
		Status:     Ptr("s"),
		Conclusion: Ptr("c"),
		WorkflowID: Ptr(int64(1)),
		URL:        Ptr("u"),
		HTMLURL:    Ptr("h"),
		PullRequests: []*PullRequest{
			{
				URL:    Ptr("u"),
				ID:     Ptr(int64(1)),
				Number: Ptr(1),
				Head: &PullRequestBranch{
					Ref: Ptr("r"),
					SHA: Ptr("s"),
					Repo: &Repository{
						ID:   Ptr(int64(1)),
						URL:  Ptr("s"),
						Name: Ptr("n"),
					},
				},
				Base: &PullRequestBranch{
					Ref: Ptr("r"),
					SHA: Ptr("s"),
					Repo: &Repository{
						ID:   Ptr(int64(1)),
						URL:  Ptr("u"),
						Name: Ptr("n"),
					},
				},
			},
		},
		CreatedAt:          &Timestamp{referenceTime},
		UpdatedAt:          &Timestamp{referenceTime},
		RunStartedAt:       &Timestamp{referenceTime},
		JobsURL:            Ptr("j"),
		LogsURL:            Ptr("l"),
		CheckSuiteURL:      Ptr("c"),
		ArtifactsURL:       Ptr("a"),
		CancelURL:          Ptr("c"),
		RerunURL:           Ptr("r"),
		PreviousAttemptURL: Ptr("p"),
		HeadCommit: &HeadCommit{
			Message: Ptr("m"),
			Author: &CommitAuthor{
				Name:  Ptr("n"),
				Email: Ptr("e"),
				Login: Ptr("l"),
			},
			URL:       Ptr("u"),
			Distinct:  Ptr(false),
			SHA:       Ptr("s"),
			ID:        Ptr("i"),
			TreeID:    Ptr("tid"),
			Timestamp: &Timestamp{referenceTime},
			Committer: &CommitAuthor{
				Name:  Ptr("n"),
				Email: Ptr("e"),
				Login: Ptr("l"),
			},
		},
		WorkflowURL: Ptr("w"),
		Repository: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("u"),
			Name: Ptr("n"),
		},
		HeadRepository: &Repository{
			ID:   Ptr(int64(1)),
			URL:  Ptr("u"),
			Name: Ptr("n"),
		},
		Actor: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
			URL:             Ptr("u"),
		},
		TriggeringActor: &User{
			Login:           Ptr("l2"),
			ID:              Ptr(int64(2)),
			AvatarURL:       Ptr("a2"),
			GravatarID:      Ptr("g2"),
			Name:            Ptr("n2"),
			Company:         Ptr("c2"),
			Blog:            Ptr("b2"),
			Location:        Ptr("l2"),
			Email:           Ptr("e2"),
			Hireable:        Ptr(false),
			Bio:             Ptr("b2"),
			TwitterUsername: Ptr("t2"),
			PublicRepos:     Ptr(2),
			Followers:       Ptr(2),
			Following:       Ptr(2),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
			URL:             Ptr("u2"),
		},
		ReferencedWorkflows: []*ReferencedWorkflow{
			{
				Path: Ptr("rwfp"),
				SHA:  Ptr("rwfsha"),
				Ref:  Ptr("rwfref"),
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
		TotalCount: Ptr(1),
		WorkflowRuns: []*WorkflowRun{
			{
				ID:         Ptr(int64(1)),
				Name:       Ptr("n"),
				NodeID:     Ptr("nid"),
				HeadBranch: Ptr("hb"),
				HeadSHA:    Ptr("hs"),
				RunNumber:  Ptr(1),
				RunAttempt: Ptr(1),
				Event:      Ptr("e"),
				Status:     Ptr("s"),
				Conclusion: Ptr("c"),
				WorkflowID: Ptr(int64(1)),
				URL:        Ptr("u"),
				HTMLURL:    Ptr("h"),
				PullRequests: []*PullRequest{
					{
						URL:    Ptr("u"),
						ID:     Ptr(int64(1)),
						Number: Ptr(1),
						Head: &PullRequestBranch{
							Ref: Ptr("r"),
							SHA: Ptr("s"),
							Repo: &Repository{
								ID:   Ptr(int64(1)),
								URL:  Ptr("s"),
								Name: Ptr("n"),
							},
						},
						Base: &PullRequestBranch{
							Ref: Ptr("r"),
							SHA: Ptr("s"),
							Repo: &Repository{
								ID:   Ptr(int64(1)),
								URL:  Ptr("u"),
								Name: Ptr("n"),
							},
						},
					},
				},
				CreatedAt:          &Timestamp{referenceTime},
				UpdatedAt:          &Timestamp{referenceTime},
				RunStartedAt:       &Timestamp{referenceTime},
				JobsURL:            Ptr("j"),
				LogsURL:            Ptr("l"),
				CheckSuiteURL:      Ptr("c"),
				ArtifactsURL:       Ptr("a"),
				CancelURL:          Ptr("c"),
				RerunURL:           Ptr("r"),
				PreviousAttemptURL: Ptr("p"),
				HeadCommit: &HeadCommit{
					Message: Ptr("m"),
					Author: &CommitAuthor{
						Name:  Ptr("n"),
						Email: Ptr("e"),
						Login: Ptr("l"),
					},
					URL:       Ptr("u"),
					Distinct:  Ptr(false),
					SHA:       Ptr("s"),
					ID:        Ptr("i"),
					TreeID:    Ptr("tid"),
					Timestamp: &Timestamp{referenceTime},
					Committer: &CommitAuthor{
						Name:  Ptr("n"),
						Email: Ptr("e"),
						Login: Ptr("l"),
					},
				},
				WorkflowURL: Ptr("w"),
				Repository: &Repository{
					ID:   Ptr(int64(1)),
					URL:  Ptr("u"),
					Name: Ptr("n"),
				},
				HeadRepository: &Repository{
					ID:   Ptr(int64(1)),
					URL:  Ptr("u"),
					Name: Ptr("n"),
				},
				Actor: &User{
					Login:           Ptr("l"),
					ID:              Ptr(int64(1)),
					AvatarURL:       Ptr("a"),
					GravatarID:      Ptr("g"),
					Name:            Ptr("n"),
					Company:         Ptr("c"),
					Blog:            Ptr("b"),
					Location:        Ptr("l"),
					Email:           Ptr("e"),
					Hireable:        Ptr(true),
					Bio:             Ptr("b"),
					TwitterUsername: Ptr("t"),
					PublicRepos:     Ptr(1),
					Followers:       Ptr(1),
					Following:       Ptr(1),
					CreatedAt:       &Timestamp{referenceTime},
					SuspendedAt:     &Timestamp{referenceTime},
					URL:             Ptr("u"),
				},
				TriggeringActor: &User{
					Login:           Ptr("l2"),
					ID:              Ptr(int64(2)),
					AvatarURL:       Ptr("a2"),
					GravatarID:      Ptr("g2"),
					Name:            Ptr("n2"),
					Company:         Ptr("c2"),
					Blog:            Ptr("b2"),
					Location:        Ptr("l2"),
					Email:           Ptr("e2"),
					Hireable:        Ptr(false),
					Bio:             Ptr("b2"),
					TwitterUsername: Ptr("t2"),
					PublicRepos:     Ptr(2),
					Followers:       Ptr(2),
					Following:       Ptr(2),
					CreatedAt:       &Timestamp{referenceTime},
					SuspendedAt:     &Timestamp{referenceTime},
					URL:             Ptr("u2"),
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
		TotalMS: Ptr(int64(1)),
		Jobs:    Ptr(1),
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
			TotalMS: Ptr(int64(1)),
			Jobs:    Ptr(1),
		},
		"MACOS": &WorkflowRunBill{
			TotalMS: Ptr(int64(1)),
			Jobs:    Ptr(1),
		},
		"WINDOWS": &WorkflowRunBill{
			TotalMS: Ptr(int64(1)),
			Jobs:    Ptr(1),
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
				TotalMS: Ptr(int64(1)),
				Jobs:    Ptr(1),
			},
			"MACOS": &WorkflowRunBill{
				TotalMS: Ptr(int64(1)),
				Jobs:    Ptr(1),
			},
			"WINDOWS": &WorkflowRunBill{
				TotalMS: Ptr(int64(1)),
				Jobs:    Ptr(1),
			},
		},
		RunDurationMS: Ptr(int64(1)),
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

	ctx := t.Context()
	deployments, _, err := client.Actions.PendingDeployments(ctx, "o", "r", 399444496, input)
	if err != nil {
		t.Errorf("Actions.PendingDeployments returned error: %v", err)
	}

	want := []*Deployment{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}}
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

	ctx := t.Context()
	deployments, _, err := client.Actions.GetPendingDeployments(ctx, "o", "r", 399444496)
	if err != nil {
		t.Errorf("Actions.GetPendingDeployments returned error: %v", err)
	}

	want := []*PendingDeployment{
		{
			Environment: &PendingDeploymentEnvironment{
				ID:      Ptr(int64(1)),
				NodeID:  Ptr("nid"),
				Name:    Ptr("n"),
				URL:     Ptr("u"),
				HTMLURL: Ptr("hu"),
			},
			WaitTimer:             Ptr(int64(0)),
			WaitTimerStartedAt:    &Timestamp{referenceTime},
			CurrentUserCanApprove: Ptr(false),
			Reviewers:             []*RequiredReviewer{},
		},
		{
			Environment: &PendingDeploymentEnvironment{
				ID:      Ptr(int64(2)),
				NodeID:  Ptr("nid"),
				Name:    Ptr("n"),
				URL:     Ptr("u"),
				HTMLURL: Ptr("hu"),
			},
			WaitTimer:             Ptr(int64(13)),
			WaitTimerStartedAt:    &Timestamp{referenceTime},
			CurrentUserCanApprove: Ptr(true),
			Reviewers: []*RequiredReviewer{
				{
					Type: Ptr("User"),
					Reviewer: &User{
						Login: Ptr("l"),
					},
				},
				{
					Type: Ptr("Team"),
					Reviewer: &Team{
						Name: Ptr("t"),
						Slug: Ptr("s"),
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
