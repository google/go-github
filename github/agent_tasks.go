// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// AgentTasksService handles communication with the agent tasks
// methods of the GitHub API.
//
// GitHub API docs: https://docs.github.com/rest/agent-tasks/agent-tasks?apiVersion=2026-03-10
type AgentTasksService service

// AgentTask represents a Copilot cloud agent task.
type AgentTask struct {
	ID           string               `json:"id"`
	URL          *string              `json:"url,omitempty"`
	HTMLURL      *string              `json:"html_url,omitempty"`
	Name         *string              `json:"name,omitempty"`
	Creator      *User                `json:"creator,omitempty"`
	CreatorType  *string              `json:"creator_type,omitempty"`
	Owner        *AgentTaskOwner      `json:"owner,omitempty"`
	Repository   *AgentTaskRepository `json:"repository,omitempty"`
	State        string               `json:"state"`
	SessionCount *int                 `json:"session_count,omitempty"`
	Artifacts    []*AgentTaskArtifact `json:"artifacts,omitempty"`
	ArchivedAt   *Timestamp           `json:"archived_at,omitempty"`
	CreatedAt    Timestamp            `json:"created_at"`
	UpdatedAt    *Timestamp           `json:"updated_at,omitempty"`
	Sessions     []*AgentTaskSession  `json:"sessions,omitempty"`
}

// AgentTaskOwner represents an agent task owner.
type AgentTaskOwner struct {
	ID *int64 `json:"id,omitempty"`
}

// AgentTaskRepository represents an agent task repository.
type AgentTaskRepository struct {
	ID *int64 `json:"id,omitempty"`
}

// AgentTaskArtifact represents an artifact produced by an agent task.
type AgentTaskArtifact struct {
	Provider string          `json:"provider"`
	Type     string          `json:"type"`
	Data     json.RawMessage `json:"data"`
}

// AgentTaskSession represents a session associated with an agent task.
type AgentTaskSession struct {
	ID          string               `json:"id"`
	Name        *string              `json:"name,omitempty"`
	User        *User                `json:"user,omitempty"`
	Owner       *AgentTaskOwner      `json:"owner,omitempty"`
	Repository  *AgentTaskRepository `json:"repository,omitempty"`
	TaskID      *string              `json:"task_id,omitempty"`
	State       string               `json:"state"`
	CreatedAt   Timestamp            `json:"created_at"`
	UpdatedAt   *Timestamp           `json:"updated_at,omitempty"`
	CompletedAt *Timestamp           `json:"completed_at,omitempty"`
	Prompt      *string              `json:"prompt,omitempty"`
	HeadRef     *string              `json:"head_ref,omitempty"`
	BaseRef     *string              `json:"base_ref,omitempty"`
	Model       *string              `json:"model,omitempty"`
}

// AgentTaskList represents a list of agent tasks.
type AgentTaskList struct {
	Tasks              []*AgentTask `json:"tasks"`
	TotalActiveCount   *int         `json:"total_active_count,omitempty"`
	TotalArchivedCount *int         `json:"total_archived_count,omitempty"`
}

// AgentTaskListOptions specifies optional parameters to AgentTasksService.List.
type AgentTaskListOptions struct {
	// Sort specifies the field to sort results by. Possible values are: updated_at, created_at.
	Sort string `url:"sort,omitempty"`

	// Direction specifies the direction to sort results by. Possible values are: asc, desc.
	Direction string `url:"direction,omitempty"`

	// State is a comma-separated list of task states to filter by.
	State string `url:"state,omitempty"`

	// IsArchived filters tasks by archived status.
	IsArchived bool `url:"is_archived,omitempty"`

	// Since filters tasks updated at or after this time.
	Since *time.Time `url:"since,omitempty"`

	ListOptions
}

// AgentTaskListByRepoOptions specifies optional parameters to AgentTasksService.ListByRepo.
type AgentTaskListByRepoOptions struct {
	AgentTaskListOptions

	// CreatorIDs filters tasks by creator user IDs.
	CreatorIDs []int64 `url:"creator_id,omitempty"`
}

