// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// GetRestrictions fetches the interaction restrictions for a repository.
//
// GitHub API docs: https://developer.github.com/v3/interactions/repos/#get-interaction-restrictions-for-a-repository
func (s *InteractionsService) GetRestrictions(ctx context.Context, owner, repo string) (*InteractionRestriction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/interaction-limits", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeRepositoryInteractionsPreview)

	repositoryInteractions := new(InteractionRestriction)

	resp, err := s.client.Do(ctx, req, repositoryInteractions)
	if err != nil {
		return nil, resp, err
	}

	return repositoryInteractions, resp, nil
}

// UpdateRestrictions adds or updates the interaction restrictions for a repository.
//
// GitHub API docs: https://developer.github.com/v3/interactions/repos/#add-or-update-interaction-restrictions-for-a-repository
func (s *InteractionsService) UpdateRestrictions(ctx context.Context, owner, repo string, interaction *InteractionRestriction) (*InteractionRestriction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/interaction-limits", owner, repo)
	req, err := s.client.NewRequest("PUT", u, interaction)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeRepositoryInteractionsPreview)

	repositoryInteractions := new(InteractionRestriction)

	resp, err := s.client.Do(ctx, req, repositoryInteractions)
	if err != nil {
		return nil, resp, err
	}

	return repositoryInteractions, resp, nil
}

// RemoveRestrictions removes the interaction restrictions for a repository.
//
// GitHub API docs: https://developer.github.com/v3/interactions/repos/#remove-interaction-restrictions-for-a-repository
func (s *InteractionsService) RemoveRestrictions(ctx context.Context, owner, repo string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/interaction-limits", owner, repo)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeRepositoryInteractionsPreview)

	return s.client.Do(ctx, req, nil)
}
