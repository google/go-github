// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// IssuesService handles communication with the issue related
// methods of the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/issues/
type IssuesService struct {
	client *Client
}

// Label represents a GitHib label on an Issue
type Label struct {
	URL   string `json:"url,omitempty"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

func (label Label) String() string {
	return fmt.Sprintf(label.Name)
}

// Issue represents a GitHub issue on a repository.
type Issue struct {
	Number    int        `json:"number,omitempty"`
	State     string     `json:"state,omitempty"`
	Title     string     `json:"title,omitempty"`
	Body      string     `json:"body,omitempty"`
	User      *User      `json:"user,omitempty"`
	Labels    []*Label   `json:"labels,omitempty"`
	Assignee  *User      `json:"assignee,omitempty"`
	Comments  int        `json:"comments,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// TODO(willnorris): milestone
}

// IssueListOptions specifies the optional parameters to the IssuesService.List
// and IssuesService.ListByOrg methods.
type IssueListOptions struct {
	// Filter specifies which issues to list.  Possible values are: assigned,
	// created, mentioned, subscribed, all.  Default is "assigned".
	Filter string

	// State filters issues based on their state.  Possible values are: open,
	// closed.  Default is "open".
	State string

	// Labels filters issues based on their label.
	Labels []string

	// Sort specifies how to sort issues.  Possible values are: created, updated,
	// and comments.  Default value is "assigned".
	Sort string

	// Direction in which to sort issues.  Possible values are: asc, desc.
	// Default is "asc".
	Direction string

	// Since filters issues by time.
	Since time.Time

	// For paginated result sets, page of results to retrieve.
	Page int
}

// List the issues for the authenticated user.  If all is true, list issues
// across all the user's visible repositories including owned, member, and
// organization repositories; if false, list only owned and member
// repositories.
//
// GitHub API docs: http://developer.github.com/v3/issues/#list-issues
func (s *IssuesService) List(all bool, opt *IssueListOptions) ([]Issue, *Response, error) {
	var u string
	if all {
		u = "issues"
	} else {
		u = "user/issues"
	}
	return s.listIssues(u, opt)
}

// ListByOrg fetches the issues in the specified organization for the
// authenticated user.
//
// GitHub API docs: http://developer.github.com/v3/issues/#list-issues
func (s *IssuesService) ListByOrg(org string, opt *IssueListOptions) ([]Issue, *Response, error) {
	u := fmt.Sprintf("orgs/%v/issues", org)
	return s.listIssues(u, opt)
}

func (s *IssuesService) listIssues(u string, opt *IssueListOptions) ([]Issue, *Response, error) {
	if opt != nil {
		params := url.Values{
			"filter":    {opt.Filter},
			"state":     {opt.State},
			"labels":    {strings.Join(opt.Labels, ",")},
			"sort":      {opt.Sort},
			"direction": {opt.Direction},
			"page":      []string{strconv.Itoa(opt.Page)},
		}
		if !opt.Since.IsZero() {
			params.Add("since", opt.Since.Format(time.RFC3339))
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	issues := new([]Issue)
	resp, err := s.client.Do(req, issues)
	return *issues, resp, err
}

// IssueListByRepoOptions specifies the optional parameters to the
// IssuesService.ListByRepo method.
type IssueListByRepoOptions struct {
	// Milestone limits issues for the specified milestone.  Possible values are
	// a milestone number, "none" for issues with no milestone, "*" for issues
	// with any milestone.
	Milestone string

	// State filters issues based on their state.  Possible values are: open,
	// closed.  Default is "open".
	State string

	// Assignee filters issues based on their assignee.  Possible values are a
	// user name, "none" for issues that are not assigned, "*" for issues with
	// any assigned user.
	Assignee string

	// Assignee filters issues based on their creator.
	Creator string

	// Assignee filters issues to those mentioned a specific user.
	Mentioned string

	// Labels filters issues based on their label.
	Labels []string

	// Sort specifies how to sort issues.  Possible values are: created, updated,
	// and comments.  Default value is "assigned".
	Sort string

	// Direction in which to sort issues.  Possible values are: asc, desc.
	// Default is "asc".
	Direction string

	// Since filters issues by time.
	Since time.Time
}

// ListByRepo lists the issues for the specified repository.
//
// GitHub API docs: http://developer.github.com/v3/issues/#list-issues-for-a-repository
func (s *IssuesService) ListByRepo(owner string, repo string, opt *IssueListByRepoOptions) ([]Issue, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues", owner, repo)
	if opt != nil {
		params := url.Values{
			"milestone": {opt.Milestone},
			"state":     {opt.State},
			"assignee":  {opt.Assignee},
			"creator":   {opt.Creator},
			"mentioned": {opt.Mentioned},
			"labels":    {strings.Join(opt.Labels, ",")},
			"sort":      {opt.Sort},
			"direction": {opt.Direction},
		}
		if !opt.Since.IsZero() {
			params.Add("since", opt.Since.Format(time.RFC3339))
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	issues := new([]Issue)
	resp, err := s.client.Do(req, issues)
	return *issues, resp, err
}

// Get a single issue.
//
// GitHub API docs: http://developer.github.com/v3/issues/#get-a-single-issue
func (s *IssuesService) Get(owner string, repo string, number int) (*Issue, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%d", owner, repo, number)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	issue := new(Issue)
	resp, err := s.client.Do(req, issue)
	return issue, resp, err
}

// Create a new issue on the specified repository.
//
// GitHub API docs: http://developer.github.com/v3/issues/#create-an-issue
func (s *IssuesService) Create(owner string, repo string, issue *Issue) (*Issue, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues", owner, repo)
	req, err := s.client.NewRequest("POST", u, issue)
	if err != nil {
		return nil, nil, err
	}
	i := new(Issue)
	resp, err := s.client.Do(req, i)
	return i, resp, err
}

// Edit an issue.
//
// GitHub API docs: http://developer.github.com/v3/issues/#edit-an-issue
func (s *IssuesService) Edit(owner string, repo string, number int, issue *Issue) (*Issue, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%d", owner, repo, number)
	req, err := s.client.NewRequest("PATCH", u, issue)
	if err != nil {
		return nil, nil, err
	}
	i := new(Issue)
	resp, err := s.client.Do(req, i)
	return i, resp, err
}
