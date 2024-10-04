// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
)

func TestIssueEvents(t *testing.T) {
	events, _, err := client.Issues.ListRepositoryEvents(context.Background(), "google", "go-github", nil)
	if err != nil {
		t.Fatalf("Issues.ListRepositoryEvents returned error: %v", err)
	}

	if len(events) == 0 {
		t.Errorf("ListRepositoryEvents returned no events")
	}

	events, _, err = client.Issues.ListIssueEvents(context.Background(), "google", "go-github", 1, nil)
	if err != nil {
		t.Fatalf("Issues.ListIssueEvents returned error: %v", err)
	}

	if len(events) == 0 {
		t.Errorf("ListIssueEvents returned no events")
	}

	event, _, err := client.Issues.GetEvent(context.Background(), "google", "go-github", *events[0].ID)
	if err != nil {
		t.Fatalf("Issues.GetEvent returned error: %v", err)
	}

	if *event.URL != *events[0].URL {
		t.Fatalf("Issues.GetEvent returned event URL: %v, want %v", *event.URL, *events[0].URL)
	}
}
