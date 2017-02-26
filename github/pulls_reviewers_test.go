// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestRequestReviewers(t *testing.T) {
	setup()
	defer teardown()

	type reviewers struct {
		Reviewers []string `json:"reviewers,omitempty"`
	}
	have := reviewers{}
	mux.HandleFunc("/repos/o/r/pulls/1/requested_reviewers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypePullRequestReviewsPreview)
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("TestReviewerRequest couldn't read request body: %v", err)
		}
		if err := json.Unmarshal(b, &have); err != nil {
			return
		}
	})

	logins := []string{"octocat", "googlebot"}

	// This returns a PR, unmarshalling of which is tested elsewhere
	_, _, err := client.PullRequests.RequestReviewers(context.Background(), "o", "r", 1, logins)
	if err != nil {
		t.Errorf("PullRequests.RequestReviewers returned error: %v", err)
	}

	want := reviewers{
		Reviewers: logins,
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("PullRequests.ListReviews returned %+v, want %+v", have, want)
	}
}

func TestRemoveReviewers(t *testing.T) {
	setup()
	defer teardown()

	type reviewers struct {
		Reviewers []string `json:"reviewers,omitempty"`
	}
	have := reviewers{}
	mux.HandleFunc("/repos/o/r/pulls/1/requested_reviewers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypePullRequestReviewsPreview)
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("TestReviewerRequest couldn't read request body: %v", err)
		}
		if err := json.Unmarshal(b, &have); err != nil {
			return
		}
	})

	logins := []string{"octocat", "googlebot"}

	_, err := client.PullRequests.RemoveReviewers(context.Background(), "o", "r", 1, logins)
	if err != nil {
		t.Errorf("PullRequests.RequestReviewers returned error: %v", err)
	}

	want := reviewers{
		Reviewers: logins,
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("PullRequests.ListReviews returned %+v, want %+v", have, want)
	}
}

func TestListReviewers(t *testing.T) {
	setup()
	defer teardown()

	sampleResponse := `[
  {
    "login": "octocat",
    "id": 1
  }
]`

	mux.HandleFunc("/repos/o/r/pulls/1/requested_reviewers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypePullRequestReviewsPreview)
		fmt.Fprintf(w, sampleResponse)
	})

	// This returns a PR, unmarshalling of which is tested elsewhere
	have, _, err := client.PullRequests.ListReviewers(context.Background(), "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.RequestReviewers returned error: %v", err)
	}
	_login := "octocat"
	_id := 1

	want := []User{
		{
			Login: &_login,
			ID:    &_id,
		},
	}
	if !reflect.DeepEqual(have, &want) {
		t.Errorf("PullRequests.ListReviews returned %+v, want %+v", have, want)
	}
}
