// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestActivityService_List(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/feeds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, feedsJSON)
	})

	ctx := t.Context()
	got, _, err := client.Activity.ListFeeds(ctx)
	if err != nil {
		t.Errorf("Activity.ListFeeds returned error: %v", err)
	}
	if want := wantFeeds; !cmp.Equal(got, want) {
		t.Errorf("Activity.ListFeeds = %+v, want %+v", got, want)
	}

	const methodName = "ListFeeds"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Activity.ListFeeds(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

var feedsJSON = []byte(`{
  "timeline_url": "https://github.com/timeline",
  "user_url": "https://github.com/{user}",
  "current_user_public_url": "https://github.com/defunkt",
  "current_user_url": "https://github.com/defunkt.private?token=abc123",
  "current_user_actor_url": "https://github.com/defunkt.private.actor?token=abc123",
  "current_user_organization_url": "",
  "current_user_organization_urls": [
    "https://github.com/organizations/github/defunkt.private.atom?token=abc123"
  ],
  "_links": {
    "timeline": {
      "href": "https://github.com/timeline",
      "type": "application/atom+xml"
    },
    "user": {
      "href": "https://github.com/{user}",
      "type": "application/atom+xml"
    },
    "current_user_public": {
      "href": "https://github.com/defunkt",
      "type": "application/atom+xml"
    },
    "current_user": {
      "href": "https://github.com/defunkt.private?token=abc123",
      "type": "application/atom+xml"
    },
    "current_user_actor": {
      "href": "https://github.com/defunkt.private.actor?token=abc123",
      "type": "application/atom+xml"
    },
    "current_user_organization": {
      "href": "",
      "type": ""
    },
    "current_user_organizations": [
      {
        "href": "https://github.com/organizations/github/defunkt.private.atom?token=abc123",
        "type": "application/atom+xml"
      }
    ]
  }
}`)

var wantFeeds = &Feeds{
	TimelineURL:                new("https://github.com/timeline"),
	UserURL:                    new("https://github.com/{user}"),
	CurrentUserPublicURL:       new("https://github.com/defunkt"),
	CurrentUserURL:             new("https://github.com/defunkt.private?token=abc123"),
	CurrentUserActorURL:        new("https://github.com/defunkt.private.actor?token=abc123"),
	CurrentUserOrganizationURL: new(""),
	CurrentUserOrganizationURLs: []string{
		"https://github.com/organizations/github/defunkt.private.atom?token=abc123",
	},
	Links: &FeedLinks{
		Timeline: &FeedLink{
			HRef: new("https://github.com/timeline"),
			Type: new("application/atom+xml"),
		},
		User: &FeedLink{
			HRef: new("https://github.com/{user}"),
			Type: new("application/atom+xml"),
		},
		CurrentUserPublic: &FeedLink{
			HRef: new("https://github.com/defunkt"),
			Type: new("application/atom+xml"),
		},
		CurrentUser: &FeedLink{
			HRef: new("https://github.com/defunkt.private?token=abc123"),
			Type: new("application/atom+xml"),
		},
		CurrentUserActor: &FeedLink{
			HRef: new("https://github.com/defunkt.private.actor?token=abc123"),
			Type: new("application/atom+xml"),
		},
		CurrentUserOrganization: &FeedLink{
			HRef: new(""),
			Type: new(""),
		},
		CurrentUserOrganizations: []*FeedLink{
			{
				HRef: new("https://github.com/organizations/github/defunkt.private.atom?token=abc123"),
				Type: new("application/atom+xml"),
			},
		},
	},
}
