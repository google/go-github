// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"net/http"
	"testing"
)

func TestUsersService_Promote(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/willnorris/site_admin", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	_, err := client.Users.Promote("willnorris")
	if err != nil {
		t.Errorf("Users.Promote returned error: %v", err)
	}
}

func TestUsersService_Demote(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/willnorris/site_admin", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Users.Demote("willnorris")
	if err != nil {
		t.Errorf("Users.Demote returned error: %v", err)
	}
}

func TestUsersService_Suspend(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/willnorris/suspended", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	_, err := client.Users.Suspend("willnorris")
	if err != nil {
		t.Errorf("Users.Suspend returned error: %v", err)
	}
}

func TestUsersService_Unsuspend(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/willnorris/suspended", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Users.Unsuspend("willnorris")
	if err != nil {
		t.Errorf("Users.Unsuspend returned error: %v", err)
	}
}
