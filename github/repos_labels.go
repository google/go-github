// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
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
	u := fmt.Sprintf("repos/%v/%v/labels", stringpointertostring(repo.Owner.Login), stringpointertostring(repo.Name))
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var labels []*Label
	resp, err := s.client.Do(ctx, req, &labels)
	if err != nil {
		return nil, resp, err
	}

	return labels, resp, nil
}

// GetLabel gets a single label.
//
// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/issues/#get-a-label
func (s *RepositoriesService) GetLabel(ctx context.Context, repo *Repository, name string) (*Label, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/labels/%v", stringpointertostring(repo.Owner.Login), stringpointertostring(repo.Name), name)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	label := new(Label)
	resp, err := s.client.Do(ctx, req, label)
	if err != nil {
		return nil, resp, err
	}

	return label, resp, nil
}

// CreateLabel creates a new label on the specified repository.
//
// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/issues/#create-a-label
func (s *RepositoriesService) CreateLabel(ctx context.Context, repo *Repository, label *Label) (*Label, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/labels", stringpointertostring(repo.Owner.Login), stringpointertostring(repo.Name))
	req, err := s.client.NewRequest("POST", u, label)
	if err != nil {
		return nil, nil, err
	}

	l := new(Label)
	resp, err := s.client.Do(ctx, req, l)
	if err != nil {
		return nil, resp, err
	}

	return l, resp, nil
}

// EditLabel edits a label.
//
// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/issues/#update-a-label
func (s *RepositoriesService) EditLabel(ctx context.Context, repo *Repository, name string, label *Label) (*Label, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/labels/%v", stringpointertostring(repo.Owner.Login), stringpointertostring(repo.Name), name)
	req, err := s.client.NewRequest("PATCH", u, label)
	if err != nil {
		return nil, nil, err
	}

	l := new(Label)
	resp, err := s.client.Do(ctx, req, l)
	if err != nil {
		return nil, resp, err
	}

	return l, resp, nil
}

// DeleteLabel deletes a label.
//
// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/issues/#delete-a-label
func (s *RepositoriesService) DeleteLabel(ctx context.Context, repo *Repository, name string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/labels/%v", stringpointertostring(repo.Owner.Login), stringpointertostring(repo.Name), name)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}
