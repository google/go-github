package github

import (
	"context"
	"fmt"
)

// OrginizationCustomRoles represents custom repository roles available in specified organization
type OrginizationCustomRoles struct {
	TotalCount  *int           `json:"total_count,omitempty"`
	CustomRoles []*CustomRoles `json:"custom_roles,omitempty"`
}

// CustomRoles represents information of custom roles
type CustomRoles struct {
	ID   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// List Custome Roles in Org
//
// GitHub API docs: https://docs.github.com/en/rest/reference/orgs#custom-repository-roles
func (s *OrganizationsService) ListCustomRoles(ctx context.Context, organization_id string) (*OrginizationCustomRoles, *Response, error) {
	u := fmt.Sprintf("organizations/%v/custom_roles", organization_id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	custom_roles := new(OrginizationCustomRoles)
	resp, err := s.client.Do(ctx, req, custom_roles)
	if err != nil {
		return nil, resp, err
	}

	return custom_roles, resp, nil
}
