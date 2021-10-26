// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUsersService_PromoteSiteAdmin(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/site_admin", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Users.PromoteSiteAdmin(ctx, "u")
	if err != nil {
		t.Errorf("Users.PromoteSiteAdmin returned error: %v", err)
	}

	const methodName = "PromoteSiteAdmin"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Users.PromoteSiteAdmin(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.PromoteSiteAdmin(ctx, "u")
	})
}

func TestUsersService_DemoteSiteAdmin(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/site_admin", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Users.DemoteSiteAdmin(ctx, "u")
	if err != nil {
		t.Errorf("Users.DemoteSiteAdmin returned error: %v", err)
	}

	const methodName = "DemoteSiteAdmin"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Users.DemoteSiteAdmin(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.DemoteSiteAdmin(ctx, "u")
	})
}

func TestUsersService_Suspend(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/suspended", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Users.Suspend(ctx, "u", nil)
	if err != nil {
		t.Errorf("Users.Suspend returned error: %v", err)
	}

	const methodName = "Suspend"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Users.Suspend(ctx, "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.Suspend(ctx, "u", nil)
	})
}

func TestUsersServiceReason_Suspend(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &UserSuspendOptions{Reason: String("test")}

	mux.HandleFunc("/users/u/suspended", func(w http.ResponseWriter, r *http.Request) {
		v := new(UserSuspendOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Users.Suspend(ctx, "u", input)
	if err != nil {
		t.Errorf("Users.Suspend returned error: %v", err)
	}
}

func TestUsersService_Unsuspend(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/suspended", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Users.Unsuspend(ctx, "u")
	if err != nil {
		t.Errorf("Users.Unsuspend returned error: %v", err)
	}

	const methodName = "Unsuspend"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Users.Unsuspend(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.Unsuspend(ctx, "u")
	})
}

func TestUserSuspendOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &UserSuspendOptions{}, "{}")

	u := &UserSuspendOptions{
		Reason: String("reason"),
	}

	want := `{
		"reason": "reason"
	}`

	testJSONMarshal(t, u, want)
}
