// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// GetHookConfiguration returns the configuration for the specified organization webhook.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/webhooks?apiVersion=2022-11-28#get-a-webhook-configuration-for-an-organization
func (s *OrganizationsService) GetHookConfiguration(ctx context.Context, org string, id int64) (*HookConfig, *Response, error) {
	u, err := newURLString("orgs/%v/hooks/%v/config", org, id)
	if err != nil {
		return nil, nil, err
	}
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

// EditHookConfiguration updates the configuration for the specified organization webhook.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/webhooks?apiVersion=2022-11-28#update-a-webhook-configuration-for-an-organization
func (s *OrganizationsService) EditHookConfiguration(ctx context.Context, org string, id int64, config *HookConfig) (*HookConfig, *Response, error) {
	u, err := newURLString("orgs/%v/hooks/%v/config", org, id)
	if err != nil {
		return nil, nil, err
	}
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
