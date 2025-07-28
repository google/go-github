// Copyright 2025 The go-github AUTHORS. All rights reserved.
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

func TestUsersService_ListSocialAccounts(t *testing.T) {
	t.Parallel()

	client, mux, _ := setup(t)

	mux.HandleFunc("/user/social_accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{
			"provider": "twitter",
			"url": "https://twitter.com/github"
		}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	accounts, _, err := client.Users.ListSocialAccounts(ctx, opt)
	if err != nil {
		t.Errorf("Users.ListSocialAccounts returned error: %v", err)
	}

	want := []*SocialAccount{{Provider: Ptr("twitter"), URL: Ptr("https://twitter.com/github")}}
	if !cmp.Equal(accounts, want) {
		t.Errorf("Users.ListSocialAccounts returned %#v, want %#v", accounts, want)
	}

	const methodName = "ListSocialAccounts"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.ListSocialAccounts(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_AddSocialAccounts(t *testing.T) {
	t.Parallel()

	client, mux, _ := setup(t)

	input := []string{"https://twitter.com/github"}

	mux.HandleFunc("/user/social_accounts", func(w http.ResponseWriter, r *http.Request) {
		var v []string
		assertNilError(t, json.NewDecoder(r.Body).Decode(&v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `[{"provider":"twitter","url":"https://twitter.com/github"},{"provider":"facebook","url":"https://facebook.com/github"}]`)
	})

	ctx := context.Background()
	accounts, _, err := client.Users.AddSocialAccounts(ctx, input)
	if err != nil {
		t.Errorf("Users.AddSocialAccounts returned error: %v", err)
	}

	want := []*SocialAccount{
		{Provider: Ptr("twitter"), URL: Ptr("https://twitter.com/github")},
		{Provider: Ptr("facebook"), URL: Ptr("https://facebook.com/github")},
	}
	if !cmp.Equal(accounts, want) {
		t.Errorf("Users.AddSocialAccounts returned %#v, want %#v", accounts, want)
	}

	const methodName = "AddSocialAccounts"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.AddSocialAccounts(ctx, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_DeleteSocialAccounts(t *testing.T) {
	t.Parallel()

	client, mux, _ := setup(t)

	input := []string{"https://twitter.com/github"}

	mux.HandleFunc("/user/social_accounts", func(_ http.ResponseWriter, r *http.Request) {
		var v []string
		assertNilError(t, json.NewDecoder(r.Body).Decode(&v))

		testMethod(t, r, "DELETE")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	ctx := context.Background()
	_, err := client.Users.DeleteSocialAccounts(ctx, input)
	if err != nil {
		t.Errorf("Users.DeleteSocialAccounts returned error: %v", err)
	}

	const methodName = "DeleteSocialAccounts"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.DeleteSocialAccounts(ctx, input)
	})
}

func TestUsersService_ListUserSocialAccounts(t *testing.T) {
	t.Parallel()

	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/social_accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{
			"provider": "twitter",
			"url": "https://twitter.com/github"
		}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	accounts, _, err := client.Users.ListUserSocialAccounts(ctx, "u", opt)
	if err != nil {
		t.Errorf("Users.ListUserSocialAccounts returned error: %v", err)
	}

	want := []*SocialAccount{{Provider: Ptr("twitter"), URL: Ptr("https://twitter.com/github")}}
	if !cmp.Equal(accounts, want) {
		t.Errorf("Users.ListUserSocialAccounts returned %#v, want %#v", accounts, want)
	}

	const methodName = "ListUserSocialAccounts"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.ListUserSocialAccounts(ctx, "u", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
