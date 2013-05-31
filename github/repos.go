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

// Repository represents a GitHub repository.
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
//
// GitHub API docs: http://developer.github.com/v3/repos/#list-user-repositories
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

// ListByOrg lists the repositories for an organization.
//
// GitHub API docs: http://developer.github.com/v3/repos/#list-organization-repositories
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

// RepositoryListAllOptions specifies the optional parameters to the
// RepositoriesService.ListAll method.
type RepositoryListAllOptions struct {
	// ID of the last repository seen
	Since int
}

// ListAll lists all GitHub repositories in the order that they were created.
//
// GitHub API docs: http://developer.github.com/v3/repos/#list-all-repositories
func (s *RepositoriesService) ListAll(opt *RepositoryListAllOptions) ([]Repository, error) {
	url_ := "repositories"
	if opt != nil {
		params := url.Values{
			"since": []string{strconv.Itoa(opt.Since)},
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

// Create a new repository.  If an organization is specified, the new
// repository will be created under that org.  If the empty string is
// specified, it will be created for the authenticated user.
//
// GitHub API docs: http://developer.github.com/v3/repos/#create
func (s *RepositoriesService) Create(org string, repo *Repository) (*Repository, error) {
	var url_ string
	if org != "" {
		url_ = fmt.Sprintf("orgs/%v/repos", org)
	} else {
		url_ = "user/repos"
	}

	req, err := s.client.NewRequest("POST", url_, repo)
	if err != nil {
		return nil, err
	}

	r := new(Repository)
	_, err = s.client.Do(req, r)
	return r, err
}

// Get fetches a repository.
//
// GitHub API docs: http://developer.github.com/v3/repos/#get
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

// RepositoryListForksOptions specifies the optional parameters to the
// RepositoriesService.ListForks method.
type RepositoryListForksOptions struct {
	// How to sort the forks list.  Possible values are: newest, oldest,
	// watchers.  Default is "newest".
	Sort string
}

// ListForks lists the forks of the specified repository.
//
// GitHub API docs: http://developer.github.com/v3/repos/forks/#list-forks
func (s *RepositoriesService) ListForks(owner, repo string, opt *RepositoryListForksOptions) ([]Repository, error) {
	url_ := fmt.Sprintf("repos/%v/%v/forks", owner, repo)
	if opt != nil {
		params := url.Values{
			"sort": []string{opt.Sort},
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

// RepositoryCreateForkOptions specifies the optional parameters to the
// RepositoriesService.CreateFork method.
type RepositoryCreateForkOptions struct {
	// The organization to fork the repository into.
	Organization string
}

// CreateFork creates a fork of the specified repository.
//
// GitHub API docs: http://developer.github.com/v3/repos/forks/#list-forks
func (s *RepositoriesService) CreateFork(owner, repo string, opt *RepositoryCreateForkOptions) (*Repository, error) {
	url_ := fmt.Sprintf("repos/%v/%v/forks", owner, repo)
	if opt != nil {
		params := url.Values{
			"organization": []string{opt.Organization},
		}
		url_ += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("POST", url_, nil)
	if err != nil {
		return nil, err
	}

	fork := new(Repository)
	_, err = s.client.Do(req, fork)
	return fork, err
}

// RepoStatus represents the status of a repository at a particular reference.
type RepoStatus struct {
	ID int `json:"id,omitempty"`

	// State is the current state of the repository.  Possible values are:
	// pending, success, error, or failure.
	State string `json:"state,omitempty"`

	// TargetURL is the URL of the page representing this status.  It will be
	// linked from the GitHub UI to allow users to see the source of the status.
	TargetURL string `json:"target_url,omitempty"`

	// Description is a short high level summary of the status.
	Description string `json:"description,omitempty"`

	Creator   *User      `json:"creator,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// ListStatuses lists the statuses of a repository at the specified
// reference.  ref can be a SHA, a branch name, or a tag name.
//
// GitHub API docs: http://developer.github.com/v3/repos/statuses/#list-statuses-for-a-specific-ref
func (s *RepositoriesService) ListStatuses(owner, repo, ref string) ([]RepoStatus, error) {
	url_ := fmt.Sprintf("repos/%v/%v/statuses/%v", owner, repo, ref)
	req, err := s.client.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	statuses := new([]RepoStatus)
	_, err = s.client.Do(req, statuses)
	return *statuses, err
}

// CreateStatus creates a new status for a repository at the specified
// reference.  Ref can be a SHA, a branch name, or a tag name.
//
// GitHub API docs: http://developer.github.com/v3/repos/statuses/#create-a-status
func (s *RepositoriesService) CreateStatus(owner, repo, ref string, status *RepoStatus) (*RepoStatus, error) {
	url_ := fmt.Sprintf("repos/%v/%v/statuses/%v", owner, repo, ref)
	req, err := s.client.NewRequest("POST", url_, status)
	if err != nil {
		return nil, err
	}

	statuses := new(RepoStatus)
	_, err = s.client.Do(req, statuses)
	return statuses, err
}
