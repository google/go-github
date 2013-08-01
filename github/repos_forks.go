// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/url"
)

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
	u := fmt.Sprintf("repos/%v/%v/forks", owner, repo)
	if opt != nil {
		params := url.Values{
			"sort": []string{opt.Sort},
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
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
	u := fmt.Sprintf("repos/%v/%v/forks", owner, repo)
	if opt != nil {
		params := url.Values{
			"organization": []string{opt.Organization},
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	fork := new(Repository)
	_, err = s.client.Do(req, fork)
	return fork, err
}
