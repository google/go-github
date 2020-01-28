// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// TeamListTeamMembersOptions specifies the optional parameters to the
// TeamsService.ListTeamMembers method.
type TeamListTeamMembersOptions struct {
	// Role filters members returned by their role in the team. Possible
	// values are "all", "member", "maintainer". Default is "all".
	Role string `url:"role,omitempty"`

	ListOptions
}

// ListTeamMembersByID lists all of the users who are members of the specified
// team given the team ID and organization ID.
//
// GitHub API docs: https://developer.github.com/v3/teams/members/#list-team-members
func (s *TeamsService) ListTeamMembersByID(ctx context.Context, orgID, teamID int64, opt *TeamListTeamMembersOptions) ([]*User, *Response, error) {
	u := fmt.Sprintf("organizations/%v/team/%v/members", orgID, teamID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var members []*User
	resp, err := s.client.Do(ctx, req, &members)
	if err != nil {
		return nil, resp, err
	}

	return members, resp, nil
}

// ListTeamMembersByName lists all of the users who are members of the specified
// team given the team slug and organization name.
//
// GitHub API docs: https://developer.github.com/v3/teams/members/#list-team-members
func (s *TeamsService) ListTeamMembersByName(ctx context.Context, org, slug string, opt *TeamListTeamMembersOptions) ([]*User, *Response, error) {
	u := fmt.Sprintf("orgs/%v/teams/%v/members", org, slug)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var members []*User
	resp, err := s.client.Do(ctx, req, &members)
	if err != nil {
		return nil, resp, err
	}

	return members, resp, nil
}

// GetTeamMembershipByID returns the membership status for a user in a team
// given the team ID and organization ID.
//
// GitHub API docs: https://developer.github.com/v3/teams/members/#get-team-membership
func (s *TeamsService) GetTeamMembershipByID(ctx context.Context, orgID, teamID int64, user string) (*Membership, *Response, error) {
	u := fmt.Sprintf("organizations/%v/team/%v/memberships/%v", orgID, teamID, user)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	t := new(Membership)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// GetTeamMembershipByName returns the membership status for a user in a team
// given the team slug and organization name..
//
// GitHub API docs: https://developer.github.com/v3/teams/members/#get-team-membership
func (s *TeamsService) GetTeamMembershipByName(ctx context.Context, org, slug, user string) (*Membership, *Response, error) {
	u := fmt.Sprintf("orgs/%v/teams/%v/memberships/%v", org, slug, user)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	t := new(Membership)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// TeamAddTeamMembershipOptions specifies the optional
// parameters to the TeamsService.AddTeamMembership method.
type TeamAddTeamMembershipOptions struct {
	// Role specifies the role the user should have in the team. Possible
	// values are:
	//     member - a normal member of the team
	//     maintainer - a team maintainer. Able to add/remove other team
	//                  members, promote other team members to team
	//                  maintainer, and edit the team’s name and description
	//
	// Default value is "member".
	Role string `json:"role,omitempty"`
}

// AddTeamMembershipByID adds or invites a user to a team by the team ID
// and organization ID.
//
// In order to add a membership between a user and a team, the authenticated
// user must have 'admin' permissions to the team or be an owner of the
// organization that the team is associated with.
//
// If the user is already a part of the team's organization (meaning they're on
// at least one other team in the organization), this endpoint will add the
// user to the team.
//
// If the user is completely unaffiliated with the team's organization (meaning
// they're on none of the organization's teams), this endpoint will send an
// invitation to the user via email. This newly-created membership will be in
// the "pending" state until the user accepts the invitation, at which point
// the membership will transition to the "active" state and the user will be
// added as a member of the team.
//
// GitHub API docs: https://developer.github.com/v3/teams/members/#add-or-update-team-membership
func (s *TeamsService) AddTeamMembershipByID(ctx context.Context, orgID, teamID int64, user string, opt *TeamAddTeamMembershipOptions) (*Membership, *Response, error) {
	u := fmt.Sprintf("organizations/%v/team/%v/memberships/%v", orgID, teamID, user)
	req, err := s.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	t := new(Membership)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// AddTeamMembershipByName adds or invites a user to a team by the team slug
// and organiation name.
//
// In order to add a membership between a user and a team, the authenticated
// user must have 'admin' permissions to the team or be an owner of the
// organization that the team is associated with.
//
// If the user is already a part of the team's organization (meaning they're on
// at least one other team in the organization), this endpoint will add the
// user to the team.
//
// If the user is completely unaffiliated with the team's organization (meaning
// they're on none of the organization's teams), this endpoint will send an
// invitation to the user via email. This newly-created membership will be in
// the "pending" state until the user accepts the invitation, at which point
// the membership will transition to the "active" state and the user will be
// added as a member of the team.
//
// GitHub API docs: https://developer.github.com/v3/teams/members/#add-or-update-team-membership
func (s *TeamsService) AddTeamMembershipByName(ctx context.Context, org, slug, user string, opt *TeamAddTeamMembershipOptions) (*Membership, *Response, error) {
	u := fmt.Sprintf("orgs/%v/teams/%v/memberships/%v", org, slug, user)
	req, err := s.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	t := new(Membership)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// RemoveTeamMembershipByID removes a user from a team given the team ID and
// organization ID.
//
// GitHub API docs: https://developer.github.com/v3/teams/members/#remove-team-membership
func (s *TeamsService) RemoveTeamMembershipByID(ctx context.Context, orgID, teamID int64, user string) (*Response, error) {
	u := fmt.Sprintf("organizations/%v/team/%v/memberships/%v", orgID, teamID, user)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// RemoveTeamMembershipByName removes a user from a team given the team slug
// and organization name.
//
// GitHub API docs: https://developer.github.com/v3/teams/members/#remove-team-membership
func (s *TeamsService) RemoveTeamMembershipByName(ctx context.Context, org, slug, user string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/teams/%v/memberships/%v", org, slug, user)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// ListPendingTeamInvitationsByID get pending invitaion list in team given the
// team ID and organization ID.
// Warning: The API may change without advance notice during the preview period.
// Preview features are not supported for production use.
//
// GitHub API docs: https://developer.github.com/v3/teams/members/#list-pending-team-invitations
func (s *TeamsService) ListPendingTeamInvitationsByID(ctx context.Context, orgID, teamID int64, opt *ListOptions) ([]*Invitation, *Response, error) {
	u := fmt.Sprintf("organizations/%v/team/%v/invitations", orgID, teamID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var pendingInvitations []*Invitation
	resp, err := s.client.Do(ctx, req, &pendingInvitations)
	if err != nil {
		return nil, resp, err
	}

	return pendingInvitations, resp, nil
}

// ListPendingTeamInvitationsByName get pending invitaion list in team given the
// team slug and organization name.
// Warning: The API may change without advance notice during the preview period.
// Preview features are not supported for production use.
//
// GitHub API docs: https://developer.github.com/v3/teams/members/#list-pending-team-invitations
func (s *TeamsService) ListPendingTeamInvitationsByName(ctx context.Context, org, slug string, opt *ListOptions) ([]*Invitation, *Response, error) {
	u := fmt.Sprintf("orgs/%v/teams/%v/invitations", org, slug)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var pendingInvitations []*Invitation
	resp, err := s.client.Do(ctx, req, &pendingInvitations)
	if err != nil {
		return nil, resp, err
	}

	return pendingInvitations, resp, nil
}
