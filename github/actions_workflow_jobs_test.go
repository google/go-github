// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestActionsService_ListWorkflowJobs(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/runs/29679449/jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"jobs":[{"id":399444496,"run_id":29679449,"started_at":"2019-01-02T15:04:05Z","completed_at":"2020-01-02T15:04:05Z"},{"id":399444497,"run_id":29679449,"started_at":"2019-01-02T15:04:05Z","completed_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	jobs, _, err := client.Actions.ListWorkflowJobs(context.Background(), "o", "r", 29679449, opts)
	if err != nil {
		t.Errorf("Actions.ListWorkflowJobs returned error: %v", err)
	}

	want := &Jobs{
		TotalCount: 4,
		Jobs: []*Job{
			{ID: 399444496, RunID: 29679449, StartedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, CompletedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
			{ID: 399444497, RunID: 29679449, StartedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, CompletedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
		},
	}
	if !reflect.DeepEqual(jobs, want) {
		t.Errorf("Actions.ListWorkflowJobs returned %+v, want %+v", jobs, want)
	}
}

func TestActionsService_GetWorkflowJobByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/jobs/399444496", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":399444496,"started_at":"2019-01-02T15:04:05Z","completed_at":"2020-01-02T15:04:05Z"}`)
	})

	job, _, err := client.Actions.GetWorkflowJobByID(context.Background(), "o", "r", 399444496)
	if err != nil {
		t.Errorf("Actions.GetWorkflowJobByID returned error: %v", err)
	}

	want := &Job{
		ID:          399444496,
		StartedAt:   Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
		CompletedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
	}
	if !reflect.DeepEqual(job, want) {
		t.Errorf("Actions.GetWorkflowJobByID returned %+v, want %+v", job, want)
	}
}

func TestActionsService_ListWorkflowJobLogs(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/jobs/399444496/logs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `https://pipelines.actions.githubusercontent.com/ab1f3cCFPB34Nd6imvFxpGZH5hNlDp2wijMwl2gDoO0bcrrlJj/_apis/pipelines/1/jobs/19/signedlogcontent?urlExpires=2020-01-22T22%3A44%3A54.1389777Z&urlSigningMethod=HMACV1&urlSignature=2TUDfIg4fm36OJmfPy6km5QD5DLCOkBVzvhWZM8B%2BUY%3D`)
	})

	logFileURL, _, err := client.Actions.ListWorkflowJobLogs(context.Background(), "o", "r", 399444496)
	if err != nil {
		t.Errorf("Actions.ListWorkflowJobLogs returned error: %v", err)
	}

	want := "https://pipelines.actions.githubusercontent.com/ab1f3cCFPB34Nd6imvFxpGZH5hNlDp2wijMwl2gDoO0bcrrlJj/_apis/pipelines/1/jobs/19/signedlogcontent?urlExpires=2020-01-22T22%3A44%3A54.1389777Z&urlSigningMethod=HMACV1&urlSignature=2TUDfIg4fm36OJmfPy6km5QD5DLCOkBVzvhWZM8B%2BUY%3D"
	if !reflect.DeepEqual(logFileURL, want) {
		t.Errorf("Actions.GetWorkflowByFileName returned %+v, want %+v", logFileURL, want)
	}
}
