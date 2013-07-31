// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// SearchService handles communication with the GitHub Search API
// (http://developer.github.com/v3/search/)
type SearchService struct {
	client *Client
}

type SortOrder string

const (
	SortOrder_Desc SortOrder = "desc"
	SortOrder_Asc            = "asc"
)

type RepositorySearchOptions struct {
	Sort      string
	Order     SortOrder
	TextMatch bool
	Page      int
	PerPage   int // up to 100 results per page according to GitHub
}

type RepositorySearchResults struct {
	TotalCount int          `json:"total_count"`
	Items      []Repository `json:"items"`
}

func (s *SearchService) Repositories(query string, opt *RepositorySearchOptions) (*RepositorySearchResults, error) {
	path := "search/repositories?"
	param := url.Values{"q": []string{query}}
	path += param.Encode()
	textMatch := false
	if opt != nil {
		textMatch = opt.TextMatch
		auxParams := url.Values{}
		if opt.Sort != "" {
			auxParams["sort"] = []string{opt.Sort}
		}
		if string(opt.Order) != "" {
			auxParams["order"] = []string{string(opt.Order)}
		}
		if opt.Page > 0 {
			auxParams["page"] = []string{strconv.Itoa(opt.Page)}
		}
		if opt.PerPage > 0 {
			auxParams["per_page"] = []string{strconv.Itoa(opt.PerPage)}
		}
		if len(auxParams) > 0 {
			path += "&" + auxParams.Encode()
		}
	}

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	enableGithubExperimental(req, textMatch)

	result := new(RepositorySearchResults)
	_, err = s.client.Do(req, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

var experimentalAccept = "application/vnd.github.preview"

// This is necessary until the new GitHub search API is officially
// released, at which point this function can be deleted
func enableGithubExperimental(req *http.Request, textMatch bool) {
	if textMatch {
		req.Header.Add("Accept", fmt.Sprintf("%s.text-match+json", experimentalAccept))
	} else {
		req.Header.Add("Accept", experimentalAccept)
	}
}
