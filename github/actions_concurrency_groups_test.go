// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestActionsService_ListConcurrencyGroups(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/concurrency_groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"after":    "abc",
			"per_page": "2",
		})
		fmt.Fprint(w, `{
			"total_count": 2,
			"concurrency_groups": [
				{
					"group_name": "deploy-prod",
					"group_url": "https://api.github.com/repos/o/r/actions/concurrency_groups/deploy-prod",
					"last_acquired_at": `+referenceTimeStr+`
				},
				{
					"group_name": "ci-build",
					"group_url": "https://api.github.com/repos/o/r/actions/concurrency_groups/ci-build",
					"last_acquired_at": null
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &ListConcurrencyGroupsOptions{After: "abc", PerPage: 2}
	groups, _, err := client.Actions.ListConcurrencyGroups(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListConcurrencyGroups returned error: %v", err)
	}

	want := &ConcurrencyGroupList{
		TotalCount: 2,
		ConcurrencyGroups: []*ConcurrencyGroup{
			{
				GroupName:      "deploy-prod",
				GroupURL:       "https://api.github.com/repos/o/r/actions/concurrency_groups/deploy-prod",
				LastAcquiredAt: &referenceTimestamp,
			},
			{
				GroupName: "ci-build",
				GroupURL:  "https://api.github.com/repos/o/r/actions/concurrency_groups/ci-build",
			},
		},
	}
	if !cmp.Equal(groups, want) {
		t.Errorf("Actions.ListConcurrencyGroups returned %+v, want %+v", groups, want)
	}

	const methodName = "ListConcurrencyGroups"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListConcurrencyGroups(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListConcurrencyGroups(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetConcurrencyGroup(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/concurrency_groups/deploy-prod", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"ahead_of_job": "2",
		})
		fmt.Fprint(w, `{
			"group_name": "deploy-prod",
			"group_url": "https://api.github.com/repos/o/r/actions/concurrency_groups/deploy-prod",
			"total_count": 2,
			"group_members": [
				{
					"run_id": 1,
					"run_name": "Deploy",
					"run_url": "https://api.github.com/repos/o/r/actions/runs/1",
					"run_html_url": "https://github.com/o/r/actions/runs/1",
					"status": "in_progress"
				},
				{
					"run_id": 3,
					"run_name": "Deploy hotfix",
					"run_url": "https://api.github.com/repos/o/r/actions/runs/3",
					"run_html_url": "https://github.com/o/r/actions/runs/3",
					"job_id": 2,
					"job_name": "deploy",
					"job_url": "https://api.github.com/repos/o/r/actions/jobs/2",
					"job_html_url": "https://github.com/o/r/actions/runs/3/job/2",
					"status": "pending"
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &GetConcurrencyGroupOptions{AheadOfJob: 2}
	group, _, err := client.Actions.GetConcurrencyGroup(ctx, "o", "r", "deploy-prod", opts)
	if err != nil {
		t.Errorf("Actions.GetConcurrencyGroup returned error: %v", err)
	}

	want := &ConcurrencyGroup{
		GroupName:  "deploy-prod",
		GroupURL:   "https://api.github.com/repos/o/r/actions/concurrency_groups/deploy-prod",
		TotalCount: new(2),
		GroupMembers: []*ConcurrencyGroupMember{
			{
				RunID:      1,
				RunName:    "Deploy",
				RunURL:     new("https://api.github.com/repos/o/r/actions/runs/1"),
				RunHTMLURL: new("https://github.com/o/r/actions/runs/1"),
				Status:     "in_progress",
			},
			{
				RunID:      3,
				RunName:    "Deploy hotfix",
				RunURL:     new("https://api.github.com/repos/o/r/actions/runs/3"),
				RunHTMLURL: new("https://github.com/o/r/actions/runs/3"),
				JobID:      new(int64(2)),
				JobName:    new("deploy"),
				JobURL:     new("https://api.github.com/repos/o/r/actions/jobs/2"),
				JobHTMLURL: new("https://github.com/o/r/actions/runs/3/job/2"),
				Status:     "pending",
			},
		},
	}
	if !cmp.Equal(group, want) {
		t.Errorf("Actions.GetConcurrencyGroup returned %+v, want %+v", group, want)
	}

	const methodName = "GetConcurrencyGroup"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetConcurrencyGroup(ctx, "\n", "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetConcurrencyGroup(ctx, "o", "r", "deploy-prod", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetConcurrencyGroup_EscapeGroupName(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		groupName        string
		escapedGroupName string
	}{
		{groupName: "deploy-prod", escapedGroupName: "deploy-prod"},
		{groupName: "CI-refs/heads/main", escapedGroupName: "CI-refs%2Fheads%2Fmain"},
	} {
		t.Run(tt.groupName, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				if got, want := r.URL.EscapedPath(), "/repos/o/r/actions/concurrency_groups/"+tt.escapedGroupName; got != want {
					t.Errorf("Request path = %q, want %q", got, want)
				}
				if got, want := r.URL.Path, "/repos/o/r/actions/concurrency_groups/"+tt.groupName; got != want {
					t.Errorf("Decoded request path = %q, want %q", got, want)
				}
				fmt.Fprint(w, `{}`)
			})

			ctx := t.Context()
			_, _, err := client.Actions.GetConcurrencyGroup(ctx, "o", "r", tt.groupName, nil)
			if err != nil {
				t.Fatalf("Actions.GetConcurrencyGroup returned error: %v", err)
			}
		})
	}
}

