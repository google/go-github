// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
)

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
