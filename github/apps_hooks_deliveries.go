// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// ListHookDeliveries lists deliveries of an App webhook
//
// GitHub API docs: https://docs.github.com/en/rest/reference/apps#list-deliveries-for-an-app-webhook
func (s *AppsService) ListHookDeliveries(ctx context.Context, opts *ListCursorOptions) ([]*HookDelivery, *Response, error) {
	u, err := addOptions("app/hook/deliveries", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	deliveries := []*HookDelivery{}
	resp, err := s.client.Do(ctx, req, &deliveries)
	if err != nil {
		return nil, resp, err
	}

	return deliveries, resp, nil
}
