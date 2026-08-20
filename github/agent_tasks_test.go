// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func agentTaskJSON() string {
	return `{
		"id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		"url": "https://api.github.com/agents/repos/octocat/hello-world/tasks/a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		"html_url": "https://github.com/octocat/hello-world/copilot/tasks/a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		"name": "Fix the login button on the homepage",
		"creator": { "id": 1 },
		"creator_type": "user",
		"user_collaborators": [{ "id": 3 }],
		"owner": { "id": 2 },
		"repository": { "id": 1296269 },
		"state": "completed",
		"session_count": 1,
		"artifacts": [
			{
				"provider": "github",
				"type": "pull",
				"data": { "id": 42 }
			}
		],
		"archived_at": null,
		"created_at": "2025-03-14T00:00:00Z",
		"updated_at": "2025-03-14T01:00:00Z"
	}`
}

func agentTaskWithSessionsJSON() string {
	return `{
		"id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		"url": "https://api.github.com/agents/repos/octocat/hello-world/tasks/a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		"html_url": "https://github.com/octocat/hello-world/copilot/tasks/a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		"name": "Fix the login button on the homepage",
		"creator": { "id": 1 },
		"creator_type": "user",
		"user_collaborators": [{ "id": 3 }],
		"owner": { "id": 2 },
		"repository": { "id": 1296269 },
		"state": "failed",
		"session_count": 1,
		"artifacts": [
			{
				"provider": "github",
				"type": "pull",
				"data": { "id": 42 }
			}
		],
		"archived_at": null,
		"created_at": "2025-03-14T00:00:00Z",
		"updated_at": "2025-03-14T01:00:00Z",
		"sessions": [
			{
				"id": "s1a2b3c4-d5e6-7890-abcd-ef1234567890",
				"name": "Fix the login button on the homepage",
				"user": { "id": 1 },
				"owner": { "id": 2 },
				"repository": { "id": 1296269 },
				"task_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
				"state": "failed",
				"created_at": "2025-03-14T00:00:00Z",
				"updated_at": "2025-03-14T01:00:00Z",
				"completed_at": "2025-03-14T01:00:00Z",
				"prompt": "Fix the login button on the homepage",
				"head_ref": "copilot/fix-1",
				"base_ref": "main",
				"model": "claude-sonnet-4.6",
				"error": { "message": "agent task failed" }
			}
		]
	}`
}

