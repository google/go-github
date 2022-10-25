// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ActionsPermissionsEnterprise represents a policy for allowed actions in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/permissions
type ActionsPermissionsEnterprise struct {
	EnabledOrganizations *string `json:"enabled_organizations,omitempty"`
	AllowedActions       *string `json:"allowed_actions,omitempty"`
	SelectedActionsURL   *string `json:"selected_actions_url,omitempty"`
}

func (a ActionsPermissionsEnterprise) String() string {
	return Stringify(a)
}

// GetActionsPermissions gets the GitHub Actions permissions policy for an enterprise.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/permissions#get-github-actions-permissions-for-an-enterprise
func (s *EnterpriseService) GetActionsPermissions(ctx context.Context, enterprise string) (*ActionsPermissionsEnterprise, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/permissions", enterprise)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	permissions := new(ActionsPermissionsEnterprise)
	resp, err := s.client.Do(ctx, req, permissions)
	if err != nil {
		return nil, resp, err
	}

	return permissions, resp, nil
}

// EditActionsPermissions sets the permissions policy in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/permissions#set-github-actions-permissions-for-an-enterprise
func (s *EnterpriseService) EditActionsPermissions(ctx context.Context, enterprise string, actionsPermissionsEnterprise ActionsPermissionsEnterprise) (*ActionsPermissionsEnterprise, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/permissions", enterprise)
	req, err := s.client.NewRequest("PUT", u, actionsPermissionsEnterprise)
	if err != nil {
		return nil, nil, err
	}

	p := new(ActionsPermissionsEnterprise)
	resp, err := s.client.Do(ctx, req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}
