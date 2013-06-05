// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package github

import "fmt"

// TreesService handles communication with the tree related
// methods of the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/trees/
type TreesService struct {
	client *Client
}

type Tree struct {
	SHA   string    `json:"sha,omitempty"`
	Trees []GitTree `json:"tree,omitempty"`
}

// Tree represents a Git tree.
type GitTree struct {
	SHA  string `json:"sha,omitempty"`
	Path string `json:"path,omitempty"`
	Mode string `json:"mode,omitempty"`
	Type string `json:"type,omitempty"`
	Size int    `json:"size,omitempty"`
}

type createTree struct {
	baseTree string    `json:base_tree`
	trees    []GitTree `json:tree`
}

// Get the Tree object for a given sha hash from a users repository.
//
// GitHub API docs: http://developer.github.com/v3/git/trees/#get-a-tree
func (s *TreesService) Get(user string, repo string, sha string, recursive bool) (*Tree, error) {
	url_ := fmt.Sprintf("repos/%v/%v/git/trees/%v", user, repo, sha)

	if recursive {
		url_ += "?recursive=1"
	}

	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	var response Tree
	_, err = s.client.Do(req, &response)
	return &response, err
}

// The tree creation API will take nested entries as well.
// If both a tree and a nested path modifying that tree are specified,
// it will overwrite the contents of that tree with the new path contents and write a new tree out.
//
// GitHub API docs: http://developer.github.com/v3/git/trees/#create-a-tree
func (s *TreesService) Create(owner string, repo string, sha string, baseTreeSha string, trees []GitTree) (*Tree, error) {
	url_ := fmt.Sprintf("repos/%v/%v/git/trees/%v", owner, repo, sha)

	req, err := s.client.NewRequest("POST", url_, createTree{
		baseTree: baseTreeSha,
		trees:    trees,
	})
	if err != nil {
		return nil, err
	}

	r := new(Tree)
	_, err = s.client.Do(req, r)
	return r, err
}
