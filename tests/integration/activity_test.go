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
	sub, _, err := client.Activity.GetSubscription("google", "go-github")
	if err != nil {
		t.Fatalf("Activity.GetSubscription returned error: %v", err)
	}
	if sub != nil {
		t.Fatalf("Already watching google/go-github.  Please manually stop watching it first.")
	}

	// watch google/go-github
	sub = &github.Subscription{Subscribed: github.Bool(true)}
	_, _, err = client.Activity.SetSubscription("google", "go-github", sub)
	if err != nil {
		t.Fatalf("Activity.SetSubscription returned error: %v", err)
	}

	// check again and verify watching
	sub, _, err = client.Activity.GetSubscription("google", "go-github")
	if err != nil {
		t.Fatalf("Activity.GetSubscription returned error: %v", err)
	}
	if sub == nil || !*sub.Subscribed {
		t.Fatalf("Not watching google/go-github after setting subscription.")
	}

	// delete subscription
	_, err = client.Activity.DeleteSubscription("google", "go-github")
	if err != nil {
		t.Fatalf("Activity.DeleteSubscription returned error: %v", err)
	}

	// check again and verify not watching
	sub, _, err = client.Activity.GetSubscription("google", "go-github")
	if err != nil {
		t.Fatalf("Activity.GetSubscription returned error: %v", err)
	}
	if sub != nil {
		t.Fatalf("Still watching google/go-github after deleting subscription.")
	}
}
