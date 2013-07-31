// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "fmt"

// ListMembers lists the members for an organization.  If the authenticated
// user is an owner of the organization, this will return both concealed and
// public members, otherwise it will only return public members.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#members-list
func (s *OrganizationsService) ListMembers(org string) ([]User, *Response, error) {
	u := fmt.Sprintf("orgs/%v/members", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	members := new([]User)
	resp, err := s.client.Do(req, members)
	return *members, resp, err
}

// ListPublicMembers lists the public members for an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#public-members-list
func (s *OrganizationsService) ListPublicMembers(org string) ([]User, *Response, error) {
	u := fmt.Sprintf("orgs/%v/public_members", org)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	members := new([]User)
	resp, err := s.client.Do(req, members)
	return *members, resp, err
}

// CheckMembership checks if a user is a member of an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#check-membership
func (s *OrganizationsService) CheckMembership(org, user string) (bool, *Response, error) {
	u := fmt.Sprintf("orgs/%v/members/%v", org, user)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return false, nil, err
	}

	resp, err := s.client.Do(req, nil)
	member, err := parseBoolResponse(err)
	return member, resp, err
}

// CheckPublicMembership checks if a user is a public member of an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#check-public-membership
func (s *OrganizationsService) CheckPublicMembership(org, user string) (bool, *Response, error) {
	u := fmt.Sprintf("orgs/%v/public_members/%v", org, user)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return false, nil, err
	}

	resp, err := s.client.Do(req, nil)
	member, err := parseBoolResponse(err)
	return member, resp, err
}

// RemoveMember removes a user from all teams of an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#remove-a-member
func (s *OrganizationsService) RemoveMember(org, user string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/members/%v", org, user)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// PublicizeMembership publicizes a user's membership in an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#publicize-a-users-membership
func (s *OrganizationsService) PublicizeMembership(org, user string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/public_members/%v", org, user)
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ConcealMembership conceals a user's membership in an organization.
//
// GitHub API docs: http://developer.github.com/v3/orgs/members/#conceal-a-users-membership
func (s *OrganizationsService) ConcealMembership(org, user string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/public_members/%v", org, user)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
