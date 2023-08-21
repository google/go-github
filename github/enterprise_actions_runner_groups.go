// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ListOrganizations represents the response from the list orgs endpoints.
type ListOrganizations struct {
	TotalCount    *int            `json:"total_count,omitempty"`
	Organizations []*Organization `json:"organizations"`
}

// EnterpriseRunnerGroup represents a self-hosted runner group configured in an enterprise.
type EnterpriseRunnerGroup struct {
	ID                           *int64   `json:"id,omitempty"`
	Name                         *string  `json:"name,omitempty"`
	Visibility                   *string  `json:"visibility,omitempty"`
	Default                      *bool    `json:"default,omitempty"`
	SelectedOrganizationsURL     *string  `json:"selected_organizations_url,omitempty"`
	RunnersURL                   *string  `json:"runners_url,omitempty"`
	Inherited                    *bool    `json:"inherited,omitempty"`
	AllowsPublicRepositories     *bool    `json:"allows_public_repositories,omitempty"`
	RestrictedToWorkflows        *bool    `json:"restricted_to_workflows,omitempty"`
	SelectedWorkflows            []string `json:"selected_workflows,omitempty"`
	WorkflowRestrictionsReadOnly *bool    `json:"workflow_restrictions_read_only,omitempty"`
}

// EnterpriseRunnerGroups represents a collection of self-hosted runner groups configured for an enterprise.
type EnterpriseRunnerGroups struct {
	TotalCount   *int                     `json:"total_count"`
	RunnerGroups []*EnterpriseRunnerGroup `json:"runner_groups"`
}

// CreateEnterpriseRunnerGroupRequest represents a request to create a Runner group for an enterprise.
type CreateEnterpriseRunnerGroupRequest struct {
	Name       *string `json:"name,omitempty"`
	Visibility *string `json:"visibility,omitempty"`
	// List of organization IDs that can access the runner group.
	SelectedOrganizationIDs []int64 `json:"selected_organization_ids,omitempty"`
	// Runners represent a list of runner IDs to add to the runner group.
	Runners []int64 `json:"runners,omitempty"`
	// If set to True, public repos can use this runner group
	AllowsPublicRepositories *bool `json:"allows_public_repositories,omitempty"`
	// If true, the runner group will be restricted to running only the workflows specified in the SelectedWorkflows slice.
	RestrictedToWorkflows *bool `json:"restricted_to_workflows,omitempty"`
	// List of workflows the runner group should be allowed to run. This setting will be ignored unless RestrictedToWorkflows is set to true.
	SelectedWorkflows []string `json:"selected_workflows,omitempty"`
}

// UpdateEnterpriseRunnerGroupRequest represents a request to update a Runner group for an enterprise.
type UpdateEnterpriseRunnerGroupRequest struct {
	Name                     *string  `json:"name,omitempty"`
	Visibility               *string  `json:"visibility,omitempty"`
	AllowsPublicRepositories *bool    `json:"allows_public_repositories,omitempty"`
	RestrictedToWorkflows    *bool    `json:"restricted_to_workflows,omitempty"`
	SelectedWorkflows        []string `json:"selected_workflows,omitempty"`
}

// SetOrgAccessRunnerGroupRequest represents a request to replace the list of organizations
// that can access a self-hosted runner group configured in an enterprise.
type SetOrgAccessRunnerGroupRequest struct {
	// Updated list of organization IDs that should be given access to the runner group.
	SelectedOrganizationIDs []int64 `json:"selected_organization_ids"`
}

// ListEnterpriseRunnerGroupOptions extend ListOptions to have the optional parameters VisibleToOrganization.
type ListEnterpriseRunnerGroupOptions struct {
	ListOptions

	// Only return runner groups that are allowed to be used by this organization.
	VisibleToOrganization string `url:"visible_to_organization,omitempty"`
}

// ListRunnerGroups lists all self-hosted runner groups configured in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#list-self-hosted-runner-groups-for-an-enterprise
func (s *EnterpriseService) ListRunnerGroups(ctx context.Context, enterprise string, opts *ListEnterpriseRunnerGroupOptions) (*EnterpriseRunnerGroups, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups", enterprise)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	groups := &EnterpriseRunnerGroups{}
	resp, err := s.client.Do(ctx, req, &groups)
	if err != nil {
		return nil, resp, err
	}

	return groups, resp, nil
}

// GetRunnerGroup gets a specific self-hosted runner group for an enterprise using its RunnerGroup ID.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#get-a-self-hosted-runner-group-for-an-enterprise
func (s *EnterpriseService) GetRunnerGroup(ctx context.Context, enterprise string, groupID int64) (*EnterpriseRunnerGroup, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups/%v", enterprise, groupID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	runnerGroup := new(EnterpriseRunnerGroup)
	resp, err := s.client.Do(ctx, req, runnerGroup)
	if err != nil {
		return nil, resp, err
	}

	return runnerGroup, resp, nil
}

