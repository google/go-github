// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "context"

type Client struct{}

func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body any) error {
	return nil
}

type service struct {
	client *Client
}

type CreateRequest struct {
	Name string
}

type UpdateRequest struct {
	Name string
}

type SuspendOptions struct {
	Reason string
}

type ServerOptions struct {
	Host string
}

// opts named, value, good type: rename only.
func (s *service) Create(ctx context.Context, opts CreateRequest) error { // want `rename request body parameter "opts" to "body"`
	return s.client.NewRequest(ctx, "POST", "u", opts)
}

// request named, pointer, good type: rename and by-value.
func (s *service) Update(ctx context.Context, request *UpdateRequest) error { // want `rename request body parameter "request" to "body"` `pass request body "request" by value, not by pointer`
	return s.client.NewRequest(ctx, "PATCH", "u", request)
}

// opts named, pointer, Options-suffixed type: rename, by-value, and type suffix.
func (s *service) Suspend(ctx context.Context, opts *SuspendOptions) error { // want `rename request body parameter "opts" to "body"` `pass request body "opts" by value, not by pointer` `request body type "SuspendOptions" should use a "Request" suffix, not "Options"`
	return s.client.NewRequest(ctx, "PUT", "u", opts)
}

// Domain-specific name, value, Options-suffixed type: rename and type suffix.
func (s *service) Save(ctx context.Context, settings ServerOptions) error { // want `rename request body parameter "settings" to "body"` `request body type "ServerOptions" should use a "Request" suffix, not "Options"`
	return s.client.NewRequest(ctx, "POST", "u", settings)
}
