// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// These event types are shared between the Events API and used as Webhook payloads.

package github

// PushEvent represents a git push to a GitHub repository.
//
// GitHub API docs: http://developer.github.com/v3/activity/events/types/#pushevent
type PushEvent struct {
	PushID  *int                 `json:"push_id,omitempty"`
	Head    *string              `json:"head,omitempty"`
	Ref     *string              `json:"ref,omitempty"`
	Size    *int                 `json:"size,omitempty"`
	Commits []PushEventCommit    `json:"commits,omitempty"`
	Repo    *PushEventRepository `json:"repository,omitempty"`
}

func (p PushEvent) String() string {
	return Stringify(p)
}

// PushEventCommit represents a git commit in a GitHub PushEvent.
type PushEventCommit struct {
	SHA      *string       `json:"sha,omitempty"`
	Message  *string       `json:"message,omitempty"`
	Author   *CommitAuthor `json:"author,omitempty"`
	URL      *string       `json:"url,omitempty"`
	Distinct *bool         `json:"distinct,omitempty"`
	Added    []string      `json:"added,omitempty"`
	Removed  []string      `json:"removed,omitempty"`
	Modified []string      `json:"modified,omitempty"`
}

func (p PushEventCommit) String() string {
	return Stringify(p)
}

// PushEventRepository represents the repo object in a PushEvent payload
type PushEventRepository struct {
	ID              *int                `json:"id,omitempty"`
	Name            *string             `json:"name,omitempty"`
	FullName        *string             `json:"full_name,omitempty"`
	Owner           *PushEventRepoOwner `json:"owner,omitempty"`
	Private         *bool               `json:"private,omitempty"`
	Description     *string             `json:"description,omitempty"`
	Fork            *bool               `json:"fork,omitempty"`
	CreatedAt       *Timestamp          `json:"created_at,omitempty"`
	PushedAt        *Timestamp          `json:"pushed_at,omitempty"`
	UpdatedAt       *Timestamp          `json:"updated_at,omitempty"`
	Homepage        *string             `json:"homepage,omitempty"`
	Size            *int                `json:"size,omitempty"`
	StargazersCount *int                `json:"stargazers_count,omitempty"`
	WatchersCount   *int                `json:"watchers_count,omitempty"`
	Language        *string             `json:"language,omitempty"`
	HasIssues       *bool               `json:"has_issues,omitempty"`
	HasDownloads    *bool               `json:"has_downloads,omitempty"`
	HasWiki         *bool               `json:"has_wiki,omitempty"`
	HasPages        *bool               `json:"has_pages,omitempty"`
	ForksCount      *int                `json:"forks_count,omitempty"`
	OpenIssuesCount *int                `json:"open_issues_count,omitempty"`
	DefaultBranch   *string             `json:"default_branch,omitempty"`
	MasterBranch    *string             `json:"master_branch,omitempty"`
	Organization    *string             `json:"organization,omitempty"`
}

// PushEventRepoOwner is a basic reporesntation of user/org in a PushEvent payload
type PushEventRepoOwner struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

//PullRequestEvent represents the payload delivered by PullRequestEvent webhook
type PullRequestEvent struct {
	Action      *string      `json:"action,omitempty"`
	Number      *int         `json:"number,omitempty"`
	PullRequest *PullRequest `json:"pull_request,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo   *Repository `json:"repository,omitempty"`
	Sender *User       `json:"sender,omitempty"`
}

// IssueActivityEvent represents the payload delivered by Issue webhook
type IssueActivityEvent struct {
	Action *string `json:"action,omitempty"`
	Issue  *Issue  `json:"issue,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo   *Repository `json:"repository,omitempty"`
	Sender *User       `json:"sender,omitempty"`
}

// IssueCommentEvent represents the payload delivered by IssueComment webhook
//
// This webhook also gets fired for comments on pull requests
type IssueCommentEvent struct {
	Action  *string       `json:"action,omitempty"`
	Issue   *Issue        `json:"issue,omitempty"`
	Comment *IssueComment `json:"comment,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo   *Repository `json:"repository,omitempty"`
	Sender *User       `json:"sender,omitempty"`
}
