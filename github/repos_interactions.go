// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// RepositoryInteraction represents the interaction
// restrictions for a repository
type RepositoryInteraction struct {
	Limit     *string    `json:"limit,omitempty"`
	Origin    *string    `json:"origin,omitempty"`
	ExpiresAt *Timestamp `json:"expires_at,omitempty"`
}

// GetInteractions fetches the interaction restrictions for a repository.
//
// GitHub API docs: https://developer.github.com/v3/interactions/repos/#get-interaction-restrictions-for-a-repository
func (s *RepositoriesService) GetInteractions(ctx context.Context, owner string, repo string) (*RepositoryInteraction, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/interaction-limits", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeRepositoryInteractionsPreview)

	repositoryInteractions := new(RepositoryInteraction)

	resp, err := s.client.Do(ctx, req, repositoryInteractions)
	if err != nil {
		return nil, resp, err
	}

	return repositoryInteractions, resp, nil
}
