// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestActionsService_ListRepoVariables(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/variables", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"variables":[{"name":"A","value":"AA","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","value":"BB","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	variables, _, err := client.Actions.ListRepoVariables(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListRepoVariables returned error: %v", err)
	}

	want := &ActionsVariables{
		TotalCount: 4,
		Variables: []*ActionsVariable{
			{Name: "A", Value: "AA", CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
			{Name: "B", Value: "BB", CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
		},
	}
	if !cmp.Equal(variables, want) {
		t.Errorf("Actions.ListRepoVariables returned %+v, want %+v", variables, want)
	}

	const methodName = "ListRepoVariables"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListRepoVariables(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListRepoVariables(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_ListRepoOrgVariables(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/organization-variables", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"variables":[{"name":"A","value":"AA","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","value":"BB","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	variables, _, err := client.Actions.ListRepoOrgVariables(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListRepoOrgVariables returned error: %v", err)
	}

	want := &ActionsVariables{
		TotalCount: 4,
		Variables: []*ActionsVariable{
			{Name: "A", Value: "AA", CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
			{Name: "B", Value: "BB", CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
		},
	}
	if !cmp.Equal(variables, want) {
		t.Errorf("Actions.ListRepoOrgVariables returned %+v, want %+v", variables, want)
	}

	const methodName = "ListRepoOrgVariables"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListRepoOrgVariables(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListRepoOrgVariables(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetRepoVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/variables/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"NAME","value":"VALUE","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}`)
	})

	ctx := t.Context()
	variable, _, err := client.Actions.GetRepoVariable(ctx, "o", "r", "NAME")
	if err != nil {
		t.Errorf("Actions.GetRepoVariable returned error: %v", err)
	}

	want := &ActionsVariable{
		Name:      "NAME",
		Value:     "VALUE",
		CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)},
		UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)},
	}
	if !cmp.Equal(variable, want) {
		t.Errorf("Actions.GetRepoVariable returned %+v, want %+v", variable, want)
	}

	const methodName = "GetRepoVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetRepoVariable(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetRepoVariable(ctx, "o", "r", "NAME")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_CreateRepoVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/variables", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"name":"NAME","value":"VALUE"}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &ActionsVariable{
		Name:  "NAME",
		Value: "VALUE",
	}
	ctx := t.Context()
	_, err := client.Actions.CreateRepoVariable(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Actions.CreateRepoVariable returned error: %v", err)
	}

	const methodName = "CreateRepoVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.CreateRepoVariable(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.CreateRepoVariable(ctx, "o", "r", input)
	})
}

