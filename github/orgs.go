// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package github

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// OrganizationsService provides access to the organization related functions
// in the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/orgs/
type OrganizationsService struct {
	client *Client
}

type Organization struct {
	Login     string `json:"login,omitempty"`
	ID        int    `json:"id,omitempty"`
	URL       string `json:"url,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Location  string `json:"location,omitempty"`
}

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
func (s *OrganizationsService) List(user string, opt *ListOptions) ([]Organization, error) {
	var url_ string
	if user != "" {
		url_ = fmt.Sprintf("users/%v/orgs", user)
	} else {
		url_ = "user/orgs"
	}
	if opt != nil {
		params := url.Values{
			"page": []string{strconv.Itoa(opt.Page)},
		}
		url_ += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	orgs := new([]Organization)
	_, err = s.client.Do(req, orgs)
	return *orgs, err
}

// Get an organization.
func (s *OrganizationsService) Get(org string) (*Organization, error) {
	url_ := fmt.Sprintf("orgs/%v", org)
	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	organization := new(Organization)
	_, err = s.client.Do(req, organization)
	return organization, err
}

// Edit an organization.
func (s *OrganizationsService) Edit(name string, org *Organization) (*Organization, error) {
	url_ := fmt.Sprintf("orgs/%v", name)
	req, err := s.client.NewRequest("PATCH", url_, org)
	if err != nil {
		return nil, err
	}

	updatedOrg := new(Organization)
	_, err = s.client.Do(req, updatedOrg)
	return updatedOrg, err
}

// List the members for an organization.  If the authenticated user is an owner
// of the organization, this will return concealed and public members,
// otherwise it will only return public members.
func (s *OrganizationsService) ListMembers(org string) ([]User, error) {
	url_ := fmt.Sprintf("orgs/%v/members", org)
	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	members := new([]User)
	_, err = s.client.Do(req, members)
	return *members, err
}

// List the public members for an organization.
func (s *OrganizationsService) ListPublicMembers(org string) ([]User, error) {
	url_ := fmt.Sprintf("orgs/%v/public_members", org)
	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	members := new([]User)
	_, err = s.client.Do(req, members)
	return *members, err
}

// CheckMembership checks if a user is a member of an organization.
func (s *OrganizationsService) CheckMembership(org, user string) (bool, error) {
	url_ := fmt.Sprintf("orgs/%v/members/%v", org, user)
	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return false, err
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		if err, ok := err.(*ErrorResponse); ok && err.Response.StatusCode == http.StatusNotFound {
			// The user is not a member of the org. In this one case, we do not pass
			// the error through.
			return false, nil
		} else {
			// some other real error occurred
			return false, err
		}
	}

	return err == nil, err
}

// CheckPublicMembership checks if a user is a public member of an organization.
func (s *OrganizationsService) CheckPublicMembership(org, user string) (bool, error) {
	url_ := fmt.Sprintf("orgs/%v/public_members/%v", org, user)
	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return false, err
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		if err, ok := err.(*ErrorResponse); ok && err.Response.StatusCode == http.StatusNotFound {
			// The user is not a member of the org. In this one case, we do not pass
			// the error through.
			return false, nil
		} else {
			// some other real error occurred
			return false, err
		}
	}

	return true, nil
}

// RemoveMember removes a user from all teams of an organization.
func (s *OrganizationsService) RemoveMember(org, user string) error {
	url_ := fmt.Sprintf("orgs/%v/members/%v", org, user)
	req, err := s.client.NewRequest("DELETE", url_, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// Publicize a user's membership in an organization.
func (s *OrganizationsService) PublicizeMembership(org, user string) error {
	url_ := fmt.Sprintf("orgs/%v/public_members/%v", org, user)
	req, err := s.client.NewRequest("PUT", url_, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// Conceal a user's membership in an organization.
func (s *OrganizationsService) ConcealMembership(org, user string) error {
	url_ := fmt.Sprintf("orgs/%v/public_members/%v", org, user)
	req, err := s.client.NewRequest("DELETE", url_, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// List the teams for an organization.
func (s *OrganizationsService) ListTeams(org string) ([]Team, error) {
	url_ := fmt.Sprintf("orgs/%v/teams", org)
	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	teams := new([]Team)
	_, err = s.client.Do(req, teams)
	return *teams, err
}

// Add a user to a team.
func (s *OrganizationsService) AddTeamMember(team int, user string) error {
	url_ := fmt.Sprintf("teams/%v/members/%v", team, user)
	req, err := s.client.NewRequest("PUT", url_, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// Remove a user from a team.
func (s *OrganizationsService) RemoveTeamMember(team int, user string) error {
	url_ := fmt.Sprintf("teams/%v/members/%v", team, user)
	req, err := s.client.NewRequest("DELETE", url_, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}
