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

func TestInteractionsService_GetRestrictionsForAuthenticatedUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeInteractionRestrictionsPreview)
		fmt.Fprint(w, `{"origin":"user"}`)
	})

	ctx := t.Context()
	userInteractions, _, err := client.Interactions.GetRestrictionsForAuthenticatedUser(ctx)
	if err != nil {
		t.Errorf("Interactions.GetRestrictionsForAuthenticatedUser returned error: %v", err)
	}

	want := &InteractionRestriction{Origin: new("user")}
	if !cmp.Equal(userInteractions, want) {
		t.Errorf("Interactions.GetRestrictionsForAuthenticatedUser returned %+v, want %+v", userInteractions, want)
	}

	const methodName = "GetRestrictionsForAuthenticatedUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Interactions.GetRestrictionsForAuthenticatedUser(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestInteractionsService_UpdateRestrictionsForAuthenticatedUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &InteractionRestriction{Limit: new("existing_users")}

	mux.HandleFunc("/user/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeInteractionRestrictionsPreview)
		testJSONBody(t, r, input)
		fmt.Fprint(w, `{"origin":"user"}`)
	})

	ctx := t.Context()
	userInteractions, _, err := client.Interactions.UpdateRestrictionsForAuthenticatedUser(ctx, input.GetLimit())
	if err != nil {
		t.Errorf("Interactions.UpdateRestrictionsForAuthenticatedUser returned error: %v", err)
	}

	want := &InteractionRestriction{Origin: new("user")}
	if !cmp.Equal(userInteractions, want) {
		t.Errorf("Interactions.UpdateRestrictionsForAuthenticatedUser returned %+v, want %+v", userInteractions, want)
	}

	const methodName = "UpdateRestrictionsForAuthenticatedUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Interactions.UpdateRestrictionsForAuthenticatedUser(ctx, input.GetLimit())
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestInteractionsService_RemoveRestrictionsFromAuthenticatedUser(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/interaction-limits", func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeInteractionRestrictionsPreview)
	})

	ctx := t.Context()
	_, err := client.Interactions.RemoveRestrictionsFromAuthenticatedUser(ctx)
	if err != nil {
		t.Errorf("Interactions.RemoveRestrictionsFromAuthenticatedUser returned error: %v", err)
	}

	const methodName = "RemoveRestrictionsFromAuthenticatedUser"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Interactions.RemoveRestrictionsFromAuthenticatedUser(ctx)
	})
}
