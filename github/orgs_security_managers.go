// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

func (s *OrganizationsService) GetSecurityManagerRole(ctx context.Context, org string) (*CustomOrgRoles, *Response, error) {
	roles, resp, err := s.ListRoles(ctx, org)
	if err != nil {
		return nil, resp, err
	}

	for _, role := range roles.CustomRepoRoles {
		if *role.Name == "security_manager" {
			return role, resp, nil
		}
	}

	return nil, resp, fmt.Errorf("security manager role not found")
}

// ListSecurityManagerTeams lists all security manager teams for an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/organization-roles#list-teams-that-are-assigned-to-an-organization-role
//
//meta:operation GET /orgs/{org}/organization-roles/{security_manager_role_id}/teams
func (s *OrganizationsService) ListSecurityManagerTeams(ctx context.Context, org string) ([]*Team, *Response, error) {
	securityManagerRole, resp, err := s.GetSecurityManagerRole(ctx, org)
	if err != nil {
		return nil, resp, err
	}

	options := &ListOptions{PerPage: 100}
	securityManagerTeams := make([]*Team, 0)
	for {
		teams, resp, err := s.ListTeamsAssignedToOrgRole(ctx, org, securityManagerRole.GetID(), options)
		if err != nil {
			return nil, resp, err
		}

		securityManagerTeams = append(securityManagerTeams, teams...)
		if resp.NextPage == 0 {
			return securityManagerTeams, resp, nil
		}

		options.Page = resp.NextPage
	}
}

// AddSecurityManagerTeam adds a team to the list of security managers for an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/organization-roles#assign-an-organization-role-to-a-team
//
//meta:operation PUT /orgs/{org}/organization-roles/teams/{team_slug}/{security_manager_role_id}
func (s *OrganizationsService) AddSecurityManagerTeam(ctx context.Context, org, team string) (*Response, error) {
	securityManagerRole, resp, err := s.GetSecurityManagerRole(ctx, org)
	if err != nil {
		return resp, err
	}

	return s.AssignOrgRoleToTeam(ctx, org, team, securityManagerRole.GetID())
}

// RemoveSecurityManagerTeam removes a team from the list of security managers for an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/organization-roles#remove-an-organization-role-from-a-team
//
//meta:operation DELETE /orgs/{org}/organization-roles/teams/{team_slug}/{security_manager_role_id}
func (s *OrganizationsService) RemoveSecurityManagerTeam(ctx context.Context, org, team string) (*Response, error) {
	securityManagerRole, resp, err := s.GetSecurityManagerRole(ctx, org)
	if err != nil {
		return resp, err
	}

	return s.RemoveOrgRoleFromTeam(ctx, org, team, securityManagerRole.GetID())
}