func TestActionsService_UpdateRepoVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/variables/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"name":"NAME","value":"VALUE"}`+"\n")
		w.WriteHeader(http.StatusNoContent)
	})

	input := &ActionsVariable{
		Name:  "NAME",
		Value: "VALUE",
	}
	ctx := t.Context()
	_, err := client.Actions.UpdateRepoVariable(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Actions.UpdateRepoVariable returned error: %v", err)
	}

	const methodName = "UpdateRepoVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.UpdateRepoVariable(ctx, "o", "r", nil)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.UpdateRepoVariable(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.UpdateRepoVariable(ctx, "o", "r", input)
	})
}

func TestActionsService_DeleteRepoVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/variables/NAME", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := t.Context()
	_, err := client.Actions.DeleteRepoVariable(ctx, "o", "r", "NAME")
	if err != nil {
		t.Errorf("Actions.( returned error: %v", err)
	}

	const methodName = "DeleteRepoVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteRepoVariable(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteRepoVariable(ctx, "o", "r", "NAME")
	})
}

func TestActionsService_ListOrgVariables(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/variables", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":3,"variables":[{"name":"A","value":"AA","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"private"},{"name":"B","value":"BB","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"all"},{"name":"C","value":"CC","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/actions/variables/VAR/repositories"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	variables, _, err := client.Actions.ListOrgVariables(ctx, "o", opts)
	if err != nil {
		t.Errorf("Actions.ListOrgVariables returned error: %v", err)
	}

	want := &ActionsVariables{
		TotalCount: 3,
		Variables: []*ActionsVariable{
			{Name: "A", Value: "AA", CreatedAt: &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)}, Visibility: Ptr("private")},
			{Name: "B", Value: "BB", CreatedAt: &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)}, Visibility: Ptr("all")},
			{Name: "C", Value: "CC", CreatedAt: &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)}, Visibility: Ptr("selected"), SelectedRepositoriesURL: Ptr("https://api.github.com/orgs/octo-org/actions/variables/VAR/repositories")},
		},
	}
	if !cmp.Equal(variables, want) {
		t.Errorf("Actions.ListOrgVariables returned %+v, want %+v", variables, want)
	}

	const methodName = "ListOrgVariables"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListOrgVariables(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListOrgVariables(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/variables/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"NAME","value":"VALUE","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z","visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/actions/variables/VAR/repositories"}`)
	})

	ctx := t.Context()
	variable, _, err := client.Actions.GetOrgVariable(ctx, "o", "NAME")
	if err != nil {
		t.Errorf("Actions.GetOrgVariable returned error: %v", err)
	}

	want := &ActionsVariable{
		Name:                    "NAME",
		Value:                   "VALUE",
		CreatedAt:               &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)},
		UpdatedAt:               &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)},
		Visibility:              Ptr("selected"),
		SelectedRepositoriesURL: Ptr("https://api.github.com/orgs/octo-org/actions/variables/VAR/repositories"),
	}
	if !cmp.Equal(variable, want) {
		t.Errorf("Actions.GetOrgVariable returned %+v, want %+v", variable, want)
	}

	const methodName = "GetOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetOrgVariable(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetOrgVariable(ctx, "o", "NAME")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_CreateOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/variables", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"name":"NAME","value":"VALUE","visibility":"selected","selected_repository_ids":[1296269,1269280]}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &ActionsVariable{
		Name:                  "NAME",
		Value:                 "VALUE",
		Visibility:            Ptr("selected"),
		SelectedRepositoryIDs: &SelectedRepoIDs{1296269, 1269280},
	}
	ctx := t.Context()
	_, err := client.Actions.CreateOrgVariable(ctx, "o", input)
	if err != nil {
		t.Errorf("Actions.CreateOrgVariable returned error: %v", err)
	}

	const methodName = "CreateOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.CreateOrgVariable(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.CreateOrgVariable(ctx, "o", input)
	})
}

func TestActionsService_UpdateOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/variables/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"name":"NAME","value":"VALUE","visibility":"selected","selected_repository_ids":[1296269,1269280]}`+"\n")
		w.WriteHeader(http.StatusNoContent)
	})

	input := &ActionsVariable{
		Name:                  "NAME",
		Value:                 "VALUE",
		Visibility:            Ptr("selected"),
		SelectedRepositoryIDs: &SelectedRepoIDs{1296269, 1269280},
	}
	ctx := t.Context()
	_, err := client.Actions.UpdateOrgVariable(ctx, "o", input)
	if err != nil {
		t.Errorf("Actions.UpdateOrgVariable returned error: %v", err)
	}

	const methodName = "UpdateOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.UpdateOrgVariable(ctx, "o", nil)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.UpdateOrgVariable(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.UpdateOrgVariable(ctx, "o", input)
	})
}

func TestActionsService_ListSelectedReposForOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/variables/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count":1,"repositories":[{"id":1}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	repos, _, err := client.Actions.ListSelectedReposForOrgVariable(ctx, "o", "NAME", opts)
	if err != nil {
		t.Errorf("Actions.( returned error: %v", err)
	}

	want := &SelectedReposList{
		TotalCount: Ptr(1),
		Repositories: []*Repository{
			{ID: Ptr(int64(1))},
		},
	}
	if !cmp.Equal(repos, want) {
		t.Errorf("Actions.( returned %+v, want %+v", repos, want)
	}

	const methodName = "ListSelectedReposForOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListSelectedReposForOrgVariable(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListSelectedReposForOrgVariable(ctx, "o", "NAME", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_SetSelectedReposForOrgSVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/variables/NAME/repositories", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"selected_repository_ids":[64780797]}`+"\n")
	})

	ctx := t.Context()
	_, err := client.Actions.SetSelectedReposForOrgVariable(ctx, "o", "NAME", SelectedRepoIDs{64780797})
	if err != nil {
		t.Errorf("Actions.( returned error: %v", err)
	}

	const methodName = "SetSelectedReposForOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.SetSelectedReposForOrgVariable(ctx, "\n", "\n", SelectedRepoIDs{64780797})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.SetSelectedReposForOrgVariable(ctx, "o", "NAME", SelectedRepoIDs{64780797})
	})
}

