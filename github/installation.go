// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// Installation represents a GitHub integration installation.
type Installation struct {
	ID              *int    `json:"id,omitempty"`
	Account         *User   `json:"account,omitempty"`
	AccessTokensURL *string `json:"access_tokens_url,omitempty"`
	RepositoriesURL *string `json:"repositories_url,omitempty"`
}
