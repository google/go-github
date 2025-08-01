// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"context"
	"fmt"
)

// Blob represents a blob object.
type Blob struct {
	Content  *string `json:"content,omitempty"`
	Encoding *string `json:"encoding,omitempty"`
	SHA      *string `json:"sha,omitempty"`
	Size     *int    `json:"size,omitempty"`
	URL      *string `json:"url,omitempty"`
	NodeID   *string `json:"node_id,omitempty"`
}

// GetBlob fetches a blob from a repo given a SHA.
//
// GitHub API docs: https://docs.github.com/rest/git/blobs#get-a-blob
//
//meta:operation GET /repos/{owner}/{repo}/git/blobs/{file_sha}
func (s *GitService) GetBlob(ctx context.Context, owner, repo, sha string) (*Blob, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/git/blobs/%v", owner, repo, sha)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	blob := new(Blob)
	resp, err := s.client.Do(ctx, req, blob)
	if err != nil {
		return nil, resp, err
	}

	return blob, resp, nil
}

// GetBlobRaw fetches a blob's contents from a repo.
// Unlike GetBlob, it returns the raw bytes rather than the base64-encoded data.
//
// GitHub API docs: https://docs.github.com/rest/git/blobs#get-a-blob
//
//meta:operation GET /repos/{owner}/{repo}/git/blobs/{file_sha}
func (s *GitService) GetBlobRaw(ctx context.Context, owner, repo, sha string) ([]byte, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/git/blobs/%v", owner, repo, sha)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3.raw")

	var buf bytes.Buffer
	resp, err := s.client.Do(ctx, req, &buf)
	if err != nil {
		return nil, resp, err
	}

	return buf.Bytes(), resp, nil
}

// CreateBlob creates a blob object.
//
// GitHub API docs: https://docs.github.com/rest/git/blobs#create-a-blob
//
//meta:operation POST /repos/{owner}/{repo}/git/blobs
func (s *GitService) CreateBlob(ctx context.Context, owner, repo string, blob *Blob) (*Blob, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/git/blobs", owner, repo)
	req, err := s.client.NewRequest("POST", u, blob)
	if err != nil {
		return nil, nil, err
	}

	t := new(Blob)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}