func TestActionsService_AddSelectedRepoToOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/variables/NAME/repositories/1234", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	repo := &Repository{ID: Ptr(int64(1234))}
	ctx := t.Context()
	_, err := client.Actions.AddSelectedRepoToOrgVariable(ctx, "o", "NAME", repo)
	if err != nil {
		t.Errorf("Actions.AddSelectedRepoToOrgVariable returned error: %v", err)
	}

	const methodName = "AddSelectedRepoToOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.AddSelectedRepoToOrgVariable(ctx, "o", "NAME", nil)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.AddSelectedRepoToOrgVariable(ctx, "o", "NAME", &Repository{ID: nil})
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.AddSelectedRepoToOrgVariable(ctx, "\n", "\n", repo)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.AddSelectedRepoToOrgVariable(ctx, "o", "NAME", repo)
	})
}

func TestActionsService_RemoveSelectedRepoFromOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/variables/NAME/repositories/1234", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	repo := &Repository{ID: Ptr(int64(1234))}
	ctx := t.Context()
	_, err := client.Actions.RemoveSelectedRepoFromOrgVariable(ctx, "o", "NAME", repo)
	if err != nil {
		t.Errorf("Actions.RemoveSelectedRepoFromOrgVariable returned error: %v", err)
	}

	const methodName = "RemoveSelectedRepoFromOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RemoveSelectedRepoFromOrgVariable(ctx, "o", "NAME", nil)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RemoveSelectedRepoFromOrgVariable(ctx, "o", "NAME", &Repository{ID: nil})
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RemoveSelectedRepoFromOrgVariable(ctx, "\n", "\n", repo)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.RemoveSelectedRepoFromOrgVariable(ctx, "o", "NAME", repo)
	})
}

func TestActionsService_DeleteOrgVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/variables/NAME", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := t.Context()
	_, err := client.Actions.DeleteOrgVariable(ctx, "o", "NAME")
	if err != nil {
		t.Errorf("Actions.DeleteOrgVariable returned error: %v", err)
	}

	const methodName = "DeleteOrgVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteOrgVariable(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteOrgVariable(ctx, "o", "NAME")
	})
}

