// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// PullRequestListStacksOptions specifies the optional parameters to the
// PullRequestsService.ListStacks method.
type PullRequestListStacksOptions struct {
	// PullRequest filters stacks to the stack containing this pull request number.
	PullRequest int `url:"pull_request,omitempty"`

	ListOptions
}

// CreatePullRequestStackRequest represents a request to create a pull request stack.
type CreatePullRequestStackRequest struct {
	// PullRequests is an ordered list of pull request numbers from the bottom of the stack to the top.
	PullRequests []int `json:"pull_requests"`
}

// AddPullRequestsToStackRequest represents a request to append pull requests to a stack.
type AddPullRequestsToStackRequest struct {
	// PullRequests is an ordered list of pull request numbers to append from the current top upward.
	PullRequests []int `json:"pull_requests"`
}

// PullRequestStackDetails represents a pull request stack returned by the
// stacked pull request endpoints.
type PullRequestStackDetails struct {
	// ID is the ID of the stack.
	ID *int64 `json:"id,omitempty"`
	// Number is the number of the stack.
	Number *int `json:"number,omitempty"`
	// NodeID is the global node ID of the stack.
	NodeID *string `json:"node_id,omitempty"`
	// URL is the API URL of the stack.
	URL *string `json:"url,omitempty"`
	// Base is the branch the entire stack ultimately targets.
	Base *PullRequestStackBase `json:"base"`
	// Open reports whether the stack contains any open pull requests.
	Open *bool `json:"open,omitempty"`
	// CreatedAt is the time the stack was created.
	CreatedAt *Timestamp `json:"created_at,omitempty"`
	// PullRequests contains the pull requests in the stack, from bottom to top.
	PullRequests []*PullRequest `json:"pull_requests,omitempty"`
}

// ListStacks lists pull request stacks in a repository.
//
// GitHub API docs: https://docs.github.com/rest/pulls/stacks?apiVersion=2022-11-28#list-pull-request-stacks
//
//meta:operation GET /repos/{owner}/{repo}/stacks
func (s *PullRequestsService) ListStacks(ctx context.Context, owner, repo string, opts *PullRequestListStacksOptions) ([]*PullRequestStackDetails, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/stacks", owner, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var stacks []*PullRequestStackDetails
	resp, err := s.client.Do(req, &stacks)
	if err != nil {
		return nil, resp, err
	}

	return stacks, resp, nil
}

// CreateStack creates a pull request stack from an ordered list of pull request numbers.
//
// GitHub API docs: https://docs.github.com/rest/pulls/stacks?apiVersion=2022-11-28#create-a-pull-request-stack
//
//meta:operation POST /repos/{owner}/{repo}/stacks
func (s *PullRequestsService) CreateStack(ctx context.Context, owner, repo string, body CreatePullRequestStackRequest) (*PullRequestStackDetails, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/stacks", owner, repo)
	req, err := s.client.NewRequest(ctx, "POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var stack *PullRequestStackDetails
	resp, err := s.client.Do(req, &stack)
	if err != nil {
		return nil, resp, err
	}

	return stack, resp, nil
}

// GetStack gets a pull request stack by its stack number.
//
// GitHub API docs: https://docs.github.com/rest/pulls/stacks?apiVersion=2022-11-28#get-a-pull-request-stack
//
//meta:operation GET /repos/{owner}/{repo}/stacks/{stack_number}
func (s *PullRequestsService) GetStack(ctx context.Context, owner, repo string, stackNumber int) (*PullRequestStackDetails, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/stacks/%v", owner, repo, stackNumber)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var stack *PullRequestStackDetails
	resp, err := s.client.Do(req, &stack)
	if err != nil {
		return nil, resp, err
	}

	return stack, resp, nil
}

// AddToStack appends pull requests to a pull request stack.
//
// GitHub API docs: https://docs.github.com/rest/pulls/stacks?apiVersion=2022-11-28#add-pull-requests-to-a-pull-request-stack
//
//meta:operation POST /repos/{owner}/{repo}/stacks/{stack_number}/add
func (s *PullRequestsService) AddToStack(ctx context.Context, owner, repo string, stackNumber int, body AddPullRequestsToStackRequest) (*PullRequestStackDetails, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/stacks/%v/add", owner, repo, stackNumber)
	req, err := s.client.NewRequest(ctx, "POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var stack *PullRequestStackDetails
	resp, err := s.client.Do(req, &stack)
	if err != nil {
		return nil, resp, err
	}

	return stack, resp, nil
}

// Unstack removes the unmerged pull requests from a pull request stack. It
// returns nil when no pull requests remain and the stack is dissolved.
//
// GitHub API docs: https://docs.github.com/rest/pulls/stacks?apiVersion=2022-11-28#remove-pull-requests-from-a-pull-request-stack
//
//meta:operation POST /repos/{owner}/{repo}/stacks/{stack_number}/unstack
func (s *PullRequestsService) Unstack(ctx context.Context, owner, repo string, stackNumber int) (*PullRequestStackDetails, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/stacks/%v/unstack", owner, repo, stackNumber)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var stack *PullRequestStackDetails
	resp, err := s.client.Do(req, &stack)
	if err != nil {
		return nil, resp, err
	}

	return stack, resp, nil
}
