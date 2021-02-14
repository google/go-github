// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// RunnerGroup represents a self-hosted runner group configured in an organization.
type RunnerGroup struct {
	ID                       *int64  `json:"id,omitempty"`
	Name                     *string `json:"name,omitempty"`
	Visibility               *string `json:"visibility,omitempty"`
	Default                  *bool   `json:"default,omitempty"`
	SelectedRepositoriesURL  *string `json:"selected_repositories_url,omitempty"`
	RunnersURL               *string `json:"runners_url,omitempty"`
	Inherited                *bool   `json:"inherited,omitempty"`
	AllowsPublicRepositories *bool   `json:"allows_public_repositories,omitempty"`
}

// RunnerGroups represents a collection of self-hosted runner groups configured for an organization
type RunnerGroups struct {
	TotalCount   int            `json:"total_count"`
	RunnerGroups []*RunnerGroup `json:"runner_groups"`
}

// ListOrganizationRunnerGroups lists all the self-hosted runner groups for an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/actions#list-self-hosted-runner-groups-for-an-organization
func (s *ActionsService) ListOrganizationRunnerGroups(ctx context.Context, owner string, opts *ListOptions) (*RunnerGroups, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/runner-groups", owner)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	groups := &RunnerGroups{}
	resp, err := s.client.Do(ctx, req, &groups)
	if err != nil {
		return nil, resp, err
	}

	return groups, resp, nil
}
