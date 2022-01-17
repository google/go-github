package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

// An example of how to use ListEnvironments method with EnvironmentListOptions.
// It's runnable with the following command:
// export GITHUB_TOKEN=your_token export GITHUB_REPOSITORY_OWNER=your_owner export GITHUB_REPOSITORY_NAME=your_repo && go run .
func main() {

	var client *github.Client
	var ctx = context.Background()

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	repo := os.Getenv("GITHUB_REPOSITORY_NAME")
	owner := os.Getenv("GITHUB_REPOSITORY_OWNER")

	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	expectedPageSize := 2

	opts := &github.EnvironmentListOptions{ListOptions: github.ListOptions{PerPage: expectedPageSize}}
	envResponse, _, err := client.Repositories.ListEnvironments(ctx, owner, repo, opts)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if len(envResponse.Environments) != expectedPageSize {
		log.Fatal("Unexpected number of environments returned")
		os.Exit(1)
	}

	// The number of environments here should be equal to expectedPageSize
	fmt.Printf("%d environments returned\n", len(envResponse.Environments))
}
