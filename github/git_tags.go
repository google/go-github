// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
)

// Object represents a git object.
type Object struct {
	Type string `json:"type,omitempty"`
	SHA  string `json:"sha,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Tag represents a tag object.
type Tag struct {
	Tag     string        `json:"tag,omitempty"`
	SHA     string        `json:"sha,omitempty"`
	URL     string        `json:"url,omitempty"`
	Message string        `json:"message,omitempty"`
	Tagger  *CommitAuthor `json:"tagger,omitempty"`
	Object  *Object       `json:"object,omitempty"`
}

// GetTag fetchs a tag from a repo given a SHA.
//
// GitHub API docs: http://developer.github.com/v3/git/tags/#get-a-tag
func (s *GitService) GetTag(owner string, repo string, sha string) (*Tag, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/git/tags/%v", owner, repo, sha)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	tag := new(Tag)
	resp, err := s.client.Do(req, tag)
	return tag, resp, err
}

// CreateTag creates a tag object.
//
// GitHub API docs: http://developer.github.com/v3/git/tags/#create-a-tag-object
func (s *GitService) CreateTag(owner string, repo string, tag *Tag) (*Tag, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/git/tags", owner, repo)
	req, err := s.client.NewRequest("POST", u, tag)
	if err != nil {
		return nil, nil, err
	}

	t := new(Tag)
	resp, err := s.client.Do(req, t)
	return t, resp, err
}
