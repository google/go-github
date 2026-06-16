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

type ListOptions struct {
	Page int
}

type SearchOptions struct {
	Query string
}

type AllowedPtr struct {
	Name string
}

type AllowedOptions struct {
	Name string
}

// Already named body, value, good type: no warning.
func (s *service) Create(ctx context.Context, body CreateRequest) error {
	return s.client.NewRequest(ctx, "POST", "u", body)
}

// No body: no warning.
func (s *service) Delete(ctx context.Context, id string) error {
	return s.client.NewRequest(ctx, "DELETE", "u", nil)
}

// GET with an Options-suffixed pointer must not fire any rule.
func (s *service) List(ctx context.Context, opts *ListOptions) error {
	return s.client.NewRequest(ctx, "GET", "u", opts)
}

// Named body, value, good type: no warning.
func (s *service) Save(ctx context.Context, body CreateRequest) error {
	return s.client.NewRequest(ctx, "PUT", "u", body)
}

// Not a client.NewRequest receiver: no warning.
type other struct{}

func (o *other) NewRequest(ctx context.Context, method, urlStr string, body any) error {
	return nil
}

func (s *service) Other(ctx context.Context, opts *ListOptions) error {
	var o other
	return o.NewRequest(ctx, "POST", "u", opts)
}

func addOptions(s string, opt any) (string, error) { return s, nil }

// opts named, pointer, Options-suffixed type: no warning.
func (s *service) Search(ctx context.Context, opts *SearchOptions) ([]string, error) {
	u, err := addOptions("search", opts)
	if err != nil {
		return nil, err
	}
	_ = u
	return nil, nil
}

// Pointer body whose type is in allowed-pointer-types: by-value rule suppressed.
func (s *service) IgnoredPointer(ctx context.Context, body *AllowedPtr) error {
	return s.client.NewRequest(ctx, "PUT", "u", body)
}

// Options-suffixed body whose type is in allowed-wrong-names: suffix rule suppressed.
func (s *service) IgnoredOptions(ctx context.Context, body AllowedOptions) error {
	return s.client.NewRequest(ctx, "POST", "u", body)
}
