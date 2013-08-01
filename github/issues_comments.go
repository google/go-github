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

// IssueComment represents a comment left on an issue.
type IssueComment struct {
	ID        int        `json:"id,omitempty"`
	Body      string     `json:"body,omitempty"`
	User      *User      `json:"user,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// IssueListCommentsOptions specifies the optional parameters to the
// IssuesService.ListComments method.
type IssueListCommentsOptions struct {
	// Sort specifies how to sort comments.  Possible values are: created, updated.
	Sort string

	// Direction in which to sort comments.  Possible values are: asc, desc.
	Direction string

	// Since filters comments by time.
	Since time.Time
}

// ListComments lists all comments on the specified issue.  Specifying an issue
// number of 0 will return all comments on all issues for the repository.
//
// GitHub API docs: http://developer.github.com/v3/issues/comments/#list-comments-on-an-issue
func (s *IssuesService) ListComments(owner string, repo string, number int, opt *IssueListCommentsOptions) ([]IssueComment, error) {
	var u string
	if number == 0 {
		u = fmt.Sprintf("repos/%v/%v/issues/comments", owner, repo)
	} else {
		u = fmt.Sprintf("repos/%v/%v/issues/%d/comments", owner, repo, number)
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
	comments := new([]IssueComment)
	_, err = s.client.Do(req, comments)
	return *comments, err
}

// GetComment fetches the specified issue comment.
//
// GitHub API docs: http://developer.github.com/v3/issues/comments/#get-a-single-comment
func (s *IssuesService) GetComment(owner string, repo string, id int) (*IssueComment, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/comments/%d", owner, repo, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	comment := new(IssueComment)
	_, err = s.client.Do(req, comment)
	return comment, err
}

// CreateComment creates a new comment on the specified issue.
//
// GitHub API docs: http://developer.github.com/v3/issues/comments/#create-a-comment
func (s *IssuesService) CreateComment(owner string, repo string, number int, comment *IssueComment) (*IssueComment, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%d/comments", owner, repo, number)
	req, err := s.client.NewRequest("POST", u, comment)
	if err != nil {
		return nil, err
	}
	c := new(IssueComment)
	_, err = s.client.Do(req, c)
	return c, err
}

// EditComment updates an issue comment.
//
// GitHub API docs: http://developer.github.com/v3/issues/comments/#edit-a-comment
func (s *IssuesService) EditComment(owner string, repo string, id int, comment *IssueComment) (*IssueComment, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/comments/%d", owner, repo, id)
	req, err := s.client.NewRequest("PATCH", u, comment)
	if err != nil {
		return nil, err
	}
	c := new(IssueComment)
	_, err = s.client.Do(req, c)
	return c, err
}

// DeleteComment deletes an issue comment.
//
// GitHub API docs: http://developer.github.com/v3/issues/comments/#delete-a-comment
func (s *IssuesService) DeleteComment(owner string, repo string, id int) error {
	u := fmt.Sprintf("repos/%v/%v/issues/comments/%d", owner, repo, id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
