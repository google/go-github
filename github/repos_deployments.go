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
func (s *RepositoriesService) CreateDeployment(owner, repo string, deployment *RepositoryDeploymentRequest) (*RepositoryDeployment, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/deployments", owner, repo)

	req, err := s.client.NewRequest("POST", u, deployment)
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
