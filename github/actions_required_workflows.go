// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// OrgRequiredWorkflow represents a required workflow object
type OrgRequiredWorkflow struct {
	ID                      *int64      `json:"id,omitempty"`
	Name                    *string     `json:"name,omitempty"`
	Path                    *string     `json:"path,omitempty"`
	Scope                   *string     `json:"scope,omitempty"`
	Ref                     *string     `json:"ref,omitempty"`
	State                   *string     `json:"state,omitempty"`
	SelectedRepositoriesUrl *string     `json:"selected_repositories_url,omitempty"`
	CreatedAt               *Timestamp  `json:"created_at,omitempty"`
	UpdatedAt               *Timestamp  `json:"updated_at,omitempty"`
	Repository              *Repository `json:"repository,omitempty"`
}

type OrgRequiredWorkflows struct {
	TotalCount        *int                   `json:"total_count,omitempty"`
	RequiredWorkflows []*OrgRequiredWorkflow `json:"required_workflows,omitempty"`
}

type CreateUpdateRequiredWorkflowOptions struct {
	WorkflowFilepath      *string         `json:"workflow_file_path,omitempty"`
	RepositoryID          *int64          `json:"repository_id,omitempty"`
	Scope                 *string         `json:"scope,omitempty"`
	SelectedRepositoryIDs SelectedRepoIDs `json:"selected_repository_ids,omitempty"`
}

type RequiredWorkflowSelectedRepositories struct {
	TotalCount   *int          `json:"total_count,omitempty"`
	Repositories []*Repository `json:"repositories,omitempty"`
}

type RepoRequiredWorkflow struct {
	ID               *int64      `json:"id,omitempty"`
	NodeID           *string     `json:"node_id,omitempty"`
	Name             *string     `json:"name,omitempty"`
	Path             *string     `json:"path,omitempty"`
	State            *string     `json:"state,omitempty"`
	URL              *string     `json:"url,omitempty"`
	HtmlURL          *string     `json:"html_url,omitempty"`
	BadgeURL         *string     `json:"badge_url,omitempty"`
	CreatedAt        *Timestamp  `json:"created_at,omitempty"`
	UpdatedAt        *Timestamp  `json:"updated_at,omitempty"`
	SourceRepository *Repository `json:"source_repository,omitempty"`
}

type RepoRequiredWorkflows struct {
	TotalCount        *int                    `json:"total_count,omitempty"`
	RequiredWorkflows []*RepoRequiredWorkflow `json:"required_workflows,omitempty"`
}

