// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type T struct {
	Field string
}

type Client struct{}

func (c *Client) Do(req any, v any) (any, error) {
	return nil, nil
}

type Receiver struct {
	client *Client
}

func (s *Receiver) TestMethod(req any) {
	// Proper usage: var pointer and pass &v
	var v1 *T
	s.client.Do(req, &v1)

	// Literal with fields
	v2 := &T{Field: "something"}
	s.client.Do(req, v2)

	// new(T) but used for something else first
	v3 := new(T)
	v3.Field = "set"
	s.client.Do(req, v3)

	// Anonymous struct
	var v11 *struct {
		F string
	}
	s.client.Do(req, &v11)
}

func (s *Receiver) MethodNameToIgnore(req any) {
	v := new(T)
	s.client.Do(req, v)
}

func (s *Receiver) unexportedMethod(req any) {
	v := new(T)
	s.client.Do(req, v) // Should be ignored because unexported.
}
