// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "fmt"

// RepositoryDeployment represents a deployment in a repo
type RepositoryDeployment struct {
	Url         *string            `json:"url,omitempty"`
	ID          *int               `json:"id,omitempty"`
	SHA         *string            `json:"sha,omitempty"`
	Ref         *string            `json:"ref,omitempty"`
	Task        *string            `json:"task,omitempty"`
	Payload     *map[string]string `json:"payload,omitempty"`
	Environment *string            `json:"environment,omitempty"`
	Description *string            `json:"description,omitempty"`
	Creator     *User              `json:"creator,omitempty"`
	CreatedAt   *Timestamp         `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp         `json:"pushed_at,omitempty"`
}

// RepositoryDeploymentRequest represents a deployment request
type RepositoryDeploymentRequest struct {
	Ref              *string   `json:"ref,omitempty"`
	Task             *string   `json:"task,omitempty"`
	AutoMerge        *bool     `json:"auto_merge,omitempty"`
	RequiredContexts *[]string `json:"required_contexts,omitempty"`
	Payload          *string   `json:"payload,omitempty"`
	Environment      *string   `json:"environment,omitempty"`
	Description      *string   `json:"description,omitempty"`
}

// DeploymentsListOptions specifies the optional parameters to the
// RepositoriesService.ListDeployments method.
type DeploymentsListOptions struct {
	// SHA of the Deployment.
	SHA string `url:"sha,omitempty"`

	// List deployments for a given ref.
	Ref string `url:"ref,omitempty"`

	// List deployments for a given task.
	Task string `url:"task,omitempty"`

	// List deployments for a given environment.
	Environment string `url:"environment,omitempty"`

	ListOptions
}

// ListDeployments lists the deployments of a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/deployments/#list-deployments
func (s *RepositoriesService) ListDeployments(owner, repo string, opt *DeploymentsListOptions) ([]RepositoryDeployment, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/deployments", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	deployments := new([]RepositoryDeployment)
	resp, err := s.client.Do(req, deployments)
	if err != nil {
		return nil, resp, err
	}

	return *deployments, resp, err
}

// CreateDeployment creates a new deployment for a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/deployments/#create-a-deployment
func (s *RepositoriesService) CreateDeployment(owner, repo string, deployment_req *RepositoryDeploymentRequest) (*RepositoryDeployment, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/deployments", owner, repo)

	req, err := s.client.NewRequest("POST", u, deployment_req)
	if err != nil {
		return nil, nil, err
	}

	d := new(RepositoryDeployment)
	resp, err := s.client.Do(req, d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, err
}

// RepositoryDeploymentStatus represents the status of a
// particular deployment.
type RepositoryDeploymentStatus struct {
	ID          *int       `json:"id,omitempty"`
	State       *string    `json:"state,omitempty"`
	Creator     *User      `json:"creator,omitempty"`
	Description *string    `json:"description,omitempty"`
	TargetUrl   *string    `json:"target_url,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp `json:"pushed_at,omitempty"`
}

// RepositoryDeploymentRequest represents a deployment request
type RepositoryDeploymentStatusRequest struct {
	State       *string `json:"state,omitempty"`
	TargetUrl   *string `json:"target_url,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ListDeploymentStatuses lists the statuses of a given deployment of a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/deployments/#list-deployment-statuses
func (s *RepositoriesService) ListDeploymentStatuses(owner, repo string, deployment_id int, opt *ListOptions) ([]RepositoryDeploymentStatus, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/deployments/%v/statuses", owner, repo, deployment_id)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	statuses := new([]RepositoryDeploymentStatus)
	resp, err := s.client.Do(req, statuses)
	if err != nil {
		return nil, resp, err
	}

	return *statuses, resp, err
}

// CreateDeploymentStatus creates a new status for a deployment.
//
// GitHub API docs: https://developer.github.com/v3/repos/deployments/#create-a-deployment-status
func (s *RepositoriesService) CreateDeploymentStatus(owner, repo string, deployment_id int, status_req *RepositoryDeploymentStatusRequest) (*RepositoryDeploymentStatus, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/deployments/%v/statuses", owner, repo, deployment_id)

	req, err := s.client.NewRequest("POST", u, status_req)
	if err != nil {
		return nil, nil, err
	}

	d := new(RepositoryDeploymentStatus)
	resp, err := s.client.Do(req, d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, err
}
