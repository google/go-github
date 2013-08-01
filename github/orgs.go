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
