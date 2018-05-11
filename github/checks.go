// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// ChecksService provides access to the Checks API in the
// GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/checks/
type ChecksService service

// CheckRun represents a GitHub check run on a repository associated with a GitHub app.
type CheckRun struct {
	ID           *int64          `json:"id,omitempty"`
	HeadSHA      *string         `json:"head_sha,omitempty"`
	ExternalID   *int64          `json:"external_id,omitempty"`
	URL          *string         `json:"url,omitempty"`
	HTMLURL      *string         `json:"html_url,omitempty"`
	Status       *string         `json:"status,omitempty"`
	Conclusion   *string         `json:"conclusion,omitempty"`
	StartedAt    *Timestamp      `json:"started_at,omitempty"`
	CompletedAt  *Timestamp      `json:"completed_at,omitempty"`
	Output       *CheckRunOutput `json:"output,omitempty"`
	Name         *string         `json:"name,omitempty"`
	CheckSuite   *CheckSuite     `json:"check_suite,omitempty"`
	App          *App            `json:"app,omitempty"`
	PullRequests []*PullRequest  `json:"pull_requests,omitempty"`
}

// CheckRunOutput represents the output of a CheckRun.
type CheckRunOutput struct {
	Title            *string            `json:"title,omitempty"`
	Summary          *string            `json:"summary,omitempty"`
	Text             *string            `json:"text,omitempty"`
	AnnotationsCount *int               `json:"annotations_count,omitempty"`
	AnnotationsURL   *string            `json:"annotations_url,omitempty"`
	Annotations      []*CheckAnnotation `json:"annotations,omitempty"`
	Images           []*CheckImage      `json:"images,omitempty"`
}

// CheckAnnotation represents an annotation object for a CheckRun output.
type CheckAnnotation struct {
	FileName     *string `json:"filename,omitempty"`
	BlobHRef     *string `json:"blob_href,omitempty"`
	StartLine    *int    `json:"start_line,omitempty"`
	EndLine      *int    `json:"end_line,omitempty"`
	WarningLevel *string `json:"warning_level,omitempty"`
	Message      *string `json:"message,omitempty"`
	Title        *string `json:"title,omitempty"`
	RawDetails   *string `json:"raw_details,omitempty"`
}

// CheckImage represents an image object for a CheckRun output.
type CheckImage struct {
	Alt      *string `json:"alt,omitempty"`
	ImageURL *string `json:"image_url,omitempty"`
	Caption  *string `json:"caption,omitempty"`
}

// CheckSuite represents a suite of check runs.
type CheckSuite struct {
	ID *int64 `json:"id,omitempty"`
}

func (c CheckRun) String() string {
	return Stringify(c)
}

// GetCheckRun gets a check-run for a repository.
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/#get-a-single-check-run
func (s *ChecksService) GetCheckRun(ctx context.Context, owner string, repo string, checkRunID int64) (*CheckRun, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-runs/%v", owner, repo, checkRunID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	checkRun := new(CheckRun)
	resp, err := s.client.Do(ctx, req, checkRun)
	if err != nil {
		return nil, resp, err
	}

	return checkRun, resp, nil
}

// CreateCheckRunOptions sets up parameters need to create a CheckRun.
type CreateCheckRunOptions struct {
	Name        *string         `json:"name,omitempty"`
	HeadBranch  *string         `json:"head_branch,omitempty"`
	HeadSHA     *string         `json:"head_sha,omitempty"`
	DetailsURL  *string         `json:"details_url,omitempty"`
	ExternalID  *int64          `json:"external_id,omitempty"`
	Status      *string         `json:"status,omitempty"`
	Conclusion  *string         `json:"conclusion,omitempty"`
	StartedAt   *Timestamp      `json:"started_at,omitempty"`
	CompletedAt *Timestamp      `json:"completed_at,omitempty"`
	Output      *CheckRunOutput `json:"output,omitempty"`
}

// CreateCheckRun Creates a check run for repository.
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/#create-a-check-run
func (s *ChecksService) CreateCheckRun(ctx context.Context, owner string, repo string, opt *CreateCheckRunOptions) (*CheckRun, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-runs", owner, repo)
	req, err := s.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	checkRun := new(CheckRun)
	resp, err := s.client.Do(ctx, req, checkRun)
	if err != nil {
		return nil, resp, err
	}

	return checkRun, resp, nil
}
