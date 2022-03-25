// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// GetActionsAllowed gets the actions that are allowed in a repository.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/actions#get-allowed-actions-and-workflows-for-a-repository
func (s *RepositoriesService) GetActionsAllowed(ctx context.Context, org, repo string) (*ActionsAllowed, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/permissions/selected-actions", org, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	actionsAllowed := new(ActionsAllowed)
	resp, err := s.client.Do(ctx, req, actionsAllowed)
	if err != nil {
		return nil, resp, err
	}

	return actionsAllowed, resp, nil
}

// EditActionsAllowed sets the actions that are allowed in a repository.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/actions#set-allowed-actions-and-workflows-for-a-repository
func (s *RepositoriesService) EditActionsAllowed(ctx context.Context, org, repo string, actionsAllowed ActionsAllowed) (*ActionsAllowed, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/permissions/selected-actions", org, repo)
	req, err := s.client.NewRequest("PUT", u, actionsAllowed)
	if err != nil {
		return nil, nil, err
	}

	p := new(ActionsAllowed)
	resp, err := s.client.Do(ctx, req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}
