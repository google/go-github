// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/url"
	"strconv"
)

// ActivityService handles communication with the activity related
// methods of the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/users/
type ActivityService struct {
	client *Client
}

// ActivityListOptions specifies the optional parameters to the
// ActivityListOptions.ListStarred method.
type ActivityListOptions struct {
	// How to sort the repository list.  Possible values are: created, updated,
	// pushed, full_name.  Default is "full_name".
	Sort string

	// Direction in which to sort repositories.  Possible values are: asc, desc.
	// Default is "asc" when sort is "full_name", otherwise default is "desc".
	Direction string

	// For paginated result sets, page of results to retrieve.
	Page int
}

// ListStarred lists all the users starred repos
//
// GitHub API docs: http://developer.github.com/v3/activity/starring/#list-stargazers
func (s *ActivityService) ListStarred(user string, opt *ActivityListOptions) ([]Repository, error) {
	var u string
	if user != "" {
		u = fmt.Sprintf("users/%v/starred", user)
	} else {
		u = "user/starred"
	}
	if opt != nil {
		params := url.Values{
			"sort":      []string{opt.Sort},
			"direction": []string{opt.Direction},
			"page":      []string{strconv.Itoa(opt.Page)},
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
