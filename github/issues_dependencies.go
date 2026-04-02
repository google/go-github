// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// IssueDependencyRequest represents a request to add a dependency to an issue.
type IssueDependencyRequest struct {
	IssueID *int64 `json:"issue_id,omitempty"`
}

// ListBlockedBy lists the dependencies that block the specified issue.
//
// GitHub API docs: https://docs.github.com/rest/issues/issue-dependencies#list-dependencies-an-issue-is-blocked-by
//
//meta:operation GET /repos/{owner}/{repo}/issues/{issue_number}/dependencies/blocked_by
func (s *IssuesService) ListBlockedBy(ctx context.Context, owner, repo string, number int, opts *ListOptions) ([]*Issue, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/dependencies/blocked_by", owner, repo, number)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var issues []*Issue
	resp, err := s.client.Do(ctx, req, &issues)
	if err != nil {
		return nil, resp, err
	}

	return issues, resp, nil
}

// AddBlockedBy adds a "blocked by" dependency to the specified issue.
//
// GitHub API docs: https://docs.github.com/rest/issues/issue-dependencies#add-a-dependency-an-issue-is-blocked-by
//
//meta:operation POST /repos/{owner}/{repo}/issues/{issue_number}/dependencies/blocked_by
func (s *IssuesService) AddBlockedBy(ctx context.Context, owner, repo string, number int, issueDepReq *IssueDependencyRequest) (*Issue, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/dependencies/blocked_by", owner, repo, number)
	req, err := s.client.NewRequest("POST", u, issueDepReq)
	if err != nil {
		return nil, nil, err
	}

	var issue *Issue
	resp, err := s.client.Do(ctx, req, &issue)
	if err != nil {
		return nil, resp, err
	}

	return issue, resp, nil
}

// RemoveBlockedBy removes a "blocked by" dependency from the specified issue.
//
// GitHub API docs: https://docs.github.com/rest/issues/issue-dependencies#remove-dependency-an-issue-is-blocked-by
//
//meta:operation DELETE /repos/{owner}/{repo}/issues/{issue_number}/dependencies/blocked_by/{issue_id}
func (s *IssuesService) RemoveBlockedBy(ctx context.Context, owner, repo string, number int, issueID int64) (*Issue, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/dependencies/blocked_by/%v", owner, repo, number, issueID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var issue *Issue
	resp, err := s.client.Do(ctx, req, &issue)
	if err != nil {
		return nil, resp, err
	}

	return issue, resp, nil
}

// ListBlocking lists the issues that the specified issue is blocking.
//
// GitHub API docs: https://docs.github.com/rest/issues/issue-dependencies#list-dependencies-an-issue-is-blocking
//
//meta:operation GET /repos/{owner}/{repo}/issues/{issue_number}/dependencies/blocking
func (s *IssuesService) ListBlocking(ctx context.Context, owner, repo string, number int, opts *ListOptions) ([]*Issue, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%v/dependencies/blocking", owner, repo, number)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var issues []*Issue
	resp, err := s.client.Do(ctx, req, &issues)
	if err != nil {
		return nil, resp, err
	}

	return issues, resp, nil
}
