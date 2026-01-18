package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v81/github"
)

func main() {
	flag.Parse()
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)

	owner := "google"
	repo := "go-github"
	issue := 2618

	opts := github.IssueListCommentsOptions{
		Sort: github.Ptr("created"),
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 5,
		},
	}

	fmt.Println("Listing comments for issue", issue, "in repository", owner+"/"+repo)

	var paginatedCommentsCount int
	paginatedOpts := opts
	for {
		pageComments, resp, err := client.Issues.ListComments(ctx, owner, repo, issue, &paginatedOpts)
		if err != nil {
			log.Fatalf("ListComments failed: %v", err)
		}
		fmt.Printf("Response: %#+v\n", resp)
		for _, c := range pageComments {
			body := c.GetBody()
			if len(body) > 50 {
				body = body[:50]
			}
			fmt.Printf("Comment: %q\n", body)
		}
		paginatedCommentsCount += len(pageComments)
		if resp.NextPage == 0 {
			break
		}
		paginatedOpts.Page = resp.NextPage
	}
	fmt.Println("Paginated comments:", paginatedCommentsCount)

	var scannedCommentsCount int
	scannedOpts := opts
	for c := range github.MustIter(github.Scan2(func(p github.PaginationOption) ([]*github.IssueComment, *github.Response, error) {
		return client.Issues.ListComments(ctx, owner, repo, issue, &scannedOpts, p)
	})) {
		body := c.GetBody()
		if len(body) > 50 {
			body = body[:50]
		}
		fmt.Printf("Comment: %q\n", body)
		scannedCommentsCount++
	}
	fmt.Println("Scanned comments:", scannedCommentsCount)

	if paginatedCommentsCount != scannedCommentsCount {
		log.Fatalf("Mismatch in comment counts: paginated=%d scanned=%d", paginatedCommentsCount, scannedCommentsCount)
	}
}
