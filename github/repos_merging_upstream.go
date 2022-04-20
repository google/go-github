// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// RepositoryMergeUpstreamRequest represents a request to sync a branch of
// a forked repository to keep it up-to-date with the upstream repository.
type RepositoryMergeUpstreamRequest struct {
	Branch *string `json:"branch,omitempty"`
}

// MergeUpstreamResult represents the result of syncing a branch of
// a forked repository with the upstream repository.
type MergeUpstreamResult struct {
	Message    *string `json:"message,omitempty"`
	MergeType  *string `json:"merge_type,omitempty"`
	BaseBranch *string `json:"base_branch,omitempty"`
}

// MergeUpstream syncs a branch of a forked repository to keep it up-to-date
// with the upstream repository.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/branches#sync-a-fork-branch-with-the-upstream-repository
func (s *RepositoriesService) MergeUpstream(ctx context.Context, owner, repo string, request *RepositoryMergeUpstreamRequest) (*MergeUpstreamResult, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/merge-upstream", owner, repo)
	req, err := s.client.NewRequest("POST", u, request)
	if err != nil {
		return nil, nil, err
	}

	result := new(MergeUpstreamResult)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
