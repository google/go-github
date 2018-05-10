// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ChecksService provides access to the the Checks API in the
// GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/checks/
type ChecksService service

// CheckRun represents a GitHub check run on a repository associated with a GitHub app
type CheckRun struct {
	ID           *int64          `json:"id,omitempty"`
	HeadSHA      *string         `json:"head_sha,omitempty"`
	ExternalID   *int64          `json:"external_id,omitempty"`
	URL          *string         `json:"url,omitempty"`
	HTMLURL      *string         `json:"html_url,omitempty"`
	Status       *string         `json:"status,omitempty"`
	Conclusion   *string         `json:"conclusion,omitempty"`
	StartedAt    *time.Time      `json:"started_at,omitempty"`
	CompletedAt  *time.Time      `json:"completed_at,omitempty"`
	Output       *CheckRunOutput `json:"output,omitempty"`
	Name         *string         `json:"name,omitempty"`
	CheckSuite   *CheckSuite     `json:"check_suite,omitempty"`
	App          *App            `json:"app,omitempty"`
	PullRequests []*PullRequest  `json:"pull_requests,omitempty"`
}

// CheckRunOutput represents the output of a CheckRun
type CheckRunOutput struct {
	Title            *string `json:"title,omitempty"`
	Summary          *string `json:"summary,omitempty"`
	Text             *string `json:"text,omitempty"`
	AnnotationsCount *int64  `json:"annotations_count,omitempty"`
	AnnotationsURL   *string `json:"annotations_url,omitempty"`
}

// CheckSuite represents a suite of check runs
type CheckSuite struct {
	ID *int64 `json:"id,omitempty"`
}

func (c CheckRun) String() string {
	return Stringify(c)
}

// GetCheckRun gets a check-run for a repository
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/#get-a-single-check-run
func (s *ChecksService) GetCheckRun(ctx context.Context, owner string, repo string, id int64) (*CheckRun, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-runs/%v", owner, repo, id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	acceptHeaders := []string{mediaTypeCheckRunsPreview}
	req.Header.Set("Accept", strings.Join(acceptHeaders, ", "))

	checkRun := new(CheckRun)
	resp, err := s.client.Do(ctx, req, checkRun)
	if err != nil {
		return nil, resp, err
	}
	return checkRun, resp, nil
}
