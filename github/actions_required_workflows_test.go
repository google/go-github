// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestActionsService_ListOrgRequiredWorkflows(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/required_workflows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"required_workflows": [
			{
			  "id": 30433642,
			  "name": "Required CI",
			  "path": ".github/workflows/ci.yml",
			  "scope": "selected",
			  "ref": "refs/head/main",
			  "state": "active",
			  "selected_repositories_url": "https://api.github.com/organizations/org/actions/required_workflows/1/repositories",
			  "created_at": "2020-01-22T19:33:08Z",
			  "updated_at": "2020-01-22T19:33:08Z"
			},
			{
			  "id": 30433643,
			  "name": "Required Linter",
			  "path": ".github/workflows/lint.yml",
			  "scope": "all",
			  "ref": "refs/head/main",
			  "state": "active",
			  "created_at": "2020-01-22T19:33:08Z",
			  "updated_at": "2020-01-22T19:33:08Z"
			}
		  ]
		}`)
	})
	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := context.Background()
	jobs, _, err := client.Actions.ListOrgRequiredWorkflows(ctx, "o", opts)

	if err != nil {
		t.Errorf("Actions.ListOrgRequiredWorkflows returned error: %v", err)
	}

	want := &OrgRequiredWorkflows{
		TotalCount: Ptr(4),
		RequiredWorkflows: []*OrgRequiredWorkflow{
			{ID: Ptr(int64(30433642)), Name: Ptr("Required CI"), Path: Ptr(".github/workflows/ci.yml"), Scope: Ptr("selected"), Ref: Ptr("refs/head/main"), State: Ptr("active"), SelectedRepositoriesURL: Ptr("https://api.github.com/organizations/org/actions/required_workflows/1/repositories"), CreatedAt: &Timestamp{time.Date(2020, time.January, 22, 19, 33, 8, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 22, 19, 33, 8, 0, time.UTC)}},
			{ID: Ptr(int64(30433643)), Name: Ptr("Required Linter"), Path: Ptr(".github/workflows/lint.yml"), Scope: Ptr("all"), Ref: Ptr("refs/head/main"), State: Ptr("active"), CreatedAt: &Timestamp{time.Date(2020, time.January, 22, 19, 33, 8, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 22, 19, 33, 8, 0, time.UTC)}},
		},
	}
	if !cmp.Equal(jobs, want) {
		t.Errorf("Actions.ListOrgRequiredWorkflows returned %+v, want %+v", jobs, want)
	}
	const methodName = "ListOrgRequiredWorkflows"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListOrgRequiredWorkflows(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListOrgRequiredWorkflows(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_CreateRequiredWorkflow(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/required_workflows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"workflow_file_path":".github/workflows/ci.yaml","repository_id":53,"scope":"selected","selected_repository_ids":[32,91]}`+"\n")
		fmt.Fprint(w, `{
			"id": 2,
			"name": "Required CI",
			"path": ".github/workflows/ci.yml",
			"scope": "selected",
			"ref": "refs/head/main",
			"state": "active",
			"selected_repositories_url": "https://api.github.com/orgs/octo-org/actions/required_workflows/2/repositories",
			"created_at": "2020-01-22T19:33:08Z",
			"updated_at": "2020-01-22T19:33:08Z",
			"repository": {
				"id": 53,
				"name": "Hello-World",
				"url": "https://api.github.com/repos/o/Hello-World"}}`)
	})
	input := &CreateUpdateRequiredWorkflowOptions{
		WorkflowFilePath:      Ptr(".github/workflows/ci.yaml"),
		RepositoryID:          Ptr(int64(53)),
		Scope:                 Ptr("selected"),
		SelectedRepositoryIDs: &SelectedRepoIDs{32, 91},
	}
	ctx := context.Background()
	requiredWokflow, _, err := client.Actions.CreateRequiredWorkflow(ctx, "o", input)
	if err != nil {
		t.Errorf("Actions.CreateRequiredWorkflow returned error: %v", err)
	}
	want := &OrgRequiredWorkflow{
		ID:                      Ptr(int64(2)),
		Name:                    Ptr("Required CI"),
		Path:                    Ptr(".github/workflows/ci.yml"),
		Scope:                   Ptr("selected"),
		Ref:                     Ptr("refs/head/main"),
		State:                   Ptr("active"),
		SelectedRepositoriesURL: Ptr("https://api.github.com/orgs/octo-org/actions/required_workflows/2/repositories"),
		CreatedAt:               &Timestamp{time.Date(2020, time.January, 22, 19, 33, 8, 0, time.UTC)},
		UpdatedAt:               &Timestamp{time.Date(2020, time.January, 22, 19, 33, 8, 0, time.UTC)},
		Repository:              &Repository{ID: Ptr(int64(53)), URL: Ptr("https://api.github.com/repos/o/Hello-World"), Name: Ptr("Hello-World")},
	}

	if !cmp.Equal(requiredWokflow, want) {
		t.Errorf("Actions.CreateRequiredWorkflow returned %+v, want %+v", requiredWokflow, want)
	}

	const methodName = "CreateRequiredWorkflow"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.CreateRequiredWorkflow(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.CreateRequiredWorkflow(ctx, "o", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetRequiredWorkflowByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/required_workflows/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 12345,
			"name": "Required CI",
			"path": ".github/workflows/ci.yml",
			"scope": "selected",
			"ref": "refs/head/main",
			"state": "active",
			"selected_repositories_url": "https://api.github.com/orgs/o/actions/required_workflows/12345/repositories",
			"created_at": "2020-01-22T19:33:08Z",
			"updated_at": "2020-01-22T19:33:08Z",
			"repository":{
				"id": 1296269,
				"url": "https://api.github.com/repos/o/Hello-World",
				"name": "Hello-World"
			}
			}`)
	})
	ctx := context.Background()
	jobs, _, err := client.Actions.GetRequiredWorkflowByID(ctx, "o", 12345)

	if err != nil {
		t.Errorf("Actions.GetRequiredWorkflowByID returned error: %v", err)
	}

	want := &OrgRequiredWorkflow{
		ID: Ptr(int64(12345)), Name: Ptr("Required CI"), Path: Ptr(".github/workflows/ci.yml"), Scope: Ptr("selected"), Ref: Ptr("refs/head/main"), State: Ptr("active"), SelectedRepositoriesURL: Ptr("https://api.github.com/orgs/o/actions/required_workflows/12345/repositories"), CreatedAt: &Timestamp{time.Date(2020, time.January, 22, 19, 33, 8, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 22, 19, 33, 8, 0, time.UTC)}, Repository: &Repository{ID: Ptr(int64(1296269)), URL: Ptr("https://api.github.com/repos/o/Hello-World"), Name: Ptr("Hello-World")},
	}
	if !cmp.Equal(jobs, want) {
		t.Errorf("Actions.GetRequiredWorkflowByID returned %+v, want %+v", jobs, want)
	}
	const methodName = "GetRequiredWorkflowByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetRequiredWorkflowByID(ctx, "\n", 1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetRequiredWorkflowByID(ctx, "o", 12345)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_UpdateRequiredWorkflow(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/required_workflows/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"workflow_file_path":".github/workflows/ci.yaml","repository_id":53,"scope":"selected","selected_repository_ids":[32,91]}`+"\n")
		fmt.Fprint(w, `{
			"id": 12345,
			"name": "Required CI",
			"path": ".github/workflows/ci.yml",
			"scope": "selected",
			"ref": "refs/head/main",
			"state": "active",
			"selected_repositories_url": "https://api.github.com/orgs/octo-org/actions/required_workflows/12345/repositories",
			"created_at": "2020-01-22T19:33:08Z",
			"updated_at": "2020-01-22T19:33:08Z",
			"repository": {
				"id": 53,
				"name": "Hello-World",
				"url": "https://api.github.com/repos/o/Hello-World"}}`)
	})
	input := &CreateUpdateRequiredWorkflowOptions{
		WorkflowFilePath:      Ptr(".github/workflows/ci.yaml"),
		RepositoryID:          Ptr(int64(53)),
		Scope:                 Ptr("selected"),
		SelectedRepositoryIDs: &SelectedRepoIDs{32, 91},
	}
	ctx := context.Background()

	requiredWokflow, _, err := client.Actions.UpdateRequiredWorkflow(ctx, "o", 12345, input)

	if err != nil {
		t.Errorf("Actions.UpdateRequiredWorkflow returned error: %v", err)
	}
	want := &OrgRequiredWorkflow{
		ID:                      Ptr(int64(12345)),
		Name:                    Ptr("Required CI"),
		Path:                    Ptr(".github/workflows/ci.yml"),
		Scope:                   Ptr("selected"),
		Ref:                     Ptr("refs/head/main"),
		State:                   Ptr("active"),
		SelectedRepositoriesURL: Ptr("https://api.github.com/orgs/octo-org/actions/required_workflows/12345/repositories"),
		CreatedAt:               &Timestamp{time.Date(2020, time.January, 22, 19, 33, 8, 0, time.UTC)},
		UpdatedAt:               &Timestamp{time.Date(2020, time.January, 22, 19, 33, 8, 0, time.UTC)},
		Repository:              &Repository{ID: Ptr(int64(53)), URL: Ptr("https://api.github.com/repos/o/Hello-World"), Name: Ptr("Hello-World")},
	}

	if !cmp.Equal(requiredWokflow, want) {
		t.Errorf("Actions.UpdateRequiredWorkflow returned %+v, want %+v", requiredWokflow, want)
	}

	const methodName = "UpdateRequiredWorkflow"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.UpdateRequiredWorkflow(ctx, "\n", 12345, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.UpdateRequiredWorkflow(ctx, "o", 12345, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_DeleteRequiredWorkflow(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/required_workflows/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := context.Background()
	_, err := client.Actions.DeleteRequiredWorkflow(ctx, "o", 12345)

	if err != nil {
		t.Errorf("Actions.DeleteRequiredWorkflow returned error: %v", err)
	}

	const methodName = "DeleteRequiredWorkflow"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteRequiredWorkflow(ctx, "\n", 12345)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteRequiredWorkflow(ctx, "o", 12345)
	})
}

