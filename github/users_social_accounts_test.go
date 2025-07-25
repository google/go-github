// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
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
		if (got != nil) {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
