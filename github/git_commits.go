// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"time"
	"fmt"
)

// Commit represents a GitHub commit.
type Commit struct {
	SHA       string        `json:"sha,omitempty"`
	Author    *CommitAuthor `json:"author,omitempty"`
	Committer *CommitAuthor `json:"committer,omitempty"`
	Message   string        `json:"message,omitempty"`
	Tree      *Tree         `json:"tree,omitempty"`
	Parents   []Commit      `json:"parents,omitempty"`
}

// CommitAuthor represents the author or committer of a commit.  The commit
// author may not correspond to a GitHub User.
type CommitAuthor struct {
	Date  *time.Time `json:"date,omitempty"`
	Name  string     `json:"name,omitempty"`
	Email string     `json:"email,omitempty"`
}

// GetCommit fetchs the Commit object for a given SHA.
//
// GitHub API docs: http://developer.github.com/v3/git/commits/#get-a-commit
func (s *GitService) GetCommit(owner string, repo string, sha string) (*Commit, error) {
	u := fmt.Sprintf("repos/%v/%v/git/commits/%v", owner, repo, sha)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	c := new(Commit)
	_, err = s.client.Do(req, c)
	return c, err
}

// CreateCommit creates a new commit in a repository.
//
// The commit.Committer is optional and will be filled with the commit.Author
// data if omitted. If the commit.Author is omitted, it will be filled in with
// the authenticated user’s information and the current date.
//
// GitHub API docs: http://developer.github.com/v3/git/commits/#create-a-commit
func (s *GitService) CreateCommit(owner string, repo string, commit *Commit) (*Commit, error) {
	u := fmt.Sprintf("repos/%v/%v/git/commits", owner, repo)
	req, err := s.client.NewRequest("POST", u, commit)
	if err != nil {
		return nil, err
	}

	c := new(Commit)
	_, err = s.client.Do(req, c)
	return c, err
}
