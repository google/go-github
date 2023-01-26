// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIssuesService_ListIssueTimeline(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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
	ctx := context.Background()
	events, _, err := client.Issues.ListIssueTimeline(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("Issues.ListIssueTimeline returned error: %v", err)
	}

	want := []*Timeline{{ID: Int64(1)}}
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
	testJSONMarshal(t, &Source{}, "{}")

	u := &Source{
		ID:  Int64(1),
		URL: String("url"),
		Actor: &User{
			Login:     String("l"),
			ID:        Int64(1),
			NodeID:    String("n"),
			URL:       String("u"),
			ReposURL:  String("r"),
			EventsURL: String("e"),
			AvatarURL: String("a"),
		},
		Type:  String("type"),
		Issue: &Issue{ID: Int64(1)},
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
	testJSONMarshal(t, &Timeline{}, "{}")

	u := &Timeline{
		ID:        Int64(1),
		URL:       String("url"),
		CommitURL: String("curl"),
		Actor: &User{
			Login:           String("l"),
			ID:              Int64(1),
			URL:             String("u"),
			AvatarURL:       String("a"),
			GravatarID:      String("g"),
			Name:            String("n"),
			Company:         String("c"),
			Blog:            String("b"),
			Location:        String("l"),
			Email:           String("e"),
			Hireable:        Bool(true),
			Bio:             String("b"),
			TwitterUsername: String("t"),
			PublicRepos:     Int(1),
			Followers:       Int(1),
			Following:       Int(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		Event:     String("event"),
		CommitID:  String("cid"),
		CreatedAt: &Timestamp{referenceTime},
		Label:     &Label{ID: Int64(1)},
		Assignee: &User{
			Login:           String("l"),
			ID:              Int64(1),
			URL:             String("u"),
			AvatarURL:       String("a"),
			GravatarID:      String("g"),
			Name:            String("n"),
			Company:         String("c"),
			Blog:            String("b"),
			Location:        String("l"),
			Email:           String("e"),
			Hireable:        Bool(true),
			Bio:             String("b"),
			TwitterUsername: String("t"),
			PublicRepos:     Int(1),
			Followers:       Int(1),
			Following:       Int(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		Milestone: &Milestone{ID: Int64(1)},
		Source: &Source{
			ID:  Int64(1),
			URL: String("url"),
			Actor: &User{
				Login:     String("l"),
				ID:        Int64(1),
				NodeID:    String("n"),
				URL:       String("u"),
				ReposURL:  String("r"),
				EventsURL: String("e"),
				AvatarURL: String("a"),
			},
			Type:  String("type"),
			Issue: &Issue{ID: Int64(1)},
		},
		Rename: &Rename{
			From: String("from"),
			To:   String("to"),
		},
		ProjectCard: &ProjectCard{ID: Int64(1)},
		State:       String("state"),
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
		"project_card": {
			"id": 1
		},
		"state": "state"
	}`

	testJSONMarshal(t, u, want)
}
