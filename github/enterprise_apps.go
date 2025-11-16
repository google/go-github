// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// EnterpriseAppsService handles communication with the enterprise apps related
// methods of the GitHub API.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/enterprise-admin#apps
type EnterpriseAppsService service

// EnterpriseInstallationRepositoriesOptions specifies the parameters for
// EnterpriseAppsService.AddRepositoriesToInstallation and
// EnterpriseAppsService.RemoveRepositoriesFromInstallation.
type EnterpriseInstallationRepositoriesOptions struct {
	SelectedRepositoryIDs []int64 `json:"selected_repository_ids"`
}

// EnterpriseInstallationRepositoriesToggleOptions specifies the parameters for
// EnterpriseAppsService.ToggleInstallationRepositories.
type EnterpriseInstallationRepositoriesToggleOptions struct {
	RepositorySelection   *string `json:"repository_selection,omitempty"` // Can be "all" or "selected"
	SelectedRepositoryIDs []int64 `json:"selected_repository_ids,omitempty"`
}

// ListRepositoriesForOrgInstallation lists the repositories that an enterprise app installation
// has access to on an organization.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/organization-installations#get-the-repositories-accessible-to-a-given-github-app-installation
//
//meta:operation GET /enterprises/{enterprise}/apps/organizations/{org}/installations/{installation_id}/repositories
func (s *EnterpriseAppsService) ListRepositoriesForOrgInstallation(ctx context.Context, enterprise, org string, installationID int64, opts *ListOptions) (*ListRepositories, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/apps/organizations/%v/installations/%v/repositories", enterprise, org, installationID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r *ListRepositories
	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// ToggleInstallationRepositories changes a GitHub App installation's repository access
// between all repositories and a selected set.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/organization-installations#toggle-installation-repository-access-between-selected-and-all-repositories
//
//meta:operation PATCH /enterprises/{enterprise}/apps/organizations/{org}/installations/{installation_id}/repositories
func (s *EnterpriseAppsService) ToggleInstallationRepositories(ctx context.Context, enterprise, org string, installationID int64, opts *EnterpriseInstallationRepositoriesToggleOptions) (*ListRepositories, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/apps/organizations/%v/installations/%v/repositories", enterprise, org, installationID)
	req, err := s.client.NewRequest("PATCH", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var r *ListRepositories
	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// AddRepositoriesToInstallation grants repository access for a GitHub App installation.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/organization-installations#grant-repository-access-to-an-organization-installation
//
//meta:operation PATCH /enterprises/{enterprise}/apps/organizations/{org}/installations/{installation_id}/repositories/add
func (s *EnterpriseAppsService) AddRepositoriesToInstallation(ctx context.Context, enterprise, org string, installationID int64, opts *EnterpriseInstallationRepositoriesOptions) (*ListRepositories, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/apps/organizations/%v/installations/%v/repositories/add", enterprise, org, installationID)
	req, err := s.client.NewRequest("PATCH", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var r *ListRepositories
	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// RemoveRepositoriesFromInstallation revokes repository access from a GitHub App installation.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/organization-installations#remove-repository-access-from-an-organization-installation
//
//meta:operation PATCH /enterprises/{enterprise}/apps/organizations/{org}/installations/{installation_id}/repositories/remove
func (s *EnterpriseAppsService) RemoveRepositoriesFromInstallation(ctx context.Context, enterprise, org string, installationID int64, opts *EnterpriseInstallationRepositoriesOptions) (*ListRepositories, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/apps/organizations/%v/installations/%v/repositories/remove", enterprise, org, installationID)
	req, err := s.client.NewRequest("PATCH", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var r *ListRepositories
	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}
