// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "context"

// createOrgRequest is a subset of Organization and is used internally
// by Create to pass only the known fields for the endpoint
//
type createOrgRequest struct {
	Login *string `json:"login"`
	Admin *string `json:"admin"`
}

// CreateOrg creates a new organization in Github Enterprise.
//
// Note that only a subset of the org fields are used and org must
// not be nil
//
// Github Enterprise API docs: https://developer.github.com/enterprise/2.16/v3/enterprise-admin/orgs/#create-an-organization
func (s *AdminService) CreateOrg(ctx context.Context, org *Organization, admin string) (*Organization, *Response, error) {
	u := "admin/organizations"

	orgReq := &createOrgRequest{
		Login: org.Login,
		Admin: &admin,
	}

	req, err := s.client.NewRequest("POST", u, orgReq)
	if err != nil {
		return nil, nil, err
	}

	o := new(Organization)
	resp, err := s.client.Do(ctx, req, o)
	if err != nil {
		return nil, resp, err
	}

	return o, resp, nil
}
