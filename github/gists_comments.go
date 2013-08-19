// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"time"
)

type GistComment struct {
	ID        string     `json:"id,omitempty"`
	URL       string     `json:"url,omitempty"`
	Body      string     `json:"body,omitempty"`
	User      *User      `json:"user,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// List all comments for a gist.
//
// GitHub API docs: http://developer.github.com/v3/gists/comments/#list-comments-on-a-gist
func (s *GistsService) ListComments(gistId string) ([]GistComment, *Response, error) {
	u := fmt.Sprintf("gists/%v/comments", gistId)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	comments := new([]GistComment)
	resp, err := s.client.Do(req, comments)
	return *comments, resp, err
}

// Get a single comment from a gist.
//
// GitHub API docs: http://developer.github.com/v3/gists/comments/#get-a-single-comment
func (s *GistsService) GetComment(gistId string, commentId string) (*GistComment, *Response, error) {
	u := fmt.Sprintf("gists/%v/comments/%v", gistId, commentId)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	c := new(GistComment)
	resp, err := s.client.Do(req, c)
	return c, resp, err
}

// Create a comment for a gist.
//
// GitHub API docs: http://developer.github.com/v3/gists/comments/#create-a-comment
func (s *GistsService) CreateComment(gistId string, comment *GistComment) (*GistComment, *Response, error) {
	u := fmt.Sprintf("gists/%v/comments", gistId)
	req, err := s.client.NewRequest("POST", u, comment)
	if err != nil {
		return nil, nil, err
	}

	c := new(GistComment)
	resp, err := s.client.Do(req, c)
	return c, resp, err
}

// Edit a gist comment.
//
// GitHub API docs: http://developer.github.com/v3/gists/comments/#edit-a-comment
func (s *GistsService) EditComment(gistId string, commentId string, comment *GistComment) (*GistComment, *Response, error) {
	u := fmt.Sprintf("gists/%v/comments/%v", gistId, commentId)
	req, err := s.client.NewRequest("PATCH", u, comment)
	if err != nil {
		return nil, nil, err
	}

	c := new(GistComment)
	resp, err := s.client.Do(req, c)
	return c, resp, err
}

// Delete a gist comment.
//
// GitHub API docs: http://developer.github.com/v3/gists/comments/#delete-a-comment
func (s *GistsService) DeleteComment(gistId string, commentId string) (*Response, error) {
	u := fmt.Sprintf("gists/%v/comments/%v", gistId, commentId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
