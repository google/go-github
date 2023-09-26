// Copyright 2023 The go-github AUTHORS. All rights reserved.
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

func TestActionsService_GetActionsPermissions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/permissions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"enabled_repositories": "all", "allowed_actions": "all"}`)
	})

	ctx := context.Background()
	org, _, err := client.Actions.GetActionsPermissions(ctx, "o")
	if err != nil {
		t.Errorf("Actions.GetActionsPermissions returned error: %v", err)
	}
	want := &ActionsPermissions{EnabledRepositories: String("all"), AllowedActions: String("all")}
	if !cmp.Equal(org, want) {
		t.Errorf("Actions.GetActionsPermissions returned %+v, want %+v", org, want)
	}

	const methodName = "GetActionsPermissions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetActionsPermissions(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetActionsPermissions(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_EditActionsPermissions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &ActionsPermissions{EnabledRepositories: String("all"), AllowedActions: String("selected")}

	mux.HandleFunc("/orgs/o/actions/permissions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionsPermissions)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"enabled_repositories": "all", "allowed_actions": "selected"}`)
	})

	ctx := context.Background()
	org, _, err := client.Actions.EditActionsPermissions(ctx, "o", *input)
	if err != nil {
		t.Errorf("Actions.EditActionsPermissions returned error: %v", err)
	}

	want := &ActionsPermissions{EnabledRepositories: String("all"), AllowedActions: String("selected")}
	if !cmp.Equal(org, want) {
		t.Errorf("Actions.EditActionsPermissions returned %+v, want %+v", org, want)
	}

	const methodName = "EditActionsPermissions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.EditActionsPermissions(ctx, "\n", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.EditActionsPermissions(ctx, "o", *input)
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
		t.Errorf("Actions.ListEnabledRepos returned error: %v", err)
	}

	want := &ActionsEnabledOnOrgRepos{TotalCount: int(2), Repositories: []*Repository{
		{ID: Int64(2)},
		{ID: Int64(3)},
	}}
	if !cmp.Equal(got, want) {
		t.Errorf("Actions.ListEnabledRepos returned %+v, want %+v", got, want)
	}

	const methodName = "ListEnabledRepos"
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
		t.Errorf("Actions.SetEnabledRepos returned error: %v", err)
	}

	const methodName = "SetEnabledRepos"

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

func TestActionsService_RemoveEnabledReposInOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/permissions/repositories/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Actions.RemoveEnabledReposInOrg(ctx, "o", 123)
	if err != nil {
		t.Errorf("Actions.RemoveEnabledReposInOrg returned error: %v", err)
	}

	const methodName = "RemoveEnabledReposInOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RemoveEnabledReposInOrg(ctx, "\n", 123)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.RemoveEnabledReposInOrg(ctx, "o", 123)
	})
}

func TestActionsService_GetActionsAllowed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/permissions/selected-actions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"github_owned_allowed":true, "verified_allowed":false, "patterns_allowed":["a/b"]}`)
	})

	ctx := context.Background()
	org, _, err := client.Actions.GetActionsAllowed(ctx, "o")
	if err != nil {
		t.Errorf("Actions.GetActionsAllowed returned error: %v", err)
	}
	want := &ActionsAllowed{GithubOwnedAllowed: Bool(true), VerifiedAllowed: Bool(false), PatternsAllowed: []string{"a/b"}}
	if !cmp.Equal(org, want) {
		t.Errorf("Actions.GetActionsAllowed returned %+v, want %+v", org, want)
	}

	const methodName = "GetActionsAllowed"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetActionsAllowed(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetActionsAllowed(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_EditActionsAllowed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	input := &ActionsAllowed{GithubOwnedAllowed: Bool(true), VerifiedAllowed: Bool(false), PatternsAllowed: []string{"a/b"}}

	mux.HandleFunc("/orgs/o/actions/permissions/selected-actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionsAllowed)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"github_owned_allowed":true, "verified_allowed":false, "patterns_allowed":["a/b"]}`)
	})

	ctx := context.Background()
	org, _, err := client.Actions.EditActionsAllowed(ctx, "o", *input)
	if err != nil {
		t.Errorf("Actions.EditActionsAllowed returned error: %v", err)
	}

	want := &ActionsAllowed{GithubOwnedAllowed: Bool(true), VerifiedAllowed: Bool(false), PatternsAllowed: []string{"a/b"}}
	if !cmp.Equal(org, want) {
		t.Errorf("Actions.EditActionsAllowed returned %+v, want %+v", org, want)
	}

	const methodName = "EditActionsAllowed"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.EditActionsAllowed(ctx, "\n", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.EditActionsAllowed(ctx, "o", *input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsAllowed_Marshal(t *testing.T) {
	testJSONMarshal(t, &ActionsAllowed{}, "{}")

	u := &ActionsAllowed{
		GithubOwnedAllowed: Bool(false),
		VerifiedAllowed:    Bool(false),
		PatternsAllowed:    []string{"s"},
	}

	want := `{
		"github_owned_allowed": false,
		"verified_allowed": false,
		"patterns_allowed": [
			"s"
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestActionsPermissions_Marshal(t *testing.T) {
	testJSONMarshal(t, &ActionsPermissions{}, "{}")

	u := &ActionsPermissions{
		EnabledRepositories: String("e"),
		AllowedActions:      String("a"),
		SelectedActionsURL:  String("sau"),
	}

	want := `{
		"enabled_repositories": "e",
		"allowed_actions": "a",
		"selected_actions_url": "sau"
	}`

	testJSONMarshal(t, u, want)
}
