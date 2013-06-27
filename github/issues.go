// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/url"
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

// Issue represents a GitHub issue on a repository.
type Issue struct {
	Number    int        `json:"number,omitempty"`
	State     string     `json:"state,omitempty"`
	Title     string     `json:"title,omitempty"`
	Body      string     `json:"body,omitempty"`
	User      *User      `json:"user,omitempty"`
	Assignee  *User      `json:"assignee,omitempty"`
	Comments  int        `json:"comments,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// TODO(willnorris): labels and milestone
}

// IssueComment represents a comment left on an issue.
type IssueComment struct {
	ID        int        `json:"id,omitempty"`
	Body      string     `json:"body,omitempty"`
	User      *User      `json:"user,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
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
}

// List the issues for the authenticated user.  If all is true, list issues
// across all the user's visible repositories including owned, member, and
// organization repositories; if false, list only owned and member
// repositories.
//
// GitHub API docs: http://developer.github.com/v3/issues/#list-issues
func (s *IssuesService) List(all bool, opt *IssueListOptions) ([]Issue, error) {
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
func (s *IssuesService) ListByOrg(org string, opt *IssueListOptions) ([]Issue, error) {
	u := fmt.Sprintf("orgs/%v/issues", org)
	return s.listIssues(u, opt)
}

func (s *IssuesService) listIssues(u string, opt *IssueListOptions) ([]Issue, error) {
	if opt != nil {
		params := url.Values{
			"filter":    {opt.Filter},
			"state":     {opt.State},
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
		return nil, err
	}

	issues := new([]Issue)
	_, err = s.client.Do(req, issues)
	return *issues, err
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
func (s *IssuesService) ListByRepo(owner string, repo string, opt *IssueListByRepoOptions) ([]Issue, error) {
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
		return nil, err
	}

	issues := new([]Issue)
	_, err = s.client.Do(req, issues)
	return *issues, err
}

// Get a single issue.
//
// GitHub API docs: http://developer.github.com/v3/issues/#get-a-single-issue
func (s *IssuesService) Get(owner string, repo string, number int) (*Issue, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%d", owner, repo, number)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	issue := new(Issue)
	_, err = s.client.Do(req, issue)
	return issue, err
}

// Create a new issue on the specified repository.
//
// GitHub API docs: http://developer.github.com/v3/issues/#create-an-issue
func (s *IssuesService) Create(owner string, repo string, issue *Issue) (*Issue, error) {
	u := fmt.Sprintf("repos/%v/%v/issues", owner, repo)
	req, err := s.client.NewRequest("POST", u, issue)
	if err != nil {
		return nil, err
	}
	i := new(Issue)
	_, err = s.client.Do(req, i)
	return i, err
}

// Edit an issue.
//
// GitHub API docs: http://developer.github.com/v3/issues/#edit-an-issue
func (s *IssuesService) Edit(owner string, repo string, number int, issue *Issue) (*Issue, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%d", owner, repo, number)
	req, err := s.client.NewRequest("PATCH", u, issue)
	if err != nil {
		return nil, err
	}
	i := new(Issue)
	_, err = s.client.Do(req, i)
	return i, err
}

// ListAssignees fetches all available assignees (owners and collaborators) to
// which issues may be assigned.
//
// GitHub API docs: http://developer.github.com/v3/issues/assignees/#list-assignees
func (s *IssuesService) ListAssignees(owner string, repo string) ([]User, error) {
	u := fmt.Sprintf("repos/%v/%v/assignees", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	assignees := new([]User)
	_, err = s.client.Do(req, assignees)
	return *assignees, err
}

// CheckAssignee checks if a user is an assignee for the specified repository.
//
// GitHub API docs: http://developer.github.com/v3/issues/assignees/#check-assignee
func (s *IssuesService) CheckAssignee(owner string, repo string, user string) (bool, error) {
	u := fmt.Sprintf("repos/%v/%v/assignees/%v", owner, repo, user)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return false, err
	}
	_, err = s.client.Do(req, nil)
	return parseBoolResponse(err)
}

// IssueListCommentsOptions specifies the optional parameters to the
// IssuesService.ListComments method.
type IssueListCommentsOptions struct {
	// Sort specifies how to sort comments.  Possible values are: created, updated.
	Sort string

	// Direction in which to sort comments.  Possible values are: asc, desc.
	Direction string

	// Since filters comments by time.
	Since time.Time
}

// ListComments lists all comments on the specified issue.  Specifying an issue
// number of 0 will return all comments on all issues for the repository.
//
// GitHub API docs: http://developer.github.com/v3/issues/comments/#list-comments-on-an-issue
func (s *IssuesService) ListComments(owner string, repo string, number int, opt *IssueListCommentsOptions) ([]IssueComment, error) {
	var u string
	if number == 0 {
		u = fmt.Sprintf("repos/%v/%v/issues/comments", owner, repo)
	} else {
		u = fmt.Sprintf("repos/%v/%v/issues/%d/comments", owner, repo, number)
	}

	if opt != nil {
		params := url.Values{
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
		return nil, err
	}
	comments := new([]IssueComment)
	_, err = s.client.Do(req, comments)
	return *comments, err
}

// GetComment fetches the specified issue comment.
//
// GitHub API docs: http://developer.github.com/v3/issues/comments/#get-a-single-comment
func (s *IssuesService) GetComment(owner string, repo string, id int) (*IssueComment, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/comments/%d", owner, repo, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	comment := new(IssueComment)
	_, err = s.client.Do(req, comment)
	return comment, err
}

// CreateComment creates a new comment on the specified issue.
//
// GitHub API docs: http://developer.github.com/v3/issues/comments/#create-a-comment
func (s *IssuesService) CreateComment(owner string, repo string, number int, comment *IssueComment) (*IssueComment, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/%d/comments", owner, repo, number)
	req, err := s.client.NewRequest("POST", u, comment)
	if err != nil {
		return nil, err
	}
	c := new(IssueComment)
	_, err = s.client.Do(req, c)
	return c, err
}

// EditComment updates an issue comment.
//
// GitHub API docs: http://developer.github.com/v3/issues/comments/#edit-a-comment
func (s *IssuesService) EditComment(owner string, repo string, id int, comment *IssueComment) (*IssueComment, error) {
	u := fmt.Sprintf("repos/%v/%v/issues/comments/%d", owner, repo, id)
	req, err := s.client.NewRequest("PATCH", u, comment)
	if err != nil {
		return nil, err
	}
	c := new(IssueComment)
	_, err = s.client.Do(req, c)
	return c, err
}

// DeleteComment deletes an issue comment.
//
// GitHub API docs: http://developer.github.com/v3/issues/comments/#delete-a-comment
func (s *IssuesService) DeleteComment(owner string, repo string, id int) error {
	u := fmt.Sprintf("repos/%v/%v/issues/comments/%d", owner, repo, id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
