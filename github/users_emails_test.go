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
	"reflect"
	"testing"
	"time"
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
	if !reflect.DeepEqual(emails, want) {
		t.Errorf("Users.ListEmails returned %+v, want %+v", emails, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Users.ListEmails(ctx, nil)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' ListEmails = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' ListEmails resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' ListEmails err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Users.ListEmails(ctx, nil)
	if got != nil {
		t.Errorf("rate.Reset.Time > now ListEmails = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now ListEmails resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now ListEmails err = nil, want error")
	}
}

func TestUsersService_AddEmails(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := []string{"new@example.com"}

	mux.HandleFunc("/user/emails", func(w http.ResponseWriter, r *http.Request) {
		var v []string
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
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
	if !reflect.DeepEqual(emails, want) {
		t.Errorf("Users.AddEmails returned %+v, want %+v", emails, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.Users.AddEmails(ctx, input)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' AddEmails = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' AddEmails resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' AddEmails err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.Users.AddEmails(ctx, input)
	if got != nil {
		t.Errorf("rate.Reset.Time > now AddEmails = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now AddEmails resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now AddEmails err = nil, want error")
	}
}

func TestUsersService_DeleteEmails(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := []string{"user@example.com"}

	mux.HandleFunc("/user/emails", func(w http.ResponseWriter, r *http.Request) {
		var v []string
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "DELETE")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	ctx := context.Background()
	_, err := client.Users.DeleteEmails(ctx, input)
	if err != nil {
		t.Errorf("Users.DeleteEmails returned error: %v", err)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	resp, err := client.Users.DeleteEmails(ctx, input)
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' DeleteEmails resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' DeleteEmails err = nil, want error")
	}
}
