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

// PullRequestReviewComment represents a comment left on a pull request review.
type PullRequestReviewComment struct {
	ID             *int       `json:"id,omitempty"`
	User           *User      `json:"user,omitempty"`
	Body           *string    `json:"body,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	CommitID       *string    `json:"commit_id,omitempty"`
	HTMLURL        *string    `json:"html_url,omitempty"`
	PullRequestURL *string    `json:"pull_request_url,omitempty"`
}

// DraftReviewComment represents a comment part of the review.
type DraftReviewComment struct {
	Path     *string `json:"path,omitempty"`
	Position *int    `json:"position,omitempty"`
	Body     *string `json:"body,omitempty"`
}

// PullRequestReviewRequest represents a request to create a review.
type PullRequestReviewRequest struct {
	Body     *string              `json:"body,omitempty"`
	Event    *string              `json:"event,omitempty"`
	Comments []DraftReviewComment `json:"comments,omitempty"`
}

// PullRequestReviewDismissalRequest represents a request to dismiss a review.
type PullRequestReviewDismissalRequest struct {
	Message *string `json:"message,omitempty"`
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

// GetReview fetches the specified pull request review.
//
// GitHub API docs: https://developer.github.com/v3/pulls/reviews/#get-a-single-review
func (s *PullRequestsService) GetReview(owner string, repo string, number int, id int) (*PullRequestReview, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/%d/reviews/%d", owner, repo, number, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypePullRequestReviewsPreview)

	review := new(PullRequestReview)
	resp, err := s.client.Do(req, review)
	if err != nil {
		return nil, resp, err
	}

	return review, resp, err
}

// ListReviewComments lists all the comments for the specified review.
//
// GitHub API docs: https://developer.github.com/v3/pulls/reviews/#get-a-single-reviews-comments
func (s *PullRequestsService) ListReviewComments(owner string, repo string, number int, id int) ([]*PullRequestReviewComment, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/%d/reviews/%d/comments", owner, repo, number, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypePullRequestReviewsPreview)

	comments := new([]*PullRequestReviewComment)
	resp, err := s.client.Do(req, comments)
	if err != nil {
		return nil, resp, err
	}

	return *comments, resp, err
}

// CreateReview creates a new review on the specified pull request.
//
// GitHub API docs: https://developer.github.com/v3/pulls/reviews/#create-a-pull-request-review
func (s *PullRequestsService) CreateReview(owner string, repo string, number int, review *PullRequestReviewRequest) (*PullRequestReview, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/%d/reviews", owner, repo, number)

	req, err := s.client.NewRequest("POST", u, review)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypePullRequestReviewsPreview)

	r := new(PullRequestReview)
	resp, err := s.client.Do(req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}

// SubmitReview submits a specified review on the specified pull request.
//
// GitHub API docs: https://developer.github.com/v3/pulls/reviews/#submit-a-pull-request-review
func (s *PullRequestsService) SubmitReview(owner string, repo string, number int, id int, review *PullRequestReviewRequest) (*PullRequestReview, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/%d/reviews/%d/events", owner, repo, number, id)

	req, err := s.client.NewRequest("POST", u, review)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypePullRequestReviewsPreview)

	r := new(PullRequestReview)
	resp, err := s.client.Do(req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}

// DismissReview dismisses a specified review on the specified pull request.
//
// GitHub API docs: https://developer.github.com/v3/pulls/reviews/#dismiss-a-pull-request-review
func (s *PullRequestsService) DismissReview(owner string, repo string, number int, id int, review *PullRequestReviewDismissalRequest) (*PullRequestReview, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/pulls/%d/reviews/%d/dismissals", owner, repo, number, id)

	req, err := s.client.NewRequest("PUT", u, review)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypePullRequestReviewsPreview)

	r := new(PullRequestReview)
	resp, err := s.client.Do(req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}
