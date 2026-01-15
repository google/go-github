// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github_test

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v81/github"
)

func ExampleRepositoriesService_ListByUserIter() {
	client := github.NewClient(nil)
	ctx := context.Background()

	// List all repositories for a user using the iterator.
	// This automatically handles pagination.
	opts := &github.RepositoryListByUserOptions{Type: "public"}
	for repo, err := range client.Repositories.ListByUserIter(ctx, "octocat", opts) {
		if err != nil {
			log.Fatalf("Error listing repos: %v", err)
		}
		fmt.Println(repo.GetName())
	}
}
