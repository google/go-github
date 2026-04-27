// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

type T struct {
	Field string
}

type Client struct{}

func (c *Client) Do(req any, v any) (any, error) {
	return nil, nil
}

type Service struct {
	client *Client
}

func assertNilError(t *testing.T, err error) {}

func (s *Service) TestMethod(req any, r *http.Request, t *testing.T) {
	v1 := new(T)
	s.client.Do(req, v1) // want "use 'var v1 [*]T' and pass '&v1' instead"

	v2 := &T{}
	s.client.Do(req, v2) // want "use 'var v2 [*]T' and pass '&v2' instead"

	v3 := new(T)
	json.NewDecoder(r.Body).Decode(v3) // want "use 'var v3 [*]T' and pass '&v3' instead"

	v4 := &T{}
	json.NewDecoder(r.Body).Decode(v4) // want "use 'var v4 [*]T' and pass '&v4' instead"

	v5 := &T{}
	s.client.Do(req, &v5) // want "use 'var v5 [*]T' and pass '&v5' instead"

	v6 := new(T)
	assertNilError(t, json.NewDecoder(r.Body).Decode(v6)) // want "use 'var v6 [*]T' and pass '&v6' instead"

	v7 := &T{Field: "something"}
	s.client.Do(req, v7) // No warning

	var v8 *T
	v8 = new(T)
	s.client.Do(req, v8) // want "use 'var v8 [*]T' and pass '&v8' instead"

	// Multiple assignments in same block
	v9 := new(T)
	v10 := new(T)
	s.client.Do(req, v9)  // want "use 'var v9 [*]T' and pass '&v9' instead"
	s.client.Do(req, v10) // want "use 'var v10 [*]T' and pass '&v10' instead"

	// Anonymous struct
	v11 := new(struct {
		F string
	})
	s.client.Do(req, v11) // want "use 'var v11 [*]T' and pass '&v11' instead"

	// Anonymous struct
	var v12 *struct {
		F string
	}
	s.client.Do(req, v12) // want "pass '&v12' instead"

	var v13 *T
	s.client.Do(req, v13) // want "pass '&v13' instead"
}
