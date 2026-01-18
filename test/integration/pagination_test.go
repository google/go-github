// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build integration

package integration

import (
	"testing"

	"github.com/google/go-github/v81/github"
)

func TestScan2_Offset(t *testing.T) {
	opts := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 5,
		},
	}
	var comments []*github.IssueComment
	for c, err := range github.Scan2(func(p github.PaginationOption) ([]*github.IssueComment, *github.Response, error) {
		return client.Issues.ListComments(t.Context(), "google", "go-github", 526, opts, p)
	}) {
		if err != nil {
			t.Fatalf("Offset scan2 iterator returned error: %v", err)
		}
		comments = append(comments, c)
	}

	if got, want := len(comments), 16; got != want {
		t.Fatalf("Offset scan2 iterator returned %v comments, want %v", got, want)
	}

	if got, want := comments[0].GetID(), int64(274246625); got != want {
		t.Fatalf("Offset scan2 iterator returned first comment ID %v, want %v", got, want)
	}
}