func TestActionsService_ListRequiredWorkflowSelectedRepos(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/required_workflows/12345/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":1,
			"repositories": [{
				"id": 1296269,
				"url": "https://api.github.com/repos/o/Hello-World",
				"name": "Hello-World"
				}]
		}`)
	})
	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := context.Background()
	jobs, _, err := client.Actions.ListRequiredWorkflowSelectedRepos(ctx, "o", 12345, opts)

	if err != nil {
		t.Errorf("Actions.ListRequiredWorkflowSelectedRepositories returned error: %v", err)
	}

	want := &RequiredWorkflowSelectedRepos{
		TotalCount: Ptr(1),
		Repositories: []*Repository{
			{ID: Ptr(int64(1296269)), URL: Ptr("https://api.github.com/repos/o/Hello-World"), Name: Ptr("Hello-World")},
		},
	}
	if !cmp.Equal(jobs, want) {
		t.Errorf("Actions.ListRequiredWorkflowSelectedRepositories returned %+v, want %+v", jobs, want)
	}
	const methodName = "ListRequiredWorkflowSelectedRepositories"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListRequiredWorkflowSelectedRepos(ctx, "\n", 12345, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListRequiredWorkflowSelectedRepos(ctx, "o", 12345, opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_SetRequiredWorkflowSelectedRepos(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/required_workflows/12345/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"selected_repository_ids":[32,91]}`+"\n")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := context.Background()
	_, err := client.Actions.SetRequiredWorkflowSelectedRepos(ctx, "o", 12345, SelectedRepoIDs{32, 91})

	if err != nil {
		t.Errorf("Actions.SetRequiredWorkflowSelectedRepositories returned error: %v", err)
	}

	const methodName = "SetRequiredWorkflowSelectedRepositories"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.SetRequiredWorkflowSelectedRepos(ctx, "\n", 12345, SelectedRepoIDs{32, 91})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.SetRequiredWorkflowSelectedRepos(ctx, "o", 12345, SelectedRepoIDs{32, 91})
	})
}

