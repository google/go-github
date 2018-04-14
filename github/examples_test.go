// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github_test

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/github"
)

func ExampleClient_Markdown() {
	client := github.NewClient(nil)

	input := "# heading #\n\nLink to issue #1"
	opt := &github.MarkdownOptions{Mode: "gfm", Context: "google/go-github"}

	output, _, err := client.Markdown(context.Background(), input, opt)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(output)
}

func ExampleRepositoriesService_GetReadme() {
	client := github.NewClient(nil)

	readme, _, err := client.Repositories.GetReadme(context.Background(), "google", "go-github", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	content, err := readme.GetContent()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("google/go-github README:\n%v\n", content)
}

func ExampleRepositoriesService_List() {
	client := github.NewClient(nil)

	user := "willnorris"
	opt := &github.RepositoryListOptions{Type: "owner", Sort: "updated", Direction: "desc"}

	repos, _, err := client.Repositories.List(context.Background(), user, opt)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Recently updated repositories by %q: %v", user, github.Stringify(repos))
}

func ExampleUsersService_ListAll() {
	client := github.NewClient(nil)
	opts := &github.UserListOptions{}
	for {
		users, _, err := client.Users.ListAll(context.Background(), opts)
		if err != nil {
			log.Fatalf("error listing users: %v", err)
		}
		if len(users) == 0 {
			break
		}
		opts.Since = *users[len(users)-1].ID
		// Process users...
	}
}

func ExamplePullRequestsService_Create() {
	// In this example we're creating a PR and displaying the HTML url at the end.

	// Note that authentication is needed here as you are performing a modification
	// so you will need to modify the example to provide an oauth client to
	// github.NewClient() instead of nil. See the following documentation for more
	// information on how to authenticate with the client:
	// https://godoc.org/github.com/google/go-github/github#hdr-Authentication
	client := github.NewClient(nil)

	newPR := &github.NewPullRequest{
		Title:               github.String("My awesome pull request"),
		Head:                github.String("branch_to_merge"),
		Base:                github.String("master"),
		Body:                github.String("This is the description of the PR created with the package `github.com/google/go-github/github`"),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(context.Background(), "myOrganization", "myRepository", newPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
}
