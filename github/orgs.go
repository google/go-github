// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// OrganizationsService provides access to the organization related functions
// in the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/orgs/
type OrganizationsService struct {
	client *Client
}

// Organization represents a GitHub organization account.
type Organization struct {
	Login     string     `json:"login,omitempty"`
	ID        int        `json:"id,omitempty"`
	URL       string     `json:"url,omitempty"`
	AvatarURL string     `json:"avatar_url,omitempty"`
	Location  string     `json:"location,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// Team represents a team within a GitHub organization.  Teams are used to
// manage access to an organization's repositories.
type Team struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	Slug         string `json:"slug,omitempty"`
	Permission   string `json:"permission,omitempty"`
	MembersCount int    `json:"members_count,omitempty"`
	ReposCount   int    `json:"repos_count,omitempty"`
}

// List the organizations for a user.  Passing the empty string will list
// organizations for the authenticated user.
//
// GitHub API docs: http://developer.github.com/v3/orgs/#list-user-organizations
func (s *OrganizationsService) List(user string, opt *ListOptions) ([]Organization, error) {
	var u string
	if user != "" {
		u = fmt.Sprintf("users/%v/orgs", user)
	} else {
		u = "user/orgs"
	}
	if opt != nil {
		params := url.Values{
			"page": []string{strconv.Itoa(opt.Page)},
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	orgs := new([]Organization)
	_, err = s.client.Do(req, orgs)
	return *orgs, err
}

// Get fetches an organization by name.
//
// GitHub API docs: http://developer.github.com/v3/orgs/#get-an-organization
func (s *OrganizationsService) Get(org string) (*Organization, error) {
	u := fmt.Sprintf("orgs/%v", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	organization := new(Organization)
	_, err = s.client.Do(req, organization)
	return organization, err
}

// Edit an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/#edit-an-organization
func (s *OrganizationsService) Edit(name string, org *Organization) (*Organization, error) {
	u := fmt.Sprintf("orgs/%v", name)
	req, err := s.client.NewRequest("PATCH", u, org)
	if err != nil {
		return nil, err
	}

	o := new(Organization)
	_, err = s.client.Do(req, o)
	return o, err
}

// ListMembers lists the members for an organization.  If the authenticated
// user is an owner of the organization, this will return both concealed and
// public members, otherwise it will only return public members.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#members-list
func (s *OrganizationsService) ListMembers(org string) ([]User, error) {
	u := fmt.Sprintf("orgs/%v/members", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	members := new([]User)
	_, err = s.client.Do(req, members)
	return *members, err
}

// ListPublicMembers lists the public members for an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#public-members-list
func (s *OrganizationsService) ListPublicMembers(org string) ([]User, error) {
	u := fmt.Sprintf("orgs/%v/public_members", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	members := new([]User)
	_, err = s.client.Do(req, members)
	return *members, err
}

// CheckMembership checks if a user is a member of an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#check-membership
func (s *OrganizationsService) CheckMembership(org, user string) (bool, error) {
	u := fmt.Sprintf("orgs/%v/members/%v", org, user)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return false, err
	}

	_, err = s.client.Do(req, nil)
	return parseBoolResponse(err)
}

// CheckPublicMembership checks if a user is a public member of an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#check-public-membership
func (s *OrganizationsService) CheckPublicMembership(org, user string) (bool, error) {
	u := fmt.Sprintf("orgs/%v/public_members/%v", org, user)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return false, err
	}

	_, err = s.client.Do(req, nil)
	return parseBoolResponse(err)
}

// RemoveMember removes a user from all teams of an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#remove-a-member
func (s *OrganizationsService) RemoveMember(org, user string) error {
	u := fmt.Sprintf("orgs/%v/members/%v", org, user)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// PublicizeMembership publicizes a user's membership in an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#publicize-a-users-membership
func (s *OrganizationsService) PublicizeMembership(org, user string) error {
	u := fmt.Sprintf("orgs/%v/public_members/%v", org, user)
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// ConcealMembership conceals a user's membership in an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#conceal-a-users-membership
func (s *OrganizationsService) ConcealMembership(org, user string) error {
	u := fmt.Sprintf("orgs/%v/public_members/%v", org, user)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// ListTeams lists all of the teams for an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#list-teams
func (s *OrganizationsService) ListTeams(org string) ([]Team, error) {
	u := fmt.Sprintf("orgs/%v/teams", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	teams := new([]Team)
	_, err = s.client.Do(req, teams)
	return *teams, err
}

// GetTeam fetches a team by ID.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#get-team
func (s *OrganizationsService) GetTeam(team int) (*Team, error) {
	u := fmt.Sprintf("teams/%v", team)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	t := new(Team)
	_, err = s.client.Do(req, t)
	return t, err
}

// CreateTeam creates a new team within an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#create-team
func (s *OrganizationsService) CreateTeam(org string, team *Team) (*Team, error) {
	u := fmt.Sprintf("orgs/%v/teams", org)
	req, err := s.client.NewRequest("POST", u, team)
	if err != nil {
		return nil, err
	}

	t := new(Team)
	_, err = s.client.Do(req, t)
	return t, err
}

// EditTeam edits a team.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#edit-team
func (s *OrganizationsService) EditTeam(id int, team *Team) (*Team, error) {
	u := fmt.Sprintf("teams/%v", id)
	req, err := s.client.NewRequest("PATCH", u, team)
	if err != nil {
		return nil, err
	}

	t := new(Team)
	_, err = s.client.Do(req, t)
	return t, err
}

// DeleteTeam deletes a team.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#delete-team
func (s *OrganizationsService) DeleteTeam(team int) error {
	u := fmt.Sprintf("teams/%v", team)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// ListTeamMembers lists all of the users who are members of the specified
// team.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#list-team-members
func (s *OrganizationsService) ListTeamMembers(team int) ([]User, error) {
	u := fmt.Sprintf("teams/%v/members", team)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	members := new([]User)
	_, err = s.client.Do(req, members)
	return *members, err
}

// CheckTeamMembership checks if a user is a member of the specified team.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#get-team-member
func (s *OrganizationsService) CheckTeamMembership(team int, user string) (bool, error) {
	u := fmt.Sprintf("teams/%v/members/%v", team, user)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return false, err
	}

	_, err = s.client.Do(req, nil)
	return parseBoolResponse(err)
}

// AddTeamMember adds a user to a team.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#add-team-member
func (s *OrganizationsService) AddTeamMember(team int, user string) error {
	u := fmt.Sprintf("teams/%v/members/%v", team, user)
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// RemoveTeamMember removes a user from a team.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#remove-team-member
func (s *OrganizationsService) RemoveTeamMember(team int, user string) error {
	u := fmt.Sprintf("teams/%v/members/%v", team, user)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// ListTeamRepos lists the repositories that the specified team has access to.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#list-team-repos
func (s *OrganizationsService) ListTeamRepos(team int) ([]Repository, error) {
	u := fmt.Sprintf("teams/%v/repos", team)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	repos := new([]Repository)
	_, err = s.client.Do(req, repos)
	return *repos, err
}

// CheckTeamRepo checks if a team manages the specified repository.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#get-team-repo
func (s *OrganizationsService) CheckTeamRepo(team int, owner string, repo string) (bool, error) {
	u := fmt.Sprintf("teams/%v/repos/%v/%v", team, owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return false, err
	}

	_, err = s.client.Do(req, nil)
	return parseBoolResponse(err)
}

// AddTeamRepo adds a repository to be managed by the specified team.  The
// specified repository must be owned by the organization to which the team
// belongs, or a direct fork of a repository owned by the organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#add-team-repo
func (s *OrganizationsService) AddTeamRepo(team int, owner string, repo string) error {
	u := fmt.Sprintf("teams/%v/repos/%v/%v", team, owner, repo)
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// RemoveTeamRepo removes a repository from being managed by the specified
// team.  Note that this does not delete the repository, it just removes it
// from the team.
//
// GitHub API docs: http://developer.github.com/v3/orgs/teams/#remove-team-repo
func (s *OrganizationsService) RemoveTeamRepo(team int, owner string, repo string) error {
	u := fmt.Sprintf("teams/%v/repos/%v/%v", team, owner, repo)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}
