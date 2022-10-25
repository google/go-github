// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// GetActionsAllowed gets the actions that are allowed in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/permissions#get-allowed-actions-and-reusable-workflows-for-an-enterprise
func (s *EnterpriseService) GetActionsAllowed(ctx context.Context, enterprise string) (*ActionsAllowed, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/permissions/selected-actions", enterprise)

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

// EditActionsAllowed sets the actions that are allowed in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/permissions#set-allowed-actions-and-reusable-workflows-for-an-enterprise
func (s *EnterpriseService) EditActionsAllowed(ctx context.Context, enterprise string, actionsAllowed ActionsAllowed) (*ActionsAllowed, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/permissions/selected-actions", enterprise)
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
