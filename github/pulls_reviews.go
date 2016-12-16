// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"time"
)

// PullRequestReview represents a review of a pull request.
type PullRequestReview struct {
	ID             *int       `json:"id,omitempty"`
	User           *User      `json:"user,omitempty"`
	Body           *string    `json:"body,omitempty"`
	SubmittedAt    *time.Time `json:"submitted_at,omitempty"`
	CommitID       *string    `json:"commit_id,omitempty"`
	HTMLURL        *string    `json:"html_url,omitempty"`
	PullRequestURL *string    `json:"pull_request_url,omitempty"`

	// State can be "ACCEPTED", "DISMISSED", "CHANGES_REQUESTED" or "COMMENTED".
	State *string `json:"state,omitempty"`
}

func (p PullRequestReview) String() string {
	return Stringify(p)
}

// ListReviews lists all reviews on the specified pull request.
//
// GitHub API docs: https://developer.github.com/v3/pulls/reviews/#list-reviews-on-a-pull-request
func (s *PullRequestsService) ListReviews(owner string, repo string, number int) ([]*PullRequestReview, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/%d/reviews", owner, repo, number)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypePullRequestReviewsPreview)

	reviews := new([]*PullRequestReview)
	resp, err := s.client.Do(req, reviews)
	if err != nil {
		return nil, resp, err
	}

	return *reviews, resp, err
}
