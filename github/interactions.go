// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// InteractionsService handles communication with the repo and org related
// methods of the GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/interactions/
type InteractionsService service

// Interaction represents the interaction restrictions
// for repository and organisation
type Interaction struct {
	Limit     *string    `json:"limit,omitempty"`
	Origin    *string    `json:"origin,omitempty"`
	ExpiresAt *Timestamp `json:"expires_at,omitempty"`
}
