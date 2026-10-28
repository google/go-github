// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// GetRestrictionsForAuthenticatedUser fetches the interaction restrictions
// for the authenticated user's public repositories.
//
// GitHub API docs: https://docs.github.com/rest/interactions/user?apiVersion=2022-11-28#get-interaction-restrictions-for-your-public-repositories
//
//meta:operation GET /user/interaction-limits
func (s *InteractionsService) GetRestrictionsForAuthenticatedUser(ctx context.Context) (*InteractionRestriction, *Response, error) {
	u := "user/interaction-limits"
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeInteractionRestrictionsPreview)

	var userInteractions *InteractionRestriction
	resp, err := s.client.Do(req, &userInteractions)
	if err != nil {
		return nil, resp, err
	}

	return userInteractions, resp, nil
}

// UpdateRestrictionsForAuthenticatedUser adds or updates the interaction restrictions
// for the authenticated user's public repositories.
//
// limit specifies the group of GitHub users who can comment, open issues, or create pull requests
// in the authenticated user's public repositories.
// Possible values are: "existing_users", "contributors_only", "collaborators_only".
//
// GitHub API docs: https://docs.github.com/rest/interactions/user?apiVersion=2022-11-28#set-interaction-restrictions-for-your-public-repositories
//
//meta:operation PUT /user/interaction-limits
func (s *InteractionsService) UpdateRestrictionsForAuthenticatedUser(ctx context.Context, limit string) (*InteractionRestriction, *Response, error) {
	u := "user/interaction-limits"

	interaction := &InteractionRestriction{Limit: &limit}

	req, err := s.client.NewRequest(ctx, "PUT", u, interaction)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeInteractionRestrictionsPreview)

	var userInteractions *InteractionRestriction
	resp, err := s.client.Do(req, &userInteractions)
	if err != nil {
		return nil, resp, err
	}

	return userInteractions, resp, nil
}

// RemoveRestrictionsFromAuthenticatedUser removes the interaction restrictions
// from the authenticated user's public repositories.
//
// GitHub API docs: https://docs.github.com/rest/interactions/user?apiVersion=2022-11-28#remove-interaction-restrictions-from-your-public-repositories
//
//meta:operation DELETE /user/interaction-limits
func (s *InteractionsService) RemoveRestrictionsFromAuthenticatedUser(ctx context.Context) (*Response, error) {
	u := "user/interaction-limits"
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", mediaTypeInteractionRestrictionsPreview)

	return s.client.Do(req, nil)
}
