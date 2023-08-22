// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
)

// ListHookDeliveries lists webhook deliveries for a webhook configured in an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/webhooks#list-deliveries-for-an-organization-webhook
func (s *OrganizationsService) ListHookDeliveries(ctx context.Context, org string, id int64, opts *ListCursorOptions) ([]*HookDelivery, *Response, error) {
	u, err := newURLString("orgs/%v/hooks/%v/deliveries", org, id)
	if err != nil {
		return nil, nil, err
	}
	u, err = addOptions(u, opts)
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

// GetHookDelivery returns a delivery for a webhook configured in an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/webhooks#get-a-webhook-delivery-for-an-organization-webhook
func (s *OrganizationsService) GetHookDelivery(ctx context.Context, owner string, hookID, deliveryID int64) (*HookDelivery, *Response, error) {
	u, err := newURLString("orgs/%v/hooks/%v/deliveries/%v", owner, hookID, deliveryID)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	h := new(HookDelivery)
	resp, err := s.client.Do(ctx, req, h)
	if err != nil {
		return nil, resp, err
	}

	return h, resp, nil
}

// RedeliverHookDelivery redelivers a delivery for a webhook configured in an organization.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/webhooks#redeliver-a-delivery-for-an-organization-webhook
func (s *OrganizationsService) RedeliverHookDelivery(ctx context.Context, owner string, hookID, deliveryID int64) (*HookDelivery, *Response, error) {
	u, err := newURLString("orgs/%v/hooks/%v/deliveries/%v/attempts", owner, hookID, deliveryID)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	h := new(HookDelivery)
	resp, err := s.client.Do(ctx, req, h)
	if err != nil {
		return nil, resp, err
	}

	return h, resp, nil
}
