// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// Step represents a single task from a sequence of tasks of a job.
type Step struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	Number      int64     `json:"number"`
	StartedAt   Timestamp `json:"started_at"`
	CompletedAt Timestamp `json:"completed_at"`
}

// Job represents a repository action workflow job.
type Job struct {
	ID          int64     `json:"id"`
	RunID       int64     `json:"run_id"`
	RunURL      string    `json:"run_url"`
	NodeID      string    `json:"node_id"`
	HeadSHA     string    `json:"head_sha"`
	URL         string    `json:"url"`
	HTMLURL     string    `json:"html_url"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	StartedAt   Timestamp `json:"started_at"`
	CompletedAt Timestamp `json:"completed_at"`
	Name        string    `json:"name"`
	Steps       []*Step   `json:"steps"`
	CheckRunURL string    `json:"check_run_url"`
}

// Jobs represents a slice of repository action workflow job.
type Jobs struct {
	TotalCount int    `json:"total_count"`
	Jobs       []*Job `json:"jobs"`
}

// ListWorkflowJobs lists all jobs for a workflow run.
//
// GitHub API docs: https://developer.github.com/v3/actions/workflow_jobs/#list-jobs-for-a-workflow-run
func (s *ActionsService) ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, opts *ListOptions) (*Jobs, *Response, error) {
	u := fmt.Sprintf("repos/%s/%s/actions/runs/%v/jobs", owner, repo, runID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	jobs := new(Jobs)
	resp, err := s.client.Do(ctx, req, &jobs)
	if err != nil {
		return nil, resp, err
	}

	return jobs, resp, nil
}

// GetWorkflowJobByID gets a specific job in a workflow run by ID.
//
// GitHub API docs: https://developer.github.com/v3/actions/workflow_jobs/#list-jobs-for-a-workflow-run
func (s *ActionsService) GetWorkflowJobByID(ctx context.Context, owner, repo string, jobID int64) (*Job, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/jobs/%v", owner, repo, jobID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	job := new(Job)
	resp, err := s.client.Do(ctx, req, job)
	if err != nil {
		return nil, resp, err
	}

	return job, resp, nil
}

// ListWorkflowJobLogs gets a redirect URL to a download a plain text file of logs for a workflow job.
//
// GitHub API docs: https://developer.github.com/v3/actions/workflow_jobs/#list-workflow-job-logs
func (s *ActionsService) ListWorkflowJobLogs(ctx context.Context, owner, repo string, jobID int64) (string, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/jobs/%v/logs", owner, repo, jobID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return "", nil, err
	}

	var logFileURL string
	resp, err := s.client.Do(ctx, req, &logFileURL)

	if err != nil {
		return "", resp, err
	}

	return logFileURL, resp, nil
}
