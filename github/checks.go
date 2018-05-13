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
	ExternalID   *string         `json:"external_id,omitempty"`
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
	Title            *string               `json:"title,omitempty"`
	Summary          *string               `json:"summary,omitempty"`
	Text             *string               `json:"text,omitempty"`
	AnnotationsCount *int                  `json:"annotations_count,omitempty"`
	AnnotationsURL   *string               `json:"annotations_url,omitempty"`
	Annotations      []*CheckRunAnnotation `json:"annotations,omitempty"`
	Images           []*CheckRunImage      `json:"images,omitempty"`
}

// CheckRunAnnotation represents an annotation object for a CheckRun output.
type CheckRunAnnotation struct {
	FileName     *string `json:"filename,omitempty"`
	BlobHRef     *string `json:"blob_href,omitempty"`
	StartLine    *int    `json:"start_line,omitempty"`
	EndLine      *int    `json:"end_line,omitempty"`
	WarningLevel *string `json:"warning_level,omitempty"`
	Message      *string `json:"message,omitempty"`
	Title        *string `json:"title,omitempty"`
	RawDetails   *string `json:"raw_details,omitempty"`
}

// CheckRunImage represents an image object for a CheckRun output.
type CheckRunImage struct {
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

// CreateCheckRunOptions sets up parameters needed to create a CheckRun.
type CreateCheckRunOptions struct {
	Name        string          `json:"name"`                   // The name of the check (e.g., "code-coverage").(Required.)
	HeadBranch  string          `json:"head_branch"`            // The name of the branch to perform a check against.(Required.)
	HeadSHA     string          `json:"head_sha"`               // The SHA of the commit.(Required.)
	DetailsURL  *string         `json:"details_url,omitempty"`  // The URL of the integrator's site that has the full details of the check. (Optional.)
	ExternalID  *string         `json:"external_id,omitempty"`  // A reference for the run on the integrator's system. (Optional.)
	Status      *string         `json:"status,omitempty"`       // The current status. Can be one of queued, in_progress, or completed. Default: queued. (Optional.)
	Conclusion  *string         `json:"conclusion,omitempty"`   // Can be one of success, failure, neutral, cancelled, timed_out, or action_required.(Optional. Required if you provide a status of completed.)
	StartedAt   *Timestamp      `json:"started_at,omitempty"`   // The time that the check run began in ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ.(Optional.)
	CompletedAt *Timestamp      `json:"completed_at,omitempty"` // The time the check completed in ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ. (Optional. Required if you provide conclusion.)
	Output      *CheckRunOutput `json:"output,omitempty"`       // Provide descriptive details about the run.(Optional)
}

// CreateCheckRun Creates a check run for repository.
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/#create-a-check-run
func (s *ChecksService) CreateCheckRun(ctx context.Context, owner string, repo string, opt CreateCheckRunOptions) (*CheckRun, *Response, error) {
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

// UpdateCheckRunOptions sets up parameters needed to update a CheckRun.
type UpdateCheckRunOptions struct {
	Name        string          `json:"name"`                   // The name of the check (e.g., "code-coverage").(Required.)
	HeadBranch  *string         `json:"head_branch"`            // The name of the branch to perform a check against.(Optional.)
	HeadSHA     *string         `json:"head_sha"`               // The SHA of the commit.(Optional.)
	DetailsURL  *string         `json:"details_url,omitempty"`  // The URL of the integrator's site that has the full details of the check. (Optional.)
	ExternalID  *string         `json:"external_id,omitempty"`  // A reference for the run on the integrator's system. (Optional.)
	Status      *string         `json:"status,omitempty"`       // The current status. Can be one of queued, in_progress, or completed. Default: queued. (Optional.)
	Conclusion  *string         `json:"conclusion,omitempty"`   // Can be one of success, failure, neutral, cancelled, timed_out, or action_required.(Optional. Required if you provide a status of completed.)
	CompletedAt *Timestamp      `json:"completed_at,omitempty"` // The time the check completed in ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ. (Optional. Required if you provide conclusion.)
	Output      *CheckRunOutput `json:"output,omitempty"`       // Provide descriptive details about the run.(Optional)
}

// UpdateCheckRun Updates a check run for a specific commit in a repository.
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/#update-a-check-run
func (s *ChecksService) UpdateCheckRun(ctx context.Context, owner string, repo string, checkRunID int64, opt UpdateCheckRunOptions) (*CheckRun, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-runs/%v", owner, repo, checkRunID)
	req, err := s.client.NewRequest("PATCH", u, opt)
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

// ListCheckRunAnnotations List the annotations for a check run.
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/#list-annotations-for-a-check-run
func (s *ChecksService) ListCheckRunAnnotations(ctx context.Context, owner string, repo string, checkRunID int64, opt *ListOptions) ([]*CheckRunAnnotation, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/check-runs/%v/annotations", owner, repo, checkRunID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	var checkRunAnnotations []*CheckRunAnnotation
	resp, err := s.client.Do(ctx, req, &checkRunAnnotations)
	if err != nil {
		return nil, resp, err
	}

	return checkRunAnnotations, resp, nil
}

// ListCheckRunsOptions Optional parameters to list check runs.
type ListCheckRunsOptions struct {
	CheckName *string `url:"check_name,omitempty"` //Returns check runs with the specified name.
	Status    *string `url:"status,omitempty"`     //Returns check runs with the specified status. Can be one of queued, in_progress, or completed.
	Filter    *string `url:"filter,omitempty"`     //Filters check runs by their completed_at timestamp. Can be one of latest (returning the most recent check runs) or all. Default: latest

	ListOptions
}

// ListCheckRunsResults represents the result of a check run list.
type ListCheckRunsResults struct {
	Total     *int       `json:"total_count,omitempty"`
	CheckRuns []CheckRun `json:"check_runs,omitempty"`
}

// ListCheckRunsForRef List check runs for a specific ref.
//
// GitHub API docs: https://developer.github.com/v3/checks/runs/list-check-runs-for-a-specific-ref
func (s *ChecksService) ListCheckRunsForRef(ctx context.Context, owner string, repo string, ref string, opt *ListCheckRunsOptions) (*ListCheckRunsResults, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/commits/%v/check-runs", owner, repo, ref)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeCheckRunsPreview)

	var checkRunResults *ListCheckRunsResults
	resp, err := s.client.Do(ctx, req, &checkRunResults)
	if err != nil {
		return nil, resp, err
	}

	return checkRunResults, resp, nil
}
