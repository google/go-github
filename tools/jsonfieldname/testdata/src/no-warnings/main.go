// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type Example struct {
	GithubThing string `json:"github_thing"`      // Should not be flagged
	ID          string `json:"id,omitempty"`      // Should not be flagged
	Strings     string `json:"strings,omitempty"` // Should not be flagged
	Ref         string `json:"$ref,omitempty"`    // Should not be flagged
}
