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

// ListStacks lists pull request stacks in a repository.
//
// GitHub API docs: https://docs.github.com/rest/pulls/stacks?apiVersion=2022-11-28#list-pull-request-stacks
//
//meta:operation GET /repos/{owner}/{repo}/stacks
func (s *PullRequestsService) ListStacks(ctx context.Context, owner, repo string, opts *PullRequestListStacksOptions) ([]*PullRequestStack, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/stacks", owner, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var stacks []*PullRequestStack
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
func (s *PullRequestsService) CreateStack(ctx context.Context, owner, repo string, body CreatePullRequestStackRequest) (*PullRequestStack, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/stacks", owner, repo)
	req, err := s.client.NewRequest(ctx, "POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var stack *PullRequestStack
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
func (s *PullRequestsService) GetStack(ctx context.Context, owner, repo string, stackNumber int) (*PullRequestStack, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/stacks/%v", owner, repo, stackNumber)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var stack *PullRequestStack
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
func (s *PullRequestsService) AddToStack(ctx context.Context, owner, repo string, stackNumber int, body AddPullRequestsToStackRequest) (*PullRequestStack, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/stacks/%v/add", owner, repo, stackNumber)
	req, err := s.client.NewRequest(ctx, "POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var stack *PullRequestStack
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
func (s *PullRequestsService) Unstack(ctx context.Context, owner, repo string, stackNumber int) (*PullRequestStack, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/stacks/%v/unstack", owner, repo, stackNumber)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var stack *PullRequestStack
	resp, err := s.client.Do(req, &stack)
	if err != nil {
		return nil, resp, err
	}

	return stack, resp, nil
}