func TestActionsService_AddRepoToRequiredWorkflow(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/required_workflows/12345/repositories/32", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := context.Background()
	_, err := client.Actions.AddRepoToRequiredWorkflow(ctx, "o", 12345, 32)

	if err != nil {
		t.Errorf("Actions.AddRepoToRequiredWorkflow returned error: %v", err)
	}

	const methodName = "AddRepoToRequiredWorkflow"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.AddRepoToRequiredWorkflow(ctx, "\n", 12345, 32)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.AddRepoToRequiredWorkflow(ctx, "o", 12345, 32)
	})
}

func TestActionsService_RemoveRepoFromRequiredWorkflow(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/required_workflows/12345/repositories/32", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := context.Background()
	_, err := client.Actions.RemoveRepoFromRequiredWorkflow(ctx, "o", 12345, 32)

	if err != nil {
		t.Errorf("Actions.RemoveRepoFromRequiredWorkflow returned error: %v", err)
	}

	const methodName = "RemoveRepoFromRequiredWorkflow"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RemoveRepoFromRequiredWorkflow(ctx, "\n", 12345, 32)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.RemoveRepoFromRequiredWorkflow(ctx, "o", 12345, 32)
	})
}

func TestActionsService_ListRepoRequiredWorkflows(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/required_workflows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":1,"required_workflows": [
			{
			  "id": 30433642,
			  "node_id": "MDg6V29ya2Zsb3cxNjEzMzU=",
			  "name": "Required CI",
			  "path": ".github/workflows/ci.yml",
			  "state": "active",
			  "created_at": "2020-01-22T19:33:08Z",
			  "updated_at": "2020-01-22T19:33:08Z",
			  "url": "https://api.github.com/repos/o/r/actions/required_workflows/161335",
			  "html_url": "https://github.com/o/r/blob/master/o/hello-world/.github/workflows/required_ci.yaml",
			  "badge_url": "https://github.com/o/r/workflows/required/o/hello-world/.github/workflows/required_ci.yaml/badge.svg",
			  "source_repository":{
				"id": 1296269,
				"url": "https://api.github.com/repos/o/Hello-World",
				"name": "Hello-World"
			  }
			}
		  ]
		}`)
	})
	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := context.Background()
	jobs, _, err := client.Actions.ListRepoRequiredWorkflows(ctx, "o", "r", opts)

	if err != nil {
		t.Errorf("Actions.ListRepoRequiredWorkflows returned error: %v", err)
	}

	want := &RepoRequiredWorkflows{
		TotalCount: Ptr(1),
		RequiredWorkflows: []*RepoRequiredWorkflow{
			{ID: Ptr(int64(30433642)), NodeID: Ptr("MDg6V29ya2Zsb3cxNjEzMzU="), Name: Ptr("Required CI"), Path: Ptr(".github/workflows/ci.yml"), State: Ptr("active"), CreatedAt: &Timestamp{time.Date(2020, time.January, 22, 19, 33, 8, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 22, 19, 33, 8, 0, time.UTC)}, URL: Ptr("https://api.github.com/repos/o/r/actions/required_workflows/161335"), BadgeURL: Ptr("https://github.com/o/r/workflows/required/o/hello-world/.github/workflows/required_ci.yaml/badge.svg"), HTMLURL: Ptr("https://github.com/o/r/blob/master/o/hello-world/.github/workflows/required_ci.yaml"), SourceRepository: &Repository{ID: Ptr(int64(1296269)), URL: Ptr("https://api.github.com/repos/o/Hello-World"), Name: Ptr("Hello-World")}},
		},
	}
	if !cmp.Equal(jobs, want) {
		t.Errorf("Actions.ListRepoRequiredWorkflows returned %+v, want %+v", jobs, want)
	}
	const methodName = "ListRepoRequiredWorkflows"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListRepoRequiredWorkflows(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListRepoRequiredWorkflows(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