// ListOrgRequiredWorkflows lists the RequiredWorkflows for an org.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/required-workflows?apiVersion=2022-11-28#list-required-workflows
func (s *ActionsService) ListOrgRequiredWorkflows(ctx context.Context, org string, opts *ListOptions) (*OrgRequiredWorkflows, *Response, error) {
	url := fmt.Sprintf("orgs/%v/actions/required_workflows", org)
	u, err := addOptions(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	requiredWorkflows := new(OrgRequiredWorkflows)
	resp, err := s.client.Do(ctx, req, &requiredWorkflows)
	if err != nil {
		return nil, resp, err
	}

	return requiredWorkflows, resp, nil
}

// CreateRequiredWorkflow creates the required workflow in an org
//
// GitHub API docs: https://docs.github.com/en/rest/actions/required-workflows?apiVersion=2022-11-28#create-a-required-workflow
func (s *ActionsService) CreateRequiredWorkflow(ctx context.Context, org string, createRequiredWorkflowOptions *CreateUpdateRequiredWorkflowOptions) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/actions/required_workflows", org)
	req, err := s.client.NewRequest("PUT", url, createRequiredWorkflowOptions)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// GetRequiredWorkflowByID get the RequiredWorkflows for an org by its ID.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/required-workflows?apiVersion=2022-11-28#list-required-workflows
func (s *ActionsService) GetRequiredWorkflowByID(ctx context.Context, owner string, requiredWorkflowID int64) (*OrgRequiredWorkflow, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/required_workflows/%v", owner, requiredWorkflowID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	requiredWorkflow := new(OrgRequiredWorkflow)
	resp, err := s.client.Do(ctx, req, &requiredWorkflow)
	if err != nil {
		return nil, resp, err
	}

	return requiredWorkflow, resp, nil
}

// UpdateRequiredWorkflow updates a required workflow in an org
//
// GitHub API docs: https://docs.github.com/en/rest/actions/required-workflows?apiVersion=2022-11-28#update-a-required-workflow
func (s *ActionsService) UpdateRequiredWorkflow(ctx context.Context, org string, requiredWorkflowID int64, updateRequiredWorkflowOptions *CreateUpdateRequiredWorkflowOptions) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/actions/required_workflows/%v", org, requiredWorkflowID)
	req, err := s.client.NewRequest("PATCH", url, updateRequiredWorkflowOptions)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// DeleteRequiredWorkflow deletes a required workflow in an org
//
// GitHub API docs: https://docs.github.com/en/rest/actions/required-workflows?apiVersion=2022-11-28#update-a-required-workflow
func (s *ActionsService) DeleteRequiredWorkflow(ctx context.Context, org string, requiredWorkflowID int64) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/actions/required_workflows/%v", org, requiredWorkflowID)
	req, err := s.client.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// ListRequiredWorkflowSelectedRepositories lists the Repositories selected for a workflow.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/required-workflows?apiVersion=2022-11-28#list-selected-repositories-for-a-required-workflow
func (s *ActionsService) ListRequiredWorkflowSelectedRepositories(ctx context.Context, org string, requiredWorkflowID int64, opts *ListOptions) (*RequiredWorkflowSelectedRepositories, *Response, error) {
	url := fmt.Sprintf("orgs/%v/actions/required_workflows/%v/repositories", org, requiredWorkflowID)
	u, err := addOptions(url, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	requiredWorkflowRepositories := new(RequiredWorkflowSelectedRepositories)
	resp, err := s.client.Do(ctx, req, &requiredWorkflowRepositories)
	if err != nil {
		return nil, resp, err
	}

	return requiredWorkflowRepositories, resp, nil
}

// SetRequiredWorkflowSelectedRepositories sets the Repositories selected for a workflow.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/required-workflows?apiVersion=2022-11-28#sets-repositories-for-a-required-workflow
func (s *ActionsService) SetRequiredWorkflowSelectedRepositories(ctx context.Context, org string, requiredWorkflowID int64, ids SelectedRepoIDs) (*Response, error) {
	type repoIDs struct {
		SelectedIDs SelectedRepoIDs `json:"selected_repository_ids"`
	}
	url := fmt.Sprintf("orgs/%v/actions/required_workflows/%v/repositories", org, requiredWorkflowID)
	req, err := s.client.NewRequest("PUT", url, repoIDs{SelectedIDs: ids})
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// AddRepoToRequiredWorkflow adds the Repository to a required workflow.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/required-workflows?apiVersion=2022-11-28#add-a-repository-to-a-required-workflow
func (s *ActionsService) AddRepoToRequiredWorkflow(ctx context.Context, org string, requiredWorkflowID int64, repoID int64) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/actions/required_workflows/%v/repositories/%v", org, requiredWorkflowID, repoID)
	req, err := s.client.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// RemoveRepoFromRequiredWorkflow removes the Repository from a required workflow.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/required-workflows?apiVersion=2022-11-28#add-a-repository-to-a-required-workflow
func (s *ActionsService) RemoveRepoFromRequiredWorkflow(ctx context.Context, org string, requiredWorkflowID int64, repoID int64) (*Response, error) {
	url := fmt.Sprintf("orgs/%v/actions/required_workflows/%v/repositories/%v", org, requiredWorkflowID, repoID)
	req, err := s.client.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// ListRepoRequiredWorkflows lists the RequiredWorkflows for a repo.
//
// Github API docs:https://docs.github.com/en/rest/actions/required-workflows?apiVersion=2022-11-28#list-repository-required-workflows
func (s *ActionsService) ListRepoRequiredWorkflows(ctx context.Context, owner, repo string, opts *ListOptions) (*RepoRequiredWorkflows, *Response, error) {
	url := fmt.Sprintf("repos/%v/%v/actions/required_workflows", owner, repo)
	u, err := addOptions(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	requiredWorkflows := new(RepoRequiredWorkflows)
	resp, err := s.client.Do(ctx, req, &requiredWorkflows)
	if err != nil {
		return nil, resp, err
	}

	return requiredWorkflows, resp, nil
}
