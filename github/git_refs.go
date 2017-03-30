// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// Reference represents a GitHub reference.
type Reference struct {
	Ref    *string    `json:"ref"`
	URL    *string    `json:"url"`
	Object *GitObject `json:"object"`
}

func (r Reference) String() string {
	return Stringify(r)
}

// GitObject represents a Git object.
type GitObject struct {
	Type *string `json:"type"`
	SHA  *string `json:"sha"`
	URL  *string `json:"url"`
}

func (o GitObject) String() string {
	return Stringify(o)
}

// createRefRequest represents the payload for creating a reference.
type createRefRequest struct {
	Ref *string `json:"ref"`
	SHA *string `json:"sha"`
}

// updateRefRequest represents the payload for updating a reference.
type updateRefRequest struct {
	SHA   *string `json:"sha"`
	Force *bool   `json:"force"`
}

// GetRef fetches a single Reference object for a given Git ref.
// If a Git ref does not match exactly, this method will return a nil result.
//
// GitHub API docs: https://developer.github.com/v3/git/refs/#get-a-reference
func (s *GitService) GetRef(ctx context.Context, owner string, repo string, ref string) (*Reference, *Response, error) {
	ref = strings.TrimPrefix(ref, "refs/")
	u := fmt.Sprintf("repos/%v/%v/git/refs/%v", owner, repo, ref)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	buf := new(bytes.Buffer)
	resp, err := s.client.Do(ctx, req, buf)
	if err != nil {
		return nil, resp, err
	}

	rbuf := bytes.NewBuffer(buf.Bytes())
	r := new(Reference)

	err = json.NewDecoder(rbuf).Decode(r)
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		// Multiple refs, means there wasn't an exact match.
		return nil, resp, nil
	} else if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// GetRefs fetches a slice of Reference objects for a given Git ref.
// If a Git ref does not match exactly, GitHub will return all refs that match it as a prefix.
// Example:
// feature -> [featureA, featureB]
// featureA -> [featureA]
//
// GitHub API docs: https://developer.github.com/v3/git/refs/#get-a-reference
func (s *GitService) GetRefs(ctx context.Context, owner string, repo string, ref string) ([]*Reference, *Response, error) {
	ref = strings.TrimPrefix(ref, "refs/")
	u := fmt.Sprintf("repos/%v/%v/git/refs/%v", owner, repo, ref)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	buf := new(bytes.Buffer)
	resp, err := s.client.Do(ctx, req, buf)
	if err != nil {
		return nil, resp, err
	}

	rbuf := bytes.NewBuffer(buf.Bytes())

	var rs []*Reference
	r := new(Reference)
	// prioritize the most common case: a single returned ref
	err = json.NewDecoder(rbuf).Decode(r)
	if err == nil {
		rs = append(rs, r)
	} else if _, ok := err.(*json.UnmarshalTypeError); ok {
		// check for multiple refs returned
		rbuf = bytes.NewBuffer(buf.Bytes())
		serr := json.NewDecoder(rbuf).Decode(&rs)
		if serr != nil {
			return nil, resp, serr
		}
	} else {
		return nil, resp, err
	}

	return rs, resp, nil
}

// ReferenceListOptions specifies optional parameters to the
// GitService.ListRefs method.
type ReferenceListOptions struct {
	Type string `url:"-"`

	ListOptions
}

// ListRefs lists all refs in a repository.
//
// GitHub API docs: https://developer.github.com/v3/git/refs/#get-all-references
func (s *GitService) ListRefs(ctx context.Context, owner, repo string, opt *ReferenceListOptions) ([]*Reference, *Response, error) {
	var u string
	if opt != nil && opt.Type != "" {
		u = fmt.Sprintf("repos/%v/%v/git/refs/%v", owner, repo, opt.Type)
	} else {
		u = fmt.Sprintf("repos/%v/%v/git/refs", owner, repo)
	}
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var rs []*Reference
	resp, err := s.client.Do(ctx, req, &rs)
	if err != nil {
		return nil, resp, err
	}

	return rs, resp, nil
}

// CreateRef creates a new ref in a repository.
//
// GitHub API docs: https://developer.github.com/v3/git/refs/#create-a-reference
func (s *GitService) CreateRef(ctx context.Context, owner string, repo string, ref *Reference) (*Reference, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/git/refs", owner, repo)
	req, err := s.client.NewRequest("POST", u, &createRefRequest{
		// back-compat with previous behavior that didn't require 'refs/' prefix
		Ref: String("refs/" + strings.TrimPrefix(*ref.Ref, "refs/")),
		SHA: ref.Object.SHA,
	})
	if err != nil {
		return nil, nil, err
	}

	r := new(Reference)
	resp, err := s.client.Do(ctx, req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// UpdateRef updates an existing ref in a repository.
//
// GitHub API docs: https://developer.github.com/v3/git/refs/#update-a-reference
func (s *GitService) UpdateRef(ctx context.Context, owner string, repo string, ref *Reference, force bool) (*Reference, *Response, error) {
	refPath := strings.TrimPrefix(*ref.Ref, "refs/")
	u := fmt.Sprintf("repos/%v/%v/git/refs/%v", owner, repo, refPath)
	req, err := s.client.NewRequest("PATCH", u, &updateRefRequest{
		SHA:   ref.Object.SHA,
		Force: &force,
	})
	if err != nil {
		return nil, nil, err
	}

	r := new(Reference)
	resp, err := s.client.Do(ctx, req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// DeleteRef deletes a ref from a repository.
//
// GitHub API docs: https://developer.github.com/v3/git/refs/#delete-a-reference
func (s *GitService) DeleteRef(ctx context.Context, owner string, repo string, ref string) (*Response, error) {
	ref = strings.TrimPrefix(ref, "refs/")
	u := fmt.Sprintf("repos/%v/%v/git/refs/%v", owner, repo, ref)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
