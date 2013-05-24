package main

import (
	"fmt"

	"github.com/google/go-github"
)

func main() {
	client := github.NewClient(nil)

	fmt.Println("Recently updated repositories owned by user willnorris:")
	opt := &github.RepositoryListOptions{Type: "owner", Sort: "updated", Direction: "desc"}
	repos, err := client.Repositories.List("willnorris", opt)
	if err != nil {
		fmt.Printf("error: %v\n\n", err)
	} else {
		fmt.Printf("%#v\n\n", repos)
	}

	rate, err := client.RateLimit()
	if err != nil {
		fmt.Printf("Error fetching rate limit: %#v\n\n", err)
		return
	}

	fmt.Printf("API Rate Limit: %#v\n\n", rate)
}
