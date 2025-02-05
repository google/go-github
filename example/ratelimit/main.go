// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The ratelimit command demonstrates using the github_ratelimit.SecondaryRateLimitWaiter.
// By using the waiter, the client automatically sleeps and retry requests
// when it hits secondary rate limits.
package main

import (
	"context"
	"fmt"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v69/github"
)

func main() {
	var username string
	fmt.Print("Enter GitHub username: ")
	fmt.Scanf("%s", &username)

	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	client := github.NewClient(rateLimiter)

	// arbitrary usage of the client
	organizations, _, err := client.Organizations.List(context.Background(), username, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, organization := range organizations {
		fmt.Printf("%v. %v\n", i+1, organization.GetLogin())
	}
}
