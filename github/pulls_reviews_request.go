// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "fmt"

// PullRequestReview represents a review of a pull request.
type PullRequestReviewerRequest struct {
	Reviewers []string `json:"reviewers,omitempty"`
}

func (p PullRequestReviewerRequest) String() string {
	return Stringify(p)
}

// RequestReviewers submits a set of logins to be potential reviewers on a PR.
//
// GitHub API docs: https://developer.github.com/v3/pulls/review_requests/#create-a-review-request
func (s *PullRequestsService) RequestReviewers(owner, repo string, number int, reviewers *PullRequestReviewerRequest) (*PullRequest, *Response, error) {
	u := fmt.Sprintf("repos/%s/%s/pulls/%d/requested_reviewers", owner, repo, number)

	req, err := s.client.NewRequest("POST", u, reviewers)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypePullRequestReviewsPreview)

	r := new(PullRequest)
	resp, err := s.client.Do(req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}