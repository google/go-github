// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestUsersService_PromoteSiteAdmin(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/site_admin", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Users.PromoteSiteAdmin(context.Background(), "u")
	if err != nil {
		t.Errorf("Users.PromoteSiteAdmin returned error: %v", err)
	}
}

func TestUsersService_DemoteSiteAdmin(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/site_admin", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Users.DemoteSiteAdmin(context.Background(), "u")
	if err != nil {
		t.Errorf("Users.DemoteSiteAdmin returned error: %v", err)
	}
}

func TestUsersService_Suspend(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/suspended", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Users.Suspend(context.Background(), "u", nil)
	if err != nil {
		t.Errorf("Users.Suspend returned error: %v", err)
	}
}

func TestUsersServiceReason_Suspend(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &UserSuspendOptions{Reason: String("test")}

	mux.HandleFunc("/users/u/suspended", func(w http.ResponseWriter, r *http.Request) {
		v := new(UserSuspendOptions)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Users.Suspend(context.Background(), "u", input)
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

	_, err := client.Users.Unsuspend(context.Background(), "u")
	if err != nil {
		t.Errorf("Users.Unsuspend returned error: %v", err)
	}
}
