// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// OrganizationCustomRepoRoles represents custom repository roles available in specified organization.
type OrganizationCustomRepoRoles struct {
	TotalCount      *int               `json:"total_count,omitempty"`
	CustomRepoRoles []*CustomRepoRoles `json:"custom_roles,omitempty"`
}

// CustomRepoRoles represents custom repository roles for an organization.
// See https://docs.github.com/en/enterprise-cloud@latest/organizations/managing-peoples-access-to-your-organization-with-roles/managing-custom-repository-roles-for-an-organization
// for more information.
type CustomRepoRoles struct {
	ID   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// ListCustomRepoRoles lists the custom repository roles available in this organization.
// In order to see custom repository roles in an organization, the authenticated user must be an organization owner.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/custom-roles#list-custom-repository-roles-in-an-organization
func (s *OrganizationsService) ListCustomRepoRoles(ctx context.Context, org string) (*OrganizationCustomRepoRoles, *Response, error) {
	u := fmt.Sprintf("orgs/%v/custom_roles", org)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	customRepoRoles := new(OrganizationCustomRepoRoles)
	resp, err := s.client.Do(ctx, req, customRepoRoles)
	if err != nil {
		return nil, resp, err
	}

	return customRepoRoles, resp, nil
}
