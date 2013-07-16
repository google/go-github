// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"time"
)

// Commit represents a GitHub commit
type Commit struct {
	Author    Author     `json:"author,omitempty"`
	Committer Author     `json:"commiter,omitempty"`
	Message   string     `json:"message,omitempty"`
	SHA       string     `json:"sha,omitempty"`
	URL       string     `json:"url,omitempty"`
}

// Author represents someone who made a commit. This does not necessarily
// represent a GitHub User.
type Author struct {
	Date  *time.Time `json:"date,omitempty"`
	Name  string     `json:"name,omitempty"`
	email string     `json:"email,omitempty"`
}
