// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "fmt"

// Pages represents a GitHub Pages site configuration.
type Pages struct {
	URL       *string `json:"url,omitempty"`
	Status    *string `json:"status,omitempty"`
	CNAME     *string `json:"cname,omitempty"`
	Custom404 *bool   `json:"custom_404"`
}

// PagesError represents a build error for a GitHub Pages site.
type PagesError struct {
	Message *string `json:"message,omitempty"`
}

// PagesBuild represents the build information for a GitHub Pages site.
type PagesBuild struct {
	URL       *string      `json:"url,omitempty"`
	Status    *string      `json:"status,omitempty"`
	Error     *PagesError  `json:"error,omitempty"`
	Pusher    *Contributor `json:"pusher,omitempty"`
	Commit    *string      `json:"commit,omitempty"`
	Duration  *int         `json:"duration,omitempty"`
	CreatedAt *Timestamp   `json:"created_at,omitempty"`
	UpdatedAt *Timestamp   `json:"created_at,omitempty"`
}

// GetPages fetches information about a GitHub Pages site.
//
// GitHub API docs: https://developer.github.com/v3/repos/pages/#get-information-about-a-pages-site
func (s *RepositoriesService) GetPages(owner string, repo string) (*Pages, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pages", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	site := new(Pages)
	resp, err := s.client.Do(req, site)
	if err != nil {
		return nil, resp, err
	}

	return site, resp, err
}

// ListPagesBuilds lists the builds for a GitHub Pages site.
//
// GitHub API docs: https://developer.github.com/v3/repos/pages/#list-pages-builds
func (s *RepositoriesService) ListPagesBuilds(owner string, repo string) ([]PagesBuild, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pages/builds", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	pages := new([]PagesBuild)
	resp, err := s.client.Do(req, pages)
	if err != nil {
		return nil, resp, err
	}

	return *pages, resp, err
}

// ListLatestPagesBuilds fetches the latest build information for a GitHub pages site.
//
// GitHub API docs: https://developer.github.com/v3/repos/pages/#list-latest-pages-build
func (s *RepositoriesService) ListLatestPagesBuilds(owner string, repo string) (*PagesBuild, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pages/builds/latest", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	build := new(PagesBuild)
	resp, err := s.client.Do(req, build)
	if err != nil {
		return nil, resp, err
	}

	return build, resp, err
}
