// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/url"
	"time"
)

// GistsService handles communication with the Gist related
// methods of the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/gists/
type GistsService struct {
	client *Client
}

// Gist represents a GitHub's gist.
type Gist struct {
	ID          string                    `json:"id,omitempty"`
	Description string                    `json:"description,omitempty"`
	Public      bool                      `json:"public,omitempty"`
	User        *User                     `json:"user,omitempty"`
	Files       map[GistFilename]GistFile `json:"files,omitempty"`
	Comments    int                       `json:"comments,omitempty"`
	HTMLURL     string                    `json:"html_url,omitempty"`
	GitPullURL  string                    `json:"git_pull_url,omitempty"`
	GitPushURL  string                    `json:"git_push_url,omitempty"`
	CreatedAt   *time.Time                `json:"created_at,omitempty"`
}

// GistFilename represents filename on a gist.
type GistFilename string

// GistFile represents a file on a gist.
type GistFile struct {
	Size     int          `json:"size,omitempty"`
	Filename GistFilename `json:"filename,omitempty"`
	RawURL   string       `json:"raw_url,omitempty"`
	Content  string       `json:"content,omitempty"`
}

// GistListOptions specifies the optional parameters to the
// GistsService.List, GistsService.ListAll, and GistsService.ListStarred methods.
type GistListOptions struct {
	// Since filters Gists by time.
	Since time.Time
}

// List gists for a user. Passing the empty string will list
// all public gists if called anonymously. However, if the call
// is authenticated, it will returns all gists for the authenticated
// user.
//
// GitHub API docs: http://developer.github.com/v3/gists/#list-gists
func (s *GistsService) List(user string, opt *GistListOptions) ([]Gist, error) {
	var u string
	if user != "" {
		u = fmt.Sprintf("users/%v/gists", user)
	} else {
		u = "gists"
	}
	if opt != nil {
		params := url.Values{}
		if !opt.Since.IsZero() {
			params.Add("since", opt.Since.Format(time.RFC3339))
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	gists := new([]Gist)
	_, err = s.client.Do(req, gists)
	return *gists, err
}

// ListAll lists all public gists.
//
// GitHub API docs: http://developer.github.com/v3/gists/#list-gists
func (s *GistsService) ListAll(opt *GistListOptions) ([]Gist, error) {
	u := "gists/public"
	if opt != nil {
		params := url.Values{}
		if !opt.Since.IsZero() {
			params.Add("since", opt.Since.Format(time.RFC3339))
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	gists := new([]Gist)
	_, err = s.client.Do(req, gists)
	return *gists, err
}

// ListStarred lists starred gists of authenticated user.
//
// GitHub API docs: http://developer.github.com/v3/gists/#list-gists
func (s *GistsService) ListStarred(opt *GistListOptions) ([]Gist, error) {
	u := "gists/starred"
	if opt != nil {
		params := url.Values{}
		if !opt.Since.IsZero() {
			params.Add("since", opt.Since.Format(time.RFC3339))
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	gists := new([]Gist)
	_, err = s.client.Do(req, gists)
	return *gists, err
}

// Get a single gist.
//
// GitHub API docs: http://developer.github.com/v3/gists/#get-a-single-gist
func (s *GistsService) Get(id string) (*Gist, error) {
	u := fmt.Sprintf("gists/%v", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	gist := new(Gist)
	_, err = s.client.Do(req, gist)
	return gist, err
}

// Create a gist for authenticated user.
//
// GitHub API docs: http://developer.github.com/v3/gists/#create-a-gist
func (s *GistsService) Create(gist *Gist) (*Gist, error) {
	u := "gists"
	req, err := s.client.NewRequest("POST", u, gist)
	if err != nil {
		return nil, err
	}
	g := new(Gist)
	_, err = s.client.Do(req, g)
	return g, err
}

// Edit a gist.
//
// GitHub API docs: http://developer.github.com/v3/gists/#edit-a-gist
func (s *GistsService) Edit(id string, gist *Gist) (*Gist, error) {
	u := fmt.Sprintf("gists/%v", id)
	req, err := s.client.NewRequest("PATCH", u, gist)
	if err != nil {
		return nil, err
	}
	g := new(Gist)
	_, err = s.client.Do(req, g)
	return g, err
}

// Delete a gist.
//
// GitHub API docs: http://developer.github.com/v3/gists/#delete-a-gist
func (s *GistsService) Delete(id string) error {
	u := fmt.Sprintf("gists/%v", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Star a gist on behalf of authenticated user.
//
// GitHub API docs: http://developer.github.com/v3/gists/#star-a-gist
func (s *GistsService) Star(id string) error {
	u := fmt.Sprintf("gists/%v/star", id)
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Unstar a gist on a behalf of authenticated user.
//
// Github API docs: http://developer.github.com/v3/gists/#unstar-a-gist
func (s *GistsService) Unstar(id string) error {
	u := fmt.Sprintf("gists/%v/star", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Starred checks if a gist is starred by authenticated user.
//
// GitHub API docs: http://developer.github.com/v3/gists/#check-if-a-gist-is-starred
func (s *GistsService) Starred(id string) (bool, error) {
	u := fmt.Sprintf("gists/%v/star", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return false, err
	}
	_, err = s.client.Do(req, nil)
	return parseBoolResponse(err)
}

// Fork a gist.
//
// GitHub API docs: http://developer.github.com/v3/gists/#fork-a-gist
func (s *GistsService) Fork(id string) (*Gist, error) {
	u := fmt.Sprintf("gists/%v/forks", id)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}
	g := new(Gist)
	_, err = s.client.Do(req, g)
	return g, err
}
