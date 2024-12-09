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

func TestActionsService_GetActionsPermissionsInEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/actions/permissions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"enabled_organizations": "all", "allowed_actions": "all"}`)
	})

	ctx := context.Background()
	ent, _, err := client.Actions.GetActionsPermissionsInEnterprise(ctx, "e")
	if err != nil {
		t.Errorf("Actions.GetActionsPermissionsInEnterprise returned error: %v", err)
	}
	want := &ActionsPermissionsEnterprise{EnabledOrganizations: Ptr("all"), AllowedActions: Ptr("all")}
	if !cmp.Equal(ent, want) {
		t.Errorf("Actions.GetActionsPermissionsInEnterprise returned %+v, want %+v", ent, want)
	}

	const methodName = "GetActionsPermissionsInEnterprise"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetActionsPermissionsInEnterprise(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetActionsPermissionsInEnterprise(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_EditActionsPermissionsInEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &ActionsPermissionsEnterprise{EnabledOrganizations: Ptr("all"), AllowedActions: Ptr("selected")}

	mux.HandleFunc("/enterprises/e/actions/permissions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionsPermissionsEnterprise)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"enabled_organizations": "all", "allowed_actions": "selected"}`)
	})

	ctx := context.Background()
	ent, _, err := client.Actions.EditActionsPermissionsInEnterprise(ctx, "e", *input)
	if err != nil {
		t.Errorf("Actions.EditActionsPermissionsInEnterprise returned error: %v", err)
	}

	want := &ActionsPermissionsEnterprise{EnabledOrganizations: Ptr("all"), AllowedActions: Ptr("selected")}
	if !cmp.Equal(ent, want) {
		t.Errorf("Actions.EditActionsPermissionsInEnterprise returned %+v, want %+v", ent, want)
	}

	const methodName = "EditActionsPermissionsInEnterprise"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.EditActionsPermissionsInEnterprise(ctx, "\n", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.EditActionsPermissionsInEnterprise(ctx, "e", *input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_ListEnabledOrgsInEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/actions/permissions/organizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "1",
		})
		fmt.Fprint(w, `{"total_count":2,"organizations":[{"id":2}, {"id":3}]}`)
	})

	ctx := context.Background()
	opt := &ListOptions{
		Page: 1,
	}
	got, _, err := client.Actions.ListEnabledOrgsInEnterprise(ctx, "e", opt)
	if err != nil {
		t.Errorf("Actions.ListEnabledOrgsInEnterprise returned error: %v", err)
	}

	want := &ActionsEnabledOnEnterpriseRepos{TotalCount: int(2), Organizations: []*Organization{
		{ID: Ptr(int64(2))},
		{ID: Ptr(int64(3))},
	}}
	if !cmp.Equal(got, want) {
		t.Errorf("Actions.ListEnabledOrgsInEnterprise returned %+v, want %+v", got, want)
	}

	const methodName = "ListEnabledOrgsInEnterprise"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListEnabledOrgsInEnterprise(ctx, "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListEnabledOrgsInEnterprise(ctx, "e", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_SetEnabledOrgsInEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/actions/permissions/organizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"selected_organization_ids":[123,1234]}`+"\n")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Actions.SetEnabledOrgsInEnterprise(ctx, "e", []int64{123, 1234})
	if err != nil {
		t.Errorf("Actions.SetEnabledOrgsInEnterprise returned error: %v", err)
	}

	const methodName = "SetEnabledOrgsInEnterprise"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.SetEnabledOrgsInEnterprise(ctx, "\n", []int64{123, 1234})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.SetEnabledOrgsInEnterprise(ctx, "e", []int64{123, 1234})
	})
}

func TestActionsService_AddEnabledOrgInEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/actions/permissions/organizations/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Actions.AddEnabledOrgInEnterprise(ctx, "e", 123)
	if err != nil {
		t.Errorf("Actions.AddEnabledOrgInEnterprise returned error: %v", err)
	}

	const methodName = "AddEnabledOrgInEnterprise"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.AddEnabledOrgInEnterprise(ctx, "\n", 123)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.AddEnabledOrgInEnterprise(ctx, "e", 123)
	})
}

func TestActionsService_RemoveEnabledOrgInEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/actions/permissions/organizations/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Actions.RemoveEnabledOrgInEnterprise(ctx, "e", 123)
	if err != nil {
		t.Errorf("Actions.RemoveEnabledOrgInEnterprise returned error: %v", err)
	}

	const methodName = "RemoveEnabledOrgInEnterprise"

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.RemoveEnabledOrgInEnterprise(ctx, "\n", 123)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.RemoveEnabledOrgInEnterprise(ctx, "e", 123)
	})
}

