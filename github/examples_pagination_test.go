// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github_test

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/google/go-github/v81/github"
)

func ExampleIssuesService_ListComments_offset_pagination_scan() {
	client := github.NewClient(nil)
	ctx := context.Background()
	opts := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 5,
		},
	}

	it, hasErr := github.Scan(func(p github.PaginationOption) ([]*github.IssueComment, *github.Response, error) {
		return client.Issues.ListComments(ctx, "google", "go-github", 526, opts, p)
	})

	comments := slices.Collect(it)
	if err := hasErr(); err != nil {
		log.Fatalf("Scan iterator returned error: %v", err)
	}

	fmt.Println("Total comments:", len(comments))
}

func ExampleIssuesService_ListComments_offset_pagination_scan2() {
	client := github.NewClient(nil)
	ctx := context.Background()
	opts := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 5,
		},
	}

	var comments []*github.IssueComment
	for c, err := range github.Scan2(func(p github.PaginationOption) ([]*github.IssueComment, *github.Response, error) {
		return client.Issues.ListComments(ctx, "google", "go-github", 526, opts, p)
	}) {
		if err != nil {
			log.Fatalf("Scan2 iterator returned error: %v", err)
		}
		comments = append(comments, c)
	}

	fmt.Println("Total comments:", len(comments))
}

func ExampleIssuesService_ListComments_offset_pagination_scanAndCollect() {
	client := github.NewClient(nil)
	ctx := context.Background()
	opts := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 5,
		},
	}

	comments, err := github.ScanAndCollect(func(p github.PaginationOption) ([]*github.IssueComment, *github.Response, error) {
		return client.Issues.ListComments(ctx, "google", "go-github", 526, opts, p)
	})
	if err != nil {
		log.Fatalf("ScanAndCollect returned error: %v", err)
	}

	fmt.Println("Total comments:", len(comments))
}

func ExampleIssuesService_ListComments_offset_pagination_scan2MustIter() {
	client := github.NewClient(nil)
	ctx := context.Background()
	opts := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 5,
		},
	}

	var comments []*github.IssueComment
	for c := range github.MustIter(github.Scan2(func(p github.PaginationOption) ([]*github.IssueComment, *github.Response, error) {
		return client.Issues.ListComments(ctx, "google", "go-github", 526, opts, p)
	})) {
		comments = append(comments, c)
	}

	fmt.Println("Total comments:", len(comments))
}

func ExampleSecurityAdvisoriesService_ListRepositorySecurityAdvisoriesForOrg_after_pagination_scan2MustIter() {
	client := github.NewClient(nil)
	ctx := context.Background()

	opts := &github.ListRepositorySecurityAdvisoriesOptions{
		ListCursorOptions: github.ListCursorOptions{
			PerPage: 1,
		},
	}

	var advisories []*github.SecurityAdvisory
	for a := range github.MustIter(github.Scan2(func(p github.PaginationOption) ([]*github.SecurityAdvisory, *github.Response, error) {
		return client.SecurityAdvisories.ListRepositorySecurityAdvisoriesForOrg(ctx, "alexandear-org", opts, p)
	})) {
		advisories = append(advisories, a)
	}

	fmt.Println("Total advisories:", len(advisories))
}
