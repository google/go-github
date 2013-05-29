// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package github

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// RepositoriesService handles communication with the repository related
// methods of the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/repos/
type RepositoriesService struct {
	client *Client
}

type Repository struct {
	ID          int        `json:"id,omitempty"`
	Owner       *User      `json:"owner,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	PushedAt    *time.Time `json:"pushed_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// RepositoryListOptions specifies the optional parameters to the
// RepositoriesService.List method.
type RepositoryListOptions struct {
	// Type of repositories to list.  Possible values are: all, owner, public,
	// private, member.  Default is "all".
	Type string

	// How to sort the repository list.  Possible values are: created, updated,
	// pushed, full_name.  Default is "full_name".
	Sort string

	// Direction in which to sort repositories.  Possible values are: asc, desc.
	// Default is "asc" when sort is "full_name", otherwise default is "desc".
	Direction string

	// For paginated result sets, page of results to retrieve.
	Page int
}

// List the repositories for a user.  Passing the empty string will list
// repositories for the authenticated user.
func (s *RepositoriesService) List(user string, opt *RepositoryListOptions) ([]Repository, error) {
	var url_ string
	if user != "" {
		url_ = fmt.Sprintf("users/%v/repos", user)
	} else {
		url_ = "user/repos"
	}
	if opt != nil {
		params := url.Values{
			"type":      []string{opt.Type},
			"sort":      []string{opt.Sort},
			"direction": []string{opt.Direction},
			"page":      []string{strconv.Itoa(opt.Page)},
		}
		url_ += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	repos := new([]Repository)
	_, err = s.client.Do(req, repos)
	return *repos, err
}

// RepositoryListByOrgOptions specifies the optional parameters to the
// RepositoriesService.ListByOrg method.
type RepositoryListByOrgOptions struct {
	// Type of repositories to list.  Possible values are: all, public, private,
	// forks, sources, member.  Default is "all".
	Type string

	// For paginated result sets, page of results to retrieve.
	Page int
}

// List the repositories for an organization.
func (s *RepositoriesService) ListByOrg(org string, opt *RepositoryListByOrgOptions) ([]Repository, error) {
	url_ := fmt.Sprintf("orgs/%v/repos", org)
	if opt != nil {
		params := url.Values{
			"type": []string{opt.Type},
			"page": []string{strconv.Itoa(opt.Page)},
		}
		url_ += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	repos := new([]Repository)
	_, err = s.client.Do(req, repos)
	return *repos, err
}

// Get fetches a repository.
func (s *RepositoriesService) Get(owner, repo string) (*Repository, error) {
	url_ := fmt.Sprintf("repos/%v/%v", owner, repo)
	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}
	repository := new(Repository)
	_, err = s.client.Do(req, repository)
	return repository, err
}
