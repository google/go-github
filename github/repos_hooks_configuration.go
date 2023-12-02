// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// GetHookConfiguration returns the configuration for the specified repository webhook.
//
// GitHub API docs: https://docs.github.com/rest/repos/webhooks#get-a-webhook-configuration-for-a-repository
//
//meta:operation GET /repos/{owner}/{repo}/hooks/{hook_id}/config
func (s *RepositoriesService) GetHookConfiguration(ctx context.Context, owner, repo string, id int64) (*HookConfig, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/hooks/%v/config", owner, repo, id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	config := new(HookConfig)
	resp, err := s.client.Do(ctx, req, config)
	if err != nil {
		return nil, resp, err
	}

	return config, resp, nil
}

// EditHookConfiguration updates the configuration for the specified repository webhook.
//
// GitHub API docs: https://docs.github.com/rest/repos/webhooks#update-a-webhook-configuration-for-a-repository
//
//meta:operation PATCH /repos/{owner}/{repo}/hooks/{hook_id}/config
func (s *RepositoriesService) EditHookConfiguration(ctx context.Context, owner, repo string, id int64, config *HookConfig) (*HookConfig, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/hooks/%v/config", owner, repo, id)
	req, err := s.client.NewRequest("PATCH", u, config)
	if err != nil {
		return nil, nil, err
	}

	c := new(HookConfig)
	resp, err := s.client.Do(ctx, req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}
