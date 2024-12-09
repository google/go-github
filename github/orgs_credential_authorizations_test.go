// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_ListCredentialAuthorizations(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/credential-authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{"per_page": "2", "page": "2", "login": "l"})
		fmt.Fprint(w, `[
			{
				"login": "l",
				"credential_id": 1,
				"credential_type": "t",
				"credential_authorized_at": "2017-01-21T00:00:00Z",
				"credential_accessed_at": "2017-01-21T00:00:00Z",
				"authorized_credential_id": 1
			}
		]`)
	})

	opts := &CredentialAuthorizationsListOptions{
		ListOptions: ListOptions{Page: 2, PerPage: 2},
		Login:       "l",
	}

	ctx := context.Background()
	creds, _, err := client.Organizations.ListCredentialAuthorizations(ctx, "o", opts)
	if err != nil {
		t.Errorf("Organizations.ListCredentialAuthorizations returned error: %v", err)
	}

	ts := time.Date(2017, time.January, 21, 0, 0, 0, 0, time.UTC)
	want := []*CredentialAuthorization{
		{
			Login:                  Ptr("l"),
			CredentialID:           Ptr(int64(1)),
			CredentialType:         Ptr("t"),
			CredentialAuthorizedAt: &Timestamp{ts},
			CredentialAccessedAt:   &Timestamp{ts},
			AuthorizedCredentialID: Ptr(int64(1)),
		},
	}
	if !cmp.Equal(creds, want) {
		t.Errorf("Organizations.ListCredentialAuthorizations returned %+v, want %+v", creds, want)
	}

	const methodName = "ListCredentialAuthorizations"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListCredentialAuthorizations(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.Organizations.ListCredentialAuthorizations(ctx, "o", opts)
		return resp, err
	})
}

func TestOrganizationsService_RemoveCredentialAuthorization(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/credential-authorizations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.Organizations.RemoveCredentialAuthorization(ctx, "o", 1)
	if err != nil {
		t.Errorf("Organizations.RemoveCredentialAuthorization returned error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Organizations.RemoveCredentialAuthorization returned %v, want %v", resp.StatusCode, http.StatusNoContent)
	}

	const methodName = "RemoveCredentialAuthorization"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.RemoveCredentialAuthorization(ctx, "\n", 0)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.RemoveCredentialAuthorization(ctx, "o", 1)
	})
}
