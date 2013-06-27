// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUsersService_Get_authenticatedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	user, err := client.Users.Get("")
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := &User{ID: 1}
	if !reflect.DeepEqual(user, want) {
		t.Errorf("Users.Get returned %+v, want %+v", user, want)
	}
}

func TestUsersService_Get_specifiedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	user, err := client.Users.Get("u")
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := &User{ID: 1}
	if !reflect.DeepEqual(user, want) {
		t.Errorf("Users.Get returned %+v, want %+v", user, want)
	}
}

func TestUsersService_Get_invalidUser(t *testing.T) {
	_, err := client.Users.Get("%")
	testURLParseError(t, err)
}

func TestUsersService_Edit(t *testing.T) {
	setup()
	defer teardown()

	input := &User{Name: "n"}

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		v := new(User)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	user, err := client.Users.Edit(input)
	if err != nil {
		t.Errorf("Users.Edit returned error: %v", err)
	}

	want := &User{ID: 1}
	if !reflect.DeepEqual(user, want) {
		t.Errorf("Users.Edit returned %+v, want %+v", user, want)
	}
}

func TestUsersService_ListAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"since": "1"})
		fmt.Fprint(w, `[{"id":2}]`)
	})

	opt := &UserListOptions{1}
	users, err := client.Users.ListAll(opt)
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := []User{User{ID: 2}}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("Users.ListAll returned %+v, want %+v", users, want)
	}
}

func TestUsersService_ListEmails(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/user/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `["user@example.com"]`)
	})

	emails, err := client.Users.ListEmails()
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

	emails, err := client.Users.AddEmails(input)
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

	err := client.Users.DeleteEmails(input)
	if err != nil {
		t.Errorf("Users.DeleteEmails returned error: %v", err)
	}
}
