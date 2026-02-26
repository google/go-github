// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIssuesService_ListIssueTimeline(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	wantAcceptHeaders := []string{mediaTypeTimelinePreview, mediaTypeProjectCardDetailsPreview}
	mux.HandleFunc("/repos/o/r/issues/1/timeline", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	ctx := t.Context()
	events, _, err := client.Issues.ListIssueTimeline(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("Issues.ListIssueTimeline returned error: %v", err)
	}

	want := []*Timeline{{ID: Ptr(int64(1))}}
	if !cmp.Equal(events, want) {
		t.Errorf("Issues.ListIssueTimeline = %+v, want %+v", events, want)
	}

	const methodName = "ListIssueTimeline"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.ListIssueTimeline(ctx, "\n", "\n", -1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.ListIssueTimeline(ctx, "o", "r", 1, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSource_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Source{}, "{}")

	u := &Source{
		ID:  Ptr(int64(1)),
		URL: Ptr("url"),
		Actor: &User{
			Login:     Ptr("l"),
			ID:        Ptr(int64(1)),
			NodeID:    Ptr("n"),
			URL:       Ptr("u"),
			ReposURL:  Ptr("r"),
			EventsURL: Ptr("e"),
			AvatarURL: Ptr("a"),
		},
		Type:  Ptr("type"),
		Issue: &Issue{ID: Ptr(int64(1))},
	}

	want := `{
		"id": 1,
		"url": "url",
		"actor": {
			"login": "l",
			"id": 1,
			"node_id": "n",
			"avatar_url": "a",
			"url": "u",
			"events_url": "e",
			"repos_url": "r"
		},
		"type": "type",
		"issue": {
			"id": 1
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestTimeline_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Timeline{}, "{}")

	u := &Timeline{
		ID:        Ptr(int64(1)),
		URL:       Ptr("url"),
		CommitURL: Ptr("curl"),
		Actor: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		Event:     Ptr("event"),
		CommitID:  Ptr("cid"),
		CreatedAt: &Timestamp{referenceTime},
		Label:     &Label{ID: Ptr(int64(1))},
		Assignee: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		Milestone: &Milestone{ID: Ptr(int64(1))},
		Source: &Source{
			ID:  Ptr(int64(1)),
			URL: Ptr("url"),
			Actor: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			Type:  Ptr("type"),
			Issue: &Issue{ID: Ptr(int64(1))},
		},
		Rename: &Rename{
			From: Ptr("from"),
			To:   Ptr("to"),
		},
		State: Ptr("state"),
	}

	want := `{
		"id": 1,
		"url": "url",
		"commit_url": "curl",
		"actor": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"event": "event",
		"commit_id": "cid",
		"created_at": ` + referenceTimeStr + `,
		"label": {
			"id": 1
		},
		"assignee": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"milestone": {
			"id": 1
		},
		"source": {
			"id": 1,
			"url": "url",
			"actor": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"url": "u",
				"events_url": "e",
				"repos_url": "r"
			},
			"type": "type",
			"issue": {
				"id": 1
			}
		},
		"rename": {
			"from": "from",
			"to": "to"
		},
		"state": "state"
	}`

	testJSONMarshal(t, u, want)
}

func TestTimeline_ReviewRequests(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/example-org/example-repo/issues/3/timeline", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
		    "id": 1234567890,
		    "url": "http://example.com/timeline/1",
		    "actor": {
			"login": "actor-user",
			"id": 1
		    },
		    "event": "review_requested",
		    "created_at": "2006-01-02T15:04:05Z",
		    "requested_reviewer": {
			"login": "reviewer-user",
			"id": 2
		    },
		    "review_requester": {
			"login": "requester-user",
			"id": 1
		    }
		},
		{
		    "id": 1234567891,
		    "url": "http://example.com/timeline/2",
		    "actor": {
			"login": "actor-user",
			"id": 1
		    },
		    "event": "review_request_removed",
		    "created_at": "2006-01-02T15:04:05Z",
		    "requested_reviewer": {
			"login": "reviewer-user",
			"id": 2
		    }
		}]`)
	})

	ctx := t.Context()
	events, _, err := client.Issues.ListIssueTimeline(ctx, "example-org", "example-repo", 3, nil)
	if err != nil {
		t.Errorf("Issues.ListIssueTimeline returned error: %v", err)
	}

	want := []*Timeline{
		{
			ID:  Ptr(int64(1234567890)),
			URL: Ptr("http://example.com/timeline/1"),
			Actor: &User{
				Login: Ptr("actor-user"),
				ID:    Ptr(int64(1)),
			},
			Event:     Ptr("review_requested"),
			CreatedAt: &Timestamp{referenceTime},
			Reviewer: &User{
				Login: Ptr("reviewer-user"),
				ID:    Ptr(int64(2)),
			},
			Requester: &User{
				Login: Ptr("requester-user"),
				ID:    Ptr(int64(1)),
			},
		},
		{
			ID:  Ptr(int64(1234567891)),
			URL: Ptr("http://example.com/timeline/2"),
			Actor: &User{
				Login: Ptr("actor-user"),
				ID:    Ptr(int64(1)),
			},
			Event:     Ptr("review_request_removed"),
			CreatedAt: &Timestamp{referenceTime},
			Reviewer: &User{
				Login: Ptr("reviewer-user"),
				ID:    Ptr(int64(2)),
			},
		},
	}

	if !cmp.Equal(events, want) {
		t.Errorf("Issues.ListIssueTimeline review request events = %+v, want %+v", events, want)
		diff := cmp.Diff(events, want)
		t.Errorf("Difference: %v", diff)
	}
}
