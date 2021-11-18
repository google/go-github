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
	testJSONMarshal(t, &ReviewersRequest{}, "{}")

	u := &ReviewersRequest{
		NodeID:        String("n"),
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
	testJSONMarshal(t, &Reviewers{}, "{}")

	u := &Reviewers{
		Users: []*User{{
			Login:       String("l"),
			ID:          Int64(1),
			AvatarURL:   String("a"),
			GravatarID:  String("g"),
			Name:        String("n"),
			Company:     String("c"),
			Blog:        String("b"),
			Location:    String("l"),
			Email:       String("e"),
			Hireable:    Bool(true),
			PublicRepos: Int(1),
			Followers:   Int(1),
			Following:   Int(1),
			CreatedAt:   &Timestamp{referenceTime},
			URL:         String("u"),
		}},
		Teams: []*Team{{
			ID:              Int64(1),
			NodeID:          String("node"),
			Name:            String("n"),
			Description:     String("d"),
			URL:             String("u"),
			Slug:            String("s"),
			Permission:      String("p"),
			Privacy:         String("priv"),
			MembersCount:    Int(1),
			ReposCount:      Int(1),
			Organization:    nil,
			MembersURL:      String("m"),
			RepositoriesURL: String("r"),
			Parent:          nil,
			LDAPDN:          String("l"),
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/requested_reviewers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"reviewers":["octocat","googlebot"],"team_reviewers":["justice-league","injustice-league"]}`+"\n")
		fmt.Fprint(w, `{"number":1}`)
	})

	// This returns a PR, unmarshalling of which is tested elsewhere
	ctx := context.Background()
	got, _, err := client.PullRequests.RequestReviewers(ctx, "o", "r", 1, ReviewersRequest{Reviewers: []string{"octocat", "googlebot"}, TeamReviewers: []string{"justice-league", "injustice-league"}})
	if err != nil {
		t.Errorf("PullRequests.RequestReviewers returned error: %v", err)
	}
	want := &PullRequest{Number: Int(1)}
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
	client, mux, _, teardown := setup()
	defer teardown()

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

func TestListReviewers(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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
				Login: String("octocat"),
				ID:    Int64(1),
			},
		},
		Teams: []*Team{
			{
				ID:   Int64(1),
				Name: String("Justice League"),
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
	client, mux, _, teardown := setup()
	defer teardown()

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