func agentTask() *AgentTask {
	return &AgentTask{
		ID:          "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		URL:         Ptr("https://api.github.com/agents/repos/octocat/hello-world/tasks/a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
		HTMLURL:     Ptr("https://github.com/octocat/hello-world/copilot/tasks/a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
		Name:        Ptr("Fix the login button on the homepage"),
		Creator:     &AgentTaskCreator{ID: Ptr(int64(1))},
		CreatorType: Ptr("user"),
		UserCollaborators: []*User{
			{ID: Ptr(int64(3))},
		},
		Owner:        &AgentTaskOwner{ID: Ptr(int64(2))},
		Repository:   &AgentTaskRepository{ID: Ptr(int64(1296269))},
		State:        "completed",
		SessionCount: Ptr(1),
		Artifacts: []*AgentTaskArtifact{
			{
				Provider: "github",
				Type:     "pull",
				Data:     &AgentTaskArtifactData{ID: Ptr(int64(42))},
			},
		},
		CreatedAt: Timestamp{time.Date(2025, time.March, 14, 0, 0, 0, 0, time.UTC)},
		UpdatedAt: &Timestamp{time.Date(2025, time.March, 14, 1, 0, 0, 0, time.UTC)},
	}
}

func agentTaskWithSessions() *AgentTask {
	task := agentTask()
	task.State = "failed"

	task.Sessions = []*AgentTaskSession{
		{
			ID:          "s1a2b3c4-d5e6-7890-abcd-ef1234567890",
			Name:        Ptr("Fix the login button on the homepage"),
			User:        &User{ID: Ptr(int64(1))},
			Owner:       &AgentTaskOwner{ID: Ptr(int64(2))},
			Repository:  &AgentTaskRepository{ID: Ptr(int64(1296269))},
			TaskID:      Ptr("a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
			State:       "failed",
			Prompt:      Ptr("Fix the login button on the homepage"),
			HeadRef:     Ptr("copilot/fix-1"),
			BaseRef:     Ptr("main"),
			Model:       Ptr("claude-sonnet-4.6"),
			Error:       &AgentTaskSessionError{Message: Ptr("agent task failed")},
			CreatedAt:   Timestamp{time.Date(2025, time.March, 14, 0, 0, 0, 0, time.UTC)},
			UpdatedAt:   &Timestamp{time.Date(2025, time.March, 14, 1, 0, 0, 0, time.UTC)},
			CompletedAt: &Timestamp{time.Date(2025, time.March, 14, 1, 0, 0, 0, time.UTC)},
		},
	}

	return task
}

func TestAgentTasksService_ListByRepo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/agents/repos/o/r/tasks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValuesList(t, r, url.Values{
			"creator_id":  {"1", "2"},
			"direction":   {"asc"},
			"is_archived": {"true"},
			"page":        {"2"},
			"per_page":    {"1"},
			"since":       {"2025-03-14T00:00:00Z"},
			"sort":        {"created_at"},
			"state":       {"queued,completed"},
		})
		w.Header().Set("Link", `<https://api.github.com/agents/repos/o/r/tasks?page=3>; rel="next"`)
		fmt.Fprintf(w, `{"tasks":[%v],"total_active_count":5,"total_archived_count":2}`, agentTaskJSON())
	})

	opts := &AgentTaskListByRepoOptions{
		AgentTaskListOptions: AgentTaskListOptions{
			Sort:        "created_at",
			Direction:   "asc",
			State:       "queued,completed",
			IsArchived:  true,
			Since:       Ptr(time.Date(2025, time.March, 14, 0, 0, 0, 0, time.UTC)),
			ListOptions: ListOptions{Page: 2, PerPage: 1},
		},
		CreatorID: []int64{1, 2},
	}

	ctx := t.Context()
	tasks, resp, err := client.AgentTasks.ListByRepo(ctx, "o", "r", opts)
	if err != nil {
		t.Fatalf("AgentTasks.ListByRepo returned error: %v", err)
	}

	want := &AgentTaskList{
		Tasks:              []*AgentTask{agentTask()},
		TotalActiveCount:   Ptr(5),
		TotalArchivedCount: Ptr(2),
	}
	if diff := cmp.Diff(want, tasks); diff != "" {
		t.Errorf("AgentTasks.ListByRepo mismatch (-want +got):\n%v", diff)
	}
	if got, want := resp.NextPage, 3; got != want {
		t.Errorf("AgentTasks.ListByRepo NextPage = %v, want %v", got, want)
	}

	const methodName = "ListByRepo"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.AgentTasks.ListByRepo(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.AgentTasks.ListByRepo(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentTasksService_Create(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := CreateAgentTaskRequest{
		Prompt:            "Fix the login button on the homepage",
		Model:             Ptr("gpt-5.3-codex"),
		CreatePullRequest: Ptr(true),
		HeadRef:           Ptr("copilot/fix-1"),
		BaseRef:           Ptr("main"),
	}

	mux.HandleFunc("/agents/repos/o/r/tasks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")
		testJSONBody(t, r, input)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, agentTaskJSON())
	})

	ctx := t.Context()
	task, _, err := client.AgentTasks.Create(ctx, "o", "r", input)
	if err != nil {
		t.Fatalf("AgentTasks.Create returned error: %v", err)
	}
	if diff := cmp.Diff(agentTask(), task); diff != "" {
		t.Errorf("AgentTasks.Create mismatch (-want +got):\n%v", diff)
	}

	const methodName = "Create"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.AgentTasks.Create(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.AgentTasks.Create(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentTasksService_GetByRepoAndID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	const taskID = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"

	mux.HandleFunc("/agents/repos/o/r/tasks/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, agentTaskWithSessionsJSON())
	})

	ctx := t.Context()
	task, _, err := client.AgentTasks.GetByRepoAndID(ctx, "o", "r", taskID)
	if err != nil {
		t.Fatalf("AgentTasks.GetByRepoAndID returned error: %v", err)
	}
	if diff := cmp.Diff(agentTaskWithSessions(), task); diff != "" {
		t.Errorf("AgentTasks.GetByRepoAndID mismatch (-want +got):\n%v", diff)
	}

	const methodName = "GetByRepoAndID"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.AgentTasks.GetByRepoAndID(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.AgentTasks.GetByRepoAndID(ctx, "o", "r", taskID)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentTasksService_List(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/agents/tasks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"direction":   "desc",
			"is_archived": "true",
			"page":        "2",
			"per_page":    "1",
			"since":       "2025-03-14T00:00:00Z",
			"sort":        "updated_at",
			"state":       "completed",
		})
		w.Header().Set("Link", `<https://api.github.com/agents/tasks?page=3>; rel="next"`)
		fmt.Fprintf(w, `{"tasks":[%v],"total_active_count":5,"total_archived_count":2}`, agentTaskJSON())
	})

	opts := &AgentTaskListOptions{
		Sort:        "updated_at",
		Direction:   "desc",
		State:       "completed",
		IsArchived:  true,
		Since:       Ptr(time.Date(2025, time.March, 14, 0, 0, 0, 0, time.UTC)),
		ListOptions: ListOptions{Page: 2, PerPage: 1},
	}

	ctx := t.Context()
	tasks, resp, err := client.AgentTasks.List(ctx, opts)
	if err != nil {
		t.Fatalf("AgentTasks.List returned error: %v", err)
	}

	want := &AgentTaskList{
		Tasks:              []*AgentTask{agentTask()},
		TotalActiveCount:   Ptr(5),
		TotalArchivedCount: Ptr(2),
	}
	if diff := cmp.Diff(want, tasks); diff != "" {
		t.Errorf("AgentTasks.List mismatch (-want +got):\n%v", diff)
	}
	if got, want := resp.NextPage, 3; got != want {
		t.Errorf("AgentTasks.List NextPage = %v, want %v", got, want)
	}

	const methodName = "List"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.AgentTasks.List(ctx, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentTasksService_Get(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	const taskID = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"

	mux.HandleFunc("/agents/tasks/"+taskID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, agentTaskWithSessionsJSON())
	})

	ctx := t.Context()
	task, _, err := client.AgentTasks.Get(ctx, taskID)
	if err != nil {
		t.Fatalf("AgentTasks.Get returned error: %v", err)
	}
	if diff := cmp.Diff(agentTaskWithSessions(), task); diff != "" {
		t.Errorf("AgentTasks.Get mismatch (-want +got):\n%v", diff)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() error {
		_, _, err := client.AgentTasks.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.AgentTasks.Get(ctx, taskID)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
