// Copyright 2015 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
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

	got, _, err := client.Authorizations.Check(context.Background(), "id", "a")
	if err != nil {
		t.Errorf("Authorizations.Check returned error: %v", err)
	}

	want := &Authorization{ID: Int64(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authorizations.Check returned auth %+v, want %+v", got, want)
	}
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

	got, _, err := client.Authorizations.Reset(context.Background(), "id", "a")
	if err != nil {
		t.Errorf("Authorizations.Reset returned error: %v", err)
	}

	want := &Authorization{ID: Int64(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authorizations.Reset returned auth %+v, want %+v", got, want)
	}
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

	_, err := client.Authorizations.Revoke(context.Background(), "id", "a")
	if err != nil {
		t.Errorf("Authorizations.Revoke returned error: %v", err)
	}
}

func TestDeleteGrant(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/applications/id/grant", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBody(t, r, `{"access_token":"a"}`+"\n")
		testHeader(t, r, "Accept", mediaTypeOAuthAppPreview)
	})

	_, err := client.Authorizations.DeleteGrant(context.Background(), "id", "a")
	if err != nil {
		t.Errorf("OAuthAuthorizations.DeleteGrant returned error: %v", err)
	}
}

func TestAuthorizationsService_CreateImpersonation(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/users/u/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":1}`)
	})

	req := &AuthorizationRequest{Scopes: []Scope{ScopePublicRepo}}
	got, _, err := client.Authorizations.CreateImpersonation(context.Background(), "u", req)
	if err != nil {
		t.Errorf("Authorizations.CreateImpersonation returned error: %+v", err)
	}

	want := &Authorization{ID: Int64(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authorizations.CreateImpersonation returned %+v, want %+v", *got.ID, *want.ID)
	}
}

func TestAuthorizationsService_DeleteImpersonation(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/users/u/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Authorizations.DeleteImpersonation(context.Background(), "u")
	if err != nil {
		t.Errorf("Authorizations.DeleteImpersonation returned error: %+v", err)
	}
}
