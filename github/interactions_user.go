// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// UpdateUserInteractionRestrictionsRequest represents the request body for
// setting the interaction restrictions for the authenticated user's public repositories.
type UpdateUserInteractionRestrictionsRequest struct {
	// Limit specifies the group of GitHub users who can comment, open issues, or create pull requests
	// in the authenticated user's public repositories.
	// Possible values are: "existing_users", "contributors_only", "collaborators_only".
	Limit string `json:"limit"`

	// Expiry specifies the duration of the interaction restriction.
	// Possible values are: "one_day", "three_days", "one_week", "one_month", "six_months".
	// The default is "one_day".
	Expiry *string `json:"expiry,omitempty"`
}

// GetUserRestrictions fetches the interaction restrictions
// for the authenticated user's public repositories.
//
// GitHub API docs: https://docs.github.com/rest/interactions/user?apiVersion=2022-11-28#get-interaction-restrictions-for-your-public-repositories
//
//meta:operation GET /user/interaction-limits
func (s *InteractionsService) GetUserRestrictions(ctx context.Context) (*InteractionRestriction, *Response, error) {
	u := "user/interaction-limits"
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var userInteractions *InteractionRestriction
	resp, err := s.client.Do(req, &userInteractions)
	if err != nil {
		return nil, resp, err
	}

	return userInteractions, resp, nil
}

// UpdateUserRestrictions adds or updates the interaction restrictions
// for the authenticated user's public repositories.
//
// GitHub API docs: https://docs.github.com/rest/interactions/user?apiVersion=2022-11-28#set-interaction-restrictions-for-your-public-repositories
//
//meta:operation PUT /user/interaction-limits
func (s *InteractionsService) UpdateUserRestrictions(ctx context.Context, body UpdateUserInteractionRestrictionsRequest) (*InteractionRestriction, *Response, error) {
	u := "user/interaction-limits"
	req, err := s.client.NewRequest(ctx, "PUT", u, body)
	if err != nil {
		return nil, nil, err
	}

	var userInteractions *InteractionRestriction
	resp, err := s.client.Do(req, &userInteractions)
	if err != nil {
		return nil, resp, err
	}

	return userInteractions, resp, nil
}

// RemoveUserRestrictions removes the interaction restrictions
// from the authenticated user's public repositories.
//
// GitHub API docs: https://docs.github.com/rest/interactions/user?apiVersion=2022-11-28#remove-interaction-restrictions-from-your-public-repositories
//
//meta:operation DELETE /user/interaction-limits
func (s *InteractionsService) RemoveUserRestrictions(ctx context.Context) (*Response, error) {
	u := "user/interaction-limits"
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