// DeleteRunnerGroup deletes a self-hosted runner group from an enterprise.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#delete-a-self-hosted-runner-group-from-an-enterprise
func (s *EnterpriseService) DeleteRunnerGroup(ctx context.Context, enterprise string, groupID int64) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups/%v", enterprise, groupID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// CreateRunnerGroup creates a new self-hosted runner group for an enterprise.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#create-a-self-hosted-runner-group-for-an-enterprise
func (s *EnterpriseService) CreateRunnerGroup(ctx context.Context, enterprise string, createReq CreateEnterpriseRunnerGroupRequest) (*EnterpriseRunnerGroup, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups", enterprise)
	req, err := s.client.NewRequest("POST", u, createReq)
	if err != nil {
		return nil, nil, err
	}

	runnerGroup := new(EnterpriseRunnerGroup)
	resp, err := s.client.Do(ctx, req, runnerGroup)
	if err != nil {
		return nil, resp, err
	}

	return runnerGroup, resp, nil
}

// UpdateRunnerGroup updates a self-hosted runner group for an enterprise.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#update-a-self-hosted-runner-group-for-an-enterprise
func (s *EnterpriseService) UpdateRunnerGroup(ctx context.Context, enterprise string, groupID int64, updateReq UpdateEnterpriseRunnerGroupRequest) (*EnterpriseRunnerGroup, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups/%v", enterprise, groupID)
	req, err := s.client.NewRequest("PATCH", u, updateReq)
	if err != nil {
		return nil, nil, err
	}

	runnerGroup := new(EnterpriseRunnerGroup)
	resp, err := s.client.Do(ctx, req, runnerGroup)
	if err != nil {
		return nil, resp, err
	}

	return runnerGroup, resp, nil
}

// ListOrganizationAccessRunnerGroup lists the organizations with access to a self-hosted runner group configured in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#list-organization-access-to-a-self-hosted-runner-group-in-an-enterprise
func (s *EnterpriseService) ListOrganizationAccessRunnerGroup(ctx context.Context, enterprise string, groupID int64, opts *ListOptions) (*ListOrganizations, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups/%v/organizations", enterprise, groupID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	orgs := &ListOrganizations{}
	resp, err := s.client.Do(ctx, req, &orgs)
	if err != nil {
		return nil, resp, err
	}

	return orgs, resp, nil
}

// SetOrganizationAccessRunnerGroup replaces the list of organizations that have access to a self-hosted runner group configured in an enterprise
// with a new List of organizations.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#set-organization-access-for-a-self-hosted-runner-group-in-an-enterprise
func (s *EnterpriseService) SetOrganizationAccessRunnerGroup(ctx context.Context, enterprise string, groupID int64, ids SetOrgAccessRunnerGroupRequest) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups/%v/organizations", enterprise, groupID)

	req, err := s.client.NewRequest("PUT", u, ids)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// AddOrganizationAccessRunnerGroup adds an organization to the list of selected organizations that can access a self-hosted runner group.
// The runner group must have visibility set to 'selected'.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#add-organization-access-to-a-self-hosted-runner-group-in-an-enterprise
func (s *EnterpriseService) AddOrganizationAccessRunnerGroup(ctx context.Context, enterprise string, groupID, orgID int64) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups/%v/organizations/%v", enterprise, groupID, orgID)

	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// RemoveOrganizationAccessRunnerGroup removes an organization from the list of selected organizations that can access a self-hosted runner group.
// The runner group must have visibility set to 'selected'.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#remove-organization-access-to-a-self-hosted-runner-group-in-an-enterprise
func (s *EnterpriseService) RemoveOrganizationAccessRunnerGroup(ctx context.Context, enterprise string, groupID, orgID int64) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups/%v/organizations/%v", enterprise, groupID, orgID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// ListEnterpriseRunnerGroupRunners lists self-hosted runners that are in a specific enterprise group.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#list-self-hosted-runners-in-a-group-for-an-enterprise
func (s *EnterpriseService) ListEnterpriseRunnerGroupRunners(ctx context.Context, enterprise string, groupID int64, opts *ListOptions) (*Runners, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups/%v/runners", enterprise, groupID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	runners := &Runners{}
	resp, err := s.client.Do(ctx, req, &runners)
	if err != nil {
		return nil, resp, err
	}

	return runners, resp, nil
}

// SetEnterpriseRunnerGroupRunners replaces the list of self-hosted runners that are part of an enterprise runner group
// with a new list of runners.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#set-self-hosted-runners-in-a-group-for-an-enterprise
func (s *EnterpriseService) SetEnterpriseRunnerGroupRunners(ctx context.Context, enterprise string, groupID int64, ids SetRunnerGroupRunnersRequest) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups/%v/runners", enterprise, groupID)

	req, err := s.client.NewRequest("PUT", u, ids)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// AddEnterpriseRunnerGroupRunners adds a self-hosted runner to a runner group configured in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#add-a-self-hosted-runner-to-a-group-for-an-enterprise
func (s *EnterpriseService) AddEnterpriseRunnerGroupRunners(ctx context.Context, enterprise string, groupID, runnerID int64) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups/%v/runners/%v", enterprise, groupID, runnerID)

	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// RemoveEnterpriseRunnerGroupRunners removes a self-hosted runner from a group configured in an enterprise.
// The runner is then returned to the default group.
//
// GitHub API docs: https://docs.github.com/en/enterprise-cloud@latest/rest/actions/self-hosted-runner-groups?apiVersion=2022-11-28#remove-a-self-hosted-runner-from-a-group-for-an-enterprise
func (s *EnterpriseService) RemoveEnterpriseRunnerGroupRunners(ctx context.Context, enterprise string, groupID, runnerID int64) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/actions/runner-groups/%v/runners/%v", enterprise, groupID, runnerID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
