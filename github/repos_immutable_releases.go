// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// checkImmutableReleases represents whether immutable releases are enabled for a repository.
type checkImmutableReleases struct {
	Enabled bool `json:"enabled,omitempty"`
}

// EnableImmutableReleases enables immutable releases for a repository.
//
// GitHub API docs: https://docs.github.com/rest/repos/repos#enable-immutable-releases-for-a-repository
//
//meta:operation PUT /repos/{owner}/{repo}/immutable-releases
func (s *RepositoriesService) EnableImmutableReleases(ctx context.Context, owner, repo string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/immutable-releases", owner, repo)

	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DisableImmutableReleases disables immutable releases for a repository.
//
// GitHub API docs: https://docs.github.com/rest/repos/repos#disable-immutable-releases-for-a-repository
//
//meta:operation DELETE /repos/{owner}/{repo}/immutable-releases
func (s *RepositoriesService) DisableImmutableReleases(ctx context.Context, owner, repo string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/immutable-releases", owner, repo)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// IsImmutableReleasesEnabled checks if immutable releases are enabled for
// the repository and returns a boolean indicating the status.
//
// GitHub API docs: https://docs.github.com/rest/repos/repos#check-if-immutable-releases-are-enabled-for-a-repository
//
//meta:operation GET /repos/{owner}/{repo}/immutable-releases
func (s *RepositoriesService) IsImmutableReleasesEnabled(ctx context.Context, owner, repo string) (bool, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/immutable-releases", owner, repo)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return false, nil, err
	}

	immutableReleases := new(checkImmutableReleases)
	resp, err := s.client.Do(ctx, req, immutableReleases)
	return immutableReleases.Enabled, resp, err
}
