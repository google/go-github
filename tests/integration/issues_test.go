// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import "testing"

func TestIssueEvents(t *testing.T) {
	events, _, err := client.Issues.ListRepositoryEvents("google", "go-github", nil)
	if err != nil {
		t.Fatalf("Issues.ListRepositoryEvents returned error: %v", err)
	}

	if len(events) == 0 {
		t.Errorf("ListRepositoryEvents returned no events")
	}

	events, _, err = client.Issues.ListIssueEvents("google", "go-github", 1, nil)
	if err != nil {
		t.Fatalf("Issues.ListIssueEvents returned error: %v", err)
	}

	if len(events) == 0 {
		t.Errorf("ListIssueEvents returned no events")
	}

	event, _, err := client.Issues.GetEvent("google", "go-github", *events[0].ID)
	if err != nil {
		t.Fatalf("Issues.GetEvent returned error: %v", err)
	}

	if *event.URL != *events[0].URL {
		t.Fatalf("Issues.GetEvent returned event URL: %v, want %v", *event.URL, *events[0].URL)
	}
}

func TestGetIssueByURL(t *testing.T) {
	i, _, err := client.Issues.GetByURL("https://api.github.com/repos/google/go-github/issues/1")

	if err != nil {
		t.Fatalf("Issues.GetByURL returned error: %v", err)
	}

	if *i.Number != 1 {
		t.Errorf("expected Issues.GetByURL to return a representation of the issue")
	}
}

func TestGetIssueCommentByURL(t *testing.T) {
	c, _, err := client.Issues.GetCommentByURL("https://api.github.com/repos/google/go-github/issues/comments/19136005")

	if err != nil {
		t.Fatalf("Issues.GetCommentByURL returned error: %v", err)
	}

	if *c.ID != 19136005 {
		t.Errorf("expected Issues.GetCommentByURL to return a representation of the specified comment")
	}
}
