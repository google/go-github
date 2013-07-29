// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "time"

// WebHookPayload represents the data that is received from GitHub when a push
// event hook is triggered.  The format of these payloads pre-date most of the
// GitHub v3 API, so there are lots of minor incompatibilities with the types
// defined in the rest of the API.  Therefore, several types are duplicated
// here to account for these differences.
//
// GitHub API docs: https://help.github.com/articles/post-receive-hooks
type WebHookPayload struct {
	After      string          `json:"after,omitempty"`
	Before     string          `json:"before,omitempty"`
	Commits    []WebHookCommit `json:"commits,omitempty"`
	Compare    string          `json:"compare,omitempty"`
	Created    bool            `json:"created,omitempty"`
	Deleted    bool            `json:"deleted,omitempty"`
	Forced     bool            `json:"forced,omitempty"`
	HeadCommit *WebHookCommit  `json:"head_commit,omitempty"`
	Pusher     *User           `json:"pusher,omitempty"`
	Ref        string          `json:"ref,omitempty"`
	Repo       *Repository     `json:"repository,omitempty"`
}

// WebHookCommit represents the commit variant we receive from GitHub in a
// WebHookPayload.
type WebHookCommit struct {
	Added     []string      `json:"added,omitempty"`
	Author    WebHookAuthor `json:"author,omitempty"`
	Committer WebHookAuthor `json:"committer,omitempty"`
	Distinct  bool          `json:"distinct,omitempty"`
	ID        string        `json:"id,omitempty"`
	Message   string        `json:"message,omitempty"`
	Modified  []string      `json:"modified,omitempty"`
	Removed   []string      `json:"removed,omitempty"`
	Timestamp *time.Time    `json:"timestamp,omitempty"`
}

// WebHookAuthor represents the author or committer of a commit, as specified
// in a WebHookCommit.  The commit author may not correspond to a GitHub User.
type WebHookAuthor struct {
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
}
