// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAgentsService_ListRepoVariables(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/agents/variables", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"variables":[{"name":"A","value":"AA","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`},{"name":"B","value":"BB","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	variables, _, err := client.Agents.ListRepoVariables(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Agents.ListRepoVariables returned error: %v", err)
	}

	want := &ActionsVariables{
		TotalCount: 4,
		Variables: []*ActionsVariable{
			{Name: "A", Value: "AA", CreatedAt: &referenceTimestamp, UpdatedAt: &referenceTimestamp},
			{Name: "B", Value: "BB", CreatedAt: &referenceTimestamp, UpdatedAt: &referenceTimestamp},
		},
	}
	if !cmp.Equal(variables, want) {
		t.Errorf("Agents.ListRepoVariables returned %+v, want %+v", variables, want)
	}

	const methodName = "ListRepoVariables"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.ListRepoVariables(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.ListRepoVariables(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_ListRepoOrgVariables(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/agents/organization-variables", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"variables":[{"name":"A","value":"AA","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`},{"name":"B","value":"BB","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	variables, _, err := client.Agents.ListRepoOrgVariables(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Agents.ListRepoOrgVariables returned error: %v", err)
	}

	want := &ActionsVariables{
		TotalCount: 4,
		Variables: []*ActionsVariable{
			{Name: "A", Value: "AA", CreatedAt: &referenceTimestamp, UpdatedAt: &referenceTimestamp},
			{Name: "B", Value: "BB", CreatedAt: &referenceTimestamp, UpdatedAt: &referenceTimestamp},
		},
	}
	if !cmp.Equal(variables, want) {
		t.Errorf("Agents.ListRepoOrgVariables returned %+v, want %+v", variables, want)
	}

	const methodName = "ListRepoOrgVariables"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.ListRepoOrgVariables(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.ListRepoOrgVariables(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_GetRepoVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/agents/variables/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		fmt.Fprint(w, `{"name":"NAME","value":"VALUE","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`}`)
	})

	ctx := t.Context()
	variable, _, err := client.Agents.GetRepoVariable(ctx, "o", "r", "NAME")
	if err != nil {
		t.Errorf("Agents.GetRepoVariable returned error: %v", err)
	}

	want := &ActionsVariable{
		Name:      "NAME",
		Value:     "VALUE",
		CreatedAt: &referenceTimestamp,
		UpdatedAt: &referenceTimestamp,
	}
	if !cmp.Equal(variable, want) {
		t.Errorf("Agents.GetRepoVariable returned %+v, want %+v", variable, want)
	}

	const methodName = "GetRepoVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.GetRepoVariable(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.GetRepoVariable(ctx, "o", "r", "NAME")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_CreateRepoVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := ActionsCreateVariableRequest{
		Name:  "NAME",
		Value: "VALUE",
	}

	mux.HandleFunc("/repos/o/r/agents/variables", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testHeader(t, r, "Content-Type", "application/json")
		testJSONBody(t, r, input)
		w.WriteHeader(http.StatusCreated)
	})

	ctx := t.Context()
	_, err := client.Agents.CreateRepoVariable(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Agents.CreateRepoVariable returned error: %v", err)
	}

	const methodName = "CreateRepoVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.CreateRepoVariable(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.CreateRepoVariable(ctx, "o", "r", input)
	})
}

func TestAgentsService_UpdateRepoVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	name := "NAME"
	input := ActionsUpdateVariableRequest{
		Name:  &name,
		Value: new("VALUE"),
	}

	mux.HandleFunc("/repos/o/r/agents/variables/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testHeader(t, r, "Content-Type", "application/json")
		testJSONBody(t, r, input)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	_, err := client.Agents.UpdateRepoVariable(ctx, "o", "r", name, input)
	if err != nil {
		t.Errorf("Agents.UpdateRepoVariable returned error: %v", err)
	}

	const methodName = "UpdateRepoVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.UpdateRepoVariable(ctx, "\n", "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.UpdateRepoVariable(ctx, "o", "r", name, input)
	})
}

