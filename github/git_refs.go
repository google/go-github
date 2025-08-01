// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Reference represents a GitHub reference.
type Reference struct {
	// The name of the fully qualified reference, i.e.: `refs/heads/master`.
	Ref    *string    `json:"ref"`
	URL    *string    `json:"url"`
	Object *GitObject `json:"object"`
	NodeID *string    `json:"node_id,omitempty"`
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

// GetRef fetches a single reference in a repository.
// The ref must be formatted as `heads/<branch name>` for branches and `tags/<tag name>` for tags.
//
// GitHub API docs: https://docs.github.com/rest/git/refs#get-a-reference
//
//meta:operation GET /repos/{owner}/{repo}/git/ref/{ref}
func (s *GitService) GetRef(ctx context.Context, owner, repo, ref string) (*Reference, *Response, error) {
	ref = strings.TrimPrefix(ref, "refs/")
	u := fmt.Sprintf("repos/%v/%v/git/ref/%v", owner, repo, refURLEscape(ref))
	req, err := s.client.NewRequest("GET", u, nil)
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

// refURLEscape escapes every path segment of the given ref. Those must
// not contain escaped "/" - as "%2F" - or github will not recognize it.
func refURLEscape(ref string) string {
	parts := strings.Split(ref, "/")
	for i, s := range parts {
		parts[i] = url.PathEscape(s)
	}
	return strings.Join(parts, "/")
}

// ReferenceListOptions specifies optional parameters to the
// GitService.ListMatchingRefs method.
type ReferenceListOptions struct {
	// The ref must be formatted as `heads/<branch name>` for branches and `tags/<tag name>` for tags.
	Ref string `url:"-"`

	ListOptions
}

// ListMatchingRefs lists references in a repository that match a supplied ref.
// Use an empty ref to list all references.
//
// GitHub API docs: https://docs.github.com/rest/git/refs#list-matching-references
//
//meta:operation GET /repos/{owner}/{repo}/git/matching-refs/{ref}
func (s *GitService) ListMatchingRefs(ctx context.Context, owner, repo string, opts *ReferenceListOptions) ([]*Reference, *Response, error) {
	var ref string
	if opts != nil {
		ref = strings.TrimPrefix(opts.Ref, "refs/")
	}
	u := fmt.Sprintf("repos/%v/%v/git/matching-refs/%v", owner, repo, refURLEscape(ref))
	u, err := addOptions(u, opts)
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
// GitHub API docs: https://docs.github.com/rest/git/refs#create-a-reference
//
//meta:operation POST /repos/{owner}/{repo}/git/refs
func (s *GitService) CreateRef(ctx context.Context, owner, repo string, ref *Reference) (*Reference, *Response, error) {
	if ref == nil {
		return nil, nil, errors.New("reference must be provided")
	}
	if ref.Ref == nil {
		return nil, nil, errors.New("ref must be provided")
	}

	u := fmt.Sprintf("repos/%v/%v/git/refs", owner, repo)
	req, err := s.client.NewRequest("POST", u, &createRefRequest{
		// back-compat with previous behavior that didn't require 'refs/' prefix
		Ref: Ptr("refs/" + strings.TrimPrefix(*ref.Ref, "refs/")),
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
// GitHub API docs: https://docs.github.com/rest/git/refs#update-a-reference
//
//meta:operation PATCH /repos/{owner}/{repo}/git/refs/{ref}
func (s *GitService) UpdateRef(ctx context.Context, owner, repo string, ref *Reference, force bool) (*Reference, *Response, error) {
	if ref == nil {
		return nil, nil, errors.New("reference must be provided")
	}
	if ref.Ref == nil {
		return nil, nil, errors.New("ref must be provided")
	}

	refPath := strings.TrimPrefix(*ref.Ref, "refs/")
	u := fmt.Sprintf("repos/%v/%v/git/refs/%v", owner, repo, refURLEscape(refPath))
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
// GitHub API docs: https://docs.github.com/rest/git/refs#delete-a-reference
//
//meta:operation DELETE /repos/{owner}/{repo}/git/refs/{ref}
func (s *GitService) DeleteRef(ctx context.Context, owner, repo, ref string) (*Response, error) {
	ref = strings.TrimPrefix(ref, "refs/")
	u := fmt.Sprintf("repos/%v/%v/git/refs/%v", owner, repo, refURLEscape(ref))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
