// Copyright 2023 The go-github AUTHORS. All rights reserved.
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
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_ListFineGrainedPersonalAccessTokens(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/personal-access-tokens", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		expectedQuery := map[string][]string{
			"per_page":  {"2"},
			"page":      {"2"},
			"sort":      {"created_at"},
			"direction": {"desc"},
			"owner[]":   {"octocat", "octodog", "otherbot"},
		}

		query := r.URL.Query()
		for key, expectedValues := range expectedQuery {
			actualValues := query[key]
			if len(actualValues) != len(expectedValues) {
				t.Errorf("Expected %d values for query param %s, got %d", len(expectedValues), key, len(actualValues))
			}
			for i, expectedValue := range expectedValues {
				if actualValues[i] != expectedValue {
					t.Errorf("Expected query param %s to be %s, got %s", key, expectedValue, actualValues[i])
				}
			}
		}

		fmt.Fprint(w, `
		[
			{
				"id": 25381,
				"owner": {
					"login": "octocat",
					"id": 1,
					"node_id": "MDQ6VXNlcjE=",
					"avatar_url": "https://github.com/images/error/octocat_happy.gif",
					"gravatar_id": "",
					"url": "https://api.github.com/users/octocat",
					"html_url": "https://github.com/octocat",
					"followers_url": "https://api.github.com/users/octocat/followers",
					"following_url": "https://api.github.com/users/octocat/following{/other_user}",
					"gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
					"starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
					"subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
					"organizations_url": "https://api.github.com/users/octocat/orgs",
					"repos_url": "https://api.github.com/users/octocat/repos",
					"events_url": "https://api.github.com/users/octocat/events{/privacy}",
					"received_events_url": "https://api.github.com/users/octocat/received_events",
					"type": "User",
					"site_admin": false
				},
				"repository_selection": "all",
				"repositories_url": "https://api.github.com/organizations/652551/personal-access-tokens/25381/repositories",
				"permissions": {
					"organization": {
						"members": "read"
					},
					"repository": {
						"metadata": "read"
					}
				},
				"access_granted_at": "2023-05-16T08:47:09.000-07:00",
				"token_expired": false,
				"token_expires_at": "2023-11-16T08:47:09.000-07:00",
				"token_last_used_at": null
			}
		]`)
	})

	opts := &ListFineGrainedPATOptions{
		ListOptions: ListOptions{Page: 2, PerPage: 2},
		Sort:        "created_at",
		Direction:   "desc",
		Owner:       []string{"octocat", "octodog", "otherbot"},
	}
	ctx := context.Background()
	tokens, resp, err := client.Organizations.ListFineGrainedPersonalAccessTokens(ctx, "o", opts)
	if err != nil {
		t.Errorf("Organizations.ListFineGrainedPersonalAccessTokens returned error: %v", err)
	}

	want := []*PersonalAccessToken{
		{
			ID: Ptr(int64(25381)),
			Owner: &User{
				Login:             Ptr("octocat"),
				ID:                Ptr(int64(1)),
				NodeID:            Ptr("MDQ6VXNlcjE="),
				AvatarURL:         Ptr("https://github.com/images/error/octocat_happy.gif"),
				GravatarID:        Ptr(""),
				URL:               Ptr("https://api.github.com/users/octocat"),
				HTMLURL:           Ptr("https://github.com/octocat"),
				FollowersURL:      Ptr("https://api.github.com/users/octocat/followers"),
				FollowingURL:      Ptr("https://api.github.com/users/octocat/following{/other_user}"),
				GistsURL:          Ptr("https://api.github.com/users/octocat/gists{/gist_id}"),
				StarredURL:        Ptr("https://api.github.com/users/octocat/starred{/owner}{/repo}"),
				SubscriptionsURL:  Ptr("https://api.github.com/users/octocat/subscriptions"),
				OrganizationsURL:  Ptr("https://api.github.com/users/octocat/orgs"),
				ReposURL:          Ptr("https://api.github.com/users/octocat/repos"),
				EventsURL:         Ptr("https://api.github.com/users/octocat/events{/privacy}"),
				ReceivedEventsURL: Ptr("https://api.github.com/users/octocat/received_events"),
				Type:              Ptr("User"),
				SiteAdmin:         Ptr(false),
			},
			RepositorySelection: Ptr("all"),
			RepositoriesURL:     Ptr("https://api.github.com/organizations/652551/personal-access-tokens/25381/repositories"),
			Permissions: &PersonalAccessTokenPermissions{
				Org:  map[string]string{"members": "read"},
				Repo: map[string]string{"metadata": "read"},
			},
			AccessGrantedAt: &Timestamp{time.Date(2023, time.May, 16, 8, 47, 9, 0, time.FixedZone("PDT", -7*60*60))},
			TokenExpired:    Ptr(false),
			TokenExpiresAt:  &Timestamp{time.Date(2023, time.November, 16, 8, 47, 9, 0, time.FixedZone("PDT", -7*60*60))},
			TokenLastUsedAt: nil,
		},
	}
	if !cmp.Equal(tokens, want) {
		t.Errorf("Organizations.ListFineGrainedPersonalAccessTokens returned %+v, want %+v", tokens, want)
	}

	if resp == nil {
		t.Error("Organizations.ListFineGrainedPersonalAccessTokens returned nil response")
	}

	const methodName = "ListFineGrainedPersonalAccessTokens"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListFineGrainedPersonalAccessTokens(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListFineGrainedPersonalAccessTokens(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_ReviewPersonalAccessTokenRequest(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := ReviewPersonalAccessTokenRequestOptions{
		Action: "a",
		Reason: Ptr("r"),
	}

	mux.HandleFunc("/orgs/o/personal-access-token-requests/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(ReviewPersonalAccessTokenRequestOptions)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, http.MethodPost)
		if !cmp.Equal(v, &input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	res, err := client.Organizations.ReviewPersonalAccessTokenRequest(ctx, "o", 1, input)
	if err != nil {
		t.Errorf("Organizations.ReviewPersonalAccessTokenRequest returned error: %v", err)
	}

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Organizations.ReviewPersonalAccessTokenRequest returned %v, want %v", res.StatusCode, http.StatusNoContent)
	}

	const methodName = "ReviewPersonalAccessTokenRequest"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.ReviewPersonalAccessTokenRequest(ctx, "\n", 0, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.ReviewPersonalAccessTokenRequest(ctx, "o", 1, input)
	})
}

func TestReviewPersonalAccessTokenRequestOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ReviewPersonalAccessTokenRequestOptions{}, "{}")

	u := &ReviewPersonalAccessTokenRequestOptions{
		Action: "a",
		Reason: Ptr("r"),
	}

	want := `{
		"action": "a",
		"reason": "r"
	}`

	testJSONMarshal(t, u, want)
}
