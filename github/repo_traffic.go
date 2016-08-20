// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "fmt"

// Referrer represent a referrer information
type Referrer struct {
	Referrer *string `json:"referrer,omitempty"`
	Count    *int    `json:"count,omitempty"`
	Uniques  *int    `json:"uniques,omitempty"`
}

// Path represent a referrer information
type Path struct {
	Path    *string `json:"path,omitempty"`
	Title   *string `json:"title,omitempty"`
	Count   *int    `json:"count,omitempty"`
	Uniques *int    `json:"uniques,omitempty"`
}

// ListReferrers list the top 10 referrers over the last 14 days.
//
// GitHub API docs: https://developer.github.com/v3/repos/traffic/#list-referrers
func (s *RepositoriesService) ListReferrers(owner, repo string, opt *ListOptions) ([]*Referrer, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/traffic/popular/referrers", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeTrafficPreview)

	var referrers []*Referrer
	resp, err := s.client.Do(req, &referrers)
	if err != nil {
		return nil, resp, err
	}

	return referrers, resp, err
}

// ListPaths list the top 10 popular content over the last 14 days.
//
// GitHub API docs: https://developer.github.com/v3/repos/traffic/#list-paths
func (s *RepositoriesService) ListPaths(owner, repo string, opt *ListOptions) ([]*Path, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/traffic/popular/paths", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeTrafficPreview)

	var paths []*Path
	resp, err := s.client.Do(req, &paths)
	if err != nil {
		return nil, resp, err
	}

	return paths, resp, err
}
