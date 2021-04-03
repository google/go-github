// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
)

// Environment represents a single environment in a repository.
type Environment struct {
	Owner           *string         `json:"owner,omitempty"`
	Repo            *string         `json:"repo,omitempty"`
	EnvironmentName *string         `json:"environment_name,omitempty"`
	WaitTimer       *int            `json:"wait_timer,omitempty"`
	Reviewers       []*EnvReviewers `json:"reviewers,omitempty"`
	BranchPolicy    *BranchPolicy   `json:"deployment_branch_policy,omitempty"`
	// Return/response only values
	ID              *int              `json:"id,omitempty"`
	NodeID          *string           `json:"node_id,omitempty"`
	Name            *string           `json:"name,omitempty"`
	URL             *string           `json:"url,omitempty"`
	HTMLURL         *string           `json:"html_url,omitempty"`
	CreatedAt       *Timestamp        `json:"created_at,omitempty"`
	UpdatedAt       *Timestamp        `json:"updated_at,omitempty"`
	ProtectionRules []*ProtectionRule `json:"protection_rules,omitempty"`
}

// EnvReviewers represents a single environment reviewer entry
type EnvReviewers struct {
	Type *string `json:"type,omitempty"`
	ID   *int    `json:"id,omitempty"`
}

// BranchPolicy represents the options for whether a branch deployment policy is applied to this environment
type BranchPolicy struct {
	ProtectedBranches    *bool `json:"protected_branches,omitempty"`
	CustomBranchPolicies *bool `json:"custom_branch_policies,omitempty"` // Unable to set the actual custom policy at the moment so unsure how useful this is
}

// EnvResponse represents the slightly different format of response that comes back when you list an environment
type EnvResponse struct {
	TotalCount   *int           `json:"total_count,omitempty"`
	Environments []*Environment `json:"environments,omitempty"`
}

// ProtectionRule represents a single protection rule applied to the environment
type ProtectionRule struct {
	ID        *int               `json:"id,omitempty"`
	NodeID    *string            `json:"node_id,omitempty"`
	Type      *string            `json:"type,omitempty"`
	WaitTimer *int               `json:"wait_timer,omitempty"`
	Reviewers []RequiredReviewer `json:"reviewers,omitempty"`
}

type RequiredReviewer struct {
	Type     *string     `json:"type,omitempty"`
	Reviewer interface{} `json:"reviewer,omitempty"`
}

func (r *RequiredReviewer) UnmarshalJSON(data []byte) error {
	type aliasReviewer RequiredReviewer
	var reviewer aliasReviewer
	if err := json.Unmarshal(data, &reviewer); err != nil {
		return err
	}

	r.Type = reviewer.Type

	switch *reviewer.Type {
	case "User":
		reviewer.Reviewer = &User{}
		if err := json.Unmarshal(data, &reviewer); err != nil {
			return err
		}
		r.Reviewer = reviewer.Reviewer
	case "Team":
		reviewer.Reviewer = &Team{}
		if err := json.Unmarshal(data, &reviewer); err != nil {
			return err
		}
		r.Reviewer = reviewer.Reviewer
	}

	return nil
}

// ListEnvironments lists all environments for a repository
// https://docs.github.com/en/rest/reference/repos#get-all-environments
func (s *RepositoriesService) ListEnvironments(ctx context.Context, owner, repo string) (*EnvResponse, *Response, error) {
	u := fmt.Sprintf("repos/%s/%s/environments", owner, repo)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var list *EnvResponse
	resp, err := s.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}
	return list, resp, nil
}

// GetEnvironment get a single environment for a repository
// https://docs.github.com/en/rest/reference/repos#get-an-environment
func (s *RepositoriesService) GetEnvironment(ctx context.Context, owner, repo, name string) (*Environment, *Response, error) {
	u := fmt.Sprintf("repos/%s/%s/environments/%s", owner, repo, name)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var env *Environment
	resp, err := s.client.Do(ctx, req, &env)
	if err != nil {
		return nil, resp, err
	}
	return env, resp, nil
}

type CreateUpdateEnvironment struct {
	WaitTimer    *int                     `json:"wait_timer,omitempty"`
	Reviewers    []*CreateUpdateReviewers `json:"reviewers,omitempty"`
	BranchPolicy *BranchPolicy            `json:"deployment_branch_policy,omitempty"`
}

type CreateUpdateReviewers struct {
	Type *string `json:"type,omitempty"`
	ID   *int    `json:"id,omitempty"`
}

// CreateUpdateEnvironment create or update a new environment for a repository
// https://docs.github.com/en/rest/reference/repos#create-or-update-an-environment
func (s *RepositoriesService) CreateUpdateEnvironment(ctx context.Context, owner, repo, name string, environment *CreateUpdateEnvironment) (*Environment, *Response, error) {
	u := fmt.Sprintf("repos/%s/%s/environments/%s", owner, repo, name)

	req, err := s.client.NewRequest("PUT", u, environment)
	if err != nil {
		return nil, nil, err
	}

	e := new(Environment)
	resp, err := s.client.Do(ctx, req, e)
	if err != nil {
		return nil, resp, err
	}
	return e, resp, nil
}

// DeleteEnvironment delete an environment from a repository
// https://docs.github.com/en/rest/reference/repos#delete-an-environment
func (s *RepositoriesService) DeleteEnvironment(ctx context.Context, owner, repo, name string) (*Response, error) {
	u := fmt.Sprintf("repos/%s/%s/environments/%s", owner, repo, name)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}
