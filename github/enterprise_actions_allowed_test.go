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

func TestEnterpriseService_GetActionsAllowed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/enterprises/e/actions/permissions/selected-actions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"github_owned_allowed":true, "verified_allowed":false, "patterns_allowed":["a/b"]}`)
	})

	ctx := context.Background()
	ent, _, err := client.Enterprise.GetActionsAllowed(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.GetActionsAllowed returned error: %v", err)
	}
	want := &ActionsAllowed{GithubOwnedAllowed: Bool(true), VerifiedAllowed: Bool(false), PatternsAllowed: []string{"a/b"}}
	if !cmp.Equal(ent, want) {
		t.Errorf("Enterprise.GetActionsAllowed returned %+v, want %+v", ent, want)
	}

	const methodName = "GetActionsAllowed"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.GetActionsAllowed(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetActionsAllowed(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_EditActionsAllowed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	input := &ActionsAllowed{GithubOwnedAllowed: Bool(true), VerifiedAllowed: Bool(false), PatternsAllowed: []string{"a/b"}}

	mux.HandleFunc("/enterprises/e/actions/permissions/selected-actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionsAllowed)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"github_owned_allowed":true, "verified_allowed":false, "patterns_allowed":["a/b"]}`)
	})

	ctx := context.Background()
	ent, _, err := client.Enterprise.EditActionsAllowed(ctx, "e", *input)
	if err != nil {
		t.Errorf("Enterprise.EditActionsAllowed returned error: %v", err)
	}

	want := &ActionsAllowed{GithubOwnedAllowed: Bool(true), VerifiedAllowed: Bool(false), PatternsAllowed: []string{"a/b"}}
	if !cmp.Equal(ent, want) {
		t.Errorf("Enterprise.EditActionsAllowed returned %+v, want %+v", ent, want)
	}

	const methodName = "EditActionsAllowed"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.EditActionsAllowed(ctx, "\n", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.EditActionsAllowed(ctx, "e", *input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}