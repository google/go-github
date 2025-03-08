// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
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

	ctx := context.Background()
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
	TimelineURL:                Ptr("https://github.com/timeline"),
	UserURL:                    Ptr("https://github.com/{user}"),
	CurrentUserPublicURL:       Ptr("https://github.com/defunkt"),
	CurrentUserURL:             Ptr("https://github.com/defunkt.private?token=abc123"),
	CurrentUserActorURL:        Ptr("https://github.com/defunkt.private.actor?token=abc123"),
	CurrentUserOrganizationURL: Ptr(""),
	CurrentUserOrganizationURLs: []string{
		"https://github.com/organizations/github/defunkt.private.atom?token=abc123",
	},
	Links: &FeedLinks{
		Timeline: &FeedLink{
			HRef: Ptr("https://github.com/timeline"),
			Type: Ptr("application/atom+xml"),
		},
		User: &FeedLink{
			HRef: Ptr("https://github.com/{user}"),
			Type: Ptr("application/atom+xml"),
		},
		CurrentUserPublic: &FeedLink{
			HRef: Ptr("https://github.com/defunkt"),
			Type: Ptr("application/atom+xml"),
		},
		CurrentUser: &FeedLink{
			HRef: Ptr("https://github.com/defunkt.private?token=abc123"),
			Type: Ptr("application/atom+xml"),
		},
		CurrentUserActor: &FeedLink{
			HRef: Ptr("https://github.com/defunkt.private.actor?token=abc123"),
			Type: Ptr("application/atom+xml"),
		},
		CurrentUserOrganization: &FeedLink{
			HRef: Ptr(""),
			Type: Ptr(""),
		},
		CurrentUserOrganizations: []*FeedLink{
			{
				HRef: Ptr("https://github.com/organizations/github/defunkt.private.atom?token=abc123"),
				Type: Ptr("application/atom+xml"),
			},
		},
	},
}

func TestFeedLink_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &FeedLink{}, "{}")

	u := &FeedLink{
		HRef: Ptr("h"),
		Type: Ptr("t"),
	}

	want := `{
		"href": "h",
		"type": "t"
	}`

	testJSONMarshal(t, u, want)
}

func TestFeeds_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Feeds{}, "{}")

	u := &Feeds{
		TimelineURL:                 Ptr("t"),
		UserURL:                     Ptr("u"),
		CurrentUserPublicURL:        Ptr("cupu"),
		CurrentUserURL:              Ptr("cuu"),
		CurrentUserActorURL:         Ptr("cuau"),
		CurrentUserOrganizationURL:  Ptr("cuou"),
		CurrentUserOrganizationURLs: []string{"a"},
		Links: &FeedLinks{
			Timeline: &FeedLink{
				HRef: Ptr("h"),
				Type: Ptr("t"),
			},
			User: &FeedLink{
				HRef: Ptr("h"),
				Type: Ptr("t"),
			},
			CurrentUserPublic: &FeedLink{
				HRef: Ptr("h"),
				Type: Ptr("t"),
			},
			CurrentUser: &FeedLink{
				HRef: Ptr("h"),
				Type: Ptr("t"),
			},
			CurrentUserActor: &FeedLink{
				HRef: Ptr("h"),
				Type: Ptr("t"),
			},
			CurrentUserOrganization: &FeedLink{
				HRef: Ptr("h"),
				Type: Ptr("t"),
			},
			CurrentUserOrganizations: []*FeedLink{
				{
					HRef: Ptr("h"),
					Type: Ptr("t"),
				},
			},
		},
	}

	want := `{
		"timeline_url": "t",
		"user_url": "u",
		"current_user_public_url": "cupu",
		"current_user_url": "cuu",
		"current_user_actor_url": "cuau",
		"current_user_organization_url": "cuou",
		"current_user_organization_urls": ["a"],
		"_links": {
			"timeline": {
				"href": "h",
				"type": "t"
				},
			"user": {
				"href": "h",
				"type": "t"
			},
			"current_user_public": {
				"href": "h",
				"type": "t"
			},
			"current_user": {
				"href": "h",
				"type": "t"
			},
			"current_user_actor": {
				"href": "h",
				"type": "t"
			},
			"current_user_organization": {
				"href": "h",
				"type": "t"
			},
			"current_user_organizations": [
				{
					"href": "h",
					"type": "t"
				}
			]
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestFeedLinks_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &FeedLinks{}, "{}")

	u := &FeedLinks{
		Timeline: &FeedLink{
			HRef: Ptr("h"),
			Type: Ptr("t"),
		},
		User: &FeedLink{
			HRef: Ptr("h"),
			Type: Ptr("t"),
		},
		CurrentUserPublic: &FeedLink{
			HRef: Ptr("h"),
			Type: Ptr("t"),
		},
		CurrentUser: &FeedLink{
			HRef: Ptr("h"),
			Type: Ptr("t"),
		},
		CurrentUserActor: &FeedLink{
			HRef: Ptr("h"),
			Type: Ptr("t"),
		},
		CurrentUserOrganization: &FeedLink{
			HRef: Ptr("h"),
			Type: Ptr("t"),
		},
		CurrentUserOrganizations: []*FeedLink{
			{
				HRef: Ptr("h"),
				Type: Ptr("t"),
			},
		},
	}

	want := `{
		"timeline": {
			"href": "h",
			"type": "t"
		},
		"user": {
			"href": "h",
			"type": "t"
		},
		"current_user_public": {
			"href": "h",
			"type": "t"
		},
		"current_user": {
			"href": "h",
			"type": "t"
		},
		"current_user_actor": {
			"href": "h",
			"type": "t"
		},
		"current_user_organization": {
			"href": "h",
			"type": "t"
		},
		"current_user_organizations": [
			{
				"href": "h",
				"type": "t"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}
