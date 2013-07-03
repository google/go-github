// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "fmt"

// GitService handles communication with the git data related
// methods of the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/git/
type GitService struct {
	client *Client
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

// createTree represents the body of a CreateTree request.
type createTree struct {
	BaseTree string      `json:base_tree`
	Entries  []TreeEntry `json:tree`
}

// GetTree fetches the Tree object for a given sha hash from a users repository.
//
// GitHub API docs: http://developer.github.com/v3/git/trees/#get-a-tree
func (s *GitService) GetTree(user string, repo string, sha string, recursive bool) (*Tree, error) {
	u := fmt.Sprintf("repos/%v/%v/git/trees/%v", user, repo, sha)
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

// CreateTree creates a new tree in a repository.  If both a tree and a nested
// path modifying that tree are specified, it will overwrite the contents of
// that tree with the new path contents and write a new tree out.
//
// GitHub API docs: http://developer.github.com/v3/git/trees/#create-a-tree
func (s *GitService) CreateTree(owner string, repo string, sha string, baseTree string, entries []TreeEntry) (*Tree, error) {
	u := fmt.Sprintf("repos/%v/%v/git/trees/%v", owner, repo, sha)

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
