// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "fmt"

// RepositoryRelease represents a GitHub release in a repository.
type RepositoryRelease struct {
	ID              *int       `json:"id,omitempty"`
	TagName         *string    `json:"tag_name,omitempty"`
	TargetCommitish *string    `json:"target_commitish,omitempty"`
	Name            *string    `json:"name,omitempty"`
	Body            *string    `json:"body,omitempty"`
	Draft           *bool      `json:"draft"`
	Prerelease      *bool      `json:"prerelease"`
	CreatedAt       *Timestamp `json:"created_at,omitempty"`
	PublishedAt     *Timestamp `json:"published_at,omitempty"`
	URL             *string    `json:"url,omitempty"`
	HTMLURL         *string    `json:"html_url,omitempty"`
	AssertsURL      *string    `json:"assets_url,omitempty"`
	UploadURL       *string    `json:"upload_url,omitempty"`
}

func (r RepositoryRelease) String() string {
	return Stringify(r)
}

// ListReleases lists the releases for a repository.
//
// GitHub API docs: http://developer.github.com/v3/repos/releases/#list-releases-for-a-repository
func (s *RepositoriesService) ListReleases(owner, repo string) ([]RepositoryRelease, *Response, error) {
	u := fmt.Sprintf("repos/%s/%s/releases", owner, repo)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", mimeReleasePreview)

	releases := new([]RepositoryRelease)
	resp, err := s.client.Do(req, releases)
	if err != nil {
		return nil, resp, err
	}
	return *releases, resp, err
}

// GetRelease fetches a single release.
//
// GitHub API docs: http://developer.github.com/v3/repos/releases/#get-a-single-release
func (s *RepositoriesService) GetRelease(owner, repo string, id int) (*RepositoryRelease, *Response, error) {
	u := fmt.Sprintf("repos/%s/%s/releases/%d", owner, repo, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", mimeReleasePreview)

	release := new(RepositoryRelease)
	resp, err := s.client.Do(req, release)
	if err != nil {
		return nil, resp, err
	}
	return release, resp, err
}

// CreateRelease adds a new release for a repository.
//
// GitHub API docs : http://developer.github.com/v3/repos/releases/#create-a-release
func (s *RepositoriesService) CreateRelease(owner, repo string, release *RepositoryRelease) (*RepositoryRelease, *Response, error) {
	u := fmt.Sprintf("repos/%s/%s/releases", owner, repo)

	req, err := s.client.NewRequest("POST", u, release)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", mimeReleasePreview)

	r := new(RepositoryRelease)
	resp, err := s.client.Do(req, r)
	if err != nil {
		return nil, resp, err
	}
	return r, resp, err
}

// EditRelease edits a repository release.
//
// GitHub API docs : http://developer.github.com/v3/repos/releases/#edit-a-release
func (s *RepositoriesService) EditRelease(owner, repo string, id int, release *RepositoryRelease) (*RepositoryRelease, *Response, error) {
	u := fmt.Sprintf("repos/%s/%s/releases/%d", owner, repo, id)

	req, err := s.client.NewRequest("PATCH", u, release)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", mimeReleasePreview)

	r := new(RepositoryRelease)
	resp, err := s.client.Do(req, r)
	if err != nil {
		return nil, resp, err
	}
	return r, resp, err
}

// DeleteRelease delete a single release from a repository.
//
// GitHub API docs : http://developer.github.com/v3/repos/releases/#delete-a-release
func (s *RepositoriesService) DeleteRelease(owner, repo string, id int) (*Response, error) {
	u := fmt.Sprintf("repos/%s/%s/releases/%d", owner, repo, id)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", mimeReleasePreview)
	return s.client.Do(req, nil)
}
