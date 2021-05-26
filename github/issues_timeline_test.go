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
