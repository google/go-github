// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// CodeQualityService handles communication with the code quality related
// methods of the GitHub API.
//
// GitHub API docs: https://docs.github.com/rest/code-quality/code-quality?apiVersion=2026-03-10
type CodeQualityService service

// CodeQualitySetupConfiguration represents a code quality setup configuration for a repository.
type CodeQualitySetupConfiguration struct {
	State       *string    `json:"state,omitempty"`
	Languages   []string   `json:"languages,omitempty"`
	RunnerType  *string    `json:"runner_type,omitempty"`
	RunnerLabel *string    `json:"runner_label,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
	Schedule    *string    `json:"schedule,omitempty"`
}

// UpdateCodeQualitySetupOptions specifies parameters to the
// CodeQualityService.UpdateSetup method.
type UpdateCodeQualitySetupOptions struct {
	State       *string  `json:"state,omitempty"`
	RunnerType  *string  `json:"runner_type,omitempty"`
	RunnerLabel *string  `json:"runner_label,omitempty"`
	Languages   []string `json:"languages,omitempty"`
}

// UpdateCodeQualitySetupResponse represents a response from updating a code quality setup configuration.
type UpdateCodeQualitySetupResponse struct {
	RunID  *int64  `json:"run_id,omitempty"`
	RunURL *string `json:"run_url,omitempty"`
}

// GetSetup gets a code quality setup configuration for a repository.
//
// GitHub API docs: https://docs.github.com/rest/code-quality/code-quality?apiVersion=2022-11-28#get-a-code-quality-setup-configuration
//
//meta:operation GET /repos/{owner}/{repo}/code-quality/setup
func (s *CodeQualityService) GetSetup(ctx context.Context, owner, repo string) (*CodeQualitySetupConfiguration, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/code-quality/setup", owner, repo)

	req, err := s.client.NewRequest(ctx, "GET", u, nil, WithVersion(api20260310))
	if err != nil {
		return nil, nil, err
	}

	var cfg *CodeQualitySetupConfiguration
	resp, err := s.client.Do(req, &cfg)
	if err != nil {
		return nil, resp, err
	}

	return cfg, resp, nil
}

// UpdateSetup updates a code quality setup configuration for a repository.
//
// This method might return an AcceptedError and a status code of 202. This is because this is the status that GitHub
// returns to signify that it has now scheduled the update in a background task.
//
// GitHub API docs: https://docs.github.com/rest/code-quality/code-quality?apiVersion=2022-11-28#update-a-code-quality-setup-configuration
//
//meta:operation PATCH /repos/{owner}/{repo}/code-quality/setup
func (s *CodeQualityService) UpdateSetup(ctx context.Context, owner, repo string, opts *UpdateCodeQualitySetupOptions) (*UpdateCodeQualitySetupResponse, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/code-quality/setup", owner, repo)

	req, err := s.client.NewRequest(ctx, "PATCH", u, opts, WithVersion(api20260310))
	if err != nil {
		return nil, nil, err
	}

	var result *UpdateCodeQualitySetupResponse
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
