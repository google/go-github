// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"time"
)

// Milestone represents a Github repository milestone.
type Milestone struct {
	URL          *string    `json:"url,omitempty"`
	Number       *int       `json:"number,omitempty"`
	State        *string    `json:"state,omitempty"`
	Title        *string    `json:"title,omitempty"`
	Description  *string    `json:"description,omitempty"`
	Creator      *User      `json:"creator,omitempty"`
	OpenIssues   *int       `json:"open_issues,omitempty"`
	ClosedIssues *int       `json:"closed_issues,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	DueOne       *time.Time `json:"due_on,omitempty"`
}

func (m Milestone) String() string {
	return Stringify(m)
}

// MilestoneListOptions specifies the optional parameters to the
// IssuesService.ListMilestones method.
type MilestoneListOptions struct {
	// State filters milestones based on their state. Possible values are:
	// open, closed. Default is "open".
	State string `url:"state,omitempty"`

	// Sort specifies how to sort milestones. Possible values are: due_date, completeness.
	// Default value is "due_date".
	Sort string `url:"sort,omitempty"`

	// Direction in which to sort milestones. Possible values are: asc, desc.
	// Default is "asc".
	Direction string `url:"direction,omitempty"`
}

// ListMilestones lists all milestones for a repository.
//
// GitHub API docs: https://developer.github.com/v3/issues/milestones/#list-milestones-for-a-repository
func (s *IssuesService) ListMilestones(owner string, repo string, opt *MilestoneListOptions) ([]Milestone, *Response, error) {
	u := fmt.Sprintf("/repos/%v/%v/milestones", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	milestones := new([]Milestone)
	resp, err := s.client.Do(req, milestones)
	if err != nil {
		return nil, resp, err
	}

	return *milestones, resp, err
}

// GetMilestone gets a single milestone.
//
// GitHub API docs: https://developer.github.com/v3/issues/milestones/#get-a-single-milestone
func (s *IssuesService) GetMilestone(owner string, repo string, number int) (*Milestone, *Response, error) {
	u := fmt.Sprintf("/repos/%v/%v/milestones/%d", owner, repo, number)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	milestone := new(Milestone)
	resp, err := s.client.Do(req, milestone)
	if err != nil {
		return nil, resp, err
	}

	return milestone, resp, err
}

