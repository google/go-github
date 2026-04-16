// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The commitpr command utilizes go-github as a CLI tool for
// pushing files to a branch and creating a pull request from it.
// It takes an auth token as an environment variable and creates
// the commit and the PR under the account affiliated with that token.
//
// The purpose of this example is to show how to use refs, trees and commits to
// create commits and pull requests.
//
// Note, if you want to push a single file, you probably prefer to use the
// content API. An example is available here:
// https://pkg.go.dev/github.com/google/go-github/v84/github#example-RepositoriesService-CreateFile
//
// Note, for this to work at least 1 commit is needed, so you if you use this
// after creating a repository you might want to make sure you set `AutoInit` to
// `true`.
package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/go-github/v84/github"
)

// downloadContents downloads the contents of a file in a repository and returns it as a byte slice.
func downloadContents(ctx context.Context, client *github.Client, owner, repo, path, ref string) ([]byte, error) {
	rc, _, err := client.Repositories.DownloadContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{Ref: ref})
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	by, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Downloaded %v/%v/%v as %d bytes\n", owner, repo, path, len(by))
	return by, nil
}

func main() {
	client := github.NewClient(nil)

	t := []struct {
		owner string
		repo  string
		path  string
		ref   string
	}{
		{"google", "go-github", "README.md", "master"},
		{"github", "rest-api-description", "descriptions/api.github.com/api.github.com.2026-03-10.yaml", "main"},
		{"ScoopInstaller", "Main", "bucket/yq.json", "master"},
	}

	for _, v := range t {
		if _, err := downloadContents(context.Background(), client, v.owner, v.repo, v.path, v.ref); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}
