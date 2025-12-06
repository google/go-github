// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

type EnterpriseTeam struct {
	ID                        int64     `json:"id"`
	URL                       string    `json:"url"`
	MemberURL                 string    `json:"member_url"`
	Name                      string    `json:"name"`
	Description               *string   `json:"description,omitempty"`
	HTMLURL                   string    `json:"html_url"`
	Slug                      string    `json:"slug"`
	CreatedAt                 Timestamp `json:"created_at"`
	UpdatedAt                 Timestamp `json:"updated_at"`
	GroupID                   int64     `json:"group_id"`
	OrganizationSelectionType *string   `json:"organization_selection_type,omitempty"`
	SyncToOrganizations       *string   `json:"sync_to_organizations,omitempty"`
}

type EnterpriseTeamCreateRequest struct {
	Name                      string  `json:"name"`
	Description               *string `json:"description,omitempty"`
	SyncToOrganizations       *string `json:"sync_to_organizations,omitempty"`
	OrganizationSelectionType *string `json:"organization_selection_type,omitempty"`
	GroupID                   *int64  `json:"group_id,omitempty"`
}

// ListTeams lists all teams in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/rest/enterprise-teams/enterprise-teams#list-enterprise-teams
//
//meta:operation GET /enterprises/{enterprise}/teams
func (s *EnterpriseService) ListTeams(ctx context.Context, enterprise string, opt *ListOptions) ([]*EnterpriseTeam, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/teams", enterprise)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var teams []*EnterpriseTeam
	resp, err := s.client.Do(ctx, req, &teams)
	if err != nil {
		return nil, resp, err
	}

	return teams, resp, nil
}

// CreateTeam creates a new team in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/rest/enterprise-teams/enterprise-teams#create-an-enterprise-team
//
//meta:operation POST /enterprises/{enterprise}/teams
func (s *EnterpriseService) CreateTeam(ctx context.Context, enterprise string, team EnterpriseTeamCreateRequest) (*EnterpriseTeam, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/teams", enterprise)

	req, err := s.client.NewRequest("POST", u, team)
	if err != nil {
		return nil, nil, err
	}

	var createdTeam *EnterpriseTeam
	resp, err := s.client.Do(ctx, req, &createdTeam)
	if err != nil {
		return nil, resp, err
	}

	return createdTeam, resp, nil
}

// GetTeam retrieves a team in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/rest/enterprise-teams/enterprise-teams#get-an-enterprise-team
//
//meta:operation GET /enterprises/{enterprise}/teams/{team_slug}
func (s *EnterpriseService) GetTeam(ctx context.Context, enterprise, teamSlug string) (*EnterpriseTeam, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/teams/%v", enterprise, teamSlug)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var team *EnterpriseTeam
	resp, err := s.client.Do(ctx, req, &team)
	if err != nil {
		return nil, resp, err
	}

	return team, resp, nil
}

// UpdateTeam updates a team in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/rest/enterprise-teams/enterprise-teams#update-an-enterprise-team
//
//meta:operation PATCH /enterprises/{enterprise}/teams/{team_slug}
func (s *EnterpriseService) UpdateTeam(ctx context.Context, enterprise, teamSlug string, team EnterpriseTeamCreateRequest) (*EnterpriseTeam, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/teams/%v", enterprise, teamSlug)

	req, err := s.client.NewRequest("PATCH", u, team)
	if err != nil {
		return nil, nil, err
	}

	var updatedTeam *EnterpriseTeam
	resp, err := s.client.Do(ctx, req, &updatedTeam)
	if err != nil {
		return nil, resp, err
	}

	return updatedTeam, resp, nil
}

// DeleteTeam deletes a team in an enterprise.
//
// GitHub API docs: https://docs.github.com/en/rest/enterprise-teams/enterprise-teams#delete-an-enterprise-team
//
//meta:operation DELETE /enterprises/{enterprise}/teams/{team_slug}
func (s *EnterpriseService) DeleteTeam(ctx context.Context, enterprise, teamSlug string) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/teams/%v", enterprise, teamSlug)

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