// CreateAgentTaskRequest represents the parameters for creating an agent task.
type CreateAgentTaskRequest struct {
	// Prompt is the user's prompt for the agent.
	Prompt string `json:"prompt"`

	// Model is the model to use for this task.
	Model *string `json:"model,omitempty"`

	// CreatePullRequest indicates whether to create a pull request.
	CreatePullRequest *bool `json:"create_pull_request,omitempty"`

	// HeadRef is the head ref for the new branch or pull request.
	HeadRef *string `json:"head_ref,omitempty"`

	// BaseRef is the base ref for the new branch or pull request.
	BaseRef *string `json:"base_ref,omitempty"`
}

// ListByRepo lists tasks for a repository.
//
// Note: This endpoint is in public preview and is subject to change.
//
// GitHub API docs: https://docs.github.com/rest/agent-tasks/agent-tasks?apiVersion=2026-03-10#list-tasks-for-repository
//
//meta:operation GET /agents/repos/{owner}/{repo}/tasks
func (s *AgentTasksService) ListByRepo(ctx context.Context, owner, repo string, opts *AgentTaskListByRepoOptions) (*AgentTaskList, *Response, error) {
	u := fmt.Sprintf("agents/repos/%v/%v/tasks", owner, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil, WithVersion(api20260310))
	if err != nil {
		return nil, nil, err
	}

	var tasks *AgentTaskList
	resp, err := s.client.Do(req, &tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, nil
}

// Create starts a new Copilot cloud agent task for a repository.
//
// Note: This endpoint is in public preview and is subject to change.
//
// GitHub API docs: https://docs.github.com/rest/agent-tasks/agent-tasks?apiVersion=2026-03-10#start-a-task
//
//meta:operation POST /agents/repos/{owner}/{repo}/tasks
func (s *AgentTasksService) Create(ctx context.Context, owner, repo string, req CreateAgentTaskRequest) (*AgentTask, *Response, error) {
	u := fmt.Sprintf("agents/repos/%v/%v/tasks", owner, repo)

	req, err := s.client.NewRequest(ctx, "POST", u, opts, WithVersion(api20260310))
	if err != nil {
		return nil, nil, err
	}

	var task *AgentTask
	resp, err := s.client.Do(req, &task)
	if err != nil {
		return nil, resp, err
	}

	return task, resp, nil
}

// GetByRepoAndID gets a repository task by ID.
//
// Note: This endpoint is in public preview and is subject to change.
//
// GitHub API docs: https://docs.github.com/rest/agent-tasks/agent-tasks?apiVersion=2026-03-10#get-a-task-by-repo
//
//meta:operation GET /agents/repos/{owner}/{repo}/tasks/{task_id}
func (s *AgentTasksService) GetByRepoAndID(ctx context.Context, owner, repo, taskID string) (*AgentTask, *Response, error) {
	u := fmt.Sprintf("agents/repos/%v/%v/tasks/%v", owner, repo, taskID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil, WithVersion(api20260310))
	if err != nil {
		return nil, nil, err
	}

	var task *AgentTask
	resp, err := s.client.Do(req, &task)
	if err != nil {
		return nil, resp, err
	}

	return task, resp, nil
}

// List lists tasks for the authenticated user.
//
// Note: This endpoint is in public preview and is subject to change.
//
// GitHub API docs: https://docs.github.com/rest/agent-tasks/agent-tasks?apiVersion=2026-03-10#list-tasks
//
//meta:operation GET /agents/tasks
func (s *AgentTasksService) List(ctx context.Context, opts *AgentTaskListOptions) (*AgentTaskList, *Response, error) {
	u := "agents/tasks"
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil, WithVersion(api20260310))
	if err != nil {
		return nil, nil, err
	}

	var tasks *AgentTaskList
	resp, err := s.client.Do(req, &tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, nil
}

// Get gets a task by ID for the authenticated user.
//
// Note: This endpoint is in public preview and is subject to change.
//
// GitHub API docs: https://docs.github.com/rest/agent-tasks/agent-tasks?apiVersion=2026-03-10#get-a-task-by-id
//
//meta:operation GET /agents/tasks/{task_id}
func (s *AgentTasksService) Get(ctx context.Context, taskID string) (*AgentTask, *Response, error) {
	u := fmt.Sprintf("agents/tasks/%v", taskID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil, WithVersion(api20260310))
	if err != nil {
		return nil, nil, err
	}

	var task *AgentTask
	resp, err := s.client.Do(req, &task)
	if err != nil {
		return nil, resp, err
	}

	return task, resp, nil
}
