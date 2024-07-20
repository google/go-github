// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
)

// PersonalAccessToken represents the minimal representation of an organization programmatic access grant.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/personal-access-tokens?apiVersion=2022-11-28
type PersonalAccessToken struct {
	// "Unique identifier of the fine-grained personal access token.
	// The `pat_id` used to get details about an approved fine-grained personal access token.
	ID *int64 `json:"id"`

	// Owner is the GitHub user associated with the token.
	Owner *User `json:"owner"`

	// RepositorySelection is the type of repository selection requested.
	// Possible values are: "none", "all", "subset".
	RepositorySelection *string `json:"repository_selection"`

	// URL to the list of repositories the fine-grained personal access token can access.
	// Only follow when `repository_selection` is `subset`.
	RepositoriesURL *string `json:"repositories_url"`

	// Permissions are the permissions requested, categorized by type.
	Permissions *PersonalAccessTokenPermissions `json:"permissions"`

	// Date and time when the fine-grained personal access token was approved to access the organization.
	AccessGrantedAt *Timestamp `json:"access_granted_at"`

	// Whether the associated fine-grained personal access token has expired.
	TokenExpired *bool `json:"token_expired"`

	// Date and time when the associated fine-grained personal access token expires.
	TokenExpiresAt *Timestamp `json:"token_expires_at"`

	// Date and time when the associated fine-grained personal access token was last used for authentication.
	TokenLastUsedAt *Timestamp `json:"token_last_used_at"`
}

// ListFineGrainedPersonalAccessTokens lists approved fine-grained personal access tokens owned by organization members that can access organization resources.
// Only GitHub Apps can call this API, using the `Personal access tokens` organization permissions (read).
//
// GitHub API docs: https://docs.github.com/rest/orgs/personal-access-tokens#list-fine-grained-personal-access-tokens-with-access-to-organization-resources
//
//meta:operation GET /orgs/{org}/personal-access-tokens
func (s *OrganizationsService) ListFineGrainedPersonalAccessTokens(ctx context.Context, org string, opts *ListOptions) ([]*PersonalAccessToken, *Response, error) {
	u := fmt.Sprintf("orgs/%v/personal-access-tokens", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, &opts)
	if err != nil {
		return nil, nil, err
	}

	var pats []*PersonalAccessToken

	resp, err := s.client.Do(ctx, req, &pats)
	if err != nil {
		return nil, resp, err
	}

	return pats, resp, nil
}

// ReviewPersonalAccessTokenRequestOptions specifies the parameters to the ReviewPersonalAccessTokenRequest method.
type ReviewPersonalAccessTokenRequestOptions struct {
	Action string  `json:"action"`
	Reason *string `json:"reason,omitempty"`
}

// ReviewPersonalAccessTokenRequest approves or denies a pending request to access organization resources via a fine-grained personal access token.
// Only GitHub Apps can call this API, using the `organization_personal_access_token_requests: write` permission.
// `action` can be one of `approve` or `deny`.
//
// GitHub API docs: https://docs.github.com/rest/orgs/personal-access-tokens#review-a-request-to-access-organization-resources-with-a-fine-grained-personal-access-token
//
//meta:operation POST /orgs/{org}/personal-access-token-requests/{pat_request_id}
func (s *OrganizationsService) ReviewPersonalAccessTokenRequest(ctx context.Context, org string, requestID int64, opts ReviewPersonalAccessTokenRequestOptions) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/personal-access-token-requests/%v", org, requestID)

	req, err := s.client.NewRequest(http.MethodPost, u, &opts)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
