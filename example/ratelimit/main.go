// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The ratelimit command demonstrates using the github_ratelimit as well as github_pagination.
// By using the waiter, the client automatically sleeps and retry requests
// when it hits secondary rate limits.
// It also prevents the client from abusing the API in case of a primary rate limit.
package main

import (
	"context"
	"fmt"

	"github.com/gofri/go-github-pagination/githubpagination"
	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit"
	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit/github_primary_ratelimit"
	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit/github_secondary_ratelimit"
	"github.com/google/go-github/v84/github"
)

func main() {
	var username string
	fmt.Print("Enter GitHub username: ")
	fmt.Scanf("%s", &username)

	rateLimiter := github_ratelimit.New(nil,
		github_primary_ratelimit.WithLimitDetectedCallback(func(ctx *github_primary_ratelimit.CallbackContext) {
			fmt.Printf("Primary rate limit detected: category %v, reset time: %v\n", ctx.Category, ctx.ResetTime)
		}),
		github_secondary_ratelimit.WithLimitDetectedCallback(func(ctx *github_secondary_ratelimit.CallbackContext) {
			fmt.Printf("Secondary rate limit detected: reset time: %v, total sleep time: %v\n", ctx.ResetTime, ctx.TotalSleepTime)
		}),
	)

	paginator := githubpagination.NewClient(rateLimiter,
		githubpagination.WithPerPage(100), // default to 100 results per page
	)
	client := github.NewClient(paginator)

	// arbitrary usage of the client
	repos, _, err := client.Repositories.ListByUser(context.Background(), username, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, repo := range repos {
		fmt.Printf("%v. %v\n", i+1, repo.GetName())
	}
}
