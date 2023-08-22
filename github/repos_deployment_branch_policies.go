// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// DeploymentBranchPolicy represents a single deployment branch policy for an environment.
type DeploymentBranchPolicy struct {
	Name   *string `json:"name,omitempty"`
	ID     *int64  `json:"id,omitempty"`
	NodeID *string `json:"node_id,omitempty"`
}

// DeploymentBranchPolicyResponse represents the slightly different format of response that comes back when you list deployment branch policies.
type DeploymentBranchPolicyResponse struct {
	TotalCount     *int                      `json:"total_count,omitempty"`
	BranchPolicies []*DeploymentBranchPolicy `json:"branch_policies,omitempty"`
}

// DeploymentBranchPolicyRequest represents a deployment branch policy request.
type DeploymentBranchPolicyRequest struct {
	Name *string `json:"name,omitempty"`
}

// ListDeploymentBranchPolicies lists the deployment branch policies for an environment.
//
// GitHub API docs: https://docs.github.com/en/rest/deployments/branch-policies#list-deployment-branch-policies
func (s *RepositoriesService) ListDeploymentBranchPolicies(ctx context.Context, owner, repo, environment string) (*DeploymentBranchPolicyResponse, *Response, error) {
	u, err := newURLString("repos/%v/%v/environments/%v/deployment-branch-policies", owner, repo, environment)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var list *DeploymentBranchPolicyResponse
	resp, err := s.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

// GetDeploymentBranchPolicy gets a deployment branch policy for an environment.
//
// GitHub API docs: https://docs.github.com/en/rest/deployments/branch-policies#get-a-deployment-branch-policy
func (s *RepositoriesService) GetDeploymentBranchPolicy(ctx context.Context, owner, repo, environment string, branchPolicyID int64) (*DeploymentBranchPolicy, *Response, error) {
	u, err := newURLString("repos/%v/%v/environments/%v/deployment-branch-policies/%v", owner, repo, environment, branchPolicyID)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var policy *DeploymentBranchPolicy
	resp, err := s.client.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}

	return policy, resp, nil
}

// CreateDeploymentBranchPolicy creates a deployment branch policy for an environment.
//
// GitHub API docs: https://docs.github.com/en/rest/deployments/branch-policies#create-a-deployment-branch-policy
func (s *RepositoriesService) CreateDeploymentBranchPolicy(ctx context.Context, owner, repo, environment string, request *DeploymentBranchPolicyRequest) (*DeploymentBranchPolicy, *Response, error) {
	u, err := newURLString("repos/%v/%v/environments/%v/deployment-branch-policies", owner, repo, environment)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, request)
	if err != nil {
		return nil, nil, err
	}

	var policy *DeploymentBranchPolicy
	resp, err := s.client.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}

	return policy, resp, nil
}

// UpdateDeploymentBranchPolicy updates a deployment branch policy for an environment.
//
// GitHub API docs: https://docs.github.com/en/rest/deployments/branch-policies#update-a-deployment-branch-policy
func (s *RepositoriesService) UpdateDeploymentBranchPolicy(ctx context.Context, owner, repo, environment string, branchPolicyID int64, request *DeploymentBranchPolicyRequest) (*DeploymentBranchPolicy, *Response, error) {
	u, err := newURLString("repos/%v/%v/environments/%v/deployment-branch-policies/%v", owner, repo, environment, branchPolicyID)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("PUT", u, request)
	if err != nil {
		return nil, nil, err
	}

	var policy *DeploymentBranchPolicy
	resp, err := s.client.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}

	return policy, resp, nil
}

// DeleteDeploymentBranchPolicy deletes a deployment branch policy for an environment.
//
// GitHub API docs: https://docs.github.com/en/rest/deployments/branch-policies#delete-a-deployment-branch-policy
func (s *RepositoriesService) DeleteDeploymentBranchPolicy(ctx context.Context, owner, repo, environment string, branchPolicyID int64) (*Response, error) {
	u, err := newURLString("repos/%v/%v/environments/%v/deployment-branch-policies/%v", owner, repo, environment, branchPolicyID)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