func TestAgentsService_DeleteRepoVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/agents/variables/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	_, err := client.Agents.DeleteRepoVariable(ctx, "o", "r", "NAME")
	if err != nil {
		t.Errorf("Agents.DeleteRepoVariable returned error: %v", err)
	}

	const methodName = "DeleteRepoVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.DeleteRepoVariable(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.DeleteRepoVariable(ctx, "o", "r", "NAME")
	})
}

func TestAgentsService_ListOrgVariables(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/variables", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":3,"variables":[{"name":"A","value":"AA","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`,"visibility":"private"},{"name":"B","value":"BB","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`,"visibility":"all"},{"name":"C","value":"CC","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`,"visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/agents/variables/VAR/repositories"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	variables, _, err := client.Agents.ListOrgVariables(ctx, "o", opts)
	if err != nil {
		t.Errorf("Agents.ListOrgVariables returned error: %v", err)
	}

	want := &ActionsVariables{
		TotalCount: 3,
		Variables: []*ActionsVariable{
			{Name: "A", Value: "AA", CreatedAt: &referenceTimestamp, UpdatedAt: &referenceTimestamp, Visibility: new("private")},
			{Name: "B", Value: "BB", CreatedAt: &referenceTimestamp, UpdatedAt: &referenceTimestamp, Visibility: new("all")},
			{Name: "C", Value: "CC", CreatedAt: &referenceTimestamp, UpdatedAt: &referenceTimestamp, Visibility: new("selected"), SelectedRepositoriesURL: new("https://api.github.com/orgs/octo-org/agents/variables/VAR/repositories")},
		},
	}
	if !cmp.Equal(variables, want) {
		t.Errorf("Agents.ListOrgVariables returned %+v, want %+v", variables, want)
	}

	const methodName = "ListOrgVariables"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.ListOrgVariables(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.ListOrgVariables(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_GetOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/variables/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		fmt.Fprint(w, `{"name":"NAME","value":"VALUE","created_at":`+referenceTimeStr+`,"updated_at":`+referenceTimeStr+`,"visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/agents/variables/VAR/repositories"}`)
	})

	ctx := t.Context()
	variable, _, err := client.Agents.GetOrgVariable(ctx, "o", "NAME")
	if err != nil {
		t.Errorf("Agents.GetOrgVariable returned error: %v", err)
	}

	want := &ActionsVariable{
		Name:                    "NAME",
		Value:                   "VALUE",
		CreatedAt:               &referenceTimestamp,
		UpdatedAt:               &referenceTimestamp,
		Visibility:              new("selected"),
		SelectedRepositoriesURL: new("https://api.github.com/orgs/octo-org/agents/variables/VAR/repositories"),
	}
	if !cmp.Equal(variable, want) {
		t.Errorf("Agents.GetOrgVariable returned %+v, want %+v", variable, want)
	}

	const methodName = "GetOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.GetOrgVariable(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.GetOrgVariable(ctx, "o", "NAME")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_CreateOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := ActionsCreateOrgVariableRequest{
		Name:                  "NAME",
		Value:                 "VALUE",
		Visibility:            "selected",
		SelectedRepositoryIDs: []int64{1296269, 1269280},
	}

	mux.HandleFunc("/orgs/o/agents/variables", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testHeader(t, r, "Content-Type", "application/json")
		testJSONBody(t, r, input)
		w.WriteHeader(http.StatusCreated)
	})

	ctx := t.Context()
	_, err := client.Agents.CreateOrgVariable(ctx, "o", input)
	if err != nil {
		t.Errorf("Agents.CreateOrgVariable returned error: %v", err)
	}

	const methodName = "CreateOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.CreateOrgVariable(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.CreateOrgVariable(ctx, "o", input)
	})
}

func TestAgentsService_UpdateOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	name := "NAME"
	input := ActionsUpdateOrgVariableRequest{
		Name:                  &name,
		Value:                 new("VALUE"),
		Visibility:            new("selected"),
		SelectedRepositoryIDs: []int64{1296269, 1269280},
	}

	mux.HandleFunc("/orgs/o/agents/variables/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testHeader(t, r, "Content-Type", "application/json")
		testJSONBody(t, r, input)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	_, err := client.Agents.UpdateOrgVariable(ctx, "o", name, input)
	if err != nil {
		t.Errorf("Agents.UpdateOrgVariable returned error: %v", err)
	}

	const methodName = "UpdateOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.UpdateOrgVariable(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.UpdateOrgVariable(ctx, "o", name, input)
	})
}

