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

func TestInteractionsService_GetUserRestrictions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"origin":"user"}`)
	})

	ctx := t.Context()
	userInteractions, _, err := client.Interactions.GetUserRestrictions(ctx)
	if err != nil {
		t.Errorf("Interactions.GetUserRestrictions returned error: %v", err)
	}

	want := &InteractionRestriction{Origin: new("user")}
	if !cmp.Equal(userInteractions, want) {
		t.Errorf("Interactions.GetUserRestrictions returned %+v, want %+v", userInteractions, want)
	}

	const methodName = "GetUserRestrictions"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Interactions.GetUserRestrictions(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestInteractionsService_UpdateUserRestrictions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := UpdateUserInteractionRestrictionsRequest{
		Limit:  "collaborators_only",
		Expiry: new("one_month"),
	}

	mux.HandleFunc("/user/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{"limit":"collaborators_only","origin":"user"}`)
	})

	ctx := t.Context()
	userInteractions, _, err := client.Interactions.UpdateUserRestrictions(ctx, input)
	if err != nil {
		t.Errorf("Interactions.UpdateUserRestrictions returned error: %v", err)
	}

	want := &InteractionRestriction{Limit: new("collaborators_only"), Origin: new("user")}
	if !cmp.Equal(userInteractions, want) {
		t.Errorf("Interactions.UpdateUserRestrictions returned %+v, want %+v", userInteractions, want)
	}

	const methodName = "UpdateUserRestrictions"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Interactions.UpdateUserRestrictions(ctx, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestInteractionsService_RemoveUserRestrictions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/interaction-limits", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := t.Context()
	_, err := client.Interactions.RemoveUserRestrictions(ctx)
	if err != nil {
		t.Errorf("Interactions.RemoveUserRestrictions returned error: %v", err)
	}

	const methodName = "RemoveUserRestrictions"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Interactions.RemoveUserRestrictions(ctx)
	})
}
