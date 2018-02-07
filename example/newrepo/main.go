// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The newrepo command utilizes go-github as a cli tool for
// creating new repositories. It takes an auth token as
// an enviroment variable and creates the new repo under
// the account affiliated with that token.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	name := strings.Join(os.Args[1:], "-")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	if name == "" {
		log.Fatal("No name: New repos must be given a name")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	r := &github.Repository{Name: github.String(name)}

	repo, _, err := client.Repositories.Create(ctx, "", r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully created new repo: %v\n", repo.GetName())
}
