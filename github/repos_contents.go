// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Repository contents API methods.
// http://developer.github.com/v3/repos/contents/

package github

import (
	"encoding/base64"
	"errors"
	"fmt"
)

// RepositoryContent represents a file or directory in a github repository.
type RepositoryContent struct {
	Type     *string `json:"type,omitempty"`
	Encoding *string `json:"encoding,omitempty"`
	Size     *int    `json:"size,omitempty"`
	Name     *string `json:"name,omitempty"`
	Path     *string `json:"path,omitempty"`
	Content  *string `json:"content,omitempty"`
	SHA      *string `json:"sha,omitempty"`
	URL      *string `json:"url,omitempty"`
	GitURL   *string `json:"giturl,omitempty"`
	HTMLURL  *string `json:"htmlurl,omitempty"`
}

func (r RepositoryContent) String() string {
	return Stringify(r)
}

// Decode decodes the file content if it is base64 encoded.
func (c *RepositoryContent) Decode() ([]byte, error) {
	if *c.Encoding != "base64" {
		return nil, errors.New("Cannot decode non-base64")
	}
	o, err := base64.StdEncoding.DecodeString(*c.Content)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// GetReadme gets the Readme file for the repository.
//
// GitHub API docs: http://developer.github.com/v3/repos/contents/#get-the-readme
func (s *RepositoriesService) GetReadme(owner, repo string) (*RepositoryContent, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/readme", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	readme := new(RepositoryContent)
	resp, err := s.client.Do(req, readme)
	if err != nil {
		return nil, resp, err
	}
	return readme, resp, err
}
