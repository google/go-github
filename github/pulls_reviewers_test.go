// Copyright 2017 The go-github AUTHORS. All rights reserved.
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

func TestReviewersRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ReviewersRequest{}, "{}")

	u := &ReviewersRequest{
		NodeID:        Ptr("n"),
		Reviewers:     []string{"r"},
		TeamReviewers: []string{"t"},
	}

	want := `{
		"node_id": "n",
		"reviewers": [
			"r"
		],
		"team_reviewers" : [
			"t"
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestReviewers_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Reviewers{}, "{}")

	u := &Reviewers{
		Users: []*User{{
			Login:       Ptr("l"),
			ID:          Ptr(int64(1)),
			AvatarURL:   Ptr("a"),
			GravatarID:  Ptr("g"),
			Name:        Ptr("n"),
			Company:     Ptr("c"),
			Blog:        Ptr("b"),
			Location:    Ptr("l"),
			Email:       Ptr("e"),
			Hireable:    Ptr(true),
			PublicRepos: Ptr(1),
			Followers:   Ptr(1),
			Following:   Ptr(1),
			CreatedAt:   &Timestamp{referenceTime},
			URL:         Ptr("u"),
		}},
		Teams: []*Team{{
			ID:              Ptr(int64(1)),
			NodeID:          Ptr("node"),
			Name:            Ptr("n"),
			Description:     Ptr("d"),
			URL:             Ptr("u"),
			Slug:            Ptr("s"),
			Permission:      Ptr("p"),
			Privacy:         Ptr("priv"),
			MembersCount:    Ptr(1),
			ReposCount:      Ptr(1),
			Organization:    nil,
			MembersURL:      Ptr("m"),
			RepositoriesURL: Ptr("r"),
			Parent:          nil,
			LDAPDN:          Ptr("l"),
		}},
	}

	want := `{
		"users" : [
			{
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
				"public_repos": 1,
				"followers": 1,
				"following": 1,
				"created_at": ` + referenceTimeStr + `,
				"url": "u"
			}
		],
		"teams" : [
			{
				"id": 1,
				"node_id": "node",
				"name": "n",
				"description": "d",
				"url": "u",
				"slug": "s",
				"permission": "p",
				"privacy": "priv",
				"members_count": 1,
				"repos_count": 1,
				"members_url": "m",
				"repositories_url": "r",
				"ldap_dn": "l"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestRequestReviewers(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1/requested_reviewers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"reviewers":["octocat","googlebot"],"team_reviewers":["justice-league","injustice-league"]}`+"\n")
		fmt.Fprint(w, `{"number":1}`)
	})

	// This returns a PR, unmarshaling of which is tested elsewhere
	ctx := context.Background()
	got, _, err := client.PullRequests.RequestReviewers(ctx, "o", "r", 1, ReviewersRequest{Reviewers: []string{"octocat", "googlebot"}, TeamReviewers: []string{"justice-league", "injustice-league"}})
	if err != nil {
		t.Errorf("PullRequests.RequestReviewers returned error: %v", err)
	}
	want := &PullRequest{Number: Ptr(1)}
	if !cmp.Equal(got, want) {
		t.Errorf("PullRequests.RequestReviewers returned %+v, want %+v", got, want)
	}

	const methodName = "RequestReviewers"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.RequestReviewers(ctx, "o", "r", 1, ReviewersRequest{Reviewers: []string{"octocat", "googlebot"}, TeamReviewers: []string{"justice-league", "injustice-league"}})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRemoveReviewers(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1/requested_reviewers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBody(t, r, `{"reviewers":["octocat","googlebot"],"team_reviewers":["justice-league"]}`+"\n")
	})

	ctx := context.Background()
	_, err := client.PullRequests.RemoveReviewers(ctx, "o", "r", 1, ReviewersRequest{Reviewers: []string{"octocat", "googlebot"}, TeamReviewers: []string{"justice-league"}})
	if err != nil {
		t.Errorf("PullRequests.RemoveReviewers returned error: %v", err)
	}

	const methodName = "RemoveReviewers"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.PullRequests.RemoveReviewers(ctx, "o", "r", 1, ReviewersRequest{Reviewers: []string{"octocat", "googlebot"}, TeamReviewers: []string{"justice-league"}})
	})
}

func TestRemoveReviewers_teamsOnly(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1/requested_reviewers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBody(t, r, `{"reviewers":[],"team_reviewers":["justice-league"]}`+"\n")
	})

	ctx := context.Background()
	_, err := client.PullRequests.RemoveReviewers(ctx, "o", "r", 1, ReviewersRequest{TeamReviewers: []string{"justice-league"}})
	if err != nil {
		t.Errorf("PullRequests.RemoveReviewers returned error: %v", err)
	}

	const methodName = "RemoveReviewers"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.PullRequests.RemoveReviewers(ctx, "o", "r", 1, ReviewersRequest{TeamReviewers: []string{"justice-league"}})
	})
}

func TestListReviewers(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1/requested_reviewers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"users":[{"login":"octocat","id":1}],"teams":[{"id":1,"name":"Justice League"}]}`)
	})

	ctx := context.Background()
	got, _, err := client.PullRequests.ListReviewers(ctx, "o", "r", 1, nil)
	if err != nil {
		t.Errorf("PullRequests.ListReviewers returned error: %v", err)
	}

	want := &Reviewers{
		Users: []*User{
			{
				Login: Ptr("octocat"),
				ID:    Ptr(int64(1)),
			},
		},
		Teams: []*Team{
			{
				ID:   Ptr(int64(1)),
				Name: Ptr("Justice League"),
			},
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("PullRequests.ListReviewers returned %+v, want %+v", got, want)
	}

	const methodName = "ListReviewers"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.ListReviewers(ctx, "o", "r", 1, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestListReviewers_withOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/pulls/1/requested_reviewers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `{}`)
	})

	ctx := context.Background()
	_, _, err := client.PullRequests.ListReviewers(ctx, "o", "r", 1, &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("PullRequests.ListReviewers returned error: %v", err)
	}

	const methodName = "ListReviewers"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PullRequests.ListReviewers(ctx, "\n", "\n", 1, &ListOptions{Page: 2})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PullRequests.ListReviewers(ctx, "o", "r", 1, &ListOptions{Page: 2})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
