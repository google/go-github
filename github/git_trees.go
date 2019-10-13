// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// Tree represents a GitHub tree.
type Tree struct {
	SHA     *string     `json:"sha,omitempty"`
	Entries []TreeEntry `json:"tree,omitempty"`

	// Truncated is true if the number of items in the tree
	// exceeded GitHub's maximum limit and the Entries were truncated
	// in the response. Only populated for requests that fetch
	// trees like Git.GetTree.
	Truncated *bool `json:"truncated,omitempty"`
}

func (t Tree) String() string {
	return Stringify(t)
}

type ITreeEntry interface {
	Tree() *TreeEntryBase
}

type TreeEntryBase struct {
	Path *string `json:"path,omitempty"`
	Mode *string `json:"mode,omitempty"`
	Type *string `json:"type,omitempty"`
}

// TreeEntry represents the contents of a tree structure. TreeEntry can
// represent either a blob, a commit (in the case of a submodule), or another
// tree.
type TreeEntry struct {
	SHA     *string `json:"sha,omitempty"`
	Path    *string `json:"path,omitempty"`
	Mode    *string `json:"mode,omitempty"`
	Type    *string `json:"type,omitempty"`
	Size    *int    `json:"size,omitempty"`
	Content *string `json:"content,omitempty"`
	URL     *string `json:"url,omitempty"`
}

func (t TreeEntry) String() string {
	return Stringify(t)
}

func (t TreeEntry) Tree() *TreeEntryBase {
	return &TreeEntryBase{
		t.Path,
		t.Mode,
		t.Type,
	}
}

// TreeDeleteEntry represents a file deletion operation. Leave SHA empty
type TreeDeleteEntry struct {
	SHA  *string `json:"sha"`
	Path *string `json:"path,omitempty"`
	Mode *string `json:"mode,omitempty"`
	Type *string `json:"type,omitempty"`
}

func (t TreeDeleteEntry) String() string {
	return Stringify(t)
}

func (t TreeDeleteEntry) Tree() *TreeEntryBase {
	return &TreeEntryBase{
		t.Path,
		t.Mode,
		t.Type,
	}
}

// GetTree fetches the Tree object for a given sha hash from a repository.
//
// GitHub API docs: https://developer.github.com/v3/git/trees/#get-a-tree
func (s *GitService) GetTree(ctx context.Context, owner string, repo string, sha string, recursive bool) (*Tree, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/git/trees/%v", owner, repo, sha)
	if recursive {
		u += "?recursive=1"
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	t := new(Tree)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// createTree represents the body of a CreateTree request.
type createTree struct {
	BaseTree string       `json:"base_tree,omitempty"`
	Entries  []ITreeEntry `json:"tree"`
}

type jsonObject map[string]interface{}

func (h jsonObject) GetString(path string) *string {
	value, hasPath := h[path]
	if !hasPath {
		return nil
	}
	valueStr, isStr := value.(string)
	if !isStr {
		return nil
	}
	return &valueStr
}

func (h jsonObject) GetInt(path string) *int {
	value, hasPath := h[path]
	if !hasPath {
		return nil
	}
	valueInt, isInt := value.(int)
	if !isInt {
		return nil
	}
	return &valueInt
}

func (s *createTree) UnmarshalJSON(data []byte) error {
	decoded := map[string]interface{}{}
	decodeErr := json.Unmarshal(data, &decoded)
	if decodeErr != nil {
		return fmt.Errorf("couldn't decode to map %v", decodeErr)
	}
	var typeOk bool
	if s.BaseTree, typeOk = decoded["base_tree"].(string); !typeOk {
		return errors.New("base_tree not string")
	}
	tree, isArray := decoded["tree"].([]interface{})
	if !isArray {
		return errors.New("tree is not array")
	}

	for _, item := range tree {
		itemAsMap, isMap := item.(map[string]interface{})
		if !isMap {
			return errors.New("tree needs to be array of map")
		}
		itemAsJson := jsonObject(itemAsMap)
		sha := itemAsJson.GetString("sha")
		content := itemAsJson.GetString("content")

		if sha != nil || content != nil {
			s.Entries = append(s.Entries, TreeEntry{
				SHA:     sha,
				Path:    itemAsJson.GetString("path"),
				Mode:    itemAsJson.GetString("mode"),
				Type:    itemAsJson.GetString("type"),
				Size:    itemAsJson.GetInt("size"),
				Content: content,
				URL:     itemAsJson.GetString("url"),
			})
		} else {
			s.Entries = append(s.Entries, TreeDeleteEntry{
				SHA:  nil,
				Path: itemAsJson.GetString("path"),
				Mode: itemAsJson.GetString("mode"),
				Type: itemAsJson.GetString("type"),
			})
		}
	}
	return nil
}

// CreateTree creates a new tree in a repository. If both a tree and a nested
// path modifying that tree are specified, it will overwrite the contents of
// that tree with the new path contents and write a new tree out.
//
// GitHub API docs: https://developer.github.com/v3/git/trees/#create-a-tree
func (s *GitService) CreateTree(ctx context.Context, owner string, repo string, baseTree string, entries []ITreeEntry) (*Tree, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/git/trees", owner, repo)

	body := &createTree{
		BaseTree: baseTree,
		Entries:  entries,
	}
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	t := new(Tree)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}
