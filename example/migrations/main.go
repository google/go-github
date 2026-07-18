// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// migrations demonstrates the functionality of
// the user data migration API for the authenticated GitHub
// user and lists all the user migrations.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v89/github"
)

func fetchAllUserMigrations() ([]*github.UserMigration, error) {
	ctx := context.Background()
	client, err := github.NewClient(github.WithAuthToken("<GITHUB_AUTH_TOKEN>"))
	if err != nil {
		return nil, err
	}

	migrations, _, err := client.Migrations.ListUserMigrations(ctx, &github.ListOptions{Page: 1})
	return migrations, err
}

func main() {
	migrations, err := fetchAllUserMigrations()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for i, m := range migrations {
		fmt.Printf("%v. %v", i+1, m.GetID())
	}
}
