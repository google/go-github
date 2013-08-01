// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
)

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
