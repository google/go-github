// Copyright 2015 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAuthorizationsService_Check(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/applications/id/token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"access_token":"a"}`+"\n")
		testHeader(t, r, "Accept", mediaTypeOAuthAppPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	got, _, err := client.Authorizations.Check(ctx, "id", "a")
	if err != nil {
		t.Errorf("Authorizations.Check returned error: %v", err)
	}

	want := &Authorization{ID: Int64(1)}
	if !cmp.Equal(got, want) {
		t.Errorf("Authorizations.Check returned auth %+v, want %+v", got, want)
	}

	const methodName = "Check"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Authorizations.Check(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Authorizations.Check(ctx, "id", "a")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAuthorizationsService_Reset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/applications/id/token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testBody(t, r, `{"access_token":"a"}`+"\n")
		testHeader(t, r, "Accept", mediaTypeOAuthAppPreview)
		fmt.Fprint(w, `{"ID":1}`)
	})

	ctx := context.Background()
	got, _, err := client.Authorizations.Reset(ctx, "id", "a")
	if err != nil {
		t.Errorf("Authorizations.Reset returned error: %v", err)
	}

	want := &Authorization{ID: Int64(1)}
	if !cmp.Equal(got, want) {
		t.Errorf("Authorizations.Reset returned auth %+v, want %+v", got, want)
	}

	const methodName = "Reset"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Authorizations.Reset(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Authorizations.Reset(ctx, "id", "a")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAuthorizationsService_Revoke(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/applications/id/token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBody(t, r, `{"access_token":"a"}`+"\n")
		testHeader(t, r, "Accept", mediaTypeOAuthAppPreview)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Authorizations.Revoke(ctx, "id", "a")
	if err != nil {
		t.Errorf("Authorizations.Revoke returned error: %v", err)
	}

	const methodName = "Revoke"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Authorizations.Revoke(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Authorizations.Revoke(ctx, "id", "a")
	})
}

func TestDeleteGrant(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/applications/id/grant", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBody(t, r, `{"access_token":"a"}`+"\n")
		testHeader(t, r, "Accept", mediaTypeOAuthAppPreview)
	})

	ctx := context.Background()
	_, err := client.Authorizations.DeleteGrant(ctx, "id", "a")
	if err != nil {
		t.Errorf("OAuthAuthorizations.DeleteGrant returned error: %v", err)
	}

	const methodName = "DeleteGrant"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Authorizations.DeleteGrant(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Authorizations.DeleteGrant(ctx, "id", "a")
	})
}

func TestAuthorizationsService_CreateImpersonation(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/users/u/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":1}`)
	})

	req := &AuthorizationRequest{Scopes: []Scope{ScopePublicRepo}}
	ctx := context.Background()
	got, _, err := client.Authorizations.CreateImpersonation(ctx, "u", req)
	if err != nil {
		t.Errorf("Authorizations.CreateImpersonation returned error: %+v", err)
	}

	want := &Authorization{ID: Int64(1)}
	if !cmp.Equal(got, want) {
		t.Errorf("Authorizations.CreateImpersonation returned %+v, want %+v", *got.ID, *want.ID)
	}

	const methodName = "CreateImpersonation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Authorizations.CreateImpersonation(ctx, "\n", req)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Authorizations.CreateImpersonation(ctx, "u", req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAuthorizationsService_DeleteImpersonation(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/users/u/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Authorizations.DeleteImpersonation(ctx, "u")
	if err != nil {
		t.Errorf("Authorizations.DeleteImpersonation returned error: %+v", err)
	}

	const methodName = "DeleteImpersonation"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Authorizations.DeleteImpersonation(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Authorizations.DeleteImpersonation(ctx, "u")
	})
}
