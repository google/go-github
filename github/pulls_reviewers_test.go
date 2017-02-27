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

type reviewersRequest struct {
	Reviewers []string `json:"reviewers,omitempty"`
}

func TestRequestReviewers(t *testing.T) {
	setup()
	defer teardown()

	logins := []string{"octocat", "googlebot"}

	mux.HandleFunc("/repos/o/r/pulls/1/requested_reviewers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypePullRequestReviewsPreview)
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("TestReviewerRequest couldn't read request body: %v", err)
		}

		reviewers := reviewersRequest{}
		if err := json.Unmarshal(b, &reviewers); err != nil {
			return
		}
		want := reviewersRequest{
			Reviewers: logins,
		}
		if !reflect.DeepEqual(reviewers, want) {
			t.Errorf("PullRequests.RequestReviewers returned %+v, want %+v", reviewers, want)
		}
	})

	// This returns a PR, unmarshalling of which is tested elsewhere
	_, _, err := client.PullRequests.RequestReviewers(context.Background(), "o", "r", 1, logins)
	if err != nil {
		t.Errorf("PullRequests.RequestReviewers returned error: %v", err)
	}
}

func TestRemoveReviewers(t *testing.T) {
	setup()
	defer teardown()
	logins := []string{"octocat", "googlebot"}

	mux.HandleFunc("/repos/o/r/pulls/1/requested_reviewers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypePullRequestReviewsPreview)
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("TestReviewerRequest couldn't read request body: %v", err)
		}

		reviewers := reviewersRequest{}
		if err := json.Unmarshal(b, &reviewers); err != nil {
			return
		}

		want := reviewersRequest{
			Reviewers: logins,
		}
		if !reflect.DeepEqual(reviewers, want) {
			t.Errorf("PullRequests.RemoveReviewers returned %+v, want %+v", reviewers, want)
		}
	})

	_, err := client.PullRequests.RemoveReviewers(context.Background(), "o", "r", 1, logins)
	if err != nil {
		t.Errorf("PullRequests.RequestReviewers returned error: %v", err)
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

	reviewers, _, err := client.PullRequests.ListReviewers(context.Background(), "o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.ListReviewers error: %v", err)
	}

	want := &[]User{
		{
			Login: String("octocat"),
			ID:    Int(1),
		},
	}
	if !reflect.DeepEqual(reviewers, want) {
		t.Errorf("PullRequests.ListReviews returned %+v, want %+v", reviewers, want)
	}
}
