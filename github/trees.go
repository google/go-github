// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package github

import (
	"fmt"
	"net/url"
	"strconv"
)

// TreesService handles communication with the tree related
// methods of the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/trees/
type TreesService struct {
	client *Client
}

type Tree struct {
	SHA   string `json:"sha,omitempty"`
	URL   string `json:"url,omitempty"`
	Trees []GitTree `json:"tree,omitempty"`
}

// Tree represents a Git tree.
type GitTree struct {
	SHA  string `json:"sha,omitempty"`
	URL  string `json:"url,omitempty"`
	Path string `json:"path,omitempty"`
	Mode string `json:"mode,omitempty"`
	Type string `json:"type,omitempty"`
	Size int    `json:"size,omitempty"`
}

// TreeListOptions specifies the optional parameters to the
// TreesService.List method.
type TreeListOptions struct {
	// Type of repositories to list.  Possible values are: all, owner, public,
	// private, member.  Default is "all".
	Type string

	// How to sort the Tree list.  Possible values are: created, updated,
	// pushed, full_name.  Default is "full_name".
	Sort string

	// Direction in which to sort repositories.  Possible values are: asc, desc.
	// Default is "asc" when sort is "full_name", otherwise default is "desc".
	Direction string

	// For paginated result sets, page of results to retrieve.
	Page int

	// For fetching trees recursively
	//
	// GitHub API docs: http://developer.github.com/v3/git/trees/#get-a-tree-recursively
	Recursive int
}

// Get the Tree object for a given sha hash from a users repository.
//
// GitHub API docs: http://developer.github.com/v3/git/trees/#get-a-tree
func (s *TreesService) List(user string, repo string, sha string, opt *TreeListOptions) (*Tree, error) {
	url_ := fmt.Sprintf("repos/%v/%v/git/trees/%v", user, repo, sha)

	if opt != nil {
		params := url.Values{
			"type":      []string{opt.Type},
			"sort":      []string{opt.Sort},
			"direction": []string{opt.Direction},
			"page":      []string{strconv.Itoa(opt.Page)},
		}
		url_ += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	var response Tree
	_, err = s.client.Do(req, &response)
	return &response, err
}

type CreateTree struct {
	BaseTree string `json:base_tree`
	Tree     []GitTree `json:tree`
}

// Create a new Tree.  If an organization is specified, the new
// Tree will be created under that org.  If the empty string is
// specified, it will be created for the authenticated user.
//
// GitHub API docs: http://developer.github.com/v3/repos/#create
func (s *TreesService) Create(user string, repo string, sha string, create *CreateTree) (*Tree, error) {
	url_ := fmt.Sprintf("repos/%v/%v/git/trees/%v", user, repo, sha)

	req, err := s.client.NewRequest("POST", url_, create)
	if err != nil {
		return nil, err
	}

	r := new(Tree)
	_, err = s.client.Do(req, r)
	return r, err
}
