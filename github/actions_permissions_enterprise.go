// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ActionsEnabledOnEnterpriseOrgs represents all the repositories in an enterprise for which Actions is enabled.
type ActionsEnabledOnEnterpriseOrgs struct {
	TotalCount    int             `json:"total_count"`
	Organizations []*Organization `json:"organizations"`
}

// ActionsPermissionsEnterprise represents a policy for allowed actions in an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/actions/permissions
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
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/actions/permissions#get-github-actions-permissions-for-an-enterprise
func (s *ActionsService) GetEnterpriseActionsPermissions(ctx context.Context, enterprise string) (*ActionsPermissionsEnterprise, *Response, error) {
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
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/actions/permissions#set-github-actions-permissions-for-an-enterprise
func (s *ActionsService) EditEnterpriseActionsPermissions(ctx context.Context, enterprise string, actionsPermissionsEnterprise ActionsPermissionsEnterprise) (*ActionsPermissionsEnterprise, *Response, error) {
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

// ListEnabledOrgsInEnterprise lists the selected organizations that are enabled for GitHub Actions in an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/actions/permissions#list-selected-organizations-enabled-for-github-actions-in-an-enterprise
func (s *ActionsService) ListEnabledOrgsInEnterprise(ctx context.Context, owner string, opts *ListOptions) (*ActionsEnabledOnEnterpriseOrgs, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/permissions/organizations", owner)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	orgs := &ActionsEnabledOnEnterpriseOrgs{}
	resp, err := s.client.Do(ctx, req, orgs)
	if err != nil {
		return nil, resp, err
	}

	return orgs, resp, nil
}

// SetEnabledOrgsInEnterprise replaces the list of selected organizations that are enabled for GitHub Actions in an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/actions/permissions#set-selected-organizations-enabled-for-github-actions-in-an-enterprise
func (s *ActionsService) SetEnabledOrgsInEnterprise(ctx context.Context, owner string, organizationIDs []int64) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/permissions/organizations", owner)

	req, err := s.client.NewRequest("PUT", u, struct {
		IDs []int64 `json:"selected_organization_ids"`
	}{IDs: organizationIDs})
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// AddEnabledOrgInEnterprise adds an organization to the list of selected organizations that are enabled for GitHub Actions in an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/actions/permissions#enable-a-selected-organization-for-github-actions-in-an-enterprise
func (s *ActionsService) AddEnabledOrgInEnterprise(ctx context.Context, owner string, organizationID int64) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/permissions/organizations/%v", owner, organizationID)

	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveEnabledOrgInEnterprise removes an organization from the list of selected organizations that are enabled for GitHub Actions in an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/actions/permissions#disable-a-selected-organization-for-github-actions-in-an-enterprise
func (s *ActionsService) RemoveEnabledOrgInEnterprise(ctx context.Context, owner string, organizationID int64) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/permissions/organizations/%v", owner, organizationID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
