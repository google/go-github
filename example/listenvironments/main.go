// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// listenvironments is an example of how to use ListEnvironments method with EnvironmentListOptions.
// It's runnable with the following command:
//
//	export GITHUB_TOKEN=your_token
//	export GITHUB_REPOSITORY_OWNER=your_owner
//	export GITHUB_REPOSITORY_NAME=your_repo
//	go run .
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v84/github"
)

func main() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	repo := os.Getenv("GITHUB_REPOSITORY_NAME")
	owner := os.Getenv("GITHUB_REPOSITORY_OWNER")

	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)

	expectedPageSize := 2

	opts := &github.EnvironmentListOptions{ListOptions: github.ListOptions{PerPage: expectedPageSize}}
	envResponse, _, err := client.Repositories.ListEnvironments(ctx, owner, repo, opts)
	if err != nil {
		log.Fatal(err)
	}

	if len(envResponse.Environments) != expectedPageSize {
		log.Fatal("Unexpected number of environments returned")
	}

	// The number of environments here should be equal to expectedPageSize
	fmt.Printf("%v environments returned\n", len(envResponse.Environments))
}
