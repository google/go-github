// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestActivity_Watching(t *testing.T) {
	if !checkAuth("TestActivity_Watching") {
		return
	}

	// first, check if already watching google/go-github
	sub, _, err := client.Activity.GetRepositorySubscription("google", "go-github")
	if err != nil {
		t.Fatalf("Activity.GetRepositorySubscription returned error: %v", err)
	}
	if sub != nil {
		t.Fatalf("Already watching google/go-github.  Please manually stop watching it first.")
	}

	// watch google/go-github
	sub = &github.RepositorySubscription{Subscribed: github.Bool(true)}
	_, _, err = client.Activity.SetRepositorySubscription("google", "go-github", sub)
	if err != nil {
		t.Fatalf("Activity.SetRepositorySubscription returned error: %v", err)
	}

	// check again and verify watching
	sub, _, err = client.Activity.GetRepositorySubscription("google", "go-github")
	if err != nil {
		t.Fatalf("Activity.GetRepositorySubscription returned error: %v", err)
	}
	if sub == nil || !*sub.Subscribed {
		t.Fatalf("Not watching google/go-github after setting subscription.")
	}

	// delete subscription
	_, err = client.Activity.DeleteRepositorySubscription("google", "go-github")
	if err != nil {
		t.Fatalf("Activity.DeleteRepositorySubscription returned error: %v", err)
	}

	// check again and verify not watching
	sub, _, err = client.Activity.GetRepositorySubscription("google", "go-github")
	if err != nil {
		t.Fatalf("Activity.GetRepositorySubscription returned error: %v", err)
	}
	if sub != nil {
		t.Fatalf("Still watching google/go-github after deleting subscription.")
	}
}
