// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/url"
	"time"
)

// PullRequestsService handles communication with the pull request related
// methods of the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/pulls/
type PullRequestsService struct {
	client *Client
}

// PullRequest represents a GitHub pull request on a repository.
type PullRequest struct {
	Number       int        `json:"number,omitempty"`
	State        string     `json:"state,omitempty"`
	Title        string     `json:"title,omitempty"`
	Body         string     `json:"body,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
	MergedAt     *time.Time `json:"merged_at,omitempty"`
	User         *User      `json:"user,omitempty"`
	Merged       bool       `json:"merged,omitempty"`
	Mergeable    bool       `json:"mergeable,omitempty"`
	MergedBy     *User      `json:"merged_by,omitempty"`
	Comments     int        `json:"comments,omitempty"`
	Commits      int        `json:"commits,omitempty"`
	Additions    int        `json:"additions,omitempty"`
	Deletions    int        `json:"deletions,omitempty"`
	ChangedFiles int        `json:"changed_files,omitempty"`

	// TODO(willnorris): add head and base once we have a Commit struct defined somewhere
}

// PullRequestComment represents a comment left on a pull request.
type PullRequestComment struct {
	ID        int        `json:"id,omitempty"`
	Body      string     `json:"body,omitempty"`
	Path      string     `json:"path,omitempty"`
	Position  int        `json:"position,omitempty"`
	CommitID  string     `json:"commit_id,omitempty"`
	User      *User      `json:"user,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// PullRequestListOptions specifies the optional parameters to the
// PullRequestsService.List method.
type PullRequestListOptions struct {
	// State filters pull requests based on their state.  Possible values are:
	// open, closed.  Default is "open".
	State string

	// Head filters pull requests by head user and branch name in the format of:
	// "user:ref-name".
	Head string

	// Base filters pull requests by base branch name.
	Base string
}

// List the pull requests for the specified repository.
//
// GitHub API docs: http://developer.github.com/v3/pulls/#list-pull-requests
func (s *PullRequestsService) List(owner string, repo string, opt *PullRequestListOptions) ([]PullRequest, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls", owner, repo)
	if opt != nil {
		params := url.Values{
			"state": {opt.State},
			"head":  {opt.Head},
			"base":  {opt.Base},
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	pulls := new([]PullRequest)
	_, err = s.client.Do(req, pulls)
	return *pulls, err
}

// Get a single pull request.
//
// GitHub API docs: https://developer.github.com/v3/pulls/#get-a-single-pull-request
func (s *PullRequestsService) Get(owner string, repo string, number int) (*PullRequest, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/%d", owner, repo, number)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	pull := new(PullRequest)
	_, err = s.client.Do(req, pull)
	return pull, err
}

// Create a new pull request on the specified repository.
//
// GitHub API docs: https://developer.github.com/v3/pulls/#create-a-pull-request
func (s *PullRequestsService) Create(owner string, repo string, pull *PullRequest) (*PullRequest, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls", owner, repo)
	req, err := s.client.NewRequest("POST", u, pull)
	if err != nil {
		return nil, err
	}
	p := new(PullRequest)
	_, err = s.client.Do(req, p)
	return p, err
}

// Edit a pull request.
//
// GitHub API docs: https://developer.github.com/v3/pulls/#update-a-pull-request
func (s *PullRequestsService) Edit(owner string, repo string, number int, pull *PullRequest) (*PullRequest, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/%d", owner, repo, number)
	req, err := s.client.NewRequest("PATCH", u, pull)
	if err != nil {
		return nil, err
	}
	p := new(PullRequest)
	_, err = s.client.Do(req, p)
	return p, err
}

// PullRequestListCommentsOptions specifies the optional parameters to the
// PullRequestsService.ListComments method.
type PullRequestListCommentsOptions struct {
	// Sort specifies how to sort comments.  Possible values are: created, updated.
	Sort string

	// Direction in which to sort comments.  Possible values are: asc, desc.
	Direction string

	// Since filters comments by time.
	Since time.Time
}

// ListComments lists all comments on the specified pull request.  Specifying a
// pull request number of 0 will return all comments on all pull requests for
// the repository.
//
// GitHub API docs: https://developer.github.com/v3/pulls/comments/#list-comments-on-a-pull-request
func (s *PullRequestsService) ListComments(owner string, repo string, number int, opt *PullRequestListCommentsOptions) ([]PullRequestComment, error) {
	var u string
	if number == 0 {
		u = fmt.Sprintf("repos/%v/%v/pulls/comments", owner, repo)
	} else {
		u = fmt.Sprintf("repos/%v/%v/pulls/%d/comments", owner, repo, number)
	}

	if opt != nil {
		params := url.Values{
			"sort":      {opt.Sort},
			"direction": {opt.Direction},
		}
		if !opt.Since.IsZero() {
			params.Add("since", opt.Since.Format(time.RFC3339))
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	comments := new([]PullRequestComment)
	_, err = s.client.Do(req, comments)
	return *comments, err
}

// GetComment fetches the specified pull request comment.
//
// GitHub API docs: https://developer.github.com/v3/pulls/comments/#get-a-single-comment
func (s *PullRequestsService) GetComment(owner string, repo string, number int) (*PullRequestComment, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/comments/%d", owner, repo, number)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	comment := new(PullRequestComment)
	_, err = s.client.Do(req, comment)
	return comment, err
}

// CreateComment creates a new comment on the specified pull request.
//
// GitHub API docs: https://developer.github.com/v3/pulls/comments/#get-a-single-comment
func (s *PullRequestsService) CreateComment(owner string, repo string, number int, comment *PullRequestComment) (*PullRequestComment, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/%d/comments", owner, repo, number)
	req, err := s.client.NewRequest("POST", u, comment)
	if err != nil {
		return nil, err
	}
	c := new(PullRequestComment)
	_, err = s.client.Do(req, c)
	return c, err
}

// EditComment updates a pull request comment.
//
// GitHub API docs: https://developer.github.com/v3/pulls/comments/#edit-a-comment
func (s *PullRequestsService) EditComment(owner string, repo string, number int, comment *PullRequestComment) (*PullRequestComment, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/comments/%d", owner, repo, number)
	req, err := s.client.NewRequest("PATCH", u, comment)
	if err != nil {
		return nil, err
	}
	c := new(PullRequestComment)
	_, err = s.client.Do(req, c)
	return c, err
}

// DeleteComment deletes a pull request comment.
//
// GitHub API docs: https://developer.github.com/v3/pulls/comments/#delete-a-comment
func (s *PullRequestsService) DeleteComment(owner string, repo string, number int) error {
	u := fmt.Sprintf("repos/%v/%v/pulls/comments/%d", owner, repo, number)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
