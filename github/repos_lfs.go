// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// EnableLfs turns the LFS feature ON for selected repo.
//
// GitHub API docs: https://docs.github.com/en/rest/repos/lfs#enable-git-lfs-for-a-repository
func (s *RepositoriesService) EnableLfs(ctx context.Context, owner string, repo string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/lfs", owner, repo)

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

// DisableLfs turns the LFS feature OFF for selected repo.
//
// Github API docs: https://docs.github.com/en/rest/repos/lfs#disable-git-lfs-for-a-repository
func (s *RepositoriesService) DisableLfs(ctx context.Context, owner string, repo string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/lfs", owner, repo)

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
