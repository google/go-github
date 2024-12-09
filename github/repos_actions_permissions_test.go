// Copyright 2022 The go-github AUTHORS. All rights reserved.
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

func TestRepositoriesService_GetActionsPermissions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/permissions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"enabled": true, "allowed_actions": "all"}`)
	})

	ctx := context.Background()
	org, _, err := client.Repositories.GetActionsPermissions(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetActionsPermissions returned error: %v", err)
	}
	want := &ActionsPermissionsRepository{Enabled: Ptr(true), AllowedActions: Ptr("all")}
	if !cmp.Equal(org, want) {
		t.Errorf("Repositories.GetActionsPermissions returned %+v, want %+v", org, want)
	}

	const methodName = "GetActionsPermissions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetActionsPermissions(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetActionsPermissions(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_EditActionsPermissions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &ActionsPermissionsRepository{Enabled: Ptr(true), AllowedActions: Ptr("selected")}

	mux.HandleFunc("/repos/o/r/actions/permissions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionsPermissionsRepository)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"enabled": true, "allowed_actions": "selected"}`)
	})

	ctx := context.Background()
	org, _, err := client.Repositories.EditActionsPermissions(ctx, "o", "r", *input)
	if err != nil {
		t.Errorf("Repositories.EditActionsPermissions returned error: %v", err)
	}

	want := &ActionsPermissionsRepository{Enabled: Ptr(true), AllowedActions: Ptr("selected")}
	if !cmp.Equal(org, want) {
		t.Errorf("Repositories.EditActionsPermissions returned %+v, want %+v", org, want)
	}

	const methodName = "EditActionsPermissions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.EditActionsPermissions(ctx, "\n", "\n", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.EditActionsPermissions(ctx, "o", "r", *input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsPermissionsRepository_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ActionsPermissions{}, "{}")

	u := &ActionsPermissionsRepository{
		Enabled:            Ptr(true),
		AllowedActions:     Ptr("all"),
		SelectedActionsURL: Ptr("someURL"),
	}

	want := `{
		"enabled": true,
		"allowed_actions": "all",
		"selected_actions_url": "someURL"
	}`

	testJSONMarshal(t, u, want)
}

func TestRepositoriesService_GetDefaultWorkflowPermissions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/permissions/workflow", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{ "default_workflow_permissions": "read", "can_approve_pull_request_reviews": true }`)
	})

	ctx := context.Background()
	org, _, err := client.Repositories.GetDefaultWorkflowPermissions(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetDefaultWorkflowPermissions returned error: %v", err)
	}
	want := &DefaultWorkflowPermissionRepository{DefaultWorkflowPermissions: Ptr("read"), CanApprovePullRequestReviews: Ptr(true)}
	if !cmp.Equal(org, want) {
		t.Errorf("Repositories.GetDefaultWorkflowPermissions returned %+v, want %+v", org, want)
	}

	const methodName = "GetDefaultWorkflowPermissions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetDefaultWorkflowPermissions(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetDefaultWorkflowPermissions(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_EditDefaultWorkflowPermissions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &DefaultWorkflowPermissionRepository{DefaultWorkflowPermissions: Ptr("read"), CanApprovePullRequestReviews: Ptr(true)}

	mux.HandleFunc("/repos/o/r/actions/permissions/workflow", func(w http.ResponseWriter, r *http.Request) {
		v := new(DefaultWorkflowPermissionRepository)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{ "default_workflow_permissions": "read", "can_approve_pull_request_reviews": true }`)
	})

	ctx := context.Background()
	org, _, err := client.Repositories.EditDefaultWorkflowPermissions(ctx, "o", "r", *input)
	if err != nil {
		t.Errorf("Repositories.EditDefaultWorkflowPermissions returned error: %v", err)
	}

	want := &DefaultWorkflowPermissionRepository{DefaultWorkflowPermissions: Ptr("read"), CanApprovePullRequestReviews: Ptr(true)}
	if !cmp.Equal(org, want) {
		t.Errorf("Repositories.EditDefaultWorkflowPermissions returned %+v, want %+v", org, want)
	}

	const methodName = "EditDefaultWorkflowPermissions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.EditDefaultWorkflowPermissions(ctx, "\n", "\n", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.EditDefaultWorkflowPermissions(ctx, "o", "r", *input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
