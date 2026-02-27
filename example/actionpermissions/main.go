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

	"github.com/google/go-github/v84/github"
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
	client := github.NewClient(nil).WithAuthToken(token)

	actionsPermissionsRepository, _, err := client.Repositories.GetActionsPermissions(ctx, *owner, *name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current ActionsPermissions %v\n", actionsPermissionsRepository)

	actionsPermissionsRepository = &github.ActionsPermissionsRepository{Enabled: github.Ptr(true), AllowedActions: github.Ptr("selected")}
	_, _, err = client.Repositories.UpdateActionsPermissions(ctx, *owner, *name, *actionsPermissionsRepository)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current ActionsPermissions %v\n", actionsPermissionsRepository)

	actionsAllowed, _, err := client.Repositories.GetActionsAllowed(ctx, *owner, *name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current ActionsAllowed %v\n", actionsAllowed)

	actionsAllowed = &github.ActionsAllowed{GithubOwnedAllowed: github.Ptr(true), VerifiedAllowed: github.Ptr(false), PatternsAllowed: []string{"a/b"}}
	_, _, err = client.Repositories.EditActionsAllowed(ctx, *owner, *name, *actionsAllowed)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current ActionsAllowed %v\n", actionsAllowed)

	actionsPermissionsRepository = &github.ActionsPermissionsRepository{Enabled: github.Ptr(true), AllowedActions: github.Ptr("all")}
	_, _, err = client.Repositories.UpdateActionsPermissions(ctx, *owner, *name, *actionsPermissionsRepository)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current ActionsPermissions %v\n", actionsPermissionsRepository)
}