func TestActionsService_GetActionsAllowedInEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/actions/permissions/selected-actions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"github_owned_allowed":true, "verified_allowed":false, "patterns_allowed":["a/b"]}`)
	})

	ctx := context.Background()
	ent, _, err := client.Actions.GetActionsAllowedInEnterprise(ctx, "e")
	if err != nil {
		t.Errorf("Actions.GetActionsAllowedInEnterprise returned error: %v", err)
	}
	want := &ActionsAllowed{GithubOwnedAllowed: Ptr(true), VerifiedAllowed: Ptr(false), PatternsAllowed: []string{"a/b"}}
	if !cmp.Equal(ent, want) {
		t.Errorf("Actions.GetActionsAllowedInEnterprise returned %+v, want %+v", ent, want)
	}

	const methodName = "GetActionsAllowedInEnterprise"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetActionsAllowedInEnterprise(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetActionsAllowedInEnterprise(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_EditActionsAllowedInEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &ActionsAllowed{GithubOwnedAllowed: Ptr(true), VerifiedAllowed: Ptr(false), PatternsAllowed: []string{"a/b"}}

	mux.HandleFunc("/enterprises/e/actions/permissions/selected-actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionsAllowed)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"github_owned_allowed":true, "verified_allowed":false, "patterns_allowed":["a/b"]}`)
	})

	ctx := context.Background()
	ent, _, err := client.Actions.EditActionsAllowedInEnterprise(ctx, "e", *input)
	if err != nil {
		t.Errorf("Actions.EditActionsAllowedInEnterprise returned error: %v", err)
	}

	want := &ActionsAllowed{GithubOwnedAllowed: Ptr(true), VerifiedAllowed: Ptr(false), PatternsAllowed: []string{"a/b"}}
	if !cmp.Equal(ent, want) {
		t.Errorf("Actions.EditActionsAllowedInEnterprise returned %+v, want %+v", ent, want)
	}

	const methodName = "EditActionsAllowedInEnterprise"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.EditActionsAllowedInEnterprise(ctx, "\n", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.EditActionsAllowedInEnterprise(ctx, "e", *input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetDefaultWorkflowPermissionsInEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/actions/permissions/workflow", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{ "default_workflow_permissions": "read", "can_approve_pull_request_reviews": true }`)
	})

	ctx := context.Background()
	ent, _, err := client.Actions.GetDefaultWorkflowPermissionsInEnterprise(ctx, "e")
	if err != nil {
		t.Errorf("Actions.GetDefaultWorkflowPermissionsInEnterprise returned error: %v", err)
	}
	want := &DefaultWorkflowPermissionEnterprise{DefaultWorkflowPermissions: Ptr("read"), CanApprovePullRequestReviews: Ptr(true)}
	if !cmp.Equal(ent, want) {
		t.Errorf("Actions.GetDefaultWorkflowPermissionsInEnterprise returned %+v, want %+v", ent, want)
	}

	const methodName = "GetDefaultWorkflowPermissionsInEnterprise"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetDefaultWorkflowPermissionsInEnterprise(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetDefaultWorkflowPermissionsInEnterprise(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_EditDefaultWorkflowPermissionsInEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &DefaultWorkflowPermissionEnterprise{DefaultWorkflowPermissions: Ptr("read"), CanApprovePullRequestReviews: Ptr(true)}

	mux.HandleFunc("/enterprises/e/actions/permissions/workflow", func(w http.ResponseWriter, r *http.Request) {
		v := new(DefaultWorkflowPermissionEnterprise)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{ "default_workflow_permissions": "read", "can_approve_pull_request_reviews": true }`)
	})

	ctx := context.Background()
	ent, _, err := client.Actions.EditDefaultWorkflowPermissionsInEnterprise(ctx, "e", *input)
	if err != nil {
		t.Errorf("Actions.EditDefaultWorkflowPermissionsInEnterprise returned error: %v", err)
	}

	want := &DefaultWorkflowPermissionEnterprise{DefaultWorkflowPermissions: Ptr("read"), CanApprovePullRequestReviews: Ptr(true)}
	if !cmp.Equal(ent, want) {
		t.Errorf("Actions.EditDefaultWorkflowPermissionsInEnterprise returned %+v, want %+v", ent, want)
	}

	const methodName = "EditDefaultWorkflowPermissionsInEnterprise"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.EditDefaultWorkflowPermissionsInEnterprise(ctx, "\n", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.EditDefaultWorkflowPermissionsInEnterprise(ctx, "e", *input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
