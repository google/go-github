// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUsersService_ListEmails(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/user/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `["user@example.com"]`)
	})

	emails, _, err := client.Users.ListEmails()
	if err != nil {
		t.Errorf("Users.ListEmails returned error: %v", err)
	}

	want := []UserEmail{"user@example.com"}
	if !reflect.DeepEqual(emails, want) {
		t.Errorf("Users.ListEmails returned %+v, want %+v", emails, want)
	}
}

func TestUsersService_AddEmails(t *testing.T) {
	setup()
	defer teardown()

	input := []UserEmail{"new@example.com"}

	mux.HandleFunc("/user/emails", func(w http.ResponseWriter, r *http.Request) {
		v := new([]UserEmail)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(*v, input) {
			t.Errorf("Request body = %+v, want %+v", *v, input)
		}

		fmt.Fprint(w, `["old@example.com", "new@example.com"]`)
	})

	emails, _, err := client.Users.AddEmails(input)
	if err != nil {
		t.Errorf("Users.AddEmails returned error: %v", err)
	}

	want := []UserEmail{"old@example.com", "new@example.com"}
	if !reflect.DeepEqual(emails, want) {
		t.Errorf("Users.AddEmails returned %+v, want %+v", emails, want)
	}
}

func TestUsersService_DeleteEmails(t *testing.T) {
	setup()
	defer teardown()

	input := []UserEmail{"user@example.com"}

	mux.HandleFunc("/user/emails", func(w http.ResponseWriter, r *http.Request) {
		v := new([]UserEmail)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "DELETE")
		if !reflect.DeepEqual(*v, input) {
			t.Errorf("Request body = %+v, want %+v", *v, input)
		}
	})

	_, err := client.Users.DeleteEmails(input)
	if err != nil {
		t.Errorf("Users.DeleteEmails returned error: %v", err)
	}
}
