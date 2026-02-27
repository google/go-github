// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIssuesService_ListIssueEvents(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	events, _, err := client.Issues.ListIssueEvents(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("Issues.ListIssueEvents returned error: %v", err)
	}

	want := []*IssueEvent{{ID: Ptr(int64(1))}}
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/issues/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	ctx := t.Context()
	events, _, err := client.Issues.ListRepositoryEvents(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Issues.ListRepositoryEvents returned error: %v", err)
	}

	want := []*IssueEvent{{ID: Ptr(int64(1))}}
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
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/issues/events/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := t.Context()
	event, _, err := client.Issues.GetEvent(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Issues.GetEvent returned error: %v", err)
	}

	want := &IssueEvent{ID: Ptr(int64(1))}
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
	t.Parallel()
	testJSONMarshal(t, &Rename{}, "{}")

	u := &Rename{
		From: Ptr("from"),
		To:   Ptr("to"),
	}

	want := `{
		"from": "from",
		"to": "to"
	}`

	testJSONMarshal(t, u, want)
}

func TestDismissedReview_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &DismissedReview{}, "{}")

	u := &DismissedReview{
		State:             Ptr("state"),
		ReviewID:          Ptr(int64(1)),
		DismissalMessage:  Ptr("dm"),
		DismissalCommitID: Ptr("dcid"),
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
	t.Parallel()
	testJSONMarshal(t, &IssueEvent{}, "{}")

	u := &IssueEvent{
		ID:  Ptr(int64(1)),
		URL: Ptr("url"),
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
		CreatedAt: &Timestamp{referenceTime},
		Issue:     &Issue{ID: Ptr(int64(1))},
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
		Assigner: &User{
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
		CommitID:  Ptr("cid"),
		Milestone: &Milestone{ID: Ptr(int64(1))},
		Label:     &Label{ID: Ptr(int64(1))},
		Rename: &Rename{
			From: Ptr("from"),
			To:   Ptr("to"),
		},
		LockReason: Ptr("lr"),
		DismissedReview: &DismissedReview{
			State:             Ptr("state"),
			ReviewID:          Ptr(int64(1)),
			DismissalMessage:  Ptr("dm"),
			DismissalCommitID: Ptr("dcid"),
		},
		RequestedReviewer: &User{
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
		RequestedTeam: &Team{
			ID:              Ptr(int64(1)),
			NodeID:          Ptr("n"),
			Name:            Ptr("n"),
			Description:     Ptr("d"),
			URL:             Ptr("u"),
			Slug:            Ptr("s"),
			Permission:      Ptr("p"),
			Privacy:         Ptr("p"),
			MembersCount:    Ptr(1),
			ReposCount:      Ptr(1),
			MembersURL:      Ptr("m"),
			RepositoriesURL: Ptr("r"),
			Organization: &Organization{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				AvatarURL: Ptr("a"),
				HTMLURL:   Ptr("h"),
				Name:      Ptr("n"),
				Company:   Ptr("c"),
				Blog:      Ptr("b"),
				Location:  Ptr("l"),
				Email:     Ptr("e"),
			},
			Parent: &Team{
				ID:           Ptr(int64(1)),
				NodeID:       Ptr("n"),
				Name:         Ptr("n"),
				Description:  Ptr("d"),
				URL:          Ptr("u"),
				Slug:         Ptr("s"),
				Permission:   Ptr("p"),
				Privacy:      Ptr("p"),
				MembersCount: Ptr(1),
				ReposCount:   Ptr(1),
			},
			LDAPDN: Ptr("l"),
		},
		PerformedViaGithubApp: &App{
			ID:     Ptr(int64(1)),
			NodeID: Ptr("n"),
			Owner: &User{
				Login:     Ptr("l"),
				ID:        Ptr(int64(1)),
				NodeID:    Ptr("n"),
				URL:       Ptr("u"),
				ReposURL:  Ptr("r"),
				EventsURL: Ptr("e"),
				AvatarURL: Ptr("a"),
			},
			Name:        Ptr("n"),
			Description: Ptr("d"),
			HTMLURL:     Ptr("h"),
			ExternalURL: Ptr("u"),
		},
		ReviewRequester: &User{
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
		"requested_team": {
			"id": 1,
			"node_id": "n",
			"name": "n",
			"description": "d",
			"url": "u",
			"slug": "s",
			"permission": "p",
			"privacy": "p",
			"members_count": 1,
			"repos_count": 1,
			"members_url": "m",
			"repositories_url": "r",
			"organization": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"avatar_url": "a",
				"html_url": "h",
				"name": "n",
				"company": "c",
				"blog": "b",
				"location": "l",
				"email": "e"
			},
			"parent": {
				"id": 1,
				"node_id": "n",
				"name": "n",
				"description": "d",
				"url": "u",
				"slug": "s",
				"permission": "p",
				"privacy": "p",
				"members_count": 1,
				"repos_count": 1
			},
			"ldap_dn": "l"
		},
		"performed_via_github_app": {
			"id": 1,
			"node_id": "n",
			"owner": {
				"login": "l",
				"id": 1,
				"node_id": "n",
				"url": "u",
				"repos_url": "r",
				"events_url": "e",
				"avatar_url": "a"
			},
			"name": "n",
			"description": "d",
			"html_url": "h",
			"external_url": "u"
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
