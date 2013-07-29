// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"time"

	"fmt"
)

// GitService handles communication with the git data related
// methods of the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/git/
type GitService struct {
	client *Client
}

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

// Tree represents a GitHub tree.
type Tree struct {
	SHA     string      `json:"sha,omitempty"`
	Entries []TreeEntry `json:"tree,omitempty"`
}

// TreeEntry represents the contents of a tree structure.  TreeEntry can
// represent either a blob, a commit (in the case of a submodule), or another
// tree.
type TreeEntry struct {
	SHA  string `json:"sha,omitempty"`
	Path string `json:"path,omitempty"`
	Mode string `json:"mode,omitempty"`
	Type string `json:"type,omitempty"`
	Size int    `json:"size,omitempty"`
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

// GetTree fetches the Tree object for a given sha hash from a repository.
//
// GitHub API docs: http://developer.github.com/v3/git/trees/#get-a-tree
func (s *GitService) GetTree(owner string, repo string, sha string, recursive bool) (*Tree, error) {
	u := fmt.Sprintf("repos/%v/%v/git/trees/%v", owner, repo, sha)
	if recursive {
		u += "?recursive=1"
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	t := new(Tree)
	_, err = s.client.Do(req, t)
	return t, err
}

// createTree represents the body of a CreateTree request.
type createTree struct {
	BaseTree string      `json:base_tree`
	Entries  []TreeEntry `json:tree`
}

// CreateTree creates a new tree in a repository.  If both a tree and a nested
// path modifying that tree are specified, it will overwrite the contents of
// that tree with the new path contents and write a new tree out.
//
// GitHub API docs: http://developer.github.com/v3/git/trees/#create-a-tree
func (s *GitService) CreateTree(owner string, repo string, baseTree string, entries []TreeEntry) (*Tree, error) {
	u := fmt.Sprintf("repos/%v/%v/git/trees", owner, repo)

	body := &createTree{
		BaseTree: baseTree,
		Entries:  entries,
	}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, err
	}

	t := new(Tree)
	_, err = s.client.Do(req, t)
	return t, err
}
