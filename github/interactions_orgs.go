// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// GetRestrictionsForOrg fetches the interaction restrictions for an organization.
//
// GitHub API docs: https://developer.github.com/v3/interactions/orgs/#get-interaction-restrictions-for-an-organization
func (s *InteractionsService) GetRestrictionsForOrg(ctx context.Context, organization string) (*InteractionRestriction, *Response, error) {
	u := fmt.Sprintf("orgs/%v/interaction-limits", organization)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeInteractionRestrictionsPreview)

	organizationInteractions := new(InteractionRestriction)

	resp, err := s.client.Do(ctx, req, organizationInteractions)
	if err != nil {
		return nil, resp, err
	}

	return organizationInteractions, resp, nil
}

// UpdateRestrictionsForOrg adds or updates the interaction restrictions for an organization.
//
// GitHub API docs: https://developer.github.com/v3/interactions/orgs/#add-or-update-interaction-restrictions-for-an-organization
func (s *InteractionsService) UpdateRestrictionsForOrg(ctx context.Context, organization string, interaction *InteractionRestriction) (*InteractionRestriction, *Response, error) {
	u := fmt.Sprintf("orgs/%v/interaction-limits", organization)
	req, err := s.client.NewRequest("PUT", u, interaction)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeInteractionRestrictionsPreview)

	organizationInteractions := new(InteractionRestriction)

	resp, err := s.client.Do(ctx, req, organizationInteractions)
	if err != nil {
		return nil, resp, err
	}

	return organizationInteractions, resp, nil
}

// RemoveRestrictionsFromOrg removes the interaction restrictions for an organization.
//
// GitHub API docs: https://developer.github.com/v3/interactions/orgs/#remove-interaction-restrictions-for-an-organization
func (s *InteractionsService) RemoveRestrictionsFromOrg(ctx context.Context, organization string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/interaction-limits", organization)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeInteractionRestrictionsPreview)

	return s.client.Do(ctx, req, nil)
}
