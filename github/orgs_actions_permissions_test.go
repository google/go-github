// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_GetActionsPermissions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/permissions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"enabled_repositories": "all", "allowed_actions": "all"}`)
	})

	ctx := context.Background()
	org, _, err := client.Organizations.GetActionsPermissions(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.GetActionsPermissions returned error: %v", err)
	}
	want := &ActionsPermissions{EnabledRepositories: String("all"), AllowedActions: String("all")}
	if !cmp.Equal(org, want) {
		t.Errorf("Organizations.GetActionsPermissions returned %+v, want %+v", org, want)
	}

	const methodName = "GetActionsPermissions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetActionsPermissions(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetActionsPermissions(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_EditActionsPermissions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &ActionsPermissions{EnabledRepositories: String("all"), AllowedActions: String("selected")}

	mux.HandleFunc("/orgs/o/actions/permissions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionsPermissions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"enabled_repositories": "all", "allowed_actions": "selected"}`)
	})

	ctx := context.Background()
	org, _, err := client.Organizations.EditActionsPermissions(ctx, "o", *input)
	if err != nil {
		t.Errorf("Organizations.EditActionsPermissions returned error: %v", err)
	}

	want := &ActionsPermissions{EnabledRepositories: String("all"), AllowedActions: String("selected")}
	if !cmp.Equal(org, want) {
		t.Errorf("Organizations.EditActionsPermissions returned %+v, want %+v", org, want)
	}

	const methodName = "EditActionsPermissions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.EditActionsPermissions(ctx, "\n", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.EditActionsPermissions(ctx, "o", *input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_ListEnabledReposInOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/permissions/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "1",
		})
		fmt.Fprint(w, `{"total_count":2,"repositories":[{"id":2}, {"id": 3}]}`)
	})

	ctx := context.Background()
	opt := &ListOptions{
		Page: 1,
	}
	got, _, err := client.Actions.ListEnabledReposInOrg(ctx, "o", opt)
	if err != nil {
		t.Errorf("Actions.ListEnabledReposInOrg returned error: %v", err)
	}

	want := &ActionsEnabledOnOrgRepos{TotalCount: int(2), Repositories: []*Repository{
		{ID: Int64(2)},
		{ID: Int64(3)},
	}}
	if !cmp.Equal(got, want) {
		t.Errorf("Actions.ListEnabledReposInOrg returned %+v, want %+v", got, want)
	}

	const methodName = "ListEnabledReposInOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListEnabledReposInOrg(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListEnabledReposInOrg(ctx, "o", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_SetEnabledReposInOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/permissions/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"selected_repository_ids":[123,1234]}`+"\n")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Actions.SetEnabledReposInOrg(ctx, "o", []int64{123, 1234})
	if err != nil {
		t.Errorf("Actions.SetEnabledReposInOrg returned error: %v", err)
	}

	const methodName = "SetEnabledReposInOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.SetEnabledReposInOrg(ctx, "\n", []int64{123, 1234})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.SetEnabledReposInOrg(ctx, "o", []int64{123, 1234})
	})
}

func TestActionsService_AddEnabledReposInOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/permissions/repositories/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Actions.AddEnabledReposInOrg(ctx, "o", 123)
	if err != nil {
		t.Errorf("Actions.AddEnabledReposInOrg returned error: %v", err)
	}

	const methodName = "AddEnabledReposInOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.AddEnabledReposInOrg(ctx, "\n", 123)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.AddEnabledReposInOrg(ctx, "o", 123)
	})
}

func TestActionsService_RemoveEnabledRepoInOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/permissions/repositories/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Actions.RemoveEnabledRepoInOrg(ctx, "o", 123)
	if err != nil {
		t.Errorf("Actions.RemoveEnabledRepoInOrg returned error: %v", err)
	}

	const methodName = "RemoveEnabledRepoInOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RemoveEnabledRepoInOrg(ctx, "\n", 123)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.RemoveEnabledRepoInOrg(ctx, "o", 123)
	})
}
