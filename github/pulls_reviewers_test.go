// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PullRequests.RequestReviewers returned %+v, want %+v", got, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.PullRequests.RequestReviewers(ctx, "o", "r", 1, ReviewersRequest{Reviewers: []string{"octocat", "googlebot"}, TeamReviewers: []string{"justice-league", "injustice-league"}})
	if got != nil {
		t.Errorf("client.BaseURL.Path='' RequestReviewers = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' RequestReviewers resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' RequestReviewers err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.PullRequests.RequestReviewers(ctx, "o", "r", 1, ReviewersRequest{Reviewers: []string{"octocat", "googlebot"}, TeamReviewers: []string{"justice-league", "injustice-league"}})
	if got != nil {
		t.Errorf("rate.Reset.Time > now RequestReviewers = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now RequestReviewers resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now RequestReviewers err = nil, want error")
	}
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

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	resp, err := client.PullRequests.RemoveReviewers(ctx, "o", "r", 1, ReviewersRequest{Reviewers: []string{"octocat", "googlebot"}, TeamReviewers: []string{"justice-league"}})
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' RemoveReviewers resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' RemoveReviewers err = nil, want error")
	}
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PullRequests.ListReviewers returned %+v, want %+v", got, want)
	}

	// Test s.client.NewRequest failure
	client.BaseURL.Path = ""
	got, resp, err := client.PullRequests.ListReviewers(ctx, "o", "r", 1, nil)
	if got != nil {
		t.Errorf("client.BaseURL.Path='' ListReviewers = %#v, want nil", got)
	}
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' ListReviewers resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("client.BaseURL.Path='' ListReviewers err = nil, want error")
	}

	// Test s.client.Do failure
	client.BaseURL.Path = "/api-v3/"
	client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	got, resp, err = client.PullRequests.ListReviewers(ctx, "o", "r", 1, nil)
	if got != nil {
		t.Errorf("rate.Reset.Time > now ListReviewers = %#v, want nil", got)
	}
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		t.Errorf("rate.Reset.Time > now ListReviewers resp = %#v, want StatusCode=%v", resp.Response, want)
	}
	if err == nil {
		t.Error("rate.Reset.Time > now ListReviewers err = nil, want error")
	}
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

	_, _, err := client.PullRequests.ListReviewers(context.Background(), "o", "r", 1, &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("PullRequests.ListReviewers returned error: %v", err)
	}

	// Test addOptions failure
	_, _, err = client.PullRequests.ListReviewers(context.Background(), "\n", "\n", 1, &ListOptions{Page: 2})
	if err == nil {
		t.Error("bad options ListReviewers err = nil, want error")
	}

}