func TestActionsService_ListEnvVariables(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/usr/1/environments/e/variables", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"variables":[{"name":"A","value":"AA","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","value":"BB","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	ctx := t.Context()
	variables, _, err := client.Actions.ListEnvVariables(ctx, "usr", "1", "e", opts)
	if err != nil {
		t.Errorf("Actions.ListEnvVariables returned error: %v", err)
	}

	want := &ActionsVariables{
		TotalCount: 4,
		Variables: []*ActionsVariable{
			{Name: "A", Value: "AA", CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
			{Name: "B", Value: "BB", CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)}, UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)}},
		},
	}
	if !cmp.Equal(variables, want) {
		t.Errorf("Actions.ListEnvVariables returned %+v, want %+v", variables, want)
	}

	const methodName = "ListEnvVariables"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListEnvVariables(ctx, "usr", "0", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListEnvVariables(ctx, "usr", "1", "e", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetEnvVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/usr/1/environments/e/variables/variable", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"variable","value":"VAR","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}`)
	})

	ctx := t.Context()
	variable, _, err := client.Actions.GetEnvVariable(ctx, "usr", "1", "e", "variable")
	if err != nil {
		t.Errorf("Actions.GetEnvVariable returned error: %v", err)
	}

	want := &ActionsVariable{
		Name:      "variable",
		Value:     "VAR",
		CreatedAt: &Timestamp{time.Date(2019, time.January, 2, 15, 4, 5, 0, time.UTC)},
		UpdatedAt: &Timestamp{time.Date(2020, time.January, 2, 15, 4, 5, 0, time.UTC)},
	}
	if !cmp.Equal(variable, want) {
		t.Errorf("Actions.GetEnvVariable returned %+v, want %+v", variable, want)
	}

	const methodName = "GetEnvVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetEnvVariable(ctx, "usr", "0", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetEnvVariable(ctx, "usr", "1", "e", "variable")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_CreateEnvVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/usr/1/environments/e/variables", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"name":"variable","value":"VAR"}`+"\n")
		w.WriteHeader(http.StatusCreated)
	})

	input := &ActionsVariable{
		Name:  "variable",
		Value: "VAR",
	}
	ctx := t.Context()
	_, err := client.Actions.CreateEnvVariable(ctx, "usr", "1", "e", input)
	if err != nil {
		t.Errorf("Actions.CreateEnvVariable returned error: %v", err)
	}

	const methodName = "CreateEnvVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.CreateEnvVariable(ctx, "usr", "0", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.CreateEnvVariable(ctx, "usr", "1", "e", input)
	})
}

func TestActionsService_UpdateEnvVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/usr/1/environments/e/variables/variable", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"name":"variable","value":"VAR"}`+"\n")
		w.WriteHeader(http.StatusNoContent)
	})

	input := &ActionsVariable{
		Name:  "variable",
		Value: "VAR",
	}
	ctx := t.Context()
	_, err := client.Actions.UpdateEnvVariable(ctx, "usr", "1", "e", input)
	if err != nil {
		t.Errorf("Actions.UpdateEnvVariable returned error: %v", err)
	}

	const methodName = "UpdateEnvVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.UpdateEnvVariable(ctx, "usr", "1", "e", nil)
		return err
	})
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.UpdateEnvVariable(ctx, "usr", "1", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.UpdateEnvVariable(ctx, "usr", "1", "e", input)
	})
}

func TestActionsService_DeleteEnvVariable(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/usr/1/environments/e/variables/variable", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := t.Context()
	_, err := client.Actions.DeleteEnvVariable(ctx, "usr", "1", "e", "variable")
	if err != nil {
		t.Errorf("Actions.DeleteEnvVariable returned error: %v", err)
	}

	const methodName = "DeleteEnvVariable"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteEnvVariable(ctx, "usr", "0", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteEnvVariable(ctx, "usr", "1", "r", "variable")
	})
}

func TestActionVariable_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ActionsVariable{}, `{
		"name": "",
		"value": ""
	}`)

	av := &ActionsVariable{
		Name:                    "n",
		Value:                   "v",
		CreatedAt:               &Timestamp{referenceTime},
		UpdatedAt:               &Timestamp{referenceTime},
		Visibility:              Ptr("v"),
		SelectedRepositoriesURL: Ptr("s"),
		SelectedRepositoryIDs:   &SelectedRepoIDs{1, 2, 3},
	}

	want := fmt.Sprintf(`{
		"name": "n",
		"value": "v",
		"created_at": %v,
		"updated_at": %v,
		"visibility": "v",
		"selected_repositories_url": "s",
		"selected_repository_ids": [1,2,3]
	}`, referenceTimeStr, referenceTimeStr)

	testJSONMarshal(t, av, want)
}
