// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestPullRequestsService_ListReviews(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypePullRequestReviewsPreview)
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	reviews, _, err := client.PullRequests.ListReviews("o", "r", 1)
	if err != nil {
		t.Errorf("PullRequests.ListReviews returned error: %v", err)
	}

	want := []*PullRequestReview{
		{ID: Int(1)},
		{ID: Int(2)},
	}
	if !reflect.DeepEqual(reviews, want) {
		t.Errorf("PullRequests.ListReviews returned %+v, want %+v", reviews, want)
	}
}

func TestPullRequestsService_GetReview(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypePullRequestReviewsPreview)
		fmt.Fprint(w, `{"id":1}`)
	})

	review, _, err := client.PullRequests.GetReview("o", "r", 1, 1)
	if err != nil {
		t.Errorf("PullRequests.GetReview returned error: %v", err)
	}

	want := &PullRequestReview{ID: Int(1)}
	if !reflect.DeepEqual(review, want) {
		t.Errorf("PullRequests.GetReview returned %+v, want %+v", review, want)
	}
}

func TestPullRequestsService_ListReviewComments(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/pulls/1/reviews/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypePullRequestReviewsPreview)
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})

	comments, _, err := client.PullRequests.ListReviewComments("o", "r", 1, 1)
	if err != nil {
		t.Errorf("PullRequests.ListReviewComments returned error: %v", err)
	}

	want := []*PullRequestReviewComment{
		{ID: Int(1)},
		{ID: Int(2)},
	}
	if !reflect.DeepEqual(comments, want) {
		t.Errorf("PullRequests.ListReviewComments returned %+v, want %+v", comments, want)
	}
}
