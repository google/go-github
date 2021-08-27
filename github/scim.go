package github

import (
	"context"
	"encoding/json"
	"fmt"
)

// SCIMService provides access to SCIM related functions in the
// GitHub API.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/scim
type SCIMService service

type SCIMUserAttributes struct {
	UserName    string            `json:"user_name"`              // Configured by the admin. Could be an email, login, or username. (Required.)
	Name        SCIMUserName      `json:"name"`                   // (Required.)
	DisplayName *string           `json:"display_name,omitempty"` // The name of the user, suitable for display to end-users. (Optional.)
	Emails      []*SCIMUserEmails `json:"email"`                  // User emails. (Required.)
	Schemas     []*string         `json:"schemas,omitempty"`      // (Optional.)
	ExternalID  *string           `json:"external_id,omitempty"`  // (Optional.)
	Groups      []*string         `json:"groups,omitempty"`       // (Optional.)
	Active      bool              `json:"active,omitempty"`       // (Optional.)
}

type SCIMUserName struct {
	GivenName  string  `json:"given_name"`          // The first name of the user. (Required.)
	FamilyName string  `json:"family_name"`         // The last name of the user. (Required.)
	Formatted  *string `json:"formatted,omitempty"` // (Optional.)
}

type SCIMUserEmails struct {
	Value   string  `json:"value"`             // (Required.)
	Primary bool    `json:"primary,omitempty"` // (Optional.)
	Type    *string `json:"type,omitempty"`    // (Optional.)
}

type ListSCIMProvisionedIdentitiesOptions struct {
	StartIndex *int    `json:"start_index,omitempty"` // Used for pagination: the index of the first result to return. (Optional.)
	Count      *int    `json:"count,omitempty"`       // Used for pagination: the number of results to return. (Optional.)
	Filter     *string `json:"filter,omitempty"`      // Filters results using the equals query parameter operator (eq). You can filter results that are equal to id, userName, emails, and external_id. For example, to search for an identity with the userName Octocat, you would use this query: ?filter=userName%20eq%20\"Octocat\". To filter results for the identity with the email octocat@github.com, you would use this query: ?filter=emails%20eq%20\"octocat@github.com\". (Optional.)
}

// ListSCIMProvisionedIdentities lists SCIM provisioned identities.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/scim#list-scim-provisioned-identities
func (s *SCIMService) ListSCIMProvisionedIdentities(ctx context.Context, org string, opts *ListSCIMProvisionedIdentitiesOptions) (*Response, error) {
	u := fmt.Sprintf("scim/v2/organizations/%v/Users", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// ProvisionAndInviteSCIMUser provisions organization membership for a user, and send an activation email to the email address.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/scim#provision-and-invite-a-scim-user
func (s *SCIMService) ProvisionAndInviteSCIMUser(ctx context.Context, org string, opts *SCIMUserAttributes) (*Response, error) {
	u := fmt.Sprintf("scim/v2/organizations/%v/Users", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// GetSCIMProvisioningInfoForUser returns SCIM provisioning information for a user
//
// GitHub API docs: https://docs.github.com/en/rest/reference/scim#get-scim-provisioning-information-for-a-user
func (s *SCIMService) GetSCIMProvisioningInfoForUser(ctx context.Context, org, scimUserID string) (*Response, error) {
	u := fmt.Sprintf("scim/v2/organizations/%v/Users/%v", org, scimUserID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// UpdateProvisionedOrgMembership updates a provisioned organization membership.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/scim#update-a-provisioned-organization-membership
func (s *SCIMService) UpdateProvisionedOrgMembership(ctx context.Context, org, scimUserID string, opts *SCIMUserAttributes) (*Response, error) {
	u := fmt.Sprintf("scim/v2/organizations/%v/Users/%v", org, scimUserID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

type UpdateAttributeForSCIMUserOptions struct {
	Schemas    []*string                            `json:"schemas"`    // (Optional.)
	Operations UpdateAttributeForSCIMUserOperations `json:"operations"` // Set of operations to be performed. (Required.)
}

type UpdateAttributeForSCIMUserOperations struct {
	Op    string          `json:"op"`              // (Required.)
	Path  *string         `json:"path,omitempty"`  // (Optional.)
	Value json.RawMessage `json:"value,omitempty"` // (Optional.)
}

// UpdateAttributeForSCIMUser updates an attribute for an SCIM user.
//
// GitHub API docs: https://docs.github.com/en/rest/reference/scim#update-an-attribute-for-a-scim-user
func (s *SCIMService) UpdateAttributeForSCIMUser(ctx context.Context, org, scimUserID string, opts *UpdateAttributeForSCIMUserOptions) (*Response, error) {
	u := fmt.Sprintf("scim/v2/organizations/%v/Users/%v", org, scimUserID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest("PATCH", u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// DeleteSCIMUserFromOrg deletes SCIM user from an organization
//
// GitHub API docs: https://docs.github.com/en/rest/reference/scim#delete-a-scim-user-from-an-organization
func (s *SCIMService) DeleteSCIMUserFromOrg(ctx context.Context, org, scimUserID string) (*Response, error) {
	u := fmt.Sprintf("scim/v2/organizations/%v/Users/%v", org, scimUserID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}
