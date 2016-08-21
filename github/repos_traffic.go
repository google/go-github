// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"strconv"
	"time"
)

// Referrer represent a referrer information.
type Referrer struct {
	Referrer *string `json:"referrer,omitempty"`
	Count    *int    `json:"count,omitempty"`
	Uniques  *int    `json:"uniques,omitempty"`
}

// Path represent a referrer information.
type Path struct {
	Path    *string `json:"path,omitempty"`
	Title   *string `json:"title,omitempty"`
	Count   *int    `json:"count,omitempty"`
	Uniques *int    `json:"uniques,omitempty"`
}

// Time represents a time stamp as in a datapoint.
//
// This is necessary because github uses milliseconds timestamp.
type Time struct {
	time.Time
}

// UnmarshalJSON parse unix timestamp.
func (t *Time) UnmarshalJSON(b []byte) error {
	s := string(b)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	// We can drop the reaminder as returned values are days and it will always be 0
	*t = Time{time.Unix(i/1000, 0)}
	return nil
}

// Datapoint represent a view in views on clones.
type Datapoint struct {
	Timestamp *Time `json:"timestamp,omitempty"`
	Count     *int  `json:"count,omitempty"`
	Uniques   *int  `json:"uniques,omitempty"`
}

// Views represent information about views on the last 14 days.
type Views struct {
	Views   *[]Datapoint `json:"views,omitempty"`
	Count   *int         `json:"count,omitempty"`
	Uniques *int         `json:"uniques,omitempty"`
}

// Clones represent information about clones on the last 14 days.
type Clones struct {
	Clones  *[]Datapoint `json:"clones,omitempty"`
	Count   *int         `json:"count,omitempty"`
	Uniques *int         `json:"uniques,omitempty"`
}

// ListReferrers list the top 10 referrers over the last 14 days.
//
// GitHub API docs: https://developer.github.com/v3/repos/traffic/#list-referrers
func (s *RepositoriesService) ListReferrers(owner, repo string) ([]*Referrer, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/traffic/popular/referrers", owner, repo)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeTrafficPreview)

	referrers := new([]*Referrer)
	resp, err := s.client.Do(req, &referrers)
	if err != nil {
		return nil, resp, err
	}

	return *referrers, resp, err
}

// ListPaths list the top 10 popular content over the last 14 days.
//
// GitHub API docs: https://developer.github.com/v3/repos/traffic/#list-paths
func (s *RepositoriesService) ListPaths(owner, repo string) ([]*Path, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/traffic/popular/paths", owner, repo)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeTrafficPreview)

	var paths = new([]*Path)
	resp, err := s.client.Do(req, &paths)
	if err != nil {
		return nil, resp, err
	}

	return *paths, resp, err
}

// ListViews get total number of views for the last 14 days and breaks it down either per day or week.
//
// GitHub API docs: https://developer.github.com/v3/repos/traffic/#views
func (s *RepositoriesService) ListViews(owner, repo string, opt *TrafficBreakdownOptions) (*Views, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/traffic/views", owner, repo)
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

	views := new(Views)
	resp, err := s.client.Do(req, &views)
	if err != nil {
		return nil, resp, err
	}

	return views, resp, err
}

// ListClones get total number of clones for the last 14 days and breaks it down either per day or week for the last 14 days.
//
// GitHub API docs: https://developer.github.com/v3/repos/traffic/#views
func (s *RepositoriesService) ListClones(owner, repo string, opt *TrafficBreakdownOptions) (*Clones, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/traffic/clones", owner, repo)
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

	clones := new(Clones)
	resp, err := s.client.Do(req, &clones)
	if err != nil {
		return nil, resp, err
	}

	return clones, resp, err
}
