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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/permissions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"enabled": true, "allowed_actions": "all"}`)
	})

	ctx := context.Background()
	org, _, err := client.Repositories.GetActionsPermissions(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetActionsPermissions returned error: %v", err)
	}
	want := &ActionsPermissionsRepository{Enabled: Bool(true), AllowedActions: String("all")}
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
	client, mux, _, teardown := setup()
	defer teardown()

	input := &ActionsPermissionsRepository{Enabled: Bool(true), AllowedActions: String("selected")}

	mux.HandleFunc("/repos/o/r/actions/permissions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionsPermissionsRepository)
		json.NewDecoder(r.Body).Decode(v)

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

	want := &ActionsPermissionsRepository{Enabled: Bool(true), AllowedActions: String("selected")}
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
	testJSONMarshal(t, &ActionsPermissions{}, "{}")

	u := &ActionsPermissionsRepository{
		Enabled:            Bool(true),
		AllowedActions:     String("all"),
		SelectedActionsURL: String("someURL"),
	}

	want := `{
		"enabled": true,
		"allowed_actions": "all",
		"selected_actions_url": "someURL"
	}`

	testJSONMarshal(t, u, want)
}
