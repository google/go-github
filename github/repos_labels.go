// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

func stringpointertostring(s *string) string {
	if s != nil {
		strPointerValue := *s
		return strPointerValue
	}
	return ""
}

// ListLabels lists all labels for a repository.
//
// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/issues/#list-labels-for-a-repository
func (s *RepositoriesService) ListLabels(ctx context.Context, repo *Repository, opts *ListOptions) ([]*Label, *Response, error) {
	convertedService := &IssuesService{
		client: s.client,
	}
	return convertedService.ListLabels(
		ctx, stringpointertostring(repo.Owner.Login), stringpointertostring(repo.Name), opts,
	)
}

// GetLabel gets a single label.
//
// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/issues/#get-a-label
func (s *RepositoriesService) GetLabel(ctx context.Context, repo *Repository, name string) (*Label, *Response, error) {
	convertedService := &IssuesService{
		client: s.client,
	}
	return convertedService.GetLabel(
		ctx, stringpointertostring(repo.Owner.Login), stringpointertostring(repo.Name), name,
	)
}

// CreateLabel creates a new label on the specified repository.
//
// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/issues/#create-a-label
func (s *RepositoriesService) CreateLabel(ctx context.Context, repo *Repository, label *Label) (*Label, *Response, error) {
	convertedService := &IssuesService{
		client: s.client,
	}
	return convertedService.CreateLabel(
		ctx, stringpointertostring(repo.Owner.Login), stringpointertostring(repo.Name), label,
	)

}

// EditLabel edits a label.
//
// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/issues/#update-a-label
func (s *RepositoriesService) EditLabel(ctx context.Context, repo *Repository, name string, label *Label) (*Label, *Response, error) {
	convertedService := &IssuesService{
		client: s.client,
	}
	return convertedService.EditLabel(
		ctx, stringpointertostring(repo.Owner.Login), stringpointertostring(repo.Name), name, label,
	)
}

// DeleteLabel deletes a label.
//
// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/issues/#delete-a-label
func (s *RepositoriesService) DeleteLabel(ctx context.Context, repo *Repository, name string) (*Response, error) {
	convertedService := &IssuesService{
		client: s.client,
	}
	return convertedService.DeleteLabel(
		ctx, stringpointertostring(repo.Owner.Login), stringpointertostring(repo.Name), name,
	)
}
