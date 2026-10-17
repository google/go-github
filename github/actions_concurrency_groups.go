// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/url"
)

// ConcurrencyGroup represents a GitHub Actions concurrency group.
type ConcurrencyGroup struct {
	GroupName      string                    `json:"group_name"`
	GroupURL       string                    `json:"group_url"`
	TotalCount     *int                      `json:"total_count,omitempty"`
	LastAcquiredAt *Timestamp                `json:"last_acquired_at,omitempty"`
	GroupMembers   []*ConcurrencyGroupMember `json:"group_members,omitzero"`
}

// ConcurrencyGroupMember represents a workflow run or job that is holding or waiting for a concurrency group.
type ConcurrencyGroupMember struct {
	RunID       int64   `json:"run_id"`
	RunName     string  `json:"run_name"`
	RunURL      *string `json:"run_url,omitempty"`
	RunHTMLURL  *string `json:"run_html_url,omitempty"`
	JobID       *int64  `json:"job_id,omitempty"`
	JobName     *string `json:"job_name,omitempty"`
	JobURL      *string `json:"job_url,omitempty"`
	JobHTMLURL  *string `json:"job_html_url,omitempty"`
	Status      string  `json:"status"`
	Position    *int    `json:"position,omitempty"`
	PositionURL *string `json:"position_url,omitempty"`
}

// ConcurrencyGroupList represents a list of GitHub Actions concurrency groups.
type ConcurrencyGroupList struct {
	TotalCount        int                 `json:"total_count"`
	ConcurrencyGroups []*ConcurrencyGroup `json:"concurrency_groups"`
}

// ListConcurrencyGroupsOptions specifies the optional parameters to the
// ActionsService.ListConcurrencyGroups method.
type ListConcurrencyGroupsOptions struct {
	// A cursor, as given in the Link header. If specified, the query only searches for concurrency groups after this cursor.
	After string `url:"after,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}

// GetConcurrencyGroupOptions specifies the optional parameters to the
// ActionsService.GetConcurrencyGroup method.
type GetConcurrencyGroupOptions struct {
	// AheadOfRun filters the results to the items ahead of this workflow run ID, plus the run itself.
	AheadOfRun int64 `url:"ahead_of_run,omitempty"`

	// AheadOfJob filters the results to the items ahead of this job ID, plus the job itself.
	AheadOfJob int64 `url:"ahead_of_job,omitempty"`
}

// ListWorkflowRunConcurrencyGroupsOptions specifies the optional parameters to the
// ActionsService.ListWorkflowRunConcurrencyGroups method.
type ListWorkflowRunConcurrencyGroupsOptions struct {
	// A cursor, as given in the Link header. If specified, the query only searches for concurrency groups before this cursor.
	Before string `url:"before,omitempty"`

	// A cursor, as given in the Link header. If specified, the query only searches for concurrency groups after this cursor.
	After string `url:"after,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}

// ListConcurrencyGroups lists the active concurrency groups for a repository.
//
// GitHub API docs: https://docs.github.com/rest/actions/concurrency-groups?apiVersion=2022-11-28#list-concurrency-groups-for-a-repository
//
//meta:operation GET /repos/{owner}/{repo}/actions/concurrency_groups
func (s *ActionsService) ListConcurrencyGroups(ctx context.Context, owner, repo string, opts *ListConcurrencyGroupsOptions) (*ConcurrencyGroupList, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/concurrency_groups", owner, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var groups *ConcurrencyGroupList
	resp, err := s.client.Do(req, &groups)
	if err != nil {
		return nil, resp, err
	}

	return groups, resp, nil
}

// GetConcurrencyGroup gets a concurrency group for a repository, including the items in its queue.
//
// Note: the concurrency group name is URL path escaped for you. See: https://pkg.go.dev/net/url#PathEscape .
//
// GitHub API docs: https://docs.github.com/rest/actions/concurrency-groups?apiVersion=2022-11-28#get-a-concurrency-group-for-a-repository
//
//meta:operation GET /repos/{owner}/{repo}/actions/concurrency_groups/{concurrency_group_name}
func (s *ActionsService) GetConcurrencyGroup(ctx context.Context, owner, repo, groupName string, opts *GetConcurrencyGroupOptions) (*ConcurrencyGroup, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/concurrency_groups/%v", owner, repo, url.PathEscape(groupName))
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var group *ConcurrencyGroup
	resp, err := s.client.Do(req, &group)
	if err != nil {
		return nil, resp, err
	}

	return group, resp, nil
}

// ListWorkflowRunConcurrencyGroups lists the concurrency groups associated with a workflow run or its jobs.
//
// GitHub API docs: https://docs.github.com/rest/actions/concurrency-groups?apiVersion=2022-11-28#list-concurrency-groups-for-a-workflow-run
//
//meta:operation GET /repos/{owner}/{repo}/actions/runs/{run_id}/concurrency_groups
func (s *ActionsService) ListWorkflowRunConcurrencyGroups(ctx context.Context, owner, repo string, runID int64, opts *ListWorkflowRunConcurrencyGroupsOptions) (*ConcurrencyGroupList, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/runs/%v/concurrency_groups", owner, repo, runID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var groups *ConcurrencyGroupList
	resp, err := s.client.Do(req, &groups)
	if err != nil {
		return nil, resp, err
	}

	return groups, resp, nil
}
