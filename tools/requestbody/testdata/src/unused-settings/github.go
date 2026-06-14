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

// Still passed by pointer: allowed-pointer-types exception is still needed.
type ActivePtr struct {
	Name string
}

// No longer passed by pointer: allowed-pointer-types exception is obsolete.
type ObsoletePtr struct { // want `unused requestbody exception: type "ObsoletePtr" in allowed-pointer-types is never passed by pointer to client.NewRequest`
	Name string
}

// Still passed as a request body: allowed-wrong-names exception is still needed.
type ActiveOptions struct {
	Name string
}

// No longer passed as a request body: allowed-wrong-names exception is obsolete.
type ObsoleteOptions struct { // want `unused requestbody exception: type "ObsoleteOptions" in allowed-wrong-names is never passed as a request body to client.NewRequest`
	Name string
}

func (s *service) UseActivePtr(ctx context.Context, body *ActivePtr) error {
	return s.client.NewRequest(ctx, "POST", "u", body)
}

func (s *service) UseActiveOptions(ctx context.Context, body ActiveOptions) error {
	return s.client.NewRequest(ctx, "POST", "u", body)
}