func TestAgentsService_ListSelectedReposForOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/variables/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		fmt.Fprint(w, `{"total_count":1,"repositories":[{"id":1}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	repos, _, err := client.Agents.ListSelectedReposForOrgVariable(ctx, "o", "NAME", opts)
	if err != nil {
		t.Errorf("Agents.ListSelectedReposForOrgVariable returned error: %v", err)
	}

	want := &SelectedReposList{
		TotalCount: new(1),
		Repositories: []*Repository{
			{ID: new(int64(1))},
		},
	}
	if !cmp.Equal(repos, want) {
		t.Errorf("Agents.ListSelectedReposForOrgVariable returned %+v, want %+v", repos, want)
	}

	const methodName = "ListSelectedReposForOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Agents.ListSelectedReposForOrgVariable(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Agents.ListSelectedReposForOrgVariable(ctx, "o", "NAME", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAgentsService_SetSelectedReposForOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := []int64{64780797}

	mux.HandleFunc("/orgs/o/agents/variables/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		testHeader(t, r, "Content-Type", "application/json")
		testJSONBody(t, r, struct {
			SelectedIDs []int64 `json:"selected_repository_ids"`
		}{
			SelectedIDs: input,
		})
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	_, err := client.Agents.SetSelectedReposForOrgVariable(ctx, "o", "NAME", input)
	if err != nil {
		t.Errorf("Agents.SetSelectedReposForOrgVariable returned error: %v", err)
	}

	const methodName = "SetSelectedReposForOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.SetSelectedReposForOrgVariable(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.SetSelectedReposForOrgVariable(ctx, "o", "NAME", input)
	})
}

func TestAgentsService_AddSelectedRepoToOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/variables/NAME/repositories/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		w.WriteHeader(http.StatusNoContent)
	})

	repoID := int64(1234)
	ctx := t.Context()
	_, err := client.Agents.AddSelectedRepoToOrgVariable(ctx, "o", "NAME", repoID)
	if err != nil {
		t.Errorf("Agents.AddSelectedRepoToOrgVariable returned error: %v", err)
	}

	const methodName = "AddSelectedRepoToOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.AddSelectedRepoToOrgVariable(ctx, "o", "NAME", 0)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.AddSelectedRepoToOrgVariable(ctx, "\n", "\n", repoID)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.AddSelectedRepoToOrgVariable(ctx, "o", "NAME", repoID)
	})
}

func TestAgentsService_RemoveSelectedRepoFromOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/variables/NAME/repositories/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		w.WriteHeader(http.StatusNoContent)
	})

	repoID := int64(1234)
	ctx := t.Context()
	_, err := client.Agents.RemoveSelectedRepoFromOrgVariable(ctx, "o", "NAME", repoID)
	if err != nil {
		t.Errorf("Agents.RemoveSelectedRepoFromOrgVariable returned error: %v", err)
	}

	const methodName = "RemoveSelectedRepoFromOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.RemoveSelectedRepoFromOrgVariable(ctx, "o", "NAME", 0)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.RemoveSelectedRepoFromOrgVariable(ctx, "\n", "\n", repoID)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.RemoveSelectedRepoFromOrgVariable(ctx, "o", "NAME", repoID)
	})
}

func TestAgentsService_DeleteOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/agents/variables/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "X-Github-Api-Version", api20260310)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := t.Context()
	_, err := client.Agents.DeleteOrgVariable(ctx, "o", "NAME")
	if err != nil {
		t.Errorf("Agents.DeleteOrgVariable returned error: %v", err)
	}

	const methodName = "DeleteOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Agents.DeleteOrgVariable(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Agents.DeleteOrgVariable(ctx, "o", "NAME")
	})
}
