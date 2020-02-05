// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"time"
)

// Workflow represents a repository action workflow.
type Workflow struct {
	ID        int64     `json:"id"`
	NodeID    string    `json:"node_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
	HTMLURL   string    `json:"html_url"`
	BadgeURL  string    `json:"badge_url"`
}

// Workflows represents one item from the ListWorkflows response.
type Workflows struct {
	TotalCount int         `json:"total_count"`
	Workflows  []*Workflow `json:"workflows"`
}

// ListWorkflows lists all workflows in a repository.
//
// GitHub API docs: https://developer.github.com/v3/actions/workflows/#list-repository-workflows
func (s *ActionsService) ListWorkflows(ctx context.Context, owner, repo string, opt *ListOptions) (*Workflows, *Response, error) {
	u := fmt.Sprintf("repos/%s/%s/actions/workflows", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	workflows := new(Workflows)
	resp, err := s.client.Do(ctx, req, &workflows)
	if err != nil {
		return nil, resp, err
	}

	return workflows, resp, nil
}

// GetWorkflow gets a specific workflow.
//
// GitHub API docs: https://developer.github.com/v3/actions/workflows/#list-repository-workflows
func (s *ActionsService) GetWorkflow(ctx context.Context, owner, repo string, workflowID int64) (*Workflow, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/workflows/%v", owner, repo, workflowID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	workflow := new(Workflow)
	resp, err := s.client.Do(ctx, req, workflow)
	if err != nil {
		return nil, resp, err
	}

	return workflow, resp, nil
}
