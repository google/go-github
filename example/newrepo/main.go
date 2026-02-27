// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The newrepo command utilizes go-github as a cli tool for
// creating new repositories. It takes an auth token as
// an environment variable and creates the new repo under
// the account affiliated with that token.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v84/github"
)

var (
	name        = flag.String("name", "", "Name of repo to create in authenticated user's GitHub account.")
	description = flag.String("description", "", "Description of created repo.")
	private     = flag.Bool("private", false, "Will created repo be private.")
	autoInit    = flag.Bool("auto-init", false, "Pass true to create an initial commit with empty README.")
)

func main() {
	flag.Parse()
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	if *name == "" {
		log.Fatal("No name: New repos must be given a name")
	}
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)

	r := &github.Repository{Name: name, Private: private, Description: description, AutoInit: autoInit}
	repo, _, err := client.Repositories.Create(ctx, "", r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully created new repo: %v\n", repo.GetName())
}
