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
