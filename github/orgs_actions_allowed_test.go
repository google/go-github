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

func TestOrganizationsService_GetActionsAllowed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/permissions/selected-actions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"github_owned_allowed":true, "verified_allowed":false, "patterns_allowed":["a/b"]}`)
	})

	ctx := context.Background()
	org, _, err := client.Organizations.GetActionsAllowed(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.GetActionsAllowed returned error: %v", err)
	}
	want := &ActionsAllowed{GithubOwnedAllowed: Bool(true), VerifiedAllowed: Bool(false), PatternsAllowed: []string{"a/b"}}
	if !cmp.Equal(org, want) {
		t.Errorf("Organizations.GetActionsAllowed returned %+v, want %+v", org, want)
	}

	const methodName = "GetActionsAllowed"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetActionsAllowed(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetActionsAllowed(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_EditActionsAllowed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	input := &ActionsAllowed{GithubOwnedAllowed: Bool(true), VerifiedAllowed: Bool(false), PatternsAllowed: []string{"a/b"}}

	mux.HandleFunc("/orgs/o/actions/permissions/selected-actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionsAllowed)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"github_owned_allowed":true, "verified_allowed":false, "patterns_allowed":["a/b"]}`)
	})

	ctx := context.Background()
	org, _, err := client.Organizations.EditActionsAllowed(ctx, "o", *input)
	if err != nil {
		t.Errorf("Organizations.EditActionsAllowed returned error: %v", err)
	}

	want := &ActionsAllowed{GithubOwnedAllowed: Bool(true), VerifiedAllowed: Bool(false), PatternsAllowed: []string{"a/b"}}
	if !cmp.Equal(org, want) {
		t.Errorf("Organizations.EditActionsAllowed returned %+v, want %+v", org, want)
	}

	const methodName = "EditActionsAllowed"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.EditActionsAllowed(ctx, "\n", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.EditActionsAllowed(ctx, "o", *input)
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
