// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// iterpagination is an example of how to use pagination with Go-native iterators.
// It's runnable with the following command:
//
//	export GITHUB_AUTH_TOKEN=your_token
//	export GITHUB_REPOSITORY_OWNER=your_owner
//	export GITHUB_REPOSITORY_NAME=your_repo
//	export GITHUB_REPOSITORY_ISSUE=your_issue
//	go run .
package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/go-github/v81/github"
)

func main() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	owner := cmp.Or(os.Getenv("GITHUB_REPOSITORY_OWNER"), "google")
	repo := cmp.Or(os.Getenv("GITHUB_REPOSITORY_NAME"), "go-github")
	issue, _ := strconv.Atoi(os.Getenv("GITHUB_REPOSITORY_ISSUE"))
	issue = cmp.Or(issue, 2618)

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)

	opts := github.IssueListCommentsOptions{
		Sort: github.Ptr("created"),
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 5,
		},
	}

	fmt.Println("Listing comments for issue", issue, "in repository", owner+"/"+repo)

	scannedOpts := opts
	for c := range github.MustIter(github.Scan2(func(p github.PaginationOption) ([]*github.IssueComment, *github.Response, error) {
		return client.Issues.ListComments(ctx, owner, repo, issue, &scannedOpts, p)
	})) {
		body := c.GetBody()
		if len(body) > 50 {
			body = body[:50]
		}
		fmt.Printf("Comment: %q\n", body)
	}
}
