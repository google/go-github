// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUsersService_ListEmails(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{
			"email": "user@example.com",
			"verified": false,
			"primary": true
		}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	emails, _, err := client.Users.ListEmails(ctx, opt)
	if err != nil {
		t.Errorf("Users.ListEmails returned error: %v", err)
	}

	want := []*UserEmail{{Email: String("user@example.com"), Verified: Bool(false), Primary: Bool(true)}}
	if !cmp.Equal(emails, want) {
		t.Errorf("Users.ListEmails returned %+v, want %+v", emails, want)
	}

	const methodName = "ListEmails"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.ListEmails(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_AddEmails(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := []string{"new@example.com"}

	mux.HandleFunc("/user/emails", func(w http.ResponseWriter, r *http.Request) {
		var v []string
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `[{"email":"old@example.com"}, {"email":"new@example.com"}]`)
	})

	ctx := context.Background()
	emails, _, err := client.Users.AddEmails(ctx, input)
	if err != nil {
		t.Errorf("Users.AddEmails returned error: %v", err)
	}

	want := []*UserEmail{
		{Email: String("old@example.com")},
		{Email: String("new@example.com")},
	}
	if !cmp.Equal(emails, want) {
		t.Errorf("Users.AddEmails returned %+v, want %+v", emails, want)
	}

	const methodName = "AddEmails"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.AddEmails(ctx, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_DeleteEmails(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := []string{"user@example.com"}

	mux.HandleFunc("/user/emails", func(w http.ResponseWriter, r *http.Request) {
		var v []string
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "DELETE")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	ctx := context.Background()
	_, err := client.Users.DeleteEmails(ctx, input)
	if err != nil {
		t.Errorf("Users.DeleteEmails returned error: %v", err)
	}

	const methodName = "DeleteEmails"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.DeleteEmails(ctx, input)
	})
}

func TestUserEmail_Marshal(t *testing.T) {
	testJSONMarshal(t, &UserEmail{}, "{}")

	u := &UserEmail{
		Email:      String("qwe@qwe.qwe"),
		Primary:    Bool(false),
		Verified:   Bool(true),
		Visibility: String("yes"),
	}

	want := `{
		"email": "qwe@qwe.qwe",
		"primary": false,
		"verified": true,
		"visibility": "yes"
	}`

	testJSONMarshal(t, u, want)
}
