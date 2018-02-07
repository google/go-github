// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The newrepo command utilizes go-github as a cli tool for
// creating new repositories.

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	name := strings.Join(os.Args[1:], "-")
	if token == "" {
		fmt.Println("Unauthorized: No token present")
		return
	}
	if name == "" {
		fmt.Println("No name: New repos must be given a name")
		return
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repo := &github.Repository{Name: github.String(name)}

	repo, _, err := client.Repositories.Create(ctx, "", repo)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Sucsesfully created new repo: ", *repo.Name)
}
