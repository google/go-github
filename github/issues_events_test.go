// Copyright 2014 The go-github AUTHORS. All rights reserved.
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

func TestIssuesService_ListIssueEvents(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/1/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectCardDetailsPreview)
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	ctx := context.Background()
	events, _, err := client.Issues.ListIssueEvents(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("Issues.ListIssueEvents returned error: %v", err)
	}

	want := []*IssueEvent{{ID: Int64(1)}}
	if !cmp.Equal(events, want) {
		t.Errorf("Issues.ListIssueEvents returned %+v, want %+v", events, want)
	}

	const methodName = "ListIssueEvents"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.ListIssueEvents(ctx, "\n", "\n", -1, &ListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.ListIssueEvents(ctx, "o", "r", 1, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_ListRepositoryEvents(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	ctx := context.Background()
	events, _, err := client.Issues.ListRepositoryEvents(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Issues.ListRepositoryEvents returned error: %v", err)
	}

	want := []*IssueEvent{{ID: Int64(1)}}
	if !cmp.Equal(events, want) {
		t.Errorf("Issues.ListRepositoryEvents returned %+v, want %+v", events, want)
	}

	const methodName = "ListRepositoryEvents"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.ListRepositoryEvents(ctx, "\n", "\n", &ListOptions{})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.ListRepositoryEvents(ctx, "o", "r", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestIssuesService_GetEvent(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/events/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	event, _, err := client.Issues.GetEvent(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Issues.GetEvent returned error: %v", err)
	}

	want := &IssueEvent{ID: Int64(1)}
	if !cmp.Equal(event, want) {
		t.Errorf("Issues.GetEvent returned %+v, want %+v", event, want)
	}

	const methodName = "GetEvent"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Issues.GetEvent(ctx, "\n", "\n", -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Issues.GetEvent(ctx, "o", "r", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRename_Marshal(t *testing.T) {
	testJSONMarshal(t, &Rename{}, "{}")

	u := &Rename{
		From: String("from"),
		To:   String("to"),
	}

	want := `{
		"from": "from",
		"to": "to"
	}`

	testJSONMarshal(t, u, want)
}

func TestDismissedReview_Marshal(t *testing.T) {
	testJSONMarshal(t, &DismissedReview{}, "{}")

	u := &DismissedReview{
		State:             String("state"),
		ReviewID:          Int64(1),
		DismissalMessage:  String("dm"),
		DismissalCommitID: String("dcid"),
	}

	want := `{
		"state": "state",
		"review_id": 1,
		"dismissal_message": "dm",
		"dismissal_commit_id": "dcid"
	}`

	testJSONMarshal(t, u, want)
}

func TestIssueEvent_Marshal(t *testing.T) {
	testJSONMarshal(t, &IssueEvent{}, "{}")

	u := &IssueEvent{
		ID:  Int64(1),
		URL: String("url"),
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
		CreatedAt: &Timestamp{referenceTime},
		Issue:     &Issue{ID: Int64(1)},
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
		Assigner: &User{
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
		CommitID:  String("cid"),
		Milestone: &Milestone{ID: Int64(1)},
		Label:     &Label{ID: Int64(1)},
		Rename: &Rename{
			From: String("from"),
			To:   String("to"),
		},
		LockReason:  String("lr"),
		ProjectCard: &ProjectCard{ID: Int64(1)},
		DismissedReview: &DismissedReview{
			State:             String("state"),
			ReviewID:          Int64(1),
			DismissalMessage:  String("dm"),
			DismissalCommitID: String("dcid"),
		},
		RequestedReviewer: &User{
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
		ReviewRequester: &User{
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
	}

	want := `{
		"id": 1,
		"url": "url",
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
		"created_at": ` + referenceTimeStr + `,
		"issue": {
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
		"assigner": {
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
		"commit_id": "cid",
		"milestone": {
			"id": 1
		},
		"label": {
			"id": 1
		},
		"rename": {
			"from": "from",
			"to": "to"
		},
		"lock_reason": "lr",
		"project_card": {
			"id": 1
		},
		"dismissed_review": {
			"state": "state",
			"review_id": 1,
			"dismissal_message": "dm",
			"dismissal_commit_id": "dcid"
		},
		"requested_reviewer": {
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
		"review_requester": {
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
		}
	}`

	testJSONMarshal(t, u, want)
}
