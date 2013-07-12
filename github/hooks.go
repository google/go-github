// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// PostReceiveHook represents the data that is received from github
// when a hook is triggered.
type PostReceiveHook struct {
	After      string     `json:"after,omitempty"`
	Before     string     `json:"before,omitempty"`
	Commits    []Commit   `json:"commits,omitempty"`
	Compare    string     `json:"compare,omitempty"`
	Created    bool       `json:"created,omitempty"`
	Deleted    bool       `json:"deleted,omitempty"`
	Forced     bool       `json:"forced,omitempty"`
	HeadCommit Commit     `json:"head_commit,omitempty"`
	Pusher     User       `json:"pusher,omitempty"`
	Ref        string     `json:"ref,omitempty"`
	Repo       Repository `json:"repository,omitempty"`
}