func TestActionsService_ListWorkflowRunConcurrencyGroups(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/runs/1/concurrency_groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"before":   "xyz",
			"per_page": "2",
		})
		fmt.Fprint(w, `{
			"total_count": 2,
			"concurrency_groups": [
				{
					"group_name": "deploy-prod",
					"group_url": "https://api.github.com/repos/o/r/actions/concurrency_groups/deploy-prod",
					"group_members": [
						{
							"run_id": 1,
							"run_name": "Deploy",
							"run_url": "https://api.github.com/repos/o/r/actions/runs/1",
							"run_html_url": "https://github.com/o/r/actions/runs/1",
							"status": "in_progress",
							"position": 0,
							"position_url": "https://api.github.com/repos/o/r/actions/concurrency_groups/deploy-prod?ahead_of_run=1"
						}
					]
				},
				{
					"group_name": "ci-build",
					"group_url": "https://api.github.com/repos/o/r/actions/concurrency_groups/ci-build",
					"group_members": []
				}
			]
		}`)
	})

	ctx := t.Context()
	opts := &ListWorkflowRunConcurrencyGroupsOptions{Before: "xyz", PerPage: 2}
	groups, _, err := client.Actions.ListWorkflowRunConcurrencyGroups(ctx, "o", "r", 1, opts)
	if err != nil {
		t.Errorf("Actions.ListWorkflowRunConcurrencyGroups returned error: %v", err)
	}

	want := &ConcurrencyGroupList{
		TotalCount: 2,
		ConcurrencyGroups: []*ConcurrencyGroup{
			{
				GroupName: "deploy-prod",
				GroupURL:  "https://api.github.com/repos/o/r/actions/concurrency_groups/deploy-prod",
				GroupMembers: []*ConcurrencyGroupMember{
					{
						RunID:       1,
						RunName:     "Deploy",
						RunURL:      new("https://api.github.com/repos/o/r/actions/runs/1"),
						RunHTMLURL:  new("https://github.com/o/r/actions/runs/1"),
						Status:      "in_progress",
						Position:    new(0),
						PositionURL: new("https://api.github.com/repos/o/r/actions/concurrency_groups/deploy-prod?ahead_of_run=1"),
					},
				},
			},
			{
				GroupName:    "ci-build",
				GroupURL:     "https://api.github.com/repos/o/r/actions/concurrency_groups/ci-build",
				GroupMembers: []*ConcurrencyGroupMember{},
			},
		},
	}
	if !cmp.Equal(groups, want) {
		t.Errorf("Actions.ListWorkflowRunConcurrencyGroups returned %+v, want %+v", groups, want)
	}

	const methodName = "ListWorkflowRunConcurrencyGroups"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListWorkflowRunConcurrencyGroups(ctx, "\n", "\n", -1, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListWorkflowRunConcurrencyGroups(ctx, "o", "r", 1, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
