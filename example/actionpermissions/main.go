// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The actionpermissions command utilizes go-github as a cli tool for
// changing GitHub Actions related permission settings for a repository.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

var (
	name  = flag.String("name", "", "repo to change Actions permissions.")
	owner = flag.String("owner", "", "owner of targeted repo.")
)

func main() {
	flag.Parse()
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	if *name == "" {
		log.Fatal("No name: repo name must be given")
	}
	if *owner == "" {
		log.Fatal("No owner: owner of repo must be given")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	actionsPermissionsRepository, _, err := client.Repositories.GetActionsPermissions(ctx, *owner, *name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current ActionsPermissions %s\n", actionsPermissionsRepository.String())

	actionsPermissionsRepository = &github.ActionsPermissionsRepository{Enabled: github.Bool(true), AllowedActions: github.String("selected")}
	_, _, err = client.Repositories.EditActionsPermissions(ctx, *owner, *name, *actionsPermissionsRepository)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current ActionsPermissions %s\n", actionsPermissionsRepository.String())

	actionsAllowed, _, err := client.Repositories.GetActionsAllowed(ctx, *owner, *name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current ActionsAllowed %s\n", actionsAllowed.String())

	actionsAllowed = &github.ActionsAllowed{GithubOwnedAllowed: github.Bool(true), VerifiedAllowed: github.Bool(false), PatternsAllowed: []string{"a/b"}}
	_, _, err = client.Repositories.EditActionsAllowed(ctx, *owner, *name, *actionsAllowed)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current ActionsAllowed %s\n", actionsAllowed.String())

	actionsPermissionsRepository = &github.ActionsPermissionsRepository{Enabled: github.Bool(true), AllowedActions: github.String("all")}
	_, _, err = client.Repositories.EditActionsPermissions(ctx, *owner, *name, *actionsPermissionsRepository)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current ActionsPermissions %s\n", actionsPermissionsRepository.String())
}
