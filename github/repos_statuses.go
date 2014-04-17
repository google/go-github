// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"time"
)

// RepoStatus represents the status of a repository at a particular reference.
type RepoStatus struct {
	ID *int `json:"id,omitempty"`

	// State is the current state of the repository.  Possible values are:
	// pending, success, error, or failure.
	State *string `json:"state,omitempty"`

	// TargetURL is the URL of the page representing this status.  It will be
	// linked from the GitHub UI to allow users to see the source of the status.
	TargetURL *string `json:"target_url,omitempty"`

	// Description is a short high level summary of the status.
	Description *string `json:"description,omitempty"`

	// A string label to differentiate this status from the statuses of other systems.
	Context *string `json:"context,omitempty"`

	Creator   *User      `json:"creator,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (r RepoStatus) String() string {
	return Stringify(r)
}

// ListStatuses lists the statuses of a repository at the specified
// reference.  ref can be a SHA, a branch name, or a tag name.
//
// GitHub API docs: http://developer.github.com/v3/repos/statuses/#list-statuses-for-a-specific-ref
func (s *RepositoriesService) ListStatuses(owner, repo, ref string, opt *ListOptions) ([]RepoStatus, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/statuses/%v", owner, repo, ref)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	statuses := new([]RepoStatus)
	resp, err := s.client.Do(req, statuses)
	if err != nil {
		return nil, resp, err
	}

	return *statuses, resp, err
}

// CreateStatus creates a new status for a repository at the specified
// reference.  Ref can be a SHA, a branch name, or a tag name.
//
// GitHub API docs: http://developer.github.com/v3/repos/statuses/#create-a-status
func (s *RepositoriesService) CreateStatus(owner, repo, ref string, status *RepoStatus) (*RepoStatus, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/statuses/%v", owner, repo, ref)
	req, err := s.client.NewRequest("POST", u, status)
	if err != nil {
		return nil, nil, err
	}

	statuses := new(RepoStatus)
	resp, err := s.client.Do(req, statuses)
	if err != nil {
		return nil, resp, err
	}

	return statuses, resp, err
}
