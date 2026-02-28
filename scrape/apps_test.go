// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scrape

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v84/github"
)

func Test_AppRestrictionsEnabled(t *testing.T) {
	t.Parallel()
	tests := []struct {
		description string
		testFile    string
		org         string
		want        bool
	}{
		{
			description: "return true for enabled orgs",
			testFile:    "access-restrictions-enabled.html",
			want:        true,
		},
		{
			description: "return false for disabled orgs",
			testFile:    "access-restrictions-disabled.html",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			client, mux := setup(t)

			mux.HandleFunc("/organizations/o/settings/oauth_application_policy", func(w http.ResponseWriter, _ *http.Request) {
				copyTestFile(t, w, tt.testFile)
			})

			got, err := client.AppRestrictionsEnabled("o")
			if err != nil {
				t.Fatalf("AppRestrictionsEnabled returned err: %v", err)
			}
			if want := tt.want; got != want {
				t.Errorf("AppRestrictionsEnabled returned %t, want %t", got, want)
			}
		})
	}
}

func Test_ListOAuthApps(t *testing.T) {
	t.Parallel()
	client, mux := setup(t)

	mux.HandleFunc("/organizations/e/settings/oauth_application_policy", func(w http.ResponseWriter, _ *http.Request) {
		copyTestFile(t, w, "access-restrictions-enabled.html")
	})

	got, err := client.ListOAuthApps("e")
	if err != nil {
		t.Fatalf("ListOAuthApps(e) returned err: %v", err)
	}
	want := []*OAuthApp{
		{
			ID:          22222,
			Name:        "Coveralls",
			Description: "Test coverage history and statistics.",
			State:       OAuthAppRequested,
			RequestedBy: "willnorris",
		},
		{
			ID:    530107,
			Name:  "Google Cloud Platform",
			State: OAuthAppApproved,
		},
		{
			ID:          231424,
			Name:        "GitKraken",
			Description: "An intuitive, cross-platform Git client that doesn't suck, built by @axosoft and made with @nodegit & @ElectronJS.",
			State:       OAuthAppDenied,
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("ListOAuthApps(o) returned %v, want %v", got, want)
	}
}

func Test_CreateApp(t *testing.T) {
	t.Parallel()
	client, mux := setup(t)

	mux.HandleFunc("/settings/apps/new", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := client.CreateApp(&AppManifest{
		URL: github.Ptr("https://example.com"),
		HookAttributes: map[string]string{
			"url": "https://example.com/hook",
		},
	}, "")
	if err != nil {
		t.Fatalf("CreateApp: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("CreateApp returned status code %v, want %v", got, want)
	}
}

func Test_CreateAppWithOrg(t *testing.T) {
	t.Parallel()
	client, mux := setup(t)

	mux.HandleFunc("/organizations/example/apps/settings/new", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	if _, err := client.CreateApp(&AppManifest{
		URL: github.Ptr("https://example.com"),
		HookAttributes: map[string]string{
			"url": "https://example.com/hook",
		},
	}, "example"); err != nil {
		t.Fatalf("CreateAppWithOrg: %v", err)
	}
}
