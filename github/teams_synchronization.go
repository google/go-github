// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// IDPGroupList represents a list of external identity provider (IDP) groups.
type IDPGroupList struct {
	Groups []*IDPGroup `json:"groups"`
}

// IDPGroup represents an external identity provider (IDP) group.
type IDPGroup struct {
	GroupID          *string `json:"group_id,omitempty"`
	GroupName        *string `json:"group_name,omitempty"`
	GroupDescription *string `json:"group_description,omitempty"`
}

// ListIDPGroupsInOrganization lists IDP groups available in an organization.
//
// GitHub API docs: https://developer.github.com/v3/teams/team_sync/#list-idp-groups-in-an-organization
func (s *TeamsService) ListIDPGroupsInOrganization(ctx context.Context, org string, opt *ListOptions) (*IDPGroupList, *Response, error) {
	u := fmt.Sprintf("orgs/%v/team-sync/groups", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	groups := new(IDPGroupList)
	resp, err := s.client.Do(ctx, req, groups)
	if err != nil {
		return nil, resp, err
	}
	return groups, resp, nil
}

// ListIDPGroupsForTeamByID lists IDP groups connected to a team on GitHub
// given a team ID and organization ID.
//
// GitHub API docs: https://developer.github.com/v3/teams/team_sync/#list-idp-groups-for-a-team
func (s *TeamsService) ListIDPGroupsForTeamByID(ctx context.Context, orgID, teamID int64) (*IDPGroupList, *Response, error) {
	u := fmt.Sprintf("organizations/%v/team/%v/team-sync/group-mappings", orgID, teamID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	groups := new(IDPGroupList)
	resp, err := s.client.Do(ctx, req, groups)
	if err != nil {
		return nil, resp, err
	}
	return groups, resp, err
}

// ListIDPGroupsForTeamByName lists IDP groups connected to a team on GitHub
// given a team slug and orgnization name.
//
// GitHub API docs: https://developer.github.com/v3/teams/team_sync/#list-idp-groups-for-a-team
func (s *TeamsService) ListIDPGroupsForTeamByName(ctx context.Context, org, slug string) (*IDPGroupList, *Response, error) {
	u := fmt.Sprintf("orgs/%v/teams/%v/team-sync/group-mappings", org, slug)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	groups := new(IDPGroupList)
	resp, err := s.client.Do(ctx, req, groups)
	if err != nil {
		return nil, resp, err
	}
	return groups, resp, err
}

// CreateOrUpdateIDPGroupConnectionsByID creates, updates, or removes a connection between a team
// and an IDP group.  Identifies organization and team by ID.
//
// GitHub API docs: https://developer.github.com/v3/teams/team_sync/#create-or-update-idp-group-connections
func (s *TeamsService) CreateOrUpdateIDPGroupConnectionsByID(ctx context.Context, orgID, teamID int64, opt IDPGroupList) (*IDPGroupList, *Response, error) {
	u := fmt.Sprintf("organizations/%v/team/%v/team-sync/group-mappings", orgID, teamID)

	req, err := s.client.NewRequest("PATCH", u, opt)
	if err != nil {
		return nil, nil, err
	}

	groups := new(IDPGroupList)
	resp, err := s.client.Do(ctx, req, groups)
	if err != nil {
		return nil, resp, err
	}

	return groups, resp, nil
}

// CreateOrUpdateIDPGroupConnectionsByName creates, updates, or removes a connection between a team
// and an IDP group.  Identifies organization by name and team by slug.
//
// GitHub API docs: https://developer.github.com/v3/teams/team_sync/#create-or-update-idp-group-connections
func (s *TeamsService) CreateOrUpdateIDPGroupConnectionsByName(ctx context.Context, org, slug string, opt IDPGroupList) (*IDPGroupList, *Response, error) {
	u := fmt.Sprintf("orgs/%v/teams/%v/team-sync/group-mappings", org, slug)

	req, err := s.client.NewRequest("PATCH", u, opt)
	if err != nil {
		return nil, nil, err
	}

	groups := new(IDPGroupList)
	resp, err := s.client.Do(ctx, req, groups)
	if err != nil {
		return nil, resp, err
	}

	return groups, resp, nil
}
