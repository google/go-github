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

func TestEnterpriseService_GetActionsPermissions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/enterprises/e/actions/permissions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"enabled_organizations": "all", "allowed_actions": "all"}`)
	})

	ctx := context.Background()
	ent, _, err := client.Enterprise.GetActionsPermissions(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.GetActionsPermissions returned error: %v", err)
	}
	want := &ActionsPermissionsEnterprise{EnabledOrganizations: String("all"), AllowedActions: String("all")}
	if !cmp.Equal(ent, want) {
		t.Errorf("Enterprise.GetActionsPermissions returned %+v, want %+v", ent, want)
	}

	const methodName = "GetActionsPermissions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.GetActionsPermissions(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetActionsPermissions(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_EditActionsPermissions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &ActionsPermissionsEnterprise{EnabledOrganizations: String("all"), AllowedActions: String("selected")}

	mux.HandleFunc("/enterprises/e/actions/permissions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionsPermissionsEnterprise)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"enabled_organizations": "all", "allowed_actions": "selected"}`)
	})

	ctx := context.Background()
	ent, _, err := client.Enterprise.EditActionsPermissions(ctx, "e", *input)
	if err != nil {
		t.Errorf("Enterprise.EditActionsPermissions returned error: %v", err)
	}

	want := &ActionsPermissionsEnterprise{EnabledOrganizations: String("all"), AllowedActions: String("selected")}
	if !cmp.Equal(ent, want) {
		t.Errorf("Enterprise.EditActionsPermissions returned %+v, want %+v", ent, want)
	}

	const methodName = "EditActionsPermissions"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.EditActionsPermissions(ctx, "\n", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.EditActionsPermissions(ctx, "e", *input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
